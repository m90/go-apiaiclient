package apiaiclient

// Response describes the wrapper around an API.AI response
type Response struct {
	Result Result `json:"result"`
}

// Result describes the result contained in a response
type Result struct {
	Contexts    ContextCollection `json:"contexts"`
	Fulfillment Fulfillment       `json:"fulfillment"`
}

// Fulfillment describes the fulfillment data contained in a response
type Fulfillment struct {
	Speech   string            `json:"speech"`
	Messages MessageCollection `json:"messages"`
}

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

// Context describes a an api-ai conversation context
type Context struct {
	Name       string                  `json:"name" bson:"name"`
	Lifespan   int                     `json:"lifespan" bson:"lifespan"`
	Parameters *map[string]interface{} `json:"parameters,omitempty" bson:"parameters,omitempty"`
}

// ContextCollection proxes []Context and adds helper methods on top
type ContextCollection []Context

// GetUpdate transforms the collection into a []interface{}
func (l ContextCollection) GetUpdate() []interface{} {
	cast := make([]interface{}, len(l))
	for i, item := range l {
		cast[i] = item
	}
	return cast
}

// FilterByContextNames removes any contexts from the collection
// that are of the give name
func (l *ContextCollection) FilterByContextNames(filters ...string) bool {
	filtered := ContextCollection{}
	removal := false
	for _, ctx := range *l {
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
	*l = filtered
	return removal
}

// ContainsContextName checks if the collection contains a context of one
// of the given names without mutating the collection
func (l *ContextCollection) ContainsContextName(filters ...string) bool {
	for _, ctx := range *l {
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
func (l *ContextCollection) FilterByGenericNames(filters ...string) bool {
	filtered := ContextCollection{}
	removal := false
	for _, ctx := range *l {
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
	*l = filtered
	return removal
}

// RequestPayload describes the payload that will be sent to API.AI
type RequestPayload struct {
	Query     string            `json:"query"`
	Contexts  ContextCollection `json:"contexts"`
	SessionID string            `json:"sessionId"`
	Lang      string            `json:"lang"`
}
