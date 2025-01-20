package dto

type JSONSubSection struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Body        string           `json:"body,omitempty"`
	Items       []JSONSubSection `json:"items,omitempty"`
}
