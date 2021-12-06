package sender

import (
	model "github.com/ozonmp/lic-license-api/internal/model/license"
)

type LicenseEventSender interface {
	Send(license *model.LicenseEvent) error
}
