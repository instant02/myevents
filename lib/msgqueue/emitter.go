package msqqueue

type EventEmitter interface{
	Emit(event Event) error
}