package models

//任务
type Task struct {
	ID       int         `json:"id,omitempty"`
	Comment  interface{} `json:"comment,omitempty"`
	Parent   string      `json:"parent,omitempty"`   //依赖任务id
	Children string      `json:"children,omitempty"` //被依赖任务id
}
