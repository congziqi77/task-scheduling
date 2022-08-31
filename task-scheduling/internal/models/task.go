package models

import (
	"encoding/json"
	"errors"
	"strings"
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
	Graphs    [][]string `json:"graphs"`
	TopicID   string     `json:"topic_id"`
	TopicName string     `json:"topic_name"`
}

// 任务
type Task struct {
	ID         string   `json:"id"`
	TaskName   string   `json:"task_name"`
	Comment    string   `json:"comment"`
	ParentName []string `json:"parent_name"` //依赖任务id
	TopicName  string   `json:"topic_name"`
	Desc       string   `json:"desc"`
}

// 将新建的taskServer存入缓存当中
func (tasks *Tasks) TaskCreateServer(topicName, topicID string) error {
	for i := range tasks.TaskList {
		tasks.TaskList[i].ID = pkg.GetID()
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

// 异步获取
func makeGraphResAsync(topic Topic) {
	go func(Topic) {
		graph := GraphNew()
		var err error
		for _, t := range topic.Tasks {
			if len(t.ParentName) == 0 || t.ParentName == nil {
				graph.addNilParents(t.TaskName)
				continue
			}
			for _, parent := range t.ParentName {
				err = graph.DependOn(t.TaskName, parent)
				if err != nil {
					ResChan <- GraphResult{
						Error:     err,
						Graphs:    nil,
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
			Graphs:    taskExecuteSequence,
			TopicID:   topic.ID,
			TopicName: topic.TopicName,
		}
	}(topic)

	//启动获取ResChan 判断是否启动
	b, _ := CacheImp.GetCache([]byte(global.ISStartGetFromResChan))
	if b == nil {
		go GetTaskTopo()
		CacheImp.SetCache([]byte(global.ISStartGetFromResChan), []byte("true"), -1)
	}
}

// 同步获取
func MakeGraphResSync(topic Topic) GraphResult {
	graph := GraphNew()
	var err error
	for _, t := range topic.Tasks {
		if len(t.ParentName) == 0 || t.ParentName == nil {
			graph.addNilParents(t.TaskName)
			continue
		}
		for _, parent := range t.ParentName {
			err = graph.DependOn(t.TaskName, parent)
			if err != nil {
				return GraphResult{
					Error:     err,
					Graphs:    nil,
					TopicID:   topic.ID,
					TopicName: topic.TopicName,
				}
			}
		}
	}
	taskExecuteSequence := graph.TopoSortedLayers()
	return GraphResult{
		Error:     nil,
		Graphs:    taskExecuteSequence,
		TopicID:   topic.ID,
		TopicName: topic.TopicName,
	}
}

// 循环获取TaskTopo
func GetTaskTopo() {
	for {
		taskTopo := <-ResChan
		logger.Debug().Interface("taskTopo:%v", taskTopo).Msg("")
		if taskTopo.Error != nil {
			logger.Error().Str("get task topo topicName: ", taskTopo.TopicName).
				Str("err:", taskTopo.Error.Error()).Msg("")
		} else {
			key := taskTopo.TopicID + global.TopicTopoSuffix
			taskTopoJson, err := json.Marshal(taskTopo)
			if err != nil {
				logger.Error().Str("topic topo to json error", err.Error()).Msg("to json")
				continue
			}

			CacheImp.SetCache([]byte(key), taskTopoJson, -1)
		}
	}
}

func TaskRun(graphTopo [][]string, topicName, topicID string) (bool, error) {
	b, err := CacheImp.GetCache([]byte(TopicKey))
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

// 执行task sql
func taskRunFromSql(topic Topic, graphTopo [][]string) (bool, error) {
	if DB == nil {
		return false, errors.New("DB is not init")
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
	delTopicFromCache(MakeTopicKey(topic.TopicName, topic.ID))
	return true, nil

}

// 执行task shell
func taskRunFromShell(topic Topic, graphTopo [][]string) (bool, error) {
	return false, nil
}

// 查询执行的sql语句
func queryExecuteStatement(taskID string, topic Topic) (string, error) {
	b, err := CacheImp.GetCache([]byte(topic.TopicName + taskID))
	if err != nil {
		return "", err
	}
	if b == nil {
		err = setTopic2Cache(taskID, topic)
		if err != nil {
			return "", nil
		}
		b, err = CacheImp.GetCache([]byte(topic.TopicName + taskID))
		if err != nil {
			return "", err
		}
	}
	var task Task
	json.Unmarshal(b, &task)

	return task.Comment, nil
}

//执行task中sql
func executeSQL(taskID string, wg *sync.WaitGroup, topic Topic) {
	defer func() {
		if e := recover(); e != nil {
			logger.Panic().Interface("execute sql panic:%v", e).Msg("")
		}
	}()
	defer wg.Done()
	comment, err := queryExecuteStatement(taskID, topic)
	if err != nil {
		panic(err)
	}
	if comment == "" {
		logger.Warn().Str("", "").Msg("task comment is nil")
		return
	}
	//通过；分割sql进行执行
	comments := strings.Split(comment, ";")
	for _, sqlComment := range comments {
		sqlComment = strings.TrimSpace(sqlComment)
		//手动保证并行执行task不操作同一个表
		if err = global.DB.Exec(sqlComment).Error; err != nil {
			panic(err)
		}
	}
}

func setTopic2Cache(taskID string, topic Topic) error {
	b, err := CacheImp.GetCache([]byte(TopicKey))
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
		err = CacheImp.SetCache([]byte(topic.TopicName+task.ID), b, -1)
		if err != nil {
			return err
		}
	}
	return nil
}
