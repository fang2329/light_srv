package main

import (
	"../msg"
	"../services"
	"fmt"
	"github.com/fang2329/logger/common_log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	checkErr(err)
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {

	path := getCurrentPath()
	fmt.Println(path)
	host := "127.0.0.1:7778"

	s, err := services.NewServStruct(host)
	if err != nil {
		return
	}

	s.SetMsgHandler(HandleMessage)
	s.SetMsgConnHandler(HandleConnect)
	s.SetDisConn(HandleDisconnect)

	go NewConnect()

	timer := time.NewTimer(time.Second * 100)
	go func() {
		<-timer.C
		common_log.LOG_DEBUG("service stoped")
	}()

	common_log.LOG_DEBUG("service running on " + host)
	s.OnSrvService()
}

func HandleMessage(s *services.SessionStruct, msg *msg.MsgStruct) {
	common_log.LOG_DEBUG("receive msg: %s", msg)
	common_log.LOG_DEBUG("receive data: %s", string(msg.GetMsgContent()))
}

func HandleDisconnect(s *services.SessionStruct, err error) {
	common_log.LOG_DEBUG("%s", s.GetConn().GetName()+" lost.")
}

func HandleConnect(s *services.SessionStruct) {
	common_log.LOG_DEBUG("%s", s.GetConn().GetName()+" connected.")
}

func NewConnect() {
	host := "127.0.0.1:7778"
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
