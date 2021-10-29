package sender

import (
	"github.com/ozonmp/lic-license-api/internal/model"
)

type EventSender interface {
	Send(subdomain *model.SubdomainEvent) error
}
