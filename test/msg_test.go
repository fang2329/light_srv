package test

import (
	"../common"
	"../msg"
	"github.com/fang2329/logger/common_log"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestMsgCode(t *testing.T) {
	testMsg := &common.Testinfo{
		Name: proto.String("abc"),
		Age:  proto.Uint32(10),
		Sex:  proto.String("woman"),
		Addr: proto.String("china-XXX"),
	}
	data, _ := proto.Marshal(testMsg)
	tmpMsg := msg.NewMsg(1, data)
	msgData, err := msg.Encode(tmpMsg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tmpMsg)
	msg2, err := msg.Decode(msgData)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("Id:%d,Data:%s", msg2.GetMsgId(), string(msg2.GetMsgContent()))
	common_log.LOG_DEBUG("Id:%d,Data:%s", msg2.GetMsgId(), string(msg2.GetMsgContent()))
}
