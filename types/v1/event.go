package v1

import (
	"time"

	"github.com/imdario/mergo"
)

// Event contains the common fields in all events.
type Event struct {
	AnonymousID       string        `json:"anonymousId,omitempty"`
	Context           *Context      `json:"context,omitempty"`
	Integrations      *Integrations `json:"integrations,omitempty"`
	MessageID         string        `json:"messageId,omitempty"`
	OriginalTimestamp Timestamp     `json:"originalTimestamp,omitempty"`
	ReceivedAt        Timestamp     `json:"receivedAt,omitempty"`
	SentAt            Timestamp     `json:"sentAt,omitempty"`
	Timestamp         Timestamp     `json:"timestamp,omitempty"`
	Type              string        `json:"type,omitempty"`
	UserID            string        `json:"userId,omitempty"`
	Version           string        `json:"version,omitempty"`
	WriteKey          string        `json:"writeKey,omitempty"`
}

// GetWriteKey returns the write key
func (ev *Event) GetWriteKey() string {
	return ev.WriteKey
}

// GetMessageID returns the message's ID
func (ev *Event) GetMessageID() string {
	return ev.MessageID
}

// GetMessageID returns the message's Context
func (ev *Event) GetContext() *Context {
	return ev.Context
}

// GetMessageID returns the message's Integrations
func (ev *Event) GetIntegrations() *Integrations {
	return ev.Integrations
}

// GetType returns the write key
func (ev *Event) GetType() string {
	return ev.Type
}

// WithSentAt sets the sentAt date
func (ev *Event) WithSentAt(t time.Time) {
	ev.SentAt.Time = t
}

// WithReceivedAt sets the receivedAt date
func (ev *Event) WithReceivedAt(t time.Time) {
	ev.ReceivedAt.Time = t
}

// SkewTimestamp corrects for incorrect client clocks
// TODO: Ensure correctness
func (ev *Event) SkewTimestamp() {
	// Set client-side timestamp to OriginalTimestamp
	ev.OriginalTimestamp = ev.Timestamp

	// Set timestamp to Timestamp +- (ReveivedAt - SentAt)
	ev.Timestamp = Timestamp{ev.Timestamp.Add(ev.ReceivedAt.Sub(ev.SentAt.Time))}
}

// WithRequestMetadata adds server side information to the message
func (ev *Event) WithRequestMetadata(m RequestMetadata) {
	// Make sure we have a context to work with
	ev.ensureContext()

	// Set the IP
	ev.Context.IP = m.IP

	// Set WriteKey if we need to
	if len(ev.WriteKey) == 0 && len(m.WriteKey) > 0 {
		ev.WriteKey = m.WriteKey
	}
}

// MergeContext merges another context onto this events
func (ev *Event) MergeContext(ctx *Context) error {
	ev.ensureContext()
	return mergo.Merge(ev.Context, ctx)
}

// MergeIntegrations merges another set of integrations onto this events
func (ev *Event) MergeIntegrations(intg Integrations) error {
	ev.ensureIntegrations()
	return mergo.Map(ev.Integrations, intg)
}

func (ev *Event) ensureContext() {
	if ev.Context == nil {
		ev.Context = &Context{}
	}
}

func (ev *Event) ensureIntegrations() {
	if ev.Integrations == nil {
		ev.Integrations = &Integrations{}
	}
}
