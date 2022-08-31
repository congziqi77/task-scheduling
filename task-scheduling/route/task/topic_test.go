package task_test

import (
	"testing"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/mock"
	. "github.com/congziqi77/task-scheduling/internal/models"
	"github.com/congziqi77/task-scheduling/internal/routers"
	"github.com/congziqi77/task-scheduling/internal/utils"
	"github.com/gin-gonic/gin"
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
	// b := []byte("{\"test:1\":{\"id\":\"1\",\"topic_name\":\"test\",\"desc\":\"传输测试\",\"type\":0,\"tasks\":null}}")
	gomock.InOrder(
		mockCache.EXPECT().GetCache([]byte("1"+global.TopicTopoSuffix)).Return(nil, nil).AnyTimes(),
		mockCache.EXPECT().SetCache(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes(),
		mockCache.EXPECT().GetCache(gomock.Any()).Return([]byte("true"), nil).AnyTimes(),
	)
	CacheImp = mockCache
})

var _ = AfterSuite(func() {
	ctl.Finish()
})

var _ = Describe("Topic", Ordered, func() {

	var route *gin.Engine

	BeforeAll(func() {
		route = routers.NewRouter()
	})

	Describe("test GetTopo", func() {
		Context("withOut error", func() {
			It("should no error", func() {
				uri := "/task/topic/getTopo?topicName=1&topicID=1"
				w := utils.Get(uri, route)
				Ω(w.Code).Should(Equal(200))
			})
		})
	})

	Describe("topicList", func() {
		It("topicList request test ", func() {
			w := utils.Get("/task/topic/list", route)
			Ω(w.Code).Should(Equal(200))
		})
	})
})
