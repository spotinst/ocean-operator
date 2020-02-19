package spot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	sdkversion "github.com/operator-framework/operator-sdk/version"
	"github.com/spotinst/ocean-operator/internal/version"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/spotinst-sdk-go/spotinst/credentials"
	sdklog "github.com/spotinst/spotinst-sdk-go/spotinst/log"
	"github.com/spotinst/spotinst-sdk-go/spotinst/session"
	"github.com/spotinst/spotinst-sdk-go/spotinst/util/useragent"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("spot")

// Blank assignment to verify that clientWrapper implements Client.
var _ Client = &clientWrapper{}

type clientWrapper struct {
	client  *client.Client
	session *session.Session
	config  *spotinst.Config
}

// NewClient returns a new Spot client.
func NewClient(options ...ClientOption) Client {
	cfg := spotinst.DefaultConfig()

	// Initialize options.
	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}

	// Configure the base URL.
	if opts.BaseURL != "" {
		cfg.WithBaseURL(opts.BaseURL)
		log.V(1).Info("Configured base URL", "url", opts.BaseURL)
	}

	// Configure the credentials.
	if opts.Token != "" || opts.Account != "" {
		cfg.WithCredentials(credentials.NewStaticCredentials(opts.Token, opts.Account))
		log.V(1).Info("Configured credentials")
	}

	// Configure the user agent.
	agents := useragent.UserAgents{
		useragent.New("ocean-operator", version.String()),
		useragent.New("operator-sdk", sdkversion.Version[1:]),
	}
	userAgent := agents.String()
	if opts.UserAgent != "" && !strings.Contains(userAgent, opts.UserAgent) {
		userAgent += " " + opts.UserAgent
	}
	cfg.WithUserAgent(userAgent)
	log.V(1).Info("Configured user agent", "userAgent", cfg.UserAgent)

	// Configure the logger.
	cfg.WithLogger(sdklog.LoggerFunc(func(format string, args ...interface{}) {
		log.V(5).Info(fmt.Sprintf(format, args...))
	}))

	// Configure the SDK to use a dry-run mode.
	if opts.DryRun {
		cfg.HTTPClient.Transport = new(roundTripperMock)
		cfg.WithCredentials(credentials.NewStaticCredentials("dry-run", "dry-run"))
		log.V(1).Info("Configured dry-run mode")
	}

	return &clientWrapper{
		client:  client.New(cfg),
		session: session.New(cfg),
		config:  cfg,
	}
}

func (x *clientWrapper) Accounts() Accounts {
	return &accounts{client: x.client}
}

func (x *clientWrapper) Services() Services {
	return &services{session: x.session}
}

type roundTripperMock struct{}

func (x *roundTripperMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		Request:    req,
	}, nil
}
