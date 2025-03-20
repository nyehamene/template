package token

import (
	"fmt"
	"temlang/tem/dsa/queue"
)

type Error struct {
	Location Location
	Msg      string
}

func (e Error) String() string {
	return fmt.Sprintf("%s at %s", e.Msg, e.Location)
}

type ErrorQueue = queue.Queue[Error]
