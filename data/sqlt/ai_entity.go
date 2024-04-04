package sqlt

import (
	"gorm.io/gorm"
)

type AiDetectedObject struct {
	Score float32 `json:"score"`
	Label string  `json:"label"`
	X1    float32 `json:"x1"`
	Y1    float32 `json:"y1"`
	X2    float32 `json:"x2"`
	Y2    float32 `json:"y2"`
}

type AiEntity struct {
	gorm.Model
	BaseEntity
	AiDetectedObject
}
