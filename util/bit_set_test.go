package util

import (
	"reflect"
	"testing"
)

func TestNewBitSet(t *testing.T) {
	b := NewBitSet(0)
	if len(b.words) != 0 {
		t.Errorf("expected 0 words, got %d", len(b.words))
	}
	b = NewBitSet(1)
	if len(b.words) != 1 {
		t.Errorf("expected 1 word, got %d", len(b.words))
	}
	b = NewBitSet(65)
	if len(b.words) != 2 {
		t.Errorf("expected 2 words, got %d", len(b.words))
	}
}

func TestSetGetClear(t *testing.T) {
	b := NewBitSet(130)
	b.Set(0)
	b.Set(64)
	b.Set(129)
	if !b.Get(0) || !b.Get(64) || !b.Get(129) || b.Get(222) {
		t.Errorf("Set/Get failed")
	}
	b.Clear(64)
	if b.Get(64) {
		t.Errorf("Clear failed")
	}
	if b.Get(128) {
		t.Errorf("Get unset bit should be false")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for negative index")
		}
	}()
	b.Get(-1)
}

func TestToByteArrayAndValueOf(t *testing.T) {
	b := NewBitSet(130)
	b.Set(0)
	b.Set(64)
	b.Set(129)
	bytes := b.ToBytes()
	b2 := NewBitSetFromBytes(bytes)
	for i := 0; i < 130; i++ {
		if b.Get(i) != b2.Get(i) {
			t.Errorf("ValueOf/ToByteArray mismatch at bit %d", i)
		}
	}
	// Test empty
	b3 := NewBitSet(10)
	bytes = b3.ToBytes()
	if len(bytes) != 0 {
		t.Errorf("expected empty byte array")
	}
	b4 := NewBitSetFromBytes(bytes)
	if !b4.IsEmpty() {
		t.Errorf("expected empty BitSet")
	}
}

func TestLengthSizeIsEmpty(t *testing.T) {
	b := NewBitSet(0)
	if b.Length() != 0 || !b.IsEmpty() {
		t.Errorf("empty BitSet should have length 0 and be empty")
	}
	b.Set(10)
	if b.Length() != 11 {
		t.Errorf("expected length 11, got %d", b.Length())
	}
	if b.IsEmpty() {
		t.Errorf("should not be empty after set")
	}
	b.Clear(10)
	if !b.IsEmpty() {
		t.Errorf("should be empty after clear")
	}
	b.Set(63)
	b.Set(64)
	if b.Length() != 65 {
		t.Errorf("expected length 65, got %d", b.Length())
	}
	if b.Size() < 2 {
		t.Errorf("expected at least 2 words, got %d", b.Size())
	}
}

func TestString(t *testing.T) {
	b := NewBitSet(0)
	if b.String() != "{}" {
		t.Errorf("expected {}, got %s", b.String())
	}
	b.Set(1)
	b.Set(3)
	b.Set(5)
	s := b.String()
	expected := "{1, 3, 5}"
	if s != expected {
		t.Errorf("expected %s, got %s", expected, s)
	}
}

func TestExpandTo(t *testing.T) {
	b := NewBitSet(1)
	b.Set(130)
	if !b.Get(130) {
		t.Errorf("expandTo failed to set bit 130")
	}
	if b.Size() < 3 {
		t.Errorf("expandTo did not expand words slice correctly")
	}
}

func TestValueOfPartialBytes(t *testing.T) {
	// Test ValueOf with less than 8 bytes
	bytes := []byte{0x01, 0x00, 0x00, 0x00}
	b := NewBitSetFromBytes(bytes)
	if !b.Get(0) {
		t.Errorf("expected bit 0 to be set")
	}
	if b.Get(1) {
		t.Errorf("expected bit 1 to be unset")
	}
}

func TestToByteArrayTrimming(t *testing.T) {
	b := NewBitSet(128)
	b.Set(0)
	b.Set(127)
	arr := b.ToBytes()
	// Should be 16 bytes, no trailing zeros
	if len(arr) != 16 {
		t.Errorf("expected 16 bytes, got %d", len(arr))
	}
	// trailing zeros trimmed
	b.Clear(127)
	arr2 := b.ToBytes()
	if len(arr2) != 1 {
		t.Errorf("expected 1 byte, got %d", len(arr2))
	}
}

func TestBitSetEquality(t *testing.T) {
	b1 := NewBitSet(100)
	b2 := NewBitSet(100)
	b1.Set(10)
	b2.Set(10)
	if !reflect.DeepEqual(b1.words, b2.words) {
		t.Errorf("BitSets with same bits set should be equal")
	}
}
