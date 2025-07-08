package util

import (
	"reflect"
	"testing"
)

func TestClean(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\nworld\t", "world"},
		{"", ""},
		{"  ", ""},
	}

	for _, tt := range tests {
		result := Clean(tt.input)
		if result != tt.expected {
			t.Errorf("Clean(%q) = %q, 预期 %q", tt.input, result, tt.expected)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\nworld\t", "world"},
		{"", ""},
		{"  ", ""},
	}

	for _, tt := range tests {
		result := Trim(tt.input)
		if result != tt.expected {
			t.Errorf("Trim(%q) = %q, 预期 %q", tt.input, result, tt.expected)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", false},
		{"a", false},
	}

	for _, tt := range tests {
		result := IsEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmpty(%q) = %v, 预期 %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{" ", true},
		{"a", true},
	}

	for _, tt := range tests {
		result := IsNotEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsNotEmpty(%q) = %v, 预期 %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"  \t\n  ", true},
		{" a ", false},
		{"hello", false},
	}

	for _, tt := range tests {
		result := IsBlank(tt.input)
		if result != tt.expected {
			t.Errorf("IsBlank(%q) = %v, 预期 %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsNotBlank(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"  \t\n  ", false},
		{" a ", true},
		{"hello", true},
	}

	for _, tt := range tests {
		result := IsNotBlank(tt.input)
		if result != tt.expected {
			t.Errorf("IsNotBlank(%q) = %v, 预期 %v", tt.input, result, tt.expected)
		}
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		str1     string
		str2     string
		expected bool
	}{
		{"hello", "hello", true},
		{"hello", "HELLO", false},
		{"", "", true},
		{"", "a", false},
	}

	for _, tt := range tests {
		result := Equals(tt.str1, tt.str2)
		if result != tt.expected {
			t.Errorf("Equals(%q, %q) = %v, 预期 %v", tt.str1, tt.str2, result, tt.expected)
		}
	}
}

func TestEqualsIgnoreCase(t *testing.T) {
	tests := []struct {
		str1     string
		str2     string
		expected bool
	}{
		{"hello", "HELLO", true},
		{"Go", "go", true},
		{"", "", true},
		{"a", "b", false},
	}

	for _, tt := range tests {
		result := EqualsIgnoreCase(tt.str1, tt.str2)
		if result != tt.expected {
			t.Errorf("EqualsIgnoreCase(%q, %q) = %v, 预期 %v", tt.str1, tt.str2, result, tt.expected)
		}
	}
}

func TestGetUTF8Bytes(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"hello", []byte("hello")},
		{"中文", []byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}},
		{"", []byte{}},
	}

	for _, tt := range tests {
		result := GetUTF8Bytes(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("GetUTF8Bytes(%q) = %v, 预期 %v", tt.input, result, tt.expected)
		}
	}
}

func TestGetUTF8String(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "hello"},
		{[]byte{0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87}, "中文"},
		{[]byte{}, ""},
		{nil, ""},
	}

	for _, tt := range tests {
		result := GetUTF8String(tt.input)
		if result != tt.expected {
			t.Errorf("GetUTF8String(%v) = %q, 预期 %q", tt.input, result, tt.expected)
		}
	}
}
