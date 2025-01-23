package dto

type JSONSubSection struct {
	Id          string           `json:"id,omitempty"`
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description"`
	Type        string           `json:"type,omitempty"`
	Content     string           `json:"content,omitempty"`
	Url         string           `json:"url,omitempty"`
	Items       []JSONSubSection `json:"items,omitempty"`
	Video       string           `json:"video,omitempty"`
	Images      string           `json:"image,omitempty"`
	Price       int64            `json:"price,omitempty"`
	Tutorial    string           `json:"tutorial,omitempty"`
}
