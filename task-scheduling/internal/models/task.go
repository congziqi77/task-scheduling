package models

type Tasks struct {
	TaskList []Task `json:"task_list"`
}

//任务
type Task struct {
	ID        string         `json:"id"`
	TaskName  string      `json:"task_name"`
	Comment   interface{} `json:"comment"`
	Parent    string      `json:"parent"`   //依赖任务id
	Children  string      `json:"children"` //被依赖任务id
	TopicName string      `json:"topic_name"`
}
