package section

type AddSectionBody struct {
	Title  string `json:"title" xml:"title" form:"title"`
	Parent string `json:"parent" xml:"parent" form:"parent"`
}
