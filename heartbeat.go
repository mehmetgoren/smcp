package main

import (
	"log"
	"time"
)

type HeartbeatClient struct {
	Repository *RedisRepository
	TimeSecond int64
}

func (h HeartbeatClient) Start() {
	var dur = time.Duration(h.TimeSecond) * time.Second
	ticker := time.NewTicker(dur)
	//quit := make(chan struct{})
	for {
		select {
		case heartBeat := <-ticker.C:
			h.Repository.Heartbeat(&heartBeat)
			log.Println("Heartbeat was beaten at " + heartBeat.Format(time.ANSIC))
			//case <- quit:
			//	ticker.Stop()
			//	return
		}
	}
}
