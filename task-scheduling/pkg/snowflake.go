package pkg

import (
	"github.com/bwmarrin/snowflake"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
)

func GetID() string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("make ID error")
		return ""
	}
	id := node.Generate()
	return id.String()
}
