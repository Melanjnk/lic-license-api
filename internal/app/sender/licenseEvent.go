package sender

import (
	"github.com/ozonmp/lic-license-api/internal/model"
)

type LicenseEventSender interface {
	Send(license *model.LicenseEvent) error
}
