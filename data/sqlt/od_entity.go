package sqlt

import (
	"gorm.io/gorm"
)

type DetectedObject struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
	X1          float32 `json:"x1"`
	Y1          float32 `json:"y1"`
	X2          float32 `json:"x2"`
	Y2          float32 `json:"y2"`
}

// OdEntity todo: add metadata
type OdEntity struct {
	gorm.Model
	BaseEntity
	DetectedObject
}
