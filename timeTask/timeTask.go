package timeTask

import (
	"time"
)

type TimeTask struct {
	closeSignalChannel chan int
	consumer           func()
	interval           time.Duration
	serialization      bool
}

func NewTask() *TimeTask {
	ts := &TimeTask{
		closeSignalChannel: make(chan int),
		consumer:           func() {},
		interval:           time.Second,
	}
	ts.start()
	return ts
}
func (_this *TimeTask) Close() *TimeTask {
	select {
	case _this.closeSignalChannel <- 1:
	default:
	}
	return _this
}
func (_this *TimeTask) start()  {
	go func() {
		for {
			select {
			case <-_this.closeSignalChannel:
				return
			case <-time.After(_this.interval):
				if _this.serialization {
					_this.consumer()
				} else {
					go _this.consumer()
				}
			}
		}
	}()
}
func (_this *TimeTask) SetInterval(interval time.Duration) *TimeTask {
	_this.interval = interval
	return _this
}
func (_this *TimeTask) SetConsumer(consumer func()) *TimeTask {
	_this.consumer = consumer
	return _this
}
func (_this *TimeTask) SetSerialization(tf bool) *TimeTask {
	_this.serialization = tf
	return _this
}
