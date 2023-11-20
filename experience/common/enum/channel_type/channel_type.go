package channel_type

// ChannelType 发送渠道类型枚举
type ChannelType int

const (
	// PUSH 通知栏
	_    ChannelType = iota
	PUSH ChannelType = 10 * iota
	// SMS 短信
	SMS ChannelType = 10 * iota
	// EMAIL 邮件
	EMAIL ChannelType = 10 * iota
)

func (c ChannelType) String() string {
	switch c {
	case PUSH:
		return "PUSH"
	case SMS:
		return "SMS"
	case EMAIL:
		return "EMAIL"
	default:
		return "UNKNOWN"
	}
}

func Values() []string {
	return []string{"PUSH", "SMS", "EMAIL"}
}
