package util

import (
	"bytes"
	"math"
	"testing"
)

func TestNewByteBufferAndBasicPutGet(t *testing.T) {
	bb := NewByteBuffer(2)
	if bb.capacity != 2 {
		t.Errorf("expected capacity 2, got %d", bb.capacity)
	}
	bb.Put(0x12)
	bb.Put(0x34)
	if bb.Position() != 2 {
		t.Errorf("expected position 2, got %d", bb.Position())
	}
	bb.SetPosition(0) // 重置位置
	val, err := bb.Get()
	if err != nil || val != 0x12 {
		t.Errorf("expected 0x12, got %x, err=%v", val, err)
	}
	val, err = bb.Get()
	if err != nil || val != 0x34 {
		t.Errorf("expected 0x34, got %x, err=%v", val, err)
	}
	_, err = bb.Get()
	if err == nil {
		t.Error("expected buffer underflow error")
	}
}

func TestPutGetShortIntDouble(t *testing.T) {
	bb := NewByteBuffer(2) // 自动扩容
	bb.PutShort(0x1234)
	bb.PutInt(0x56789abc)
	bb.PutDouble(3.1415926)
	bb.SetPosition(0)
	s, err := bb.GetShort()
	if err != nil || s != 0x1234 {
		t.Errorf("expected 0x1234, got %x, err=%v", s, err)
	}
	i, err := bb.GetInt()
	if err != nil || i != 0x56789abc {
		t.Errorf("expected 0x56789abc, got %x, err=%v", i, err)
	}
	d, err := bb.GetDouble()
	if err != nil || math.Abs(d-3.1415926) > 1e-7 {
		t.Errorf("expected 3.1415926, got %f, err=%v", d, err)
	}
}

func TestEndian(t *testing.T) {
	bb := NewByteBuffer(8)
	bb.SetEndian(true) // 小端
	bb.PutShort(0x1234)
	bb.PutInt(0x56789abc)
	bb.PutDouble(1.5)
	bb.SetPosition(0)
	s, _ := bb.GetShort()
	if s != 0x1234 {
		t.Errorf("little endian short failed, got %x", s)
	}
	i, _ := bb.GetInt()
	if i != 0x56789abc {
		t.Errorf("little endian int failed, got %x", i)
	}
	d, _ := bb.GetDouble()
	if math.Abs(d-1.5) > 1e-7 {
		t.Errorf("little endian double failed, got %f", d)
	}
}

func TestPutIntAt(t *testing.T) {
	bb := NewByteBuffer(8)
	bb.PutInt(0x11111111)
	err := bb.PutIntAt(0, 0x22222222)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	bb.SetPosition(0)
	i, _ := bb.GetInt()
	if i != 0x22222222 {
		t.Errorf("PutIntAt failed, got %x", i)
	}
	err = bb.PutIntAt(100, 0x33333333)
	if err == nil {
		t.Error("expected out of range error")
	}
}

func TestPutBytesAndToBytes(t *testing.T) {
	bb := NewByteBuffer(2)
	bb.PutBytes([]byte{1, 2, 3, 4})
	out := bb.ToBytes()
	if !bytes.Equal(out, []byte{1, 2, 3, 4}) {
		t.Errorf("ToBytes failed, got %v", out)
	}
}

func TestFromBytes(t *testing.T) {
	src := []byte{9, 8, 7}
	bb := FromBytes(src)
	if bb.Position() != 3 {
		t.Errorf("expected position 3, got %d", bb.Position())
	}
	out := bb.ToBytes()
	if !bytes.Equal(out, src) {
		t.Errorf("FromBytes failed, got %v", out)
	}
}

func TestClearAndRemaining(t *testing.T) {
	bb := NewByteBuffer(4)
	bb.Put(1)
	bb.Put(2)
	bb.Clear()
	if bb.Position() != 0 {
		t.Errorf("Clear failed, position=%d", bb.Position())
	}
	if bb.Remaining() != 4 {
		t.Errorf("Remaining failed, got %d", bb.Remaining())
	}
}

func TestSlice(t *testing.T) {
	bb := NewByteBuffer(5)
	bb.PutBytes([]byte{1, 2, 3, 4, 5})
	s, err := bb.Slice(1, 4)
	if err != nil || !bytes.Equal(s, []byte{2, 3, 4}) {
		t.Errorf("Slice failed, got %v, err=%v", s, err)
	}
	_, err = bb.Slice(-1, 2)
	if err == nil {
		t.Error("expected out of range error")
	}
}

func TestCopyFrom(t *testing.T) {
	src := NewByteBuffer(4)
	src.PutBytes([]byte{1, 2, 3, 4})
	dst := NewByteBuffer(2)
	err := dst.CopyFrom(src, 1, 4)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !bytes.Equal(dst.ToBytes(), []byte{2, 3, 4}) {
		t.Errorf("CopyFrom failed, got %v", dst.ToBytes())
	}
	err = dst.CopyFrom(src, -1, 2)
	if err == nil {
		t.Error("expected out of bounds error")
	}
}

func TestSetPosition(t *testing.T) {
	bb := NewByteBuffer(3)
	err := bb.SetPosition(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	err = bb.SetPosition(-1)
	if err == nil {
		t.Error("expected out of range error")
	}
	err = bb.SetPosition(4)
	if err == nil {
		t.Error("expected out of range error")
	}
}
func TestEnsureCapacity(t *testing.T) {
	bb := NewByteBuffer(2)
	bb.Put(1)
	bb.Put(2)
	if bb.Remaining() != 0 {
		t.Errorf("expected remaining 0, got %d", bb.Remaining())
	}
	bb.Put(3) // should trigger auto-expand
	if bb.Remaining() < 1 {
		t.Errorf("expected remaining >= 1 after auto-expand, got %d", bb.Remaining())
	}
	if bb.Position() != 3 {
		t.Errorf("expected position 3 after putting 3 bytes, got %d", bb.Position())
	}
}
func TestPutGetLargeData(t *testing.T) {
	bb := NewByteBuffer(256)
	data := make([]byte, 512)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}
	bb.PutBytes(data)

	if bb.Position() != len(data) {
		t.Errorf("expected position %d, got %d", len(data), bb.Position())
	}

	out := bb.ToBytes()
	if !bytes.Equal(out, data) {
		t.Errorf("expected data to match, got %v", out)
	}
}
