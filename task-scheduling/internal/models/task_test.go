package models_test

import (
	"log"
	"testing"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/mock"
	. "github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/setting"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "task test")
}

var ctl *gomock.Controller

var _ = BeforeSuite(func() {
	ctl = gomock.NewController(&testing.T{})
	mockCache := mock.NewMockICache(ctl)
	b := []byte("{\"test:1\":{\"id\":\"1\",\"topic_name\":\"test\",\"desc\":\"传输测试\",\"type\":0,\"tasks\":null}}")
	topicKeyByte := []byte(
		`{
			"topicList":{
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
								"B",
								"C"
							],
							"topic_name":"test"
						}
					]
				}
			}
		}`)
	gomock.InOrder(
		mockCache.EXPECT().GetCache([]byte(TopicKey)).Return(b, nil),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockCache.EXPECT().GetCache(gomock.Any()).Return([]byte("true"), nil),
		mockCache.EXPECT().GetCache([]byte(TopicKey)).Return(topicKeyByte, nil).AnyTimes(),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)
	CacheImp = mockCache
	set, err := setting.NewSetting()
	Ω(err).ShouldNot(HaveOccurred())
	if err = set.ReadSection("database", &global.DbSetting); err != nil {
		Ω(err).ShouldNot(HaveOccurred())
	}
	DB, _ = NewDBImp()
})

var _ = AfterSuite(func() {
	ctl.Finish()
})

var _ = Describe("Task", Ordered, func() {
	var (
		tasks Tasks
		topic Topic
	)

	BeforeEach(func() {
		tasks = Tasks{
			[]Task{
				{
					TaskName:   "A",
					Comment:    "select * from task_a;",
					ParentName: []string{},
				},
				{
					TaskName:   "B",
					Comment:    "select * from task_a;",
					ParentName: []string{"A"},
				},
				{
					TaskName:   "C",
					Comment:    "select * from task_a;",
					ParentName: []string{"B"},
				},
				{
					TaskName:   "D",
					Comment:    "select * from task_a;",
					ParentName: []string{"B"},
				},
			},
		}

		topic = Topic{
			ID:        "1",
			TopicName: "test",
			Desc:      "test",
			Type:      0,
			Tasks:     tasks.TaskList,
		}

	})

	Describe("test async create Server", func() {
		Context("without error", func() {
			It("should be no error", func() {
				err := tasks.TaskCreateServer("test", "1")
				Ω(err).ShouldNot(HaveOccurred())
				s := <-ResChan
				log.Printf("s is %v", s.Graphs)
			})
		})
	})

	Describe("test sync create Server", func() {
		Context("withError error", func() {
			It("should be equals", func() {
				gTrue := GraphResult{
					Error:     nil,
					Graphs:    [][]string{{"A"}, {"B"}, {"C", "D"}},
					TopicID:   "1",
					TopicName: "test",
				}
				log.Println("gTrue.Graphs is", gTrue.Graphs)
				Ω(MakeGraphResSync(topic)).Should(Equal(gTrue))
			})
		})
	})

	Describe("sql exec Run test", func() {
		Context("without error", func() {
			It("should be no error", func() {
				set, err := setting.NewSetting()
				Ω(err).ShouldNot(HaveOccurred())
				err = set.ReadSection("database", &global.DbSetting)
				Ω(err).ShouldNot(HaveOccurred())
				isBool, err := TaskRun([][]string{{"A"}, {"B"}, {"C", "D"}}, "test", "1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(isBool).Should(BeTrue())
			})
		})
	})

})
