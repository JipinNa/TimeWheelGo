package main

import (
	"TimeWheelGo/task"
	"log"
	"time"
)

func main() {
	taskHandler, err := task.New(time.Second, 5, callback)
	if err != nil {
		log.Println(err.Error())
		return
	}
	taskHandler.AddTask(&task.Task{
		Data:  "123456",
		Delay: time.Second * 20,
	})
	taskHandler.Start()
	for {
		select {
		case <-time.NewTimer(100).C:
			break
		}
	}
}

func callback(data interface{}) {
	log.Println(data)
}
