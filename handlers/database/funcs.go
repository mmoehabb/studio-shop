package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v2"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db"
	"github.com/mmoehabb/studio-shop/db/photos"
	"github.com/mmoehabb/studio-shop/db/relations"
	"github.com/mmoehabb/studio-shop/db/sections"
)

func Seed(c *fiber.Ctx) error {
	defer anc.Recover(c)
	anc.Must(nil, db.Seed())
	return c.SendString("Database has been seeded.")
}

var dirs map[string]string
var images map[string]string
var sectionIdMap map[string]int

var status string = "idle"
var progress int
var totalProgress int

func Reseed(c *fiber.Ctx) error {
	defer anc.Recover(c)
	c.Response().Header.Add(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	anc.Must(nil, db.Reseed())
	if status != "idle" && status != "done" {
		return c.SendString("<div>Reseeding is still in progress... <a href='/reseed/status'>View Status<a></div>")
	}

	progress = 0
	dirs = make(map[string]string)
	images = make(map[string]string)
	sectionIdMap = make(map[string]int)

	go reseedFunc(anc.S3)
	return c.SendString("<div>Database is reseeding now... <a href='/reseed/status'>View Status<a></div>")
}

func reseedFunc(bucket anc.S3Bucket) {
	status = "retrieving data..."
	objs := anc.Must(bucket.ListObjects(context.Background())).([]types.Object)

	status = "collecting data..."
	collectData(objs)

	status = "reseeding sections..."
	newSections := []sections.DataModel{}
	for dir := range dirs {
		newSections = append(newSections, sections.DataModel{Title: dir})
		progress += 1
	}

	if err := sections.Add(newSections); err != nil {
		log.Println("couldn't add sections!")
		log.Println("error:", err)
	}

	status = "reseeding relations..."
	newRelations := []relations.DataModel{}
	for child, parent := range dirs {
		c := anc.Must(sections.GetId(child)).(int)
		sectionIdMap[child] = c

		if parent == "" {
			continue
		}
		p := anc.Must(sections.GetId(parent)).(int)
		sectionIdMap[parent] = p

		newRelations = append(newRelations, relations.DataModel{
			Parent: p,
			Child:  c,
		})
	}

	if err := relations.Add(newRelations); err != nil {
		log.Println("couldn't add relations!")
		log.Println("error:", err)
	}

	status = "reseeding photos..."
	newPhotos := []photos.DataModel{}
	for child, parent := range images {
		if parent == "" {
			continue
		}
		subdirs := strings.Split(parent, "/")
		newPhotos = append(newPhotos, photos.DataModel{
			Name:      child,
			Url:       parent + "/" + child,
			SectionId: sectionIdMap[subdirs[len(subdirs)-1]],
		})
		progress += 1
	}

	if err := photos.Add(newPhotos); err != nil {
		log.Println("couldn't add photos!")
		log.Println("error:", err)
	}

	status = "done"
}

func collectData(objs []types.Object) {
	for _, obj := range objs {
		key := *obj.Key
		isDir := key[len(key)-1] == '/'
		if isDir == true {
			parts := strings.Split(key, "/")
			for i, part := range parts {
				if part != "" && i > 0 {
					dirs[part] = parts[i-1]
					continue
				}
				if part != "" {
					dirs[part] = ""
				}
			}
		} else {
			parts := strings.Split(key, "/")
			folders := parts[:len(parts)-1]
			img := parts[len(parts)-1]
			for i, folder := range folders {
				if folder != "" && i > 0 {
					dirs[folder] = folders[i-1]
					continue
				}
				if folder != "" {
					dirs[folder] = ""
				}
			}
			images[img] = strings.Join(folders, "/")
		}
	}

	totalProgress = len(dirs) + len(images)
}

func ReseedStatus(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf(
		"status: %s\nprogress: %d / %d\n",
		status, progress, totalProgress,
	))
}
