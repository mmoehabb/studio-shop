package database

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/drive/v3"

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

var dirs []*drive.File
var images []*drive.File

func collectData(files []*drive.File) {
	for _, file := range files {
		if file.MimeType == "application/vnd.google-apps.folder" {
			dirs = append(dirs, file)
		} else if file.MimeType == "image/jpeg" {
			images = append(images, file)
		}
	}
}

func Reseed(c *fiber.Ctx) error {
	defer anc.Recover(c)
	anc.Must(nil, db.Reseed())
	dirs = nil
	images = nil

	service := anc.GetDriveService()
	go reseedFunc(service)

	return c.SendString("Database is reseeding now...")
}

func reseedFunc(service *drive.Service) {
	/* Context Structure
	  <digit> <- 1, 2, 3, ... 9
	  <name> <- any string

	  <dir-prefix> <- <digit>
	  <dir-prefix> <- <digit>.<dir-prefix>
	  <dir-name> <- <dir-prefix> <name>

	  <file-prefix> <- <digit>.
	  <file-prefix> <- <digit><file-prefix>
	  <file-name> <- <file-prefix> <name>

	* Examples
	  1. main-dir
	  2. sec-dir
	  1.1 inner-main-dir
	  2.1 inner-sec-dir
	  2.1.1 example-file
	*/
	driveRes := anc.Must(service.Files.List().Do()).(*drive.FileList)

	for driveRes.NextPageToken != "" {
		collectData(driveRes.Files)
		driveRes = anc.Must(service.Files.List().PageToken(driveRes.NextPageToken).Do()).(*drive.FileList)
	}
	collectData(driveRes.Files)

	var prefixNameMap = make(map[string]string)

	// DONE: insert sections into the Database
  log.Println("Reseeding sections...")
	for _, file := range dirs {
		var nameSlice = strings.SplitN(file.Name, " ", 2)
		if len(nameSlice) < 2 {
			continue
		}
		var dirPrefix = nameSlice[0]

		prefixParts := strings.Split(dirPrefix, ".")
		invalidDir := false
		for _, digit := range prefixParts {
			if _, err := strconv.Atoi(digit); err != nil {
				invalidDir = true
				break
			}
		}

		if invalidDir == true {
			continue
		}

		prefixNameMap[dirPrefix] = file.Name
		newSection := sections.DataModel{Title: file.Name}
		sections.Add([]sections.DataModel{newSection})
	}

	// DONE: insert relations into the Database
  log.Println("Reseeding relations...")
	for prefix, name := range prefixNameMap {
		prefixParts := strings.Split(prefix, ".")
		if len(prefixParts) < 2 {
			continue
		}

		var parentPrefix = strings.Join(prefixParts[0:len(prefixParts)-1], ".")
		if prefixNameMap[parentPrefix] == "" {
			continue
		}

		parentId := anc.Must(sections.GetId(prefixNameMap[parentPrefix])).(int)
		childId := anc.Must(sections.GetId(name)).(int)

		newRelation := relations.DataModel{
			Parent: parentId,
			Child:  childId,
		}
		relations.Add([]relations.DataModel{newRelation})
	}

	// DONE: insert photos into the Database
  log.Println("Reseeding photos...")
	for _, file := range images {
		var nameSlice = strings.SplitN(file.Name, " ", 2)
		if len(nameSlice) < 2 {
			continue
		}
		var prefix = nameSlice[0]

		invalid := false
		prefixParts := strings.Split(prefix, ".")
		for _, digit := range prefixParts {
			if _, err := strconv.Atoi(digit); err != nil {
				invalid = true
				break
			}
		}
		if invalid == true {
			continue
		}

		sectionPrefix := strings.Join(prefixParts[0:len(prefixParts)-1], ".")
		if prefixNameMap[sectionPrefix] == "" {
			log.Println("Bad File: ", file.Name)
		}

		parentId := anc.Must(sections.GetId(prefixNameMap[sectionPrefix])).(int)
		newPhoto := photos.DataModel{
			Name:      file.Name,
			Url:       "https://lh3.googleusercontent.com/d/" + file.Id,
			SectionId: parentId,
		}
		photos.Add([]photos.DataModel{newPhoto})
	}

  log.Println("Done.")
}
