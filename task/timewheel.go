package task

import (
	"container/list"
	"errors"
	"log"
	"time"
)

type Callback func(interface{})

type TimeWheel interface {
	taskSchedule
}

type timeWheel struct {
	TimeWheel
	ticker       *time.Ticker
	interval     time.Duration
	buckets      []*list.List
	bucketSize   int
	curPos       int
	callbackFunc Callback
	cancel       chan bool
}

func (tw *timeWheel) AddTask(task *Task) {
	// 任务延迟时间
	delaySeconds := int(task.Delay.Seconds())
	// tw轮询时间
	intervalSeconds := int(tw.interval.Seconds())
	// 计算轮次
	loop := delaySeconds / intervalSeconds / tw.bucketSize
	// 计算需要插入到哪个桶内
	pos := tw.curPos + delaySeconds/intervalSeconds%tw.bucketSize
	task.loop = loop
	tw.buckets[pos].PushBack(task)
}

func (tw *timeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go func() {
		for {
			select {
			case <-tw.ticker.C:
				log.Println(`1 tick`)
				tw.tickHandler()
			case <-tw.cancel:
				tw.ticker.Stop()
				break
			}
		}
	}()
}

func (tw *timeWheel) tickHandler() {
	bucket := tw.buckets[tw.curPos]
	for e := bucket.Front(); e != nil; {
		task, ok := e.Value.(*Task)
		if !ok || task == nil {
			continue
		}
		if task.loop > 0 {
			task.loop--
			e = e.Next()
			continue
		}
		go tw.callbackFunc(task.Data)
		next := e.Next()
		bucket.Remove(e)
		e = next
	}
	if tw.curPos == tw.bucketSize-1 {
		log.Println(`new loop`)
		tw.curPos = 0
	} else {
		tw.curPos++
	}
}

func (tw *timeWheel) Stop() {
	tw.cancel <- true
}

func New(interval time.Duration, bucketSize int, callbackFunc Callback) (TimeWheel, error) {
	if interval <= 0 || bucketSize <= 0 || callbackFunc == nil {
		return nil, errors.New("invalid input params")
	}
	tw := &timeWheel{
		interval:     interval,
		buckets:      make([]*list.List, bucketSize),
		bucketSize:   bucketSize,
		curPos:       0,
		callbackFunc: callbackFunc,
		cancel:       make(chan bool),
	}
	for i := range tw.buckets {
		tw.buckets[i] = list.New()
	}
	return tw, nil
}
