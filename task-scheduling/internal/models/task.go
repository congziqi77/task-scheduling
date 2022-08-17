package models

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"github.com/congziqi77/task-scheduling/pkg"
)

var ResChan = make(chan GraphResult, 10)

type Tasks struct {
	TaskList []Task `json:"task_list"`
}

type GraphResult struct {
	Error     error      `json:"error"`
	Graphres  [][]string `json:"graphres"`
	TopicID   string     `json:"topic_id"`
	TopicName string     `json:"topic_name"`
}

//任务
type Task struct {
	ID        string   `json:"id"`
	TaskName  string   `json:"task_name"`
	Comment   string   `json:"comment"`
	ParentId  []string `json:"parent_id"` //依赖任务id
	TopicName string   `json:"topic_name"`
}

//将新建的taskServer存入缓存当中
func (tasks *Tasks) TaskCreateServer(topicName, topicID string) error {
	for i := range tasks.TaskList {
		tasks.TaskList[i].ID = pkg.GetID()
		log.Printf("id:%v", tasks.TaskList[i].ID)
		tasks.TaskList[i].TopicName = topicName
	}
	maps, err := GetTopicMapFromCache()
	if err != nil {
		return err
	}
	topicKey := MakeTopicKey(topicName, topicID)
	topic, ok := maps[topicKey]
	if !ok {
		return errors.New("topicName or topicID is not correct")
	}
	topic.Tasks = append(topic.Tasks, tasks.TaskList...)
	maps[topicKey] = topic
	err = SetTopicMapToCache(maps)
	if err != nil {
		return err
	}
	makeGraphResAsync(topic)
	return nil
}

//异步获取
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
						Error:     err,
						Graphres:  nil,
						TopicID:   topic.ID,
						TopicName: topic.TopicName,
					}
					return
				}
			}
		}
		taskExecuteSequence := graph.TopoSortedLayers()
		ResChan <- GraphResult{
			Error:     nil,
			Graphres:  taskExecuteSequence,
			TopicID:   topic.ID,
			TopicName: topic.TopicName,
		}
	}(topic)

	//启动获取ResChan 判断是否启动
	b, _ := global.FreeCache.Get([]byte(global.ISStartGetFromResChan))
	if b == nil {
		go GetTaskTopo()
		global.FreeCache.Set([]byte(global.ISStartGetFromResChan), []byte("true"), -1)
	}
}

//同步获取
func MakeGraphResSync(topic Topic) GraphResult {
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
				return GraphResult{
					Error:     err,
					Graphres:  nil,
					TopicID:   topic.ID,
					TopicName: topic.TopicName,
				}
			}
		}
	}
	taskExecuteSequence := graph.TopoSortedLayers()
	return GraphResult{
		Error:     nil,
		Graphres:  taskExecuteSequence,
		TopicID:   topic.ID,
		TopicName: topic.TopicName,
	}
}

//循环获取TaskTopo
func GetTaskTopo() {
	for {
		taskTopo := <-ResChan
		logger.Debug().Interface("taskTopo:%v", taskTopo).Msg("")
		if taskTopo.Error != nil {
			logger.Error().Str("get tasktopo topicName: ", taskTopo.TopicName).
				Str("err:", taskTopo.Error.Error()).Msg("")
		} else {
			key := taskTopo.TopicID + global.TopicTopoSuffix
			taskTopoJson, err := json.Marshal(taskTopo)
			if err != nil {
				logger.Error().Str("topic topo to json error", err.Error()).Msg("to json")
				continue
			}
			global.FreeCache.Set([]byte(key), taskTopoJson, -1)
		}
	}
}

func TaskRun(graphTopo [][]string, topicName, topicID string) (bool, error) {
	b, err := global.FreeCache.Get([]byte(TopicKey))
	if err != nil {
		return false, err
	}

	topicMap := make(map[string]Topic)
	err = json.Unmarshal(b, &topicMap)
	if err != nil {
		return false, err
	}
	topicKey := MakeTopicKey(topicName, topicID)
	topic := topicMap[topicKey]
	var res bool
	if topic.Type == 0 {
		res, err = taskRunFromSql(topic, graphTopo)
	} else {
		res, err = taskRunFromShell(topic, graphTopo)
	}
	return res, err
}

//执行task sql
func taskRunFromSql(topic Topic, graphTopo [][]string) (bool, error) {
	var err error
	if global.DB == nil {
		global.DB, err = NewDBEngine()
		if err != nil {
			return false, err
		}
	}
	for _, s := range graphTopo {
		tierLen := len(s)
		var wg sync.WaitGroup
		wg.Add(tierLen)
		for i := 0; i < tierLen; i++ {
			go executeSQL(s[i], &wg, topic)
		}
		wg.Wait()
	}
	return true, nil
}

//执行task shell
func taskRunFromShell(topic Topic, graphTopo [][]string) (bool, error) {
	return false, nil
}

func queryExecuteStatement(taskID string, topic Topic) (string, error) {
	b, err := global.FreeCache.Get([]byte(topic.TopicName + taskID))
	if err != nil {
		return "", err
	}
	if b == nil {
		err = setTopic2Cache(taskID, topic)
		if err != nil {
			return "", nil
		}
		b, err = global.FreeCache.Get([]byte(topic.TopicName + taskID))
		if err != nil {
			return "", err
		}
	}
	var task Task
	json.Unmarshal(b, &task)

	return task.Comment, nil
}

func executeSQL(taskID string, wg *sync.WaitGroup, topic Topic) {
	defer func() {
		if e := recover(); e != nil {
			logger.Panic().Interface("execute sql panic:%v", e).Msg("")
		}
	}()
	defer wg.Done()
	commit, err := queryExecuteStatement(taskID, topic)
	if err != nil {
		panic(err)
	}
	if commit == "" {
		return
	}
	if err = global.DB.Exec(commit).Error; err != nil {
		panic(err)
	}
}

func setTopic2Cache(taskID string, topic Topic) error {
	b, err := global.FreeCache.Get([]byte(TopicKey))
	if err != nil {
		return err
	}
	mapTopic := make(map[string]Topic)
	json.Unmarshal(b, &mapTopic)
	topicOne := mapTopic[MakeTopicKey(topic.TopicName, topic.ID)]
	for _, task := range topicOne.Tasks {
		b, err := json.Marshal(task)
		if err != nil {
			return err
		}
		err = global.FreeCache.Set([]byte(topic.TopicName+task.ID), b, -1)
		if err != nil {
			return err
		}
	}
	return nil
}
