package spot

import (
	"os"
	"strings"
)

// ClientOptions contains options to configure a Client instance. Each option can
// be set through setter functions. See documentation for each setter function for
// an explanation of the option.
type ClientOptions struct {
	// Static user credentials.
	Token, Account string

	// BaseURL configures the default base URL of the Spot API.
	// Defaults to `https://api.spotinst.io`.
	BaseURL string

	// UserAgent configures the User-Agent HTTP header to set when invoking HTTP
	// requests.
	UserAgent string

	// DryRun configures the client to print the actions that would be executed,
	// without executing them.
	DryRun bool
}

// ClientOption allows specifying various settings configurable by the client.
type ClientOption func(*ClientOptions)

// WithCredentials specifies static credentials.
func WithCredentials(token, account string) ClientOption {
	return func(opts *ClientOptions) {
		opts.Token = strings.TrimSpace(token)
		opts.Account = strings.TrimSpace(account)
	}
}

// WithBaseURL defines the base URL of the Spot API.
func WithBaseURL(url string) ClientOption {
	return func(opts *ClientOptions) {
		opts.BaseURL = strings.TrimSpace(url)
	}
}

// WithUserAgent defines the User Agent.
func WithUserAgent(ua string) ClientOption {
	return func(opts *ClientOptions) {
		opts.UserAgent = strings.TrimSpace(ua)
	}
}

// WithDryRun toggles the dry-run mode on/off.
func WithDryRun(value bool) ClientOption {
	return func(opts *ClientOptions) {
		opts.DryRun = value
	}
}

func defaultOptions() *ClientOptions {
	return &ClientOptions{
		Token:   os.Getenv(EnvCredentialsToken),
		Account: os.Getenv(EnvCredentialsAccount),
		BaseURL: os.Getenv(EnvBaseURL),
	}
}
