package models

import (
	"encoding/json"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
)

type Topic struct {
	ID        string   `json:"id"`
	TopicName string   `json:"topic_name" binding:"required,max=50"`
	Desc      string   `json:"desc" binding:"required"`
	Type      int      `json:"type" binding:"oneof=0 1"` //0 sql 1 shell
	Tasks     []string `json:"tasks"`
}

const (
	TopicKey = "topicList"
)

//保存topic到cache中
func (topic *Topic) SaveTopic2Cache() error {
	b, _ := global.BigCache.Get(TopicKey)
	mapTopic := make(map[string]Topic)
	if b != nil {
		err := json.Unmarshal(b, &mapTopic)
		if err != nil {
			logger.Error().Str("err", err.Error()).Msg("json err")
			return err
		}
	}
	mapTopic[topic.TopicName+topic.ID] = *topic
	b, err := json.Marshal(mapTopic)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("map to json error")
		return err
	}
	global.BigCache.Set(TopicKey, b)

	return nil
}

//todo 删除cache中的topic request封装
//任务执行完删除对应topic
func delTopicFromCache(key string) error {
	d, _ := global.BigCache.Get(TopicKey)
	if d == nil {
		return nil
	}
	mapTopic := make(map[string]Topic)
	err := json.Unmarshal(d, &mapTopic)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("json err")
		return err
	}
	delete(mapTopic, key)
	d, err = json.Marshal(mapTopic)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("json err")
		return err
	}
	global.BigCache.Set(TopicKey, d)
	return nil
}
