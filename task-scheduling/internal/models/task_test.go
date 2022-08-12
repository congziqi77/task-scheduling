package models

import (
	"testing"
)

func TestTasks_TaskCreateServer(t *testing.T) {
	type fields struct {
		TaskList []Task
	}
	type args struct {
		topicName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := &Tasks{
				TaskList: tt.fields.TaskList,
			}
			if err := tasks.TaskCreateServer(tt.args.topicName); (err != nil) != tt.wantErr {
				t.Errorf("Tasks.TaskCreateServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_makeGraphResAsync(t *testing.T) {
	topic := Topic{
		ID:        "1",
		TopicName: "test",
		Desc:      "sql",
		Type:      0,
		Tasks: []Task{
			{
				ID:        "01",
				TaskName:  "task01",
				Comment:   "task01 run",
				ParentId:  nil,
				TopicName: "test1",
			},
			{
				ID:        "02",
				TaskName:  "task02",
				Comment:   "task02 run",
				ParentId:  []string{"01"},
				TopicName: "test2",
			},
			{
				ID:        "03",
				TaskName:  "task03",
				Comment:   "task03 run",
				ParentId:  []string{"01"},
				TopicName: "test3",
			},
			{
				ID:        "04",
				TaskName:  "task04",
				Comment:   "task04 run",
				ParentId:  []string{"01", "02"},
				TopicName: "test4",
			},
		},
	}

	makeGraphResAsync(topic)
	s := <-ResChan
	if s.Error != nil {
		t.Error(s.Error)
	}
	t.Log("s async is", s)

}

func Test_makeGraphResSync(t *testing.T) {
	topic := Topic{
		ID:        "1",
		TopicName: "test",
		Desc:      "sql",
		Type:      0,
		Tasks: []Task{
			{
				ID:        "01",
				TaskName:  "task01",
				Comment:   "task01 run",
				ParentId:  nil,
				TopicName: "test1",
			},
			{
				ID:        "02",
				TaskName:  "task02",
				Comment:   "task02 run",
				ParentId:  []string{"01"},
				TopicName: "test2",
			},
			{
				ID:        "03",
				TaskName:  "task03",
				Comment:   "task03 run",
				ParentId:  []string{"01"},
				TopicName: "test3",
			},
			{
				ID:        "04",
				TaskName:  "task04",
				Comment:   "task04 run",
				ParentId:  []string{"01", "02"},
				TopicName: "test4",
			},
		},
	}

	s, err := makeGraphResSync(topic)
	// s := <-ResChan
	// if s.Error != nil {
	// 	t.Error(s.Error)
	// }
	if err != nil {
		t.Log("err ----------->", err.Error())
	}
	t.Log("s is", s)
}
