package consumer

import "cloud.google.com/go/pubsub"

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

// Message is the interface that must be implemented by the Handler HandleMessage method
type Message interface {
	Done(bool)
	Data() []byte
}
