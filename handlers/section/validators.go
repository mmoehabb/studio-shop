package section

import "strconv"

func ValidateAddSectionBody(body *AddSectionBody) (bool, map[string]string) {
	ok := true
	errs := make(map[string]string)
	if body.Title == "" {
		errs["title"] = "Title cannot be empty."
		ok = false
	}
	if body.Parent != "" && body.Parent != "none" {
		_, err := strconv.Atoi(body.Parent)
		if err != nil {
			errs["parent"] = "Something wrong with parent data."
			ok = false
		}
	}
	return ok, errs
}
