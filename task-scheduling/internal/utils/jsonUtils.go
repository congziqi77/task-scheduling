package main

import (
	"encoding/json"

	"github.com/congziqi77/task-scheduling/internal/modules/logger"
)

//struct -> json
func Struct2Json(stru any) string {
	res, err := json.Marshal(stru)
	if err != nil {
		logger.NewLogger.Error().Msgf("struct to json error : %v", err)
	}
	return string(res)
}

//json -> struc any is struct point

func Json2Struct(j []byte, v any) any {
	err := json.Unmarshal(j, v)
	if err != nil {
		logger.NewLogger.Error().Msgf("json to struct error : %v", err)
	}
	return v
}

// json -> map
func Json2Map(j []byte, m *map[string]any) *map[string]any {
	err := json.Unmarshal(j, m)
	if err != nil {
		logger.NewLogger.Error().Msgf("json to map error : %v", err)
	}
	return m
}

//map -> json
func Map2Json(m *map[string]any) string {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		logger.NewLogger.Error().Msgf("map to json error : %v", err)
	}
	return string(jsonStr)
}
