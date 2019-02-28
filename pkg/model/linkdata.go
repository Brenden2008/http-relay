package model

type LinkData struct {
	Data     *Data
	BackChan chan *Data
}

func NewLinkData(data *Data) *LinkData {
	return &LinkData{
		data,
		make(chan *Data),
	}
}
