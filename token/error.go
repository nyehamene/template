package token

import (
	"fmt"
	"temlang/tem/dsa/queue"
)

func NewError(offset int, msg string) Error {
	return Error{offset: offset, msg: &msg}
}

type ErrorQueue = queue.Queue[Error]

type Error struct {
	offset int
	msg    *string
}

func (e Error) Offset() int {
	return e.offset
}

func (e Error) Message() string {
	return *e.msg
}

func (e Error) String() string {
	return fmt.Sprintf("%s at %d", *e.msg, e.offset)
}
