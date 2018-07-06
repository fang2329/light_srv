package msgSend

import (
	"../msg"
	"../services"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	host := "127.0.0.1:7777"

	s, err := services.NewServStruct(host)
	if err != nil {
		return
	}

	s.SetMsgHandler(HandleMessage)
	s.SetMsgConnHandler(HandleConnect)
	s.SetDisConn(HandleDisconnect)

	go NewConnect()

	timer := time.NewTimer(time.Second * 1)
	go func() {
		<-timer.C
		t.Log("service stoped")
	}()

	t.Log("service running on " + host)
	s.OnSrvService()
}

func HandleMessage(s *services.SessionStruct, msg *msg.MsgStruct) {
	fmt.Println("receive msg:", msg)
	fmt.Println("receive data:", string(msg.GetMsgContent()))
}

func HandleDisconnect(s *services.SessionStruct, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *services.SessionStruct) {
	fmt.Println(s.GetConn().GetName() + " connected.")
}

func NewConnect() {
	host := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	tmpCon := services.NewConnect(conn, time.Second*20)

	tmpMsg := msg.NewMsg(1, []byte("hello world"))
	tmpCon.SendMsg(tmpMsg)
}
