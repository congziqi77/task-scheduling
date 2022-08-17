package pkg

import (
	"time"

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
	time.Sleep(1 * time.Millisecond)
	return id.String()
}
