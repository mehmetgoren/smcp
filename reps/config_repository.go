package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"smcp/models"
)

type ConfigRepository struct {
	Connection *redis.Client
}

func (r *ConfigRepository) GetConfig() (*models.Config, error) {
	var config = &models.Config{}
	conn := r.Connection
	key := "config"
	data, err := conn.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return config, nil
		} else {
			log.Println("Error getting sources from redis: ", err)
			return nil, err
		}
	}

	err = json.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *ConfigRepository) SaveConfig(config *models.Config) error {
	conn := r.Connection
	key := "config"
	data, err := json.Marshal(config)
	if err != nil {
		log.Println("Error serializing json: ", err)
		return err
	}

	err = conn.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		log.Println("Error saving ml config to redis: ", err)
		return err
	}

	return nil
}

func (r *ConfigRepository) RestoreConfig() (*models.Config, error) {
	objJson := `{"ai":{"detected_folder":"/mnt/sde1/detected/"},"device":{"device_name":"gokalp-ubuntu","device_services":[1,2,4,8,16],"device_type":0},"ffmpeg":{"check_ffmpeg_record_running_process_interval":30,"check_ffmpeg_stream_running_process_interval":10,"check_leaky_ffmpeg_processes_interval":600,"check_unstopped_containers_interval":600,"event_listener_handler_type":0,"max_operation_retry_count":10000000,"start_task_wait_for_interval":3.0,"use_double_quotes_for_path":false},"handler":{"read_service_overlay":true,"save_image_extension":"jpg","save_image_folder_path":"/mnt/sde1","show_image_caption":false,"show_image_fullscreen":false,"show_image_wait_key":1},"heartbeat":{"interval":5},"jetson":{"model_name":"ssd-mobilenet-v2","threshold":0.1,"white_list":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90]},"once_detector":{"imagehash_threshold":3,"psnr_threshold":0.2,"ssim_threshold":0.2},"path":{"read":"/mnt/sde1/read","record":"/mnt/sde1/playback","stream":"/mnt/sde1/live"},"redis":{"host":"127.0.0.1","port":6379},"source_reader":{"buffer_size":2,"fps":1,"max_retry":150,"max_retry_in":6},"torch":{"model_name":"ultralytics/yolov5","model_name_specific":"yolov5x6","threshold":0.1,"white_list":[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79]}}`

	conn := r.Connection
	key := "config"
	err := conn.Set(context.Background(), key, objJson, 0).Err()
	if err != nil {
		log.Println("Error saving default ml config to redis: ", err)
		return nil, err
	}

	var config = &models.Config{}
	err = json.Unmarshal([]byte(objJson), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
