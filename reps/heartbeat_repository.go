package reps

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"
)

type HeartbeatRepository struct {
	Client      *redis.Client
	TimeSecond  int64
	ServiceName string
}

func DatetimeNow(t *time.Time) string {
	sep := "_"
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.Year()))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(int(t.Month())))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(t.Day()))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(t.Hour()))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(t.Minute()))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(t.Second()))
	sb.WriteString(sep)
	sb.WriteString(strconv.Itoa(t.Nanosecond()))

	return sb.String()
}

func (h *HeartbeatRepository) Start() {
	var dur = time.Duration(h.TimeSecond) * time.Second
	ticker := time.NewTicker(dur)
	//quit := make(chan struct{})
	for {
		select {
		case timeTicker := <-ticker.C:
			heartbeatObj := map[string]interface{}{
				"heartbeat": DatetimeNow(&timeTicker),
			}
			h.Client.HSet(context.Background(), "services:"+h.ServiceName, heartbeatObj)
			log.Println("Heartbeat was beaten at " + timeTicker.Format(time.ANSIC))
			//case <- quit:
			//	ticker.Stop()
			//	return
		}
	}
}
