package msg

import (
	"bytes"
	"encoding/binary"
)

func (this *MsgHeader) EncodeHeader() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, this)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func (this *MsgHeader) DecodeHeader(buf []byte) error {
	buffer := bytes.NewReader(buf)
	err := binary.Read(buffer, binary.LittleEndian, this)
	if err != nil {
		return err
	}
	return nil
}

func Encode(msg *MsgStruct) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.MsgId)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.MsgLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.MsgContent)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode from []byte to Message
func Decode(data []byte) (*MsgStruct, error) {
	bufReader := bytes.NewReader(data)

	dataSize := uint32(len(data))
	// msgId
	var msgID uint32
	err := binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		return nil, err
	}

	// msg data
	headLen := uint32(GetMsgHeaderLen())
	dataLength := dataSize - headLen
	dataBuff := make([]byte, dataLength)
	err = binary.Read(bufReader, binary.LittleEndian, &dataBuff)
	if err != nil {
		return nil, err
	}

	msg := &MsgStruct{}
	msg.SetMsgId(msgID)
	msg.SetMsgLen(dataLength)
	msg.MsgContent = dataBuff

	return msg, nil
}

func FormMsg(id uint32, data []byte) []byte {
	datalen := uint32(len(data))
	tmpMsg := &MsgStruct{}
	tmpMsg.SetMsgId(id)
	tmpMsg.SetMsgLen(datalen)
	tmpMsg.MsgContent = data
	msg, err := Encode(tmpMsg)
	if err != nil {
		return nil
	}
	return msg
}
