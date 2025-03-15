package photo

import "strconv"

func ValidateAddPhotoBody(body *AddPhotoBody) (bool, map[string]string) {
	ok := true
	errs := make(map[string]string)
	if body.Name == "" {
		errs["name"] = "Name cannot be empty."
		ok = false
	}
	if body.SectionId == "none" {
		errs["section-id"] = "Section Id cannot be none."
		ok = false
	} else if body.SectionId == "" {
		errs["section-id"] = "Section Id cannot be empty."
		ok = false
	} else if _, err := strconv.Atoi(body.SectionId); err != nil {
		errs["section-id"] = "Something wrong with section id."
		ok = false
	}
	return ok, errs
}
