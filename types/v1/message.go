package v1

import "time"

// Message defines the methods available on a server-side event
type Message interface {
	GetWriteKey() string
	GetMessageID() string
	GetContext() *Context
	GetIntegrations() *Integrations
	GetType() string
	WithSentAt(time.Time)
	WithReceivedAt(time.Time)
	WithRequestMetadata(RequestMetadata)
	MergeContext(*Context) error
	MergeIntegrations(Integrations) error
	SkewTimestamp()
	String() string
	Validate() error
}
