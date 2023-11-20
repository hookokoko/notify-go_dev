package handler

/**
1. baseHandler只作为各种子handler的基类，不单独使用，不暴露公共的初始化方法
2. 具体使用依赖子类的初始化
*/

type FlowControlParam struct {
}

type FlowControlFactory struct {
}

type AnchorState int

type AnchorInfo struct {
	State       AnchorState
	BizID       string
	MessageID   string
	BusinessID  string
	ReceiverIDs []string
}

type baseHandler struct {
	flowControlParam   *FlowControlParam
	handlerHolder      *HandlerHolder
	flowControlFactory *FlowControlFactory
}

func newBaseHandler(hh *HandlerHolder, fp *FlowControlParam, fc *FlowControlFactory) *baseHandler {
	return &baseHandler{
		flowControlParam:   fp,
		handlerHolder:      hh,
		flowControlFactory: fc,
	}
}
