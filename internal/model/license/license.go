package license

type License struct {
	ID    uint64
	Title string
}

type EventLicType uint8

type EventStatus uint8

const (
	Created EventLicType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type LicenseEvent struct {
	ID     uint64
	Type   EventLicType
	Status EventStatus
	Entity *License
}
