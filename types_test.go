package apiaiclient

import (
	"testing"
)

func TestFilterContextsByName(t *testing.T) {
	tests := []struct {
		collection     ContextCollection
		name           string
		expectedLen    int
		expectedResult bool
	}{
		{
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
			ContextCollection{},
			"zalgo",
			0,
			false,
		},
	}
	for _, test := range tests {
		removed := test.collection.FilterByContextNames(test.name)
		if len(test.collection) != test.expectedLen {
			t.Errorf("Expected length of %v, got %v", test.expectedLen, len(test.collection))
		}
		if removed != test.expectedResult {
			t.Errorf("Expected return value of %v, got %v", test.expectedResult, removed)
		}
	}
}

func TestFilterByGenericNames(t *testing.T) {
	tests := []struct {
		collection     ContextCollection
		genericName    string
		expectedLen    int
		expectedResult bool
	}{
		{
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
			ContextCollection{},
			"filter",
			0,
			false,
		},
	}
	for _, test := range tests {
		removed := test.collection.FilterByGenericNames(test.genericName)
		if len(test.collection) != test.expectedLen {
			t.Errorf("Expected length of %v, got %v", test.expectedLen, len(test.collection))
		}
		if removed != test.expectedResult {
			t.Errorf("Expected return value of %v, got %v", test.expectedResult, removed)
		}
	}
}

func TestSelectPlatformMessages(t *testing.T) {
	tests := []struct {
		collection  MessageCollection
		platform    string
		expectedLen int
	}{
		{
			MessageCollection{
				Message{Platform: "foo"},
				Message{Platform: "bar"},
				Message{},
			},
			"foo",
			2,
		},
		{
			MessageCollection{
				Message{},
				Message{},
				Message{},
			},
			"foo",
			3,
		},
		{
			MessageCollection{
				Message{Platform: "foo"},
				Message{Platform: "foo"},
				Message{Platform: "foo"},
			},
			"foo",
			3,
		},
		{
			MessageCollection{},
			"foo",
			0,
		},
	}
	for _, test := range tests {
		test.collection.SelectPlatformMesssages(test.platform)
		if len(test.collection) != test.expectedLen {
			t.Errorf("Expected %v messages, got %v", test.expectedLen, len(test.collection))
		}
	}
}

func TestGetUpdate(t *testing.T) {
	tests := []struct {
		collection ContextCollection
	}{
		{
			ContextCollection{
				Context{Name: "foo"},
				Context{Name: "bar"},
				Context{Name: "baz"},
			},
		},
		{
			ContextCollection{},
		},
	}
	for _, test := range tests {
		update := test.collection.GetUpdate()
		if len(update) != len(test.collection) {
			t.Error("Expected update to be of same length as collection")
		}
	}
}
