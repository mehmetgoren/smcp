package data

type AiClip struct {
	Enabled        bool   `json:"enabled" bson:"enabled"`
	FileName       string `json:"file_name" bson:"file_name"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
	Duration       int    `json:"duration" bson:"duration"`
}

func (a *AiClip) Setup(fileName string, createdAt string, lastModifiedAt string, duration int) {
	a.Enabled = true
	a.FileName = fileName
	a.CreatedAt = createdAt
	a.LastModifiedAt = lastModifiedAt
	a.Duration = duration
}
