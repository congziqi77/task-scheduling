package models

import (
	"encoding/json"
	"errors"

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

// 保存topic到cache中
func (topic *Topic) SaveTopic2Cache() error {
	topic.ID = pkg.GetID()
	mapTopic, err := GetTopicMapFromCache()
	if err != nil {
		logger.Debug().Str("get cache error", err.Error()).Msg("")
		return err
	}
	mapTopic[MakeTopicKey(topic.TopicName, topic.ID)] = *topic

	if err = SetTopicMapToCache(mapTopic); err != nil {
		logger.Debug().Str("set cache error", err.Error()).Msg("")
		return err
	}
	return nil
}

func delTopicFromCache(key string) error {
	b, _ := CacheImp.GetCache([]byte(TopicKey))
	if b == nil {
		return nil
	}
	mapTopic := make(map[string]Topic)
	if err := json.Unmarshal(b, &mapTopic); err != nil {
		return err
	}
	delete(mapTopic, key)
	err := SetTopicMapToCache(mapTopic)
	if err != nil {
		return err
	}
	return nil
}

func GetTopicMapFromCache() (map[string]Topic, error) {
	mapTopic := make(map[string]Topic)
	b, _ := CacheImp.GetCache([]byte(TopicKey))
	if b == nil {
		return mapTopic, nil
	}
	if b != nil {
		if err := json.Unmarshal(b, &mapTopic); err != nil {
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
	err2 := CacheImp.SetCache([]byte(TopicKey), b, -1)
	if err2 != nil {
		return err2
	}
	return nil
}

func GetTopicTopo(topicName, topicID string) ([][]string, error) {
	b, _ := CacheImp.GetCache([]byte(topicID + global.TopicTopoSuffix))
	var graphRes GraphResult
	//如果异步没有获取到数据或发生错误那么就从同步中获取
	if b == nil {
		logger.Print("warn: get topic topo by sync")
		topicMap := make(map[string]Topic)
		b, err := CacheImp.GetCache([]byte(TopicKey))
		if err != nil {
			return nil, err
		}
		json.Unmarshal(b, &topicMap)
		topic := topicMap[MakeTopicKey(topicName, topicID)]
		graphRes = MakeGraphResSync(topic)
		if graphRes.Error != nil {
			return nil, err
		}
		b, err = json.Marshal(graphRes)
		if err != nil {
			return nil, err
		}
		CacheImp.SetCache([]byte(topicID+global.TopicTopoSuffix), b, -1)
		return graphRes.Graphs, graphRes.Error
	}
	json.Unmarshal(b, &graphRes)
	//如果异步数据发生异常返回空集合
	if graphRes.Error != nil {
		return nil, graphRes.Error
	}
	return graphRes.Graphs, nil
}

func Run(topicName, topicID string) (bool, error) {
	b, err := CacheImp.GetCache([]byte(topicID + global.TopicTopoSuffix))
	if err != nil {
		return false, err
	}
	var graphRes GraphResult
	err = json.Unmarshal(b, &graphRes)
	if graphRes.Graphs == nil || len(graphRes.Graphs) == 0 {
		return false, errors.New("no graphs")
	}
	if err != nil {
		return false, err
	}
	result, err := TaskRun(graphRes.Graphs, topicName, topicID)
	return result, err
}

func MakeTopicKey(topicName, topicID string) string {
	return topicName + ":" + topicID
}
