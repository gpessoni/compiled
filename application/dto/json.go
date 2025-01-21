package dto

type JSONSubSection struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Content     string           `json:"content,omitempty"`
	Items       []JSONSubSection `json:"items,omitempty"`
}
