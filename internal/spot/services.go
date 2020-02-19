package spot

import (
	"fmt"

	"github.com/spotinst/spotinst-sdk-go/service/ocean"
	"github.com/spotinst/spotinst-sdk-go/spotinst/session"
)

type services struct {
	session *session.Session
}

func (x *services) Ocean(provider CloudProviderName) (Ocean, error) {
	svc := ocean.New(x.session)

	switch provider {
	case CloudProviderAWS:
		return &oceanAWS{svc.CloudProviderAWS()}, nil
	case CloudProviderGCP:
		return &oceanGCP{svc.CloudProviderGCP()}, nil
	default:
		return nil, fmt.Errorf("spotinst: unsupported cloud provider: %s", provider)
	}
}
