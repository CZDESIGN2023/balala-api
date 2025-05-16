package v8

import (
	"errors"
	"go-cs/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// nilByte is used in JSON marshal/unmarshal
	nilByte = []byte("null")

	// ErrNoClient is raised when no Elasticsearch node is available.
	ErrNoClient = errors.New("no Elasticsearch node available")

	// ErrRetry is raised when a request cannot be executed after the configured
	// number of retries.
	ErrRetry = errors.New("cannot connect after several retries")

	// ErrTimeout is raised when a request timed out, e.g. when WaitForStatus
	// didn't return in time.
	ErrTimeout = errors.New("timeout")

	// noRetries is a retrier that does not retry.
	// noRetries = NewStopRetrier()

	// noDeprecationLog is a no-op for logging deprecations.
	// noDeprecationLog = func(*http.Request, *http.Response) {}
)

type Config struct {
	Addresses              []string // A list of Elasticsearch nodes to use.
	Username               string   // Username for HTTP Basic Authentication.
	Password               string   // Password for HTTP Basic Authentication.
	EnableDebugLogger      bool     // Enable the debug logging.
	CertificateFingerprint string
}

func NewConfig(data *conf.Data) *Config {
	c := &Config{
		Addresses:              data.Es.Addresses,
		Username:               data.Es.Username,
		Password:               data.Es.Password,
		EnableDebugLogger:      data.Es.EnableDebugLogger,
		CertificateFingerprint: data.Es.CertificateFingerprint,
	}

	return c
}

func NewEsClient(config *Config, logger log.Logger) (*elasticsearch.Client, func(), error) {

	cfg := elasticsearch.Config{
		Addresses:         config.Addresses,
		Username:          config.Username,
		Password:          config.Password,
		EnableDebugLogger: config.EnableDebugLogger,
	}

	if config.CertificateFingerprint != "" {
		cfg.CertificateFingerprint = config.CertificateFingerprint
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the es client")
	}

	return client, cleanup, nil
}
