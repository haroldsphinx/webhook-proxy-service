package providers

import (
	"errors"
	"strings"
)

// Header parameters
const (
	GitlabProviderKind  = "gitlab"
	ContentTypeHeader = "Content-Type"
	DefaultContentTypeHeader = "application/json"
)

type Event string

type Provider interface {
	GetHeaderKeys() []string
	Validate(hook Hook) bool
	GetCommitter(hook Hook) string
	GetProviderName() string
}

type Hook struct {
	Payload []byte
	Headers map[string]string
	RequestMethod string
}  

func assertProviderImplementations() {
	var _ Provider = (*GitlabProvider)(nil)
}

func NewProvider(provider string, secret string) (Provider, error) {
	if len(provider) == 0 {
		return nil, errors.New("Empty provider string specified")
	}

	switch strings.ToLower(provider) {
	case GitlabProviderKind:
		return NewGitlabProvider(secret)
	default:
		return nil, errors.New("Unknown Git Provider '"+ provider + "' specified")
	}
}

