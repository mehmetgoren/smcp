package reps

import "encoding/json"

func Map(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}
