package apiaiclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Response describes the wrapper around an API.AI response
type Response struct {
	Result Result `json:"result"`
}

// Result describes the result contained in a response
type Result struct {
	Contexts    ContextCollection `json:"contexts"`
	Fulfillment Fulfillment       `json:"fulfillment"`
	Metadata    Metadata          `json:"metadata"`
}

// Fulfillment describes the fulfillment data contained in a response
type Fulfillment struct {
	Speech   string            `json:"speech"`
	Messages MessageCollection `json:"messages"`
}

// Metadata contains metadata about a result
type Metadata struct {
	IntentID                  string `json:"intentId"`
	IntentName                string `json:"intentName"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	WebhookUsed               string `json:"webhookUsed"`
	WebhookResponseTime       int64  `json:"webhookResponseTime"`
}

// ResponseTime returns the structs `WebhookResponseTime` as a time.Duration
func (m *Metadata) ResponseTime() time.Duration {
	if value, err := time.ParseDuration(fmt.Sprintf("%dms", m.WebhookResponseTime)); err == nil {
		return value
	}
	return 0
}

// constants that will populate a message's `type` field
const (
	MessageTypeText          = 0
	MessageTypeCardMessage   = 1
	MessageTypeQuickReplies  = 2
	MessageTypeImage         = 3
	MessageTypeCustomPayload = 4
)

// Message describes a message contained in a server payload
type Message struct {
	Type     int                               `json:"type"`
	Platform string                            `json:"platform,omitempty"`
	Speech   string                            `json:"speech,omitempty"`
	ImageURL string                            `json:"imageUrl,omitempty"`
	Title    string                            `json:"title,omitempty"`
	Subtitle string                            `json:"subtitle,omitempty"`
	Buttons  []Button                          `json:"buttons,omitempty"`
	Replies  []string                          `json:"replies,omitempty"`
	Payload  map[string]map[string]interface{} `json:"payload,omitempty"`
}

// Button describes a button included in a message
type Button struct {
	Text     string `json:"text"`
	Postback string `json:"postback"`
}

// MessageCollection proxies []Message and adds helper methods on top
type MessageCollection []Message

// SelectPlatformMesssages removes all messsages from the collection
// that do not match the given platform string
func (l *MessageCollection) SelectPlatformMesssages(platform string) {
	platformMessages := &MessageCollection{}
	for _, msg := range *l {
		if msg.Platform == "" || msg.Platform == platform {
			*platformMessages = append(*platformMessages, msg)
		}
	}
	*l = *platformMessages
}

// Context describes a an api.ai conversation context
type Context struct {
	Name       string                  `json:"name" bson:"name"`
	Lifespan   int                     `json:"lifespan" bson:"lifespan"`
	Parameters *map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
}

// ContextCollection adds helper methods to []Context
type ContextCollection []Context

// FilterByContextNames removes any contexts from the collection
// that are of the given name
func (c *ContextCollection) FilterByContextNames(filters ...string) bool {
	filtered := ContextCollection{}
	removal := false
	for _, ctx := range *c {
		match := false
		for _, name := range filters {
			if ctx.Name == name {
				removal = true
				match = true
				break
			}
		}
		if !match {
			filtered = append(filtered, ctx)
		}
	}
	*c = filtered
	return removal
}

// FilterParametersByKey removes all parameter values whose key contains
// one of the given tokens
func (c *ContextCollection) FilterParametersByKey(tokens ...string) bool {
	removal := false
	for _, ctx := range *c {
		if ctx.Parameters == nil {
			continue
		}
		params := map[string]interface{}{}
	param:
		for key, value := range *ctx.Parameters {
			for _, token := range tokens {
				if strings.Contains(key, token) {
					removal = true
					continue param
				}
			}
			params[key] = value
		}
		*ctx.Parameters = params
	}
	return removal
}

// ContainsContextName checks if the collection contains a context of one
// of the given names without mutating the collection
func (c *ContextCollection) ContainsContextName(filters ...string) bool {
	for _, ctx := range *c {
		for _, name := range filters {
			if ctx.Name == name {
				return true
			}
		}
	}
	return false
}

// FilterByGenericNames removes any context from the collection that is
// of name generic and the generic context contains any of the given keys
func (c *ContextCollection) FilterByGenericNames(filters ...string) bool {
	filtered := ContextCollection{}
	removal := false
	for _, ctx := range *c {
		if ctx.Name != "generic" || ctx.Parameters == nil {
			filtered = append(filtered, ctx)
			continue
		}
		match := false
		for _, key := range filters {
			if _, ok := (*ctx.Parameters)[key]; ok {
				removal = true
				match = true
				break
			}
		}
		if !match {
			filtered = append(filtered, ctx)
		}
	}
	*c = filtered
	return removal
}

// MarshalJSON returns the JSON representation of the collection
func (c *ContextCollection) MarshalJSON() ([]byte, error) {
	out := []Context{}
	for _, item := range *c {
		out = append(out, item)
	}
	return json.Marshal(out)
}

// RequestPayload describes the payload that will be sent to api.ai
type RequestPayload struct {
	Query     string         `json:"query"`
	Contexts  json.Marshaler `json:"contexts"`
	SessionID string         `json:"sessionId"`
	Lang      string         `json:"lang"`
}
