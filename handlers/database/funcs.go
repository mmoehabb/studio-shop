package database

import (
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

func Reseed(c *fiber.Ctx) error {
  defer anc.Recover(c)
  anc.Must(nil, db.Reseed())

  service := anc.GetDriveService()
  driveRes := anc.Must(service.Files.List().Do()).(*drive.FileList)

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
    2.1.1. example-file
  */
  
  var prefixNameMap = make(map[string]string)

  // DONE: insert sections into the Database
  for _, file := range driveRes.Files {
    if file.MimeType == "application/vnd.google-apps.folder" {
      var nameSlice = strings.SplitN(file.Name, " ", 2)
      if len(nameSlice) < 2 {
        continue
      }
      var dirPrefix = nameSlice[0]
      var dirName = nameSlice[1] 

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

      prefixNameMap[dirPrefix] = dirName
      newSection := sections.DataModel{ Title: dirName }
      sections.Add([]sections.DataModel{ newSection })
    }
  }

  // DONE: insert relations into the Database
  for prefix, name := range prefixNameMap {
    if len(prefix) < 3 {
      continue
    }
    var parentPrefix = prefix[0:len(prefix)-2]
    if prefixNameMap[parentPrefix] == "" {
      continue
    }

    parentId := anc.Must(sections.GetId(prefixNameMap[parentPrefix])).(int)
    childId := anc.Must(sections.GetId(name)).(int)

    newRelation := relations.DataModel{
      Parent: parentId,
      Child: childId,
    }
    relations.Add([]relations.DataModel{ newRelation })
  }

  // DONE: insert photos into the Database
  for _, file := range driveRes.Files {
    if file.MimeType == "image/jpeg" {
      var nameSlice = strings.SplitN(file.Name, " ", 2)
      if len(nameSlice) < 2 {
        continue
      }
      var prefix = nameSlice[0]
      var name = nameSlice[1] 

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

      sectionPrefix := prefix[0:len(prefix)-2]
      if prefixNameMap[sectionPrefix] == "" {
        return c.SendStatus(fiber.StatusBadRequest)
      }

      parentId := anc.Must(sections.GetId(prefixNameMap[sectionPrefix])).(int)
      newPhoto := photos.DataModel{ 
        Name: name, 
        Url: "https://lh3.googleusercontent.com/d/" + file.Id,
        SectionId: parentId, 
      }
      photos.Add([]photos.DataModel{ newPhoto })
    }
  }

  return c.SendString("Database has been seeded.")
}
