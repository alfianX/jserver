package helper

import (
	"errors"
	"strconv"
)

func CheckMsgLen(msg string) error {
	length, err := strconv.ParseInt(msg[:4], 16, 64)
	if err != nil {
		return errors.New("iso message length not match: " + err.Error())
	}

	length = length * 2
	msgLength := len(msg[4:])
	if int64(msgLength) != length {
		return errors.New("iso message length not match")
	}

	return nil
}
