package ioservice

import (
	"fmt"
	"reflect"
	"time"
)

// IOService 一个ioservice一个goroutine
type IOService interface {
	// Post rpc 方式
	Post(f func())
	// AfterPost 调用cancel可以取消本次调用，但一旦开始执行就会取消失败
	AfterPost(d time.Duration, f func()) (cancel func() bool)
	// TryPost 非阻塞 rpc 方式
	TryPost(f func()) bool

	RPCGo(f interface{}, args ...interface{})
	// AfterRPCGo 见AfterPost
	AfterRPCGo(d time.Duration, f interface{}, args ...interface{}) (cancel func() bool)
	// RPCCall 注意多返回值的时候ret是切片
	RPCCall(id interface{}, args ...interface{}) (ret interface{})
	// AfterRPCCall 见RPCCall
	AfterRPCCall(d time.Duration, f interface{}, args ...interface{}) (ret interface{})

	Init()
	// Run 开启goroutine
	Run()
	// Fini 关闭
	Fini()
}

// 最小时间单位
const rpcMinTime = 100 * time.Millisecond

// 误差百分比
const rpcTimeOffset = 100

// rpcEvent ...
type rpcEvent struct {
	f    reflect.Value
	args []interface{}

	retChan chan interface{}
}

func newRPCEvent(f interface{}, args ...interface{}) *rpcEvent {
	isOk, funcValue := isFunc(f)
	if !isOk {
		panic(fmt.Sprintf("function %v is not func ", f))
	}

	event := &rpcEvent{
		f:    funcValue,
		args: args,
	}

	return event
}

func newRPCEvent2(f interface{}, args ...interface{}) *rpcEvent {
	isOk, funcValue := isFunc(f)
	if !isOk {
		panic(fmt.Sprintf("function %v is not func ", f))
	}

	var retChan chan interface{}
	if funcValue.Type().NumOut() > 0 {
		retChan = make(chan interface{}, 1)
	}

	event := &rpcEvent{
		f:       funcValue,
		args:    args,
		retChan: retChan,
	}

	return event
}

func isFunc(f interface{}) (bool, reflect.Value) {
	funcType := reflect.TypeOf(f)
	if funcType == nil || funcType.Kind() != reflect.Func {
		return false, reflect.ValueOf(funcType)
	}

	funcValue := reflect.ValueOf(f)
	if funcValue.IsValid() == false {
		return false, funcValue
	}

	return true, funcValue
}

func (e *rpcEvent) doRPC() {
	argsCount := len(e.args)
	if argsCount != e.f.Type().NumIn() {
		panic("dorpc call args != function args")
	}

	var inargs []reflect.Value
	if argsCount > 0 {
		inargs = make([]reflect.Value, argsCount)
		for k, arg := range e.args {
			inargs[k] = reflect.ValueOf(arg)
		}
	}

	if e.retChan == nil {
		e.f.Call(inargs)
	} else {
		rets := e.f.Call(inargs)

		if len(rets) == 1 {
			e.retChan <- rets[0].Interface()
		} else if len(rets) > 1 {
			retsIf := make([]interface{}, len(rets), len(rets))
			for i, v := range rets {
				retsIf[i] = v.Interface()
			}

			e.retChan <- retsIf
		}
	}
}
