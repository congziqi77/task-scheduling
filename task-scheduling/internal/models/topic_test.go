package models

import (
	"errors"
	"testing"

	"github.com/congziqi77/task-scheduling/internal/mock"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

/*
*
goMock进行mock单元测试 配合convey
*/
func TestTopic_SaveTopic2Cache(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	//通过生成的源文件mock出一个接口的实现类
	mockCache := mock.NewMockICache(ctl)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "noError",
			wantErr: false,
		},
		{
			name:    "hasError",
			wantErr: true,
		},
	}
	topic := Topic{
		TopicName: "mockTest",
		Desc:      "mock test",
		Type:      0,
	}

	b := []byte(`{"test:1559820807885033472":{"id":"1559820807885033472","topic_name":"test","desc":"传输测试","type":0,"tasks":null}}`)
	gomock.InOrder(
		mockCache.EXPECT().GetCache(gomock.Any()).Return(b, nil),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockCache.EXPECT().GetCache(gomock.Any()).Return(b, errors.New("Get Error")).Times(1),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Set Error")),
		mockCache.EXPECT().GetCache(gomock.Any()).Return(b, nil),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Set Error")),
	)
	//将mock接口赋值给全局变量
	CacheImp = mockCache
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := topic.SaveTopic2Cache(); (err != nil) != tt.wantErr {
				t.Errorf("test is Error:%v", err.Error())
			}
		})
	}

}

func TestGetTopicTopo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	//通过生成的源文件mock出一个接口的实现类
	mockCache := mock.NewMockICache(ctl)
	b := []byte(`{
		"test:1":{
			"id":"1",
			"topic_name":"test",
			"desc":"传输测试",
			"type":0,
			"tasks":[
				{
					"id":"1",
					"task_name":"A",
					"comment":"test A",
					"parent_name":[
	
					],
					"topic_name":"test"
				},
				{
					"id":"2",
					"task_name":"B",
					"comment":"test B",
					"parent_name":[
						"A"
					],
					"topic_name":"test"
				},
				{
					"id":"3",
					"task_name":"C",
					"comment":"test C",
					"parent_name":[
						"A",
						"B"
					],
					"topic_name":"test"
				},
				{
					"id":"4",
					"task_name":"D",
					"comment":"test D",
					"parent_name":[
						"B"
					],
					"topic_name":"test"
				}
			]
		}
	}`)
	gomock.InOrder(
		mockCache.EXPECT().GetCache(gomock.Any()).Return(nil, nil),
		mockCache.EXPECT().GetCache(gomock.Any()).Return(b, nil),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)
	//将mock接口赋值给全局变量
	CacheImp = mockCache
	Convey("get topo sync", t, func() {
		s, err := GetTopicTopo("test", "1")
		So(s, ShouldResemble, [][]string{{"A"}, {"B"}, {"C", "D"}})
		So(err, ShouldBeEmpty)
	})
}
