package msqqueue

type Event interface{
	EventName() string
}