package model

type SyncData struct {
	Data     *Data
	BackChan chan *Data
}

func NewSyncData(data *Data) *SyncData {
	return &SyncData{
		data,
		make(chan *Data),
	}
}
