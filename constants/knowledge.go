package constants

type formatsStruct struct {
	Markdown string `json:"markdown"`
	JSON     string `json:"json"`
	Text     string `json:"text"`
}

var Formats = formatsStruct{
	Markdown: "markdown",
	JSON:     "json",
	Text:     "text",
}
