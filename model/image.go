package model

type ImageDBEntry struct {
	Author string
	Title  string
	Link   string
	Width  int
	Height int
}

type ImageResponse struct {
	Data []byte
}
