//
// GoogleCloudPubsub is a package that aims to give you easy to test GooglePubsub Consumer interface
//
// ## Usage
//
// First of all, you need a Handler struct that will implement the business logic for the Gcloud pubsub subscription consuming:
//
//
//		type ExampleHandler struct {
//			// you can use use fields to store messages if you want for example bulk processing
//			// messages []consumer.Message
//		}
//
//		// HandleMessage method will be called after a message had beed fetched from the pubsub.
//		// return error will Nack the message
//		func (eh *ExampleHandler) HandleMessage(msg consumer.Message) error {
//			// eh.messages = append(eh.messages, msg)
//			// single element processing can be implemented here
//			return nil
//		}
//
//		// Finish method will be called before all messages should be acked
//		// With this method, you can do bulk actions after the HandleMessage Collected all the elements
//		func (eh *ExampleHandler) Finish() error {
//			defer eh.wg.Done()
//			// heavy bulk actions can be implemented here
//			return nil
//		}
//
//		// The return value should be the consumer.Handler interface,
//		// not the actual struct pointer
//		func NewExampleHandler() consumer.Handler {
//			return &ExampleHandler{}
//		}
//
//
// Now with your new fancy struct and with it's constructor function, you can begin to use the Consumer
//
//
//	    ctx := context.Background()
//	    c := consumer.New(ctx, "example-subscription-name", NewExampleHandler)
//
//
// If you want to specify further options for the consumer, you can do so with option setters.
//
//
//		ctx := context.Background()
//		cons := consumer.New(ctx, "example", NewExampleHandler,
//				// you can configure the new consumer to use given amount of BatchSize
//				// This is the amount that will be passed for the HandleMessage method for a single Handler object
//				consumer.SetBatchSizeTo(amount),
//
//				// This will set the Google Pubsub Message Iterators MaxExtensionDuration
//				consumer.SetMaxExtensionDurationTo(10*time.Minute),
//
//				// This will configure the consumer to how manny parallel worker should pull from the subscription
//				consumer.SetWorkersCountTo(runtime.NumCPU()))
//
//
//
// ## Testing
//
// When You test your application, Before the Consumer is being initialized, you should turn on Mock mod.
// When Mock mod enabled, not the original but a Mock consumer will be created when the New method called.
// It's behavior is alike, but remove the Dependency to use Google Pubsub Emulated Host,
// and increase the speed for your tests.
//
// Make even the Benchmarking more valuable
//
//
//		func TestConsumerMockingAllPerfect(t *testing.T) {
//			consumer.TurnMockModOn()
//			defer consumer.TurnMockModOff()
//
//			// consumer creation is the same, and not required to be happen here,
//			// this is just an example , that it should be created after the mock mod enabled
//			ctx := context.Background()
//			c := consumer.New(ctx, "example-subscription-name", NewExampleHandler)
//			c.Start()
//			defer c.Stop()
//
//			// And this is how you Send Messages to the Mock Consumer
//			consumer.MockMessageFeeder["example-subscription-name"] <- []byte(`Hello World!`)
//
//			// super complex business logic testing here
//
//		}
//
package consumer
