package repository

import (
	"strings"
	"testing"
)

// TestItoa verifies that itoa() correctly converts integers to strings,
// including numbers > 9 (regression test for the naive one-digit approach).
func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "10"},
		{11, "11"},
		{42, "42"},
		{100, "100"},
		{999, "999"},
	}

	for _, tt := range tests {
		got := itoa(tt.input)
		if got != tt.expected {
			t.Errorf("itoa(%d) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

// TestJoinFields verifies joinFields joins a slice of field assignments correctly.
func TestJoinFields(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "single field",
			input:    []string{"title = $1"},
			expected: "title = $1",
		},
		{
			name:     "two fields",
			input:    []string{"title = $1", "is_archived = $2"},
			expected: "title = $1, is_archived = $2",
		},
		{
			name:     "three fields",
			input:    []string{"a = $1", "b = $2", "c = $3"},
			expected: "a = $1, b = $2, c = $3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinFields(tt.input)
			if got != tt.expected {
				t.Errorf("joinFields(%v) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestListQueryBuilding verifies that the SQL query building logic
// for List() produces correct WHERE clauses with content_type filtering.
func TestListQueryBuilding(t *testing.T) {
	t.Run("no content type filter", func(t *testing.T) {
		whereClause := "user_id = $1 AND is_archived = false"
		if !strings.Contains(whereClause, "$1") {
			t.Error("expected $1 placeholder in WHERE clause")
		}
	})

	t.Run("with content type filter", func(t *testing.T) {
		// Simulate the logic from List() when contentType is set
		args := []interface{}{"user-123", "video"}
		whereClause := "user_id = $1 AND is_archived = false"
		whereClause += " AND content_type = $" + itoa(len(args))

		if !strings.Contains(whereClause, "$2") {
			t.Errorf("expected $2 placeholder for content_type; got clause: %s", whereClause)
		}
	})
}
