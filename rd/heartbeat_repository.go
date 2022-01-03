package rd

import (
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"
)

type HeartbeatRepository struct {
	*RedisOptions
	TimeSecond int64
}

func toMyFormat(t *time.Time) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t.Year()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(t.Month())))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Day()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Hour()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Minute()))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(t.Second()))
	sb.WriteString("-")
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
				"smcp_service": toMyFormat(&timeTicker),
			}
			h.Client.HSet(context.Background(), "heartbeat", heartbeatObj)
			log.Println("Heartbeat was beaten at " + timeTicker.Format(time.ANSIC))
			//case <- quit:
			//	ticker.Stop()
			//	return
		}
	}
}
