package model

import (
	"time"
)

type McastSeq struct {
	mcastDataMap map[int]*McastData
	newSeqId     int
	oldSeqId     int
	comm
}

func NewMcastSeq(initialSeqId int) *McastSeq {
	initialMap := make(map[int]*McastData)
	initialMap[initialSeqId] = NewMcastData()
	return &McastSeq{
		mcastDataMap: initialMap,
		newSeqId:     initialSeqId,
		oldSeqId:     initialSeqId,
	}
}

func (this *McastSeq) Close() {
	this.RLock()
	defer this.RUnlock()
	for _, v := range this.mcastDataMap {
		v.Close()
	}
}

func (this *McastSeq) GetData(seqId int) (data *TeeData, ok bool) {
	this.RLock()
	defer this.RUnlock()

	if seqId < this.newSeqId {
		if mcastData, ok := this.mcastDataMap[seqId]; ok {
			return mcastData.data, ok
		}
	}

	return
}

func (this *McastSeq) Read(wantedSeqId int, closeChan <-chan struct{}) (data *TeeData, seqId int, ok bool) {
	this.AddWaiter()
	defer this.RemoveWaiter()

	this.Lock()
	if wantedSeqId == -1 {
		if this.newSeqId == this.oldSeqId {
			seqId = this.newSeqId
		} else {
			seqId = this.newSeqId - 1
		}
	} else if wantedSeqId > this.newSeqId {
		seqId = this.newSeqId
	} else if wantedSeqId < this.oldSeqId {
		seqId = this.oldSeqId
	} else {
		seqId = wantedSeqId
	}

	mcastData := this.mcastDataMap[seqId]
	this.Unlock()

	data, ok = mcastData.Read(closeChan)
	return
}

func (this *McastSeq) Write(data *TeeData) (seqId int) {
	this.Lock()
	defer this.Unlock()

	this.accessed = time.Now()
	seqId = this.newSeqId
	mcastData := this.mcastDataMap[seqId]
	mcastData.Write(data)
	this.preserveSize()
	this.newSeqId += 1
	this.mcastDataMap[this.newSeqId] = NewMcastData()
	return
}

func (this *McastSeq) preserveSize() {
	for this.size() > 11000000 { // Total allowed 11Mb while reqest limited 10Mb
		this.removeOldest()
	}
}

func (this *McastSeq) removeOldest() {
	if this.oldSeqId < this.newSeqId {
		this.remove(this.oldSeqId)
		this.oldSeqId += 1
	}
}

func (this *McastSeq) remove(seqId int) {
	if mcastData, ok := this.mcastDataMap[seqId]; ok {
		delete(this.mcastDataMap, seqId)
		mcastData.Close()
	}
}

func (this *McastSeq) size() (size int) {
	for _, v := range this.mcastDataMap {
		size += v.Size()
	}
	return
}

func (this *McastSeq) Size() int {
	this.RLock()
	defer this.RUnlock()
	return this.size()
}

func (this *McastSeq) DataCount() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.mcastDataMap)
}

func (this *McastSeq) NewSeqId() int {
	return this.newSeqId
}
