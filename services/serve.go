package services

import (
	"../msg"
	"context"
	"github.com/fang2329/logger/common_log"
	"net"
	"sync"
	"time"
)

/*status*/
const (
	ErrorSer = iota
	InitializeSer
	RunnomalSer
	StopfinishSer
)

type ServeStruct struct {
	onMsg     func(*SessionStruct, *msg.MsgStruct)
	onConn    func(*SessionStruct)
	onDisConn func(*SessionStruct, error)
	timeOut   time.Duration
	status    int
	listener  net.Listener
	localAddr string
	session   *sync.Map
	finishCh  chan error
}

func NewServStruct(addr string) (*ServeStruct, error) {
	newListen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ser := &ServeStruct{
		timeOut:   0 * time.Second,
		status:    InitializeSer,
		listener:  newListen,
		localAddr: addr,
		session:   &sync.Map{},
		finishCh:  make(chan error),
	}
	return ser, nil
}

/*set msg handler*/
func (this *ServeStruct) SetMsgHandler(msgHandler func(*SessionStruct, *msg.MsgStruct)) {
	this.onMsg = msgHandler
}

/*set connect handler*/
func (this *ServeStruct) SetMsgConnHandler(connHandler func(*SessionStruct)) {
	this.onConn = connHandler
}

/*set disconnect handler*/
func (this *ServeStruct) SetDisConn(disConnHandler func(*SessionStruct, error)) {
	this.onDisConn = disConnHandler
}

func (this *ServeStruct) connHandler(ctx context.Context, c net.Conn) {
	common_log.LOG_DEBUG("connHandler...")
	conn := NewConnect(c, this.timeOut)
	session := NewSession(conn)
	this.session.Store(session.GetSessionId(), session)
	conCtx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		conn.CloseConn()
		this.session.Delete(session.GetSessionId())
	}()

	common_log.LOG_DEBUG("routine start....")
	go conn.ReadFromGoroutine(conCtx)
	go conn.WriteToGoroutine(conCtx)

	if this.session != nil {
		common_log.LOG_DEBUG("session is not nil")
		this.onConn(session)
	}

	common_log.LOG_DEBUG("do")
	for {
		select {
		case err := <-conn.doChan:
			if this.onDisConn != nil {
				this.onDisConn(session, err)
			}
			return
		case msg := <-conn.msgChan:
			if this.onMsg != nil {
				this.onMsg(session, msg)
			}

		}
	}

}

func (this *ServeStruct) acceprHandler(ctx context.Context) {
	common_log.LOG_DEBUG("acceprHandler...")
	con, err := this.listener.Accept()
	if err != nil {
		this.finishCh <- err
		return
	}
	go this.connHandler(ctx, con)
}

func (this *ServeStruct) OnSrvService() {
	common_log.LOG_DEBUG("OnSrvService...")
	this.status = RunnomalSer
	tmpctx, tmpcancel := context.WithCancel(context.Background())
	defer func() {
		this.status = StopfinishSer
		tmpcancel()
		this.listener.Close()
	}()

	go this.acceprHandler(tmpctx)
	for {
		select {
		case <-this.finishCh:
			return
		}
	}
}

func (this *ServeStruct) GetStatus() int {
	return this.status
}

func (this *ServeStruct) SetStatus(status int) {
	this.status = status
}

func (this *ServeStruct) GetConnectCount() int {
	var count int
	this.session.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}
