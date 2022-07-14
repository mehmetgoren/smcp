package sqlt

import (
	"gorm.io/gorm"
)

type DetectedFace struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name" gorm:"index:idx_query"` //Index
	X1          float32 `json:"x1"`
	Y1          float32 `json:"y1"`
	X2          float32 `json:"x2"`
	Y2          float32 `json:"y2"`
}

type FrEntity struct {
	gorm.Model
	BaseEntity
	DetectedFace
}
