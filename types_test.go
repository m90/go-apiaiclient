package apiaiclient

import (
	"testing"
)

func TestFilterContextsByName(t *testing.T) {
	tests := []struct {
		name           string
		collection     ContextCollection
		contextName    string
		expectedLen    int
		expectedResult bool
	}{
		{
			"default",
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
				Context{Name: "baz"},
			},
			"foo",
			2,
			true,
		},
		{
			"no match",
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
				Context{Name: "baz"},
			},
			"zalgo",
			3,
			false,
		},
		{
			"empty collection",
			ContextCollection{},
			"zalgo",
			0,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			removed := test.collection.FilterByContextNames(test.contextName)
			if len(test.collection) != test.expectedLen {
				t.Errorf("Expected length of %v, got %v", test.expectedLen, len(test.collection))
			}
			if removed != test.expectedResult {
				t.Errorf("Expected return value of %v, got %v", test.expectedResult, removed)
			}
		})
	}
}

func TestFilterByGenericNames(t *testing.T) {
	tests := []struct {
		name           string
		collection     ContextCollection
		genericName    string
		expectedLen    int
		expectedResult bool
	}{
		{
			"default",
			ContextCollection{
				Context{
					Name:     "generic",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"filterme": true,
					},
				},
				Context{
					Name:     "bar",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"something": false,
					},
				},
				Context{
					Name:     "baz",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"something": 1337,
					},
				},
			},
			"filterme",
			2,
			true,
		},
		{
			"no match",
			ContextCollection{
				Context{
					Name:     "generic",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"donotfilterme": true,
					},
				},
				Context{
					Name:     "bar",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"something": false,
					},
				},
				Context{
					Name:     "baz",
					Lifespan: 99,
					Parameters: &map[string]interface{}{
						"something": 1337,
					},
				},
			},
			"baz",
			3,
			false,
		},
		{
			"empty collection",
			ContextCollection{},
			"filter",
			0,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			removed := test.collection.FilterByGenericNames(test.genericName)
			if len(test.collection) != test.expectedLen {
				t.Errorf("Expected length of %v, got %v", test.expectedLen, len(test.collection))
			}
			if removed != test.expectedResult {
				t.Errorf("Expected return value of %v, got %v", test.expectedResult, removed)
			}
		})
	}
}

func TestContainsContextName(t *testing.T) {
	tests := []struct {
		name       string
		collection ContextCollection
		filters    []string
		expected   bool
	}{
		{
			"default",
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
			},
			[]string{"bar"},
			true,
		},
		{
			"miss",
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
			},
			[]string{"baz"},
			false,
		},
		{
			"empty",
			ContextCollection{},
			[]string{"baz"},
			false,
		},
		{
			"multi arg",
			ContextCollection{
				Context{Name: "foo"},
			},
			[]string{"baz", "foo"},
			true,
		},
		{
			"nested",
			ContextCollection{
				Context{
					Name:       "bar",
					Parameters: &map[string]interface{}{"foo": 99},
				},
			},
			[]string{"foo"},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := test.collection.ContainsContextName(test.filters...); result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestSelectPlatformMessages(t *testing.T) {
	tests := []struct {
		name        string
		collection  MessageCollection
		platform    string
		expectedLen int
	}{
		{
			"mixed",
			MessageCollection{
				Message{Platform: "foo"},
				Message{Platform: "bar"},
				Message{},
			},
			"foo",
			2,
		},
		{
			"no platform match",
			MessageCollection{
				Message{},
				Message{},
				Message{},
			},
			"foo",
			3,
		},
		{
			"all match",
			MessageCollection{
				Message{Platform: "foo"},
				Message{Platform: "foo"},
				Message{Platform: "foo"},
			},
			"foo",
			3,
		},
		{
			"empty collection",
			MessageCollection{},
			"foo",
			0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.collection.SelectPlatformMesssages(test.platform)
			if len(test.collection) != test.expectedLen {
				t.Errorf("Expected %v messages, got %v", test.expectedLen, len(test.collection))
			}
		})
	}
}

func TestGetUpdate(t *testing.T) {
	tests := []struct {
		name       string
		collection ContextCollection
	}{
		{
			"default",
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
				Context{Name: "baz"},
			},
		},
		{
			"empty",
			ContextCollection{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			update := test.collection.GetUpdate()
			if len(update) != len(test.collection) {
				t.Error("Expected update to be of same length as collection")
			}
		})
	}
}
