package msg

import "unsafe"

/*packaget head*/
type MsgHeader struct {
	MsgId  uint32
	MsgLen uint32
}
/*packet*/
type MsgStruct struct {
	MsgHeader
	MsgContent  []byte
}

func NewMsgHeader(id,length uint32) (*MsgHeader) {
	return  &MsgHeader{
		MsgId:id,
		MsgLen:length,
	}
}

func GetMsgHeaderLen() uintptr  {
	tmpHeader := NewMsgHeader(0,0)
	return unsafe.Sizeof(tmpHeader)
}

func (this *MsgHeader)GetMsgId() uint32  {
	return this.MsgId
}

func (this *MsgHeader)SetMsgId(id uint32)  {
	this.MsgId = id
}

func (this *MsgHeader)GetMsgLen() uint32  {
	return this.MsgLen
}

func (this *MsgHeader)SetMsgLen(length uint32)  {
	this.MsgLen = length
}

func NewMsg(id uint32,data []byte) (*MsgStruct) {
	lenth := uint32(len(data))
	tmpMsg := &MsgStruct{}
	tmpMsg.SetMsgId(id)
	tmpMsg.SetMsgLen(lenth)
	tmpMsg.MsgContent = data
	return tmpMsg
}

func (this *MsgStruct)GetMsgContent()[]byte  {
	return this.MsgContent
}