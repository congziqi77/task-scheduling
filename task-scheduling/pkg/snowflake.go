package pkg

import (
	"github.com/bwmarrin/snowflake"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
)

func GetID() string {
	node, err := snowflake.NewNode(1)
	if err != nil {
		logger.NewLogger.Fatal().Msgf("make ID error", err)
		return ""
	}
	id := node.Generate()
	return id.String()
}
