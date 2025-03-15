package photo

type AddPhotoBody struct {
	Name      string `json:"name" xml:"name" form:"name"`
	Url       string `json:"url" xml:"url" form:"url"`
	SectionId string `json:"section-id" xml:"section-id" form:"section-id"`
}
