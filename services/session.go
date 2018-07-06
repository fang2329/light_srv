package services

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/fang2329/logger/common_log"
	"io"
)

type SessionStruct struct {
	sessionId string
	userId    uint32
	conn      *Connect
}

func GetMd5String(s string) string {
	str := md5.New()
	str.Write([]byte(s))
	return hex.EncodeToString(str.Sum(nil))
}

func UniqueId() string {
	buff := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, buff)
	if err != nil {
		common_log.LOG_ERROR("obtain uniqueid faild")
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(buff))
}

func NewSession(con *Connect) *SessionStruct {
	return &SessionStruct{
		sessionId: UniqueId(),
		userId:    0,
		conn:      con,
	}
}

func (this *SessionStruct) GetSessionId() string {
	return this.sessionId
}

func (this *SessionStruct) GetUserId() uint32 {
	return this.userId
}

func (this *SessionStruct) SetUserId(id uint32) {
	this.userId = id
}

func (this *SessionStruct) GetConn() *Connect {
	return this.conn
}

func (this *SessionStruct) SetConn(con *Connect) {
	this.conn = con
}
