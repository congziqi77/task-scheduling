package models

import (
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/pkg"
)

var ResChan = make(chan GraphResult, 1)

type Tasks struct {
	TaskList []Task `json:"task_list"`
}

type GraphResult struct {
	Error    error      `json:"error"`
	Graphres [][]string `json:"graphres"`
}

//任务
type Task struct {
	ID       string      `json:"id"`
	TaskName string      `json:"task_name"`
	Comment  interface{} `json:"comment"`
	ParentId []string    `json:"parent_id"` //依赖任务id
	// Children  string      `json:"children"` //被依赖任务id
	TopicName string `json:"topic_name"`
}

//将新建的taskServer存入缓存当中
func (tasks *Tasks) TaskCreateServer(topicName string) error {
	for i := range tasks.TaskList {
		tasks.TaskList[i].ID = pkg.GetID()
		tasks.TaskList[i].TopicName = topicName
	}
	maps, err := GetTopicMapFromCache()
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("get cache err")
		return err
	}
	topic := maps[topicName]
	topic.Tasks = append(topic.Tasks, tasks.TaskList...)
	maps[topicName] = topic
	err = SetTopicMapToCache(maps)
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg("set cache error")
		return err
	}
	makeGraphResAsync(topic)
	return nil
}

func makeGraphResAsync(topic Topic) {
	go func(Topic) {
		graph := GraphNew()
		var err error
		for _, t := range topic.Tasks {
			if t.ParentId == nil {
				graph.addNilParents(t.ID)
				continue
			}
			for _, parent := range t.ParentId {
				err = graph.DependOn(t.ID, parent)
				if err != nil {
					ResChan <- GraphResult{
						Error:    err,
						Graphres: nil,
					}
					return
				}
			}
		}
		taskExecuteSequence := graph.TopoSortedLayers()
		ResChan <- GraphResult{
			Error:    nil,
			Graphres: taskExecuteSequence,
		}
	}(topic)
}

func makeGraphResSync(topic Topic) ([][]string, error) {
	graph := GraphNew()
	var err error
	for _, t := range topic.Tasks {
		if t.ParentId == nil {
			graph.addNilParents(t.ID)
			continue
		}
		for _, parent := range t.ParentId {
			err = graph.DependOn(t.ID, parent)
			if err != nil {
				logger.Error().Str("make GraphResSync err", err.Error())
				return nil, err
			}
		}
	}
	taskExecuteSequence := graph.TopoSortedLayers()
	return taskExecuteSequence, nil
}
