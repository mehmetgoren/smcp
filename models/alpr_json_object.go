package models

type AlprJsonBaseObject struct {
	ImgWidth         int     `json:"img_width"`
	ImgHeight        int     `json:"img_height"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
	Results          []struct {
		Plate            string  `json:"plate"`
		Confidence       float64 `json:"confidence"`
		ProcessingTimeMs float64 `json:"processing_time_ms"`
		Coordinates      struct {
			X0 int `json:"x0"`
			Y0 int `json:"y0"`
			X1 int `json:"x1"`
			Y1 int `json:"y1"`
		} `json:"coordinates"`
		Candidates []struct {
			Plate      string  `json:"plate"`
			Confidence float64 `json:"confidence"`
		} `json:"candidates"`
	} `json:"results"`
	Id            string `json:"id"`
	SourceId      string `json:"source_id"`
	CreatedAt     string `json:"created_at"`
	AiClipEnabled bool   `json:"ai_clip_enabled"`
	ImageFileName string `json:"image_file_name"`
	DataFileName  string `json:"data_file_name"`
}

type AlprJsonObject struct {
	AlprResults *AlprJsonBaseObject  `json:"alpr_results"`
	Video       *VideoClipJsonObject `json:"video"`
}
