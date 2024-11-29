package util

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Recover recover panic, 写入Stderr
func Recover() {
	if e := recover(); e != nil {
		stack := debug.Stack()
		err := fmt.Sprintf("%v\n", e)
		logrus.WithFields(logrus.Fields{
			"err":   err,
			"stack": string(stack),
		}).Error("Recover")

		os.Stderr.Write([]byte(err))
		os.Stderr.Write(stack)

		//fmt.Printf("%v\n", e)
		//fmt.Printf(string(debug.Stack()))
		//_, _ = os.Stderr.Write([]byte(fmt.Sprintf("%v\n%s", e, debug.Stack())))
	}
}

// SafeGo go
func SafeGo(f func()) {
	if f != nil {
		go func() {
			defer Recover()
			f()
		}()
	}
}
