package task

import "time"

type Task struct {
	Id    interface{}
	Data  interface{}
	Delay time.Duration
	loop  int
}

type taskSchedule interface {
	AddTask(task *Task)
	Start()
	Stop()
}
