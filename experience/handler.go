package experience

import (
	"fmt"
	"sync"
)

// FlowControlParam 结构体定义
type FlowControlParam struct {
	// 在这里添加相应的字段
}

// AnchorState 枚举定义
type AnchorState int

const (
	SendSuccess AnchorState = iota
	SendFail
)

// AnchorInfo 结构体定义
type AnchorInfo struct {
	State       AnchorState
	BizID       string
	MessageID   string
	BusinessID  string
	ReceiverIDs []string
}

// Handler 接口定义
type Handler interface {
	Do(taskInfo *TaskInfo) bool
	//Recall()
}

// BaseHandler 结构体定义
type BaseHandler struct {
	channelCode int
	//Handler
	//flowControlParam   *FlowControlParam
	handlerHolder *HandlerHolder
	//flowControlFactory *FlowControlFactory
}

// DoHandler 实现Handler接口
func (b *BaseHandler) handle(taskInfo *TaskInfo) {
	//if b.flowControlParam != nil {
	//	b.flowControlFactory.FlowControl(taskInfo, b.flowControlParam)
	//}

	if b.handlerHolder.route(b.channelCode).Do(taskInfo) {
		fmt.Println(SendSuccess)
		return
	}

	fmt.Println(SendFail)
}

// HandlerHolder 结构体定义
type HandlerHolder struct {
	handlers map[int]Handler
	mu       sync.Mutex
}

// NewHandlerHolder 创建 HandlerHolder 实例
func NewHandlerHolder() *HandlerHolder {
	return &HandlerHolder{
		handlers: make(map[int]Handler),
	}
}

// PutHandler 将处理程序放入映射
func (h *HandlerHolder) route(channelCode int) Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.handlers[channelCode]
}

// PutHandler 将处理程序放入映射
func (h *HandlerHolder) PutHandler(channelCode int, handler Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[channelCode] = handler
}

type EmailHandler struct {
	//BaseHandler
	// 其他字段
}

func (e *EmailHandler) getAccountConfig() bool {
	return true
}

func (e *EmailHandler) Do(taskInfo *TaskInfo) bool {
	fmt.Println("email handler")
	return true
}

/*
*
核心流行说明
 1. 注册各种发送方式的handler，时机待定
 2. 基于渠道码构建base handler，根据不同的渠道码，路由到相应的handler，执行其的do发送
*/
func main() {
	e := EmailHandler{}
	// ...等等其它handler

	hh := NewHandlerHolder()
	hh.PutHandler(20, &e)

	NewBaseHandler(20, hh).handle(&TaskInfo{})
}

func NewBaseHandler(code int, hh *HandlerHolder) *BaseHandler {
	return &BaseHandler{
		channelCode:   code,
		handlerHolder: hh,
	}
}
