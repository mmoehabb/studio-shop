package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

var dirs map[string]*drive.File
var new_dirs []string
var images map[string]*drive.File
var new_images []string

func saveCache() {
	dirsJsonStr := anc.Must(json.Marshal(dirs)).([]byte)
	imagesJsonStr := anc.Must(json.Marshal(images)).([]byte)
	os.WriteFile("./dirs.json", dirsJsonStr, os.ModePerm)
	os.WriteFile("./images.json", imagesJsonStr, os.ModePerm)
	dirs = nil
	images = nil
	new_dirs = nil
	new_images = nil
}

func loadCache() {
	if _, err := os.Stat("./dirs.json"); err == nil {
		dirsCache := anc.Must(os.ReadFile("./dirs.json")).([]byte)
		json.Unmarshal(dirsCache, &dirs)
	} else {
		dirs = make(map[string]*drive.File)
	}
	if _, err := os.Stat("./images.json"); err == nil {
		imagesCache := anc.Must(os.ReadFile("./images.json")).([]byte)
		json.Unmarshal(imagesCache, &images)
	} else {
		images = make(map[string]*drive.File)
	}
}

var status string = "idle"
var progress int
var totalProgress int
var drivePage int

func Reseed(c *fiber.Ctx) error {
	defer anc.Recover(c)
	c.Response().Header.Add(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	if status != "idle" && status != "done" {
		return c.SendString("<div>Reseeding is still in progress... <a href='/reseed/status'>View Status<a></div>")
	}

	service := anc.GetDriveService()
	go reseedFunc(service)

	return c.SendString("<div>Database is reseeding now... <a href='/reseed/status'>View Status<a></div>")
}

func ReseedReset(c *fiber.Ctx) error {
	status = "idle"
	progress = 0
	totalProgress = 0
	anc.Must(nil, db.Reseed())
	return c.SendString("Reseed values have been reseted.")
}

func ReseedStatus(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf(
		"status: %s\nprogress: %d / %d\ndrivePage: %d\n",
		status, progress, totalProgress, drivePage,
	))
}

func collectData(files []*drive.File) {
	for _, file := range files {
		if dirs[file.Name] != nil || images[file.Name] != nil {
			drivePage = -1 // this breaks the outer for loop
			break
		}
		if file.MimeType == "application/vnd.google-apps.folder" {
			dirs[file.Name] = file
			new_dirs = append(new_dirs, file.Name)
		} else if file.MimeType == "image/jpeg" {
			images[file.Name] = file
			new_images = append(new_images, file.Name)
		}
	}
	totalProgress = len(new_dirs) + len(new_images)
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
	status = "loading data from cache..."
	loadCache()

	status = "retrieving data from drive..."
	drivePage = 0
	var driveRes = anc.Must(service.Files.List().
		Fields("nextPageToken", "files(id,name,mimeType)").
		OrderBy("mimeType desc").
		OrderBy("name desc").
		PageSize(1000).
		Do()).(*drive.FileList)

	for driveRes.NextPageToken != "" {
		collectData(driveRes.Files)

		if drivePage < 0 {
			break
		}

		driveRes = anc.Must(service.Files.List().
			Fields("nextPageToken", "files(id,name,mimeType)").
			OrderBy("mimeType desc").
			OrderBy("name desc").
			PageSize(1000).
			PageToken(driveRes.NextPageToken).
			Do()).(*drive.FileList)

		drivePage += 1
	}
	collectData(driveRes.Files)

	var prefixNameMap = make(map[string]string)

	progress = 0
	status = "reseeding data..."
	// DONE: insert sections into the Database
	log.Println("Reseeding sections...")
	for _, filename := range new_dirs {
		filename = strings.Trim(filename, " ")
		var nameSlice = strings.SplitN(filename, " ", 2)
		if len(nameSlice) < 2 {
			totalProgress -= 1
			log.Println(filename, "has been skipped!")
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
			totalProgress -= 1
			log.Println(filename, "has been skipped!")
			continue
		}

		prefixNameMap[dirPrefix] = filename
		newSection := sections.DataModel{Title: filename}
		sections.Add([]sections.DataModel{newSection})
		progress += 1
	}

	// DONE: insert relations into the Database
	log.Println("Reseeding relations...")
	for prefix, name := range prefixNameMap {
		prefixParts := strings.Split(prefix, ".")
		if len(prefixParts) < 2 {
			log.Println(prefix, name, "has been skipped!")
			continue
		}

		var parentPrefix = strings.Join(prefixParts[0:len(prefixParts)-1], ".")
		if prefixNameMap[parentPrefix] == "" {
			log.Println(prefix, name, "has been skipped!")
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
	for _, filename := range new_images {
		filename = strings.Trim(filename, " ")
		var nameSlice = strings.SplitN(filename, " ", 2)
		if len(nameSlice) < 2 {
			totalProgress -= 1
			log.Println(filename, "has been skipped!")
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
			totalProgress -= 1
			log.Println(filename, "has been skipped!")
			continue
		}

		sectionPrefix := strings.Join(prefixParts[0:len(prefixParts)-1], ".")
		if prefixNameMap[sectionPrefix] == "" {
			log.Println(filename, "has been skipped!")
			log.Println("SectionPrefix:", sectionPrefix)
		}

		parentId := anc.Must(sections.GetId(prefixNameMap[sectionPrefix])).(int)
		newPhoto := photos.DataModel{
			Name:      filename,
			Url:       "https://lh3.googleusercontent.com/d/" + images[filename].Id,
			SectionId: parentId,
		}
		photos.Add([]photos.DataModel{newPhoto})
		progress += 1
	}

	log.Println("Done.")
	status = "saving data in cache..."
	saveCache()
	status = "done"
}
