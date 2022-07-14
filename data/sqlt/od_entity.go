package sqlt

import (
	"gorm.io/gorm"
)

type DetectedObject struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name" gorm:"index:idx_query"` //Index
}

type OdEntity struct {
	gorm.Model
	BaseEntity
	DetectedObject
}
