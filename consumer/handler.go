package consumer

// Handler is the public interface that must be implemented when using the Consumer object
// With that the code that use the Consumer can focuse on the business logic
// instead of the pubsub implementation
type Handler interface {
	HandleMessage(Message) error
	Finish() error
}

// HandlerConstructor function is the required function,
// that can return a new Handler object that will consume messages during the processing
type HandlerConstructor func() Handler
