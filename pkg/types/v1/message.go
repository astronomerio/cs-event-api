package v1

import (
	"encoding/json"
	"time"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/util"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Type      string                 `json:"type,omitempty"`
	MessageID string                 `json:"messageId,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`

	Timestamp         time.Time `json:"timestamp,omitempty"`
	OriginalTimestamp time.Time `json:"originalTimestamp,omitempty"`
	SentAt            time.Time `json:"sentAt,omitempty"`
	ReceivedAt        time.Time `json:"receivedAt,omitempty"`

	AppID        string                 `json:"appId,oimtempty"`
	WriteKey     string                 `json:"writeKey,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	AnonymousID  string                 `json:"anonymousId,omitempty"`
	UserID       string                 `json:"userId,omitempty"`
	GroupID      string                 `json:"groupId,omitempty"`
	PreviousID   string                 `json:"previousId,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Action       string                 `json:"action,omitempty"`
	Channel      string                 `json:"channel,omitempty"`
	Event        string                 `json:"event,omitempty"`
	Version      string                 `json:"version,omitempty"`
}

// IsValid returns whether or not the message is valid
func (m *Message) IsValid() bool {
	valid := true

	if m.AppID == "" {
		valid = false
	}

	return valid
}

// BindRequest adds the request level fields to the Message
func (m *Message) BindRequest(c *gin.Context) {
	md := GetRequestMetadata(c)
	m.ApplyMetadata(md)
	m.ReceivedAt = util.NowUTC()
}

// ApplyMetadata adds the RequestMetadata to the Message
func (m *Message) ApplyMetadata(metadata RequestMetadata) {
	m.Context["ip"] = metadata.IP

	if metadata.AppID != "" && m.AppID == "" {
		m.AppID = metadata.AppID
	}
}

// SkewTimestamp will skew the time fields by the difference between SentAt and ReceivedAt
func (m *Message) SkewTimestamp() {
	m.Timestamp = m.Timestamp.UTC().Round(time.Millisecond)
	m.SentAt = m.SentAt.UTC().Round(time.Millisecond)

	m.OriginalTimestamp = m.Timestamp

	// SentAt *should* be at most a few seconds earlier than time.Now()
	diff := m.ReceivedAt.Sub(m.SentAt)
	m.Timestamp = m.Timestamp.Add(diff)
}

// FormatTimestamps converts client side timestamps to UTC and rounds them to the nearest
// Millisecond
func (m *Message) FormatTimestamps() {
	m.SentAt = m.SentAt.UTC().Round(time.Millisecond)
	m.Timestamp = m.Timestamp.UTC().Round(time.Millisecond)
}

// MaybeFix will add fields that should be present if they aren't
func (m *Message) MaybeFix() {
	if m.MessageID == "" {
		m.MessageID = uuid.NewV4().String()
	}

	if m.SentAt.IsZero() {
		m.SentAt = util.NowUTC()
	}

	if m.Timestamp.IsZero() {
		m.Timestamp = util.NowUTC()
	}
}

// PartitionKey returns the partition key for this message. Returns the MessageID field currently
func (m *Message) PartitionKey() string {
	return m.MessageID
}

func (m *Message) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}
