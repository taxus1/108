package message

// Sender 消息发送对象
var Sender *sender

// Sender 消息发送
type sender struct {
}

// Send 发送
func (s *sender) Send(user int64, msg Message) {
	//TODO
	// s.xxx(user, msg.Marshal())
}
