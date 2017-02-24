package consumer

import (
	"cloud.google.com/go/pubsub"
)

// MessageWrapper is the wrapper struct for *pubsub.Message so
// it can be easily tested without the need of the Proxy Server
type MessageWrapper struct {
	message *pubsub.Message
}

func newMessageWrapper(msg *pubsub.Message) *MessageWrapper {
	return &MessageWrapper{message: msg}
}

// Data is the proxy function to retrive *pubsub.Message.Data content
func (mw *MessageWrapper) Data() []byte {
	return mw.message.Data
}

// Done is the proxy function for *pubsub.Message object
func (mw *MessageWrapper) Done(ackType bool) {
	mw.message.Done(ackType)
}

// MessageInterface is the interface that must be implemented by the Handler object
type MessageInterface interface {
	Done(bool)
	Data() []byte
}

// Handler is the public interface that must be implemented when using the Consumer object
// With that the code that use the Consumer can focuse on the business logic
// instead of the pubsub implementation
type Handler interface {
	HandleMessage(MessageInterface) error
	Finish() error
}

// HandlerConstructor function is the required function,
// that can return a new Handler object that will consume messages during the processing
type HandlerConstructor func() Handler
