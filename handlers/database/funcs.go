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
var images map[string]*drive.File

func saveCache() {
  dirsJsonStr := anc.Must(json.Marshal(dirs)).([]byte)
  imagesJsonStr := anc.Must(json.Marshal(images)).([]byte)
  os.WriteFile("./dirs.json", dirsJsonStr, os.ModePerm)
  os.WriteFile("./images.json", imagesJsonStr, os.ModePerm)
}

func loadCache() {
  if _, err := os.Stat("./dirs.json"); err == nil {
    dirsCache := anc.Must(os.ReadFile("./dirs.json")).([]byte)
    json.Unmarshal(dirsCache, &dirs)
  }
  if _, err := os.Stat("./images.json"); err == nil {
    imagesCache := anc.Must(os.ReadFile("./images.json")).([]byte)
    json.Unmarshal(imagesCache, &images)
  }
}

var status string = "Idle"
var progress int
var totalProgress int

func Reseed(c *fiber.Ctx) error {
	defer anc.Recover(c)
  if progress != totalProgress {
    c.Response().Header.Add(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.SendString("<div>Reseeding is still in progress... <a href='/reseed/status'>View Status<a></div>")
  }
	anc.Must(nil, db.Reseed())
	dirs = make(map[string]*drive.File)
	images = make(map[string]*drive.File)

	service := anc.GetDriveService()
	go reseedFunc(service)

	return c.SendString("Database is reseeding now...")
}

func ReseedReset(c *fiber.Ctx) error {
  status = "Idle"
  progress = 0
  totalProgress = 0
  return c.SendString("Reseed values have been reseted.")
}

func ReseedStatus(c *fiber.Ctx) error {
  return c.SendString(fmt.Sprintf("status: %s\n progress: %d / %d\n", status, progress, totalProgress))
}

func collectData(files []*drive.File) {
	for _, file := range files {
    if dirs[file.Name] != nil || images[file.Name] != nil {
      break
    }
		if file.MimeType == "application/vnd.google-apps.folder" {
      dirs[file.Name] = file
		} else if file.MimeType == "image/jpeg" {
      images[file.Name] = file
		}
	}
  totalProgress = len(dirs) + len(images)
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
  loadCache()
	driveRes := anc.Must(service.Files.List().
    PageSize(1000).
    OrderBy("createdTime").
    Do()).(*drive.FileList)

	for driveRes.NextPageToken != "" {
		collectData(driveRes.Files)
		driveRes = anc.Must(service.Files.List().
      PageSize(1000).
      OrderBy("createdTime").
      PageToken(driveRes.NextPageToken).
      Do()).(*drive.FileList)
	}
	collectData(driveRes.Files)

	var prefixNameMap = make(map[string]string)

  progress = 0
  status = "in progress..."
	// DONE: insert sections into the Database
  log.Println("Reseeding sections...")
	for _, file := range dirs {
		var nameSlice = strings.SplitN(file.Name, " ", 2)
		if len(nameSlice) < 2 {
      totalProgress -= 1
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
			continue
		}

		prefixNameMap[dirPrefix] = file.Name
		newSection := sections.DataModel{Title: file.Name}
		sections.Add([]sections.DataModel{newSection})
    progress += 1
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
      totalProgress -= 1
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
    progress += 1
	}

  log.Println("Done.")
  status = "done"
  saveCache()
}
