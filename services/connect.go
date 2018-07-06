package services

import (
	"../msg"
	"context"
	"github.com/fang2329/logger/common_log"
	"io"
	"net"
	"time"
)

type Connect struct {
	sessionId    string
	originalConn net.Conn
	sendChan     chan []byte
	doChan       chan error
	name         string
	msgChan      chan *msg.MsgStruct
	timeOut      time.Duration
}

func NewConnect(c net.Conn, tmpTimeout time.Duration) *Connect {
	tmpconn := &Connect{
		originalConn: c,
		sendChan:     make(chan []byte, 100),
		doChan:       make(chan error),
		msgChan:      make(chan *msg.MsgStruct, 100),
		timeOut:      tmpTimeout,
	}
	tmpconn.name = c.RemoteAddr().String()
	return tmpconn
}

func (this *Connect) GetName() string {
	return this.name
}

func (this *Connect) CloseConn() {
	this.originalConn.Close()
}

func (this *Connect) SendMsg(pmsg *msg.MsgStruct) error {
	data, err := msg.Encode(pmsg)
	if err != nil {
		common_log.LOG_ERROR("encode msg failed")
		return err
	}
	common_log.LOG_DEBUG("encode success %s", data)
	this.originalConn.Write(data)

	return nil
}

func (this *Connect) WriteToGoroutine(ctx context.Context) {
	common_log.LOG_DEBUG("WriteToGoroutine...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := <-this.sendChan
			if data == nil {
				continue
			}
			_, err := this.originalConn.Write(data)
			if err != nil {
				common_log.LOG_DEBUG("failed")
				this.doChan <- err
			}
		}
	}
}

func (this *Connect) ReadFromGoroutine(ctx context.Context) {
	common_log.LOG_DEBUG("ReadFromGoroutine...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if this.timeOut > 0 {
				err := this.originalConn.SetReadDeadline(time.Now().Add(this.timeOut))
				if err != nil {
					this.doChan <- err
					continue
				}
			}
			/*header*/
			headerLen := uint32(msg.GetMsgHeaderLen())
			header := make([]byte, headerLen)
			_, err := io.ReadFull(this.originalConn, header)
			if err != nil {
				this.doChan <- err
				continue
			}

			msgHeader := &msg.MsgHeader{}
			err = msgHeader.DecodeHeader(header)
			if err != nil {
				this.doChan <- err
				continue
			}

			/*data content*/
			dataLen := msgHeader.GetMsgLen()
			tmpmsg := make([]byte, dataLen)
			_, err = io.ReadFull(this.originalConn, tmpmsg)
			if err != nil {
				this.doChan <- err
				continue
			}

			/*decode*/
			msgBuff := make([]byte, headerLen+dataLen)
			copy(msgBuff, header)
			copy(msgBuff[headerLen:], tmpmsg[0:])
			dstMsg, err := msg.Decode(tmpmsg)
			if err != nil {
				this.doChan <- err
				continue
			}
			this.msgChan <- dstMsg
		}
	}
}
