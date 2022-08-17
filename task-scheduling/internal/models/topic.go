package models

import (
	"encoding/json"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/pkg"
)

type Topic struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name" binding:"required,max=50"`
	Desc      string `json:"desc" binding:"required"`
	Type      int    `json:"type" binding:"oneof=0 1"` //0 sql 1 shell
	Tasks     []Task `json:"tasks"`
}

const (
	TopicKey = "topicList"
)

//保存topic到cache中
func (topic *Topic) SaveTopic2Cache() error {
	topic.ID = pkg.GetID()
	mapTopic, err := GetTopicMapFromCache()
	if err != nil {
		return err
	}
	mapTopic[MakeTopicKey(topic.TopicName, topic.ID)] = *topic
	err = SetTopicMapToCache(mapTopic)
	if err != nil {

		return err
	}
	return nil
}

//任务执行完删除对应topic
func delTopicFromCache(key string) error {
	b, _ := global.FreeCache.Get([]byte(TopicKey))
	if b == nil {
		return nil
	}
	mapTopic := make(map[string]Topic)
	err := json.Unmarshal(b, &mapTopic)
	if err != nil {
		return err
	}
	delete(mapTopic, key)
	err = SetTopicMapToCache(mapTopic)
	if err != nil {
		return err
	}
	return nil
}

func GetTopicMapFromCache() (map[string]Topic, error) {
	b, _ := global.FreeCache.Get([]byte(TopicKey))
	mapTopic := make(map[string]Topic)
	if b != nil {
		err := json.Unmarshal(b, &mapTopic)
		if err != nil {
			return nil, err
		}
	}
	return mapTopic, nil
}

func SetTopicMapToCache(maps map[string]Topic) error {
	b, err := json.Marshal(maps)
	if err != nil {
		return err
	}
	err = global.FreeCache.Set([]byte(TopicKey), b, -1)
	if err != nil {
		return err
	}
	return nil
}

func GetTopicTopo(topicName, topicID string) ([][]string, error) {
	b, err := global.FreeCache.Get([]byte(topicID + global.TopicTopoSuffix))
	if err != nil {
		return nil, err
	}
	var graphRes GraphResult
	//如果异步没有获取到数据那么就从同步中获取
	if b == nil {
		logger.Warn().Msg("warn: get topic topo by sync")
		topicMap := make(map[string]Topic)
		json.Unmarshal(b, &topicMap)
		topic := topicMap[MakeTopicKey(topicName, topicID)]
		graphRes = MakeGraphResSync(topic)
		if graphRes.Error != nil {
			return nil, err
		}
		b, err := json.Marshal(graphRes)
		if err != nil {
			return nil, err
		}
		global.FreeCache.Set([]byte(topicID+global.TopicTopoSuffix), b, -1)
		return graphRes.Graphres, graphRes.Error
	}

	json.Unmarshal(b, &graphRes)
	return graphRes.Graphres, graphRes.Error
}

func Run(topicName, topicID string) (bool, error) {
	b, err := global.FreeCache.Get([]byte(topicID + global.TopicTopoSuffix))
	if err != nil {
		return false, err
	}
	var graphRes GraphResult
	err = json.Unmarshal(b, &graphRes)
	if err != nil {
		return false, err
	}
	result, err := TaskRun(graphRes.Graphres, topicName, topicID)
	return result, err
}

func MakeTopicKey(topicName, topicID string) string {
	return topicName + ":" + topicID
}
