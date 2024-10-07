package args

import (
	"flag"
	"fmt"
)

const (
	configPathParam = "config"
)

// clientArgs stores cli flags passed by user at client application startup.
type clientArgs struct {
	configPath string
}

// validateClientParams verify parsed cli flags.
func (c *clientArgs) validateClientParams() error {
	type paramCheck struct {
		value string
		name  string
	}

	params := []paramCheck{
		{c.ConfigDir(), configPathParam},
	}

	for _, item := range params {
		if item.value == "" {
			flag.Usage()
			return fmt.Errorf("%q is required", item.name)
		}
	}
	return nil
}

// ConfigDir get client's config path.
func (c *clientArgs) ConfigDir() string {
	return c.configPath
}

// NewClientFlags client application flags parser, returns clientArgs
// or error if something went wrong during flag parsing.
func NewClientFlags() (clientArgs, error) {
	args := clientArgs{}

	flag.StringVar(&args.configPath, configPathParam, "", "Path to the client application config directory")
	flag.Parse()

	if err := args.validateClientParams(); err != nil {
		return clientArgs{}, err
	}

	return args, nil
}
