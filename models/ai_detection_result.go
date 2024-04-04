package models

type AiDetectionResult struct {
	Label string          `json:"label"`
	Score float32         `json:"score"`
	Box   *AiDetectionBox `json:"box"`
}
