package util

import "log"

type Logger struct {
}

func (l Logger) Print(v...interface{}) {
	log.Println(v)
}
