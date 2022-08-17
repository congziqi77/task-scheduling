package models

import (
	"reflect"
	"testing"
)

func TestTopic_SaveTopic2Cache(t *testing.T) {
	type fields struct {
		ID        string
		TopicName string
		Desc      string
		Type      int
		Tasks     []Task
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic := &Topic{
				ID:        tt.fields.ID,
				TopicName: tt.fields.TopicName,
				Desc:      tt.fields.Desc,
				Type:      tt.fields.Type,
				Tasks:     tt.fields.Tasks,
			}
			if err := topic.SaveTopic2Cache(); (err != nil) != tt.wantErr {
				t.Errorf("Topic.SaveTopic2Cache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_delTopicFromCache(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := delTopicFromCache(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("delTopicFromCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTopicMapFromCache(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]Topic
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTopicMapFromCache()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopicMapFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopicMapFromCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTopicMapToCache(t *testing.T) {
	type args struct {
		maps map[string]Topic
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetTopicMapToCache(tt.args.maps); (err != nil) != tt.wantErr {
				t.Errorf("SetTopicMapToCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTopicTopo(t *testing.T) {
	type args struct {
		topicName string
		topicID   string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTopicTopo(tt.args.topicName, tt.args.topicID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTopicTopo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTopicTopo() = %v, want %v", got, tt.want)
			}
		})
	}
}
