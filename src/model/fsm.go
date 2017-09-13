package model

// EventCurrent 状态机事件和状态
type EventCurrent interface {
	Event(e string) error

	Current() string
}
