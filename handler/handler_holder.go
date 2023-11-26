package handler

import (
	"sync"
)

var handlerHolder *HandlerHolder = NewHandlerHolder()

// HandlerHolder 结构体定义
type HandlerHolder struct {
	handlers map[string]Handler
	mu       sync.Mutex
}

// NewHandlerHolder 创建 HandlerHolder 实例
func NewHandlerHolder() *HandlerHolder {
	hh := &HandlerHolder{
		handlers: make(map[string]Handler),
	}

	return hh
}

// PutHandler 将处理程序放入映射
func (h *HandlerHolder) route(channel string) Handler {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.handlers[channel]
}

// PutHandler 将处理程序放入映射
func (h *HandlerHolder) put(channel string, handler Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[channel] = handler
}
