package model

import "time"

type TgMsg struct {
	UserId            uint64
	Username          string
	MessageText       string
	CallbackQueryData string
	TimeStamp         time.Time
}
