package rd

import (
	"golang.org/x/net/context"
	"os"
	"runtime"
	"time"
)

type ServiceRepository struct {
	*RedisOptions
}

func getKey(serviceName string) string {
	return "services:" + serviceName
}

func (r *ServiceRepository) Add(serviceName string) (int64, error) {
	now := time.Now()
	host, _ := os.Hostname()
	heartbeatObj := map[string]interface{}{
		"name":             serviceName,
		"platform":         runtime.GOOS,
		"platform_version": runtime.GOARCH,
		"hostname":         host,
		"ip_address":       "127.0.0.1",
		"mac_address":      "00:00:00:00:00:00",
		"processor":        "unknown",
		"cpu_count":        runtime.NumCPU(),
		"ram":              "unknown",
		"pid":              os.Getpid(),
		"created_at":       DatetimeNow(&now),
		"heartbeat":        "",
	}

	return r.Client.HSet(context.Background(), getKey(serviceName), heartbeatObj).Result()
}
