package util

import (
	"testing"
)

func TestParseDotSeparatedStringAndToString(t *testing.T) {
	src := "a.b.c.d"
	field := NewFieldFromPath(src)
	if field.String() != src {
		t.Errorf("expected %s, got %s", src, field.String())
	}
}

func TestHandleSingleName(t *testing.T) {
	src := "table1"
	field := NewFieldFromPath(src)
	if field.String() != src {
		t.Errorf("expected %s, got %s", src, field.String())
	}
	if field.GetName() != "table1" {
		t.Errorf("expected name table1, got %s", field.GetName())
	}
	if field.GetOwner() != nil {
		t.Errorf("expected owner nil, got %v", field.GetOwner())
	}
}

func TestHandleNestedField(t *testing.T) {
	owner := NewFieldFromPath("table1")
	field := NewFieldWithOwner("field1", owner)
	if field.String() != "table1.field1" {
		t.Errorf("expected table1.field1, got %s", field.String())
	}
	if field.GetName() != "field1" {
		t.Errorf("expected name field1, got %s", field.GetName())
	}
	if field.GetOwner() != owner {
		t.Errorf("expected owner %v, got %v", owner, field.GetOwner())
	}
}

func TestHandleRootField(t *testing.T) {
	field := NewFieldWithOwner("f1", nil)
	if field.String() != "f1" {
		t.Errorf("expected f1, got %s", field.String())
	}
	if field.GetName() != "f1" {
		t.Errorf("expected name f1, got %s", field.GetName())
	}
	if field.GetOwner() != nil {
		t.Errorf("expected owner nil, got %v", field.GetOwner())
	}
}

func TestHandleRepeatedToStringCalls(t *testing.T) {
	field := NewFieldFromPath("a.b")
	s1 := field.String()
	s2 := field.String()
	if s1 != "a.b" || s2 != "a.b" {
		t.Errorf("expected a.b, got %s and %s", s1, s2)
	}
}
