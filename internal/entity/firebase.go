package entity

type Url struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}

func (u Url) QueryEscape(filename string) any {
	panic("unimplemented")
}

type MultipleFileUploadResponse struct {
	Url []*Url `json:"url"`
}
