package models

type Task struct {
	ID       string
	Comment  interface{}
	Parent   string //依赖任务id
	Children string //被依赖任务id
	Topic    string //属于那个topic
}
