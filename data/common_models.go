package data

import (
	"time"
)

type AiClip struct {
	Enabled        bool   `json:"enabled" bson:"enabled"`
	FileName       string `json:"file_name" bson:"file_name"` //Index
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
	Duration       int    `json:"duration" bson:"duration"`
}

type SetVideoFileNameParams struct {
	SourceId      string
	Duration      int
	T1            *time.Time
	T2            *time.Time
	VideoFilename string
}
