package util

import (
	"encoding/binary"
	"errors"
	"math"
)

type ByteBuffer struct {
	buf          []byte
	position     int
	capacity     int
	littleEndian bool
}

// NewByteBuffer 创建指定容量的 ByteBuffer
func NewByteBuffer(capacity int) *ByteBuffer {
	return &ByteBuffer{
		buf:      make([]byte, capacity),
		capacity: capacity,
	}
}

// FromBytes 用已有字节创建 ByteBuffer
func FromBytes(bytes []byte) *ByteBuffer {
	bb := NewByteBuffer(len(bytes))
	copy(bb.buf, bytes)
	bb.position = len(bytes)
	return bb
}

// SetEndian 设置字节序（true=小端，false=大端，默认大端）
func (bb *ByteBuffer) SetEndian(littleEndian bool) {
	bb.littleEndian = littleEndian
}

// Position 获取当前位置
func (bb *ByteBuffer) Position() int {
	return bb.position
}

// SetPosition 设置当前位置
func (bb *ByteBuffer) SetPosition(pos int) error {
	if pos < 0 || pos > bb.capacity {
		return errors.New("position out of range")
	}
	bb.position = pos
	return nil
}

// Remaining 返回剩余空间
func (bb *ByteBuffer) Remaining() int {
	return bb.capacity - bb.position
}

// Clear 重置 position
func (bb *ByteBuffer) Clear() {
	bb.position = 0
}

// EnsureCapacity 自动扩容
func (bb *ByteBuffer) ensureCapacity(required int) {
	if bb.Remaining() >= required {
		return
	}
	newCap := bb.capacity*2 + required
	newBuf := make([]byte, newCap)
	copy(newBuf, bb.buf)
	bb.buf = newBuf
	bb.capacity = newCap
}

// Put 写入单字节
func (bb *ByteBuffer) Put(val byte) {
	bb.ensureCapacity(1)
	bb.buf[bb.position] = val
	bb.position++
}

// PutBytes 写入字节数组
func (bb *ByteBuffer) PutBytes(bytes []byte) {
	bb.ensureCapacity(len(bytes))
	copy(bb.buf[bb.position:], bytes)
	bb.position += len(bytes)
}

// PutShort 写入 int16
func (bb *ByteBuffer) PutShort(val int16) {
	bb.ensureCapacity(2)
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	order.PutUint16(bb.buf[bb.position:], uint16(val))
	bb.position += 2
}

// PutInt 写入 int32
func (bb *ByteBuffer) PutInt(val int32) {
	bb.ensureCapacity(4)
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	order.PutUint32(bb.buf[bb.position:], uint32(val))
	bb.position += 4
}

// PutIntAt 指定位置写 int32
func (bb *ByteBuffer) PutIntAt(index int, val int32) error {
	if index < 0 || index+4 > bb.capacity {
		return errors.New("index out of range")
	}
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	order.PutUint32(bb.buf[index:], uint32(val))
	return nil
}

// PutDouble 写入 float64
func (bb *ByteBuffer) PutDouble(val float64) {
	bb.ensureCapacity(8)
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	bits := math.Float64bits(val)
	order.PutUint64(bb.buf[bb.position:], bits)
	bb.position += 8
}

// Get 读取单字节
func (bb *ByteBuffer) Get() (byte, error) {
	if bb.position >= bb.capacity {
		return 0, errors.New("buffer underflow")
	}
	val := bb.buf[bb.position]
	bb.position++
	return val, nil
}

// GetShort 读取 int16
func (bb *ByteBuffer) GetShort() (int16, error) {
	if bb.position+2 > bb.capacity {
		return 0, errors.New("buffer underflow")
	}
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	val := int16(order.Uint16(bb.buf[bb.position:]))
	bb.position += 2
	return val, nil
}

// GetInt 读取 int32
func (bb *ByteBuffer) GetInt() (int32, error) {
	if bb.position+4 > bb.capacity {
		return 0, errors.New("buffer underflow")
	}
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	val := int32(order.Uint32(bb.buf[bb.position:]))
	bb.position += 4
	return val, nil
}

// GetDouble 读取 float64
func (bb *ByteBuffer) GetDouble() (float64, error) {
	if bb.position+8 > bb.capacity {
		return 0, errors.New("buffer underflow")
	}
	var order binary.ByteOrder = binary.BigEndian
	if bb.littleEndian {
		order = binary.LittleEndian
	}
	bits := order.Uint64(bb.buf[bb.position:])
	bb.position += 8
	return math.Float64frombits(bits), nil
}

// ToBytes 返回 [0, position) 区间的字节
func (bb *ByteBuffer) ToBytes() []byte {
	return bb.buf[:bb.position]
}

// Slice 返回指定区间的字节
func (bb *ByteBuffer) Slice(start, end int) ([]byte, error) {
	if start < 0 || end > bb.capacity || start > end {
		return nil, errors.New("slice out of range")
	}
	return bb.buf[start:end], nil
}

// CopyFrom 从另一个 ByteBuffer 拷贝数据
func (bb *ByteBuffer) CopyFrom(src *ByteBuffer, srcStart, srcEnd int) error {
	if srcStart < 0 || srcEnd > src.position || srcStart > srcEnd {
		return errors.New("source range out of bounds")
	}
	length := srcEnd - srcStart
	bb.ensureCapacity(length)
	copy(bb.buf[bb.position:], src.buf[srcStart:srcEnd])
	bb.position += length
	return nil
}
