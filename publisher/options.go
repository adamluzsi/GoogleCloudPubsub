package publisher

// Options is the public interface when working with options
type Options interface {
	Configure(config *config)
}

type config struct {
}

func newConfig() *config {
	return &config{}
}

func (c *config) consumeOptions(opts []Options) {
	for _, opt := range opts {
		opt.Configure(c)
	}
}
