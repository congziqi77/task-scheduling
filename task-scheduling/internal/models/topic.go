package models

type Topic struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name" binding:"required,max=50"`
	Desc      string `json:"desc" binding:"required"`
	Type      int    `json:"type" binding:"required,oneof=0 1"` //0 sql 1 shell
	Tasks     []Task `json:"tasks"`
}
