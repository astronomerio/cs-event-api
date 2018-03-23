package v1

import "time"

// Message defines the methods available on a server-side event
type Message interface {
	WithSentAt(time.Time)
	WithReceivedAt(time.Time)
	WithRequestMetadata(RequestMetadata)
	SkewTimestamp()
	MergeContext(*Context) error
	MergeIntegrations(Integrations) error
	String() string
	GetWriteKey() string
	GetMessageID() string
	GetType() string
}
