package util

import (
	"encoding/binary"
	"fmt"
)

// BitSet 兼容 Java 的 BitSet 二进制结构
type BitSet struct {
	words []uint64
}

// NewBitSet 创建一个指定 bit 数的 BitSet
func NewBitSet(nbits int) *BitSet {
	nwords := 0
	if nbits > 0 {
		nwords = (nbits-1)/64 + 1
	}
	return &BitSet{
		words: make([]uint64, nwords),
	}
}

// ValueOf 从二进制数据恢复 BitSet（兼容 Java/TS 版本的 toByteArray）
func NewBitSetFromBytes(bytes []byte) *BitSet {
	nwords := (len(bytes) + 7) / 8
	words := make([]uint64, nwords)
	for i := 0; i < nwords; i++ {
		if i*8+8 <= len(bytes) {
			words[i] = binary.LittleEndian.Uint64(bytes[i*8 : i*8+8])
		} else {
			// 剩余不足 8 字节
			var tmp [8]byte
			copy(tmp[:], bytes[i*8:])
			words[i] = binary.LittleEndian.Uint64(tmp[:])
		}
	}
	return &BitSet{words: words}
}

// Get 获取某一位
func (b *BitSet) Get(bitIndex int) bool {
	if bitIndex < 0 {
		panic(fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}
	word := bitIndex / 64
	if word >= len(b.words) {
		return false
	}
	return (b.words[word] & (1 << (bitIndex % 64))) != 0
}

// Set 设置某一位
func (b *BitSet) Set(bitIndex int) {
	if bitIndex < 0 {
		panic(fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}
	word := bitIndex / 64
	if word >= len(b.words) {
		b.expandTo(word)
	}
	b.words[word] |= 1 << (bitIndex % 64)
}

// Clear 清除某一位
func (b *BitSet) Clear(bitIndex int) {
	if bitIndex < 0 {
		panic(fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}
	word := bitIndex / 64
	if word < len(b.words) {
		b.words[word] &^= 1 << (bitIndex % 64)
	}
}

// ToByteArray 转为二进制（兼容 Java/TS 版本的 toByteArray）
func (b *BitSet) ToBytes() []byte {
	last := len(b.words) - 1
	for last >= 0 && b.words[last] == 0 {
		last--
	}
	if last < 0 {
		return []byte{}
	}
	bytes := make([]byte, (last+1)*8)
	for i := 0; i <= last; i++ {
		binary.LittleEndian.PutUint64(bytes[i*8:], b.words[i])
	}
	// 去除末尾多余的 0 字节
	realLen := len(bytes)
	for realLen > 0 && bytes[realLen-1] == 0 {
		realLen--
	}
	return bytes[:realLen]
}

// Length 返回最高 set 位的下标+1
func (b *BitSet) Length() int {
	last := len(b.words) - 1
	for last >= 0 && b.words[last] == 0 {
		last--
	}
	if last < 0 {
		return 0
	}
	word := b.words[last]
	for i := 63; i >= 0; i-- {
		if (word & (1 << uint(i))) != 0 {
			return last*64 + i + 1
		}
	}
	return 0
}

// Size 返回 word 数
func (b *BitSet) Size() int {
	return len(b.words)
}

// IsEmpty 判断是否全为 0
func (b *BitSet) IsEmpty() bool {
	for _, w := range b.words {
		if w != 0 {
			return false
		}
	}
	return true
}

func (b *BitSet) expandTo(word int) {
	for len(b.words) <= word {
		b.words = append(b.words, 0)
	}
}

func (b *BitSet) String() string {
	if b.IsEmpty() {
		return "{}"
	}
	res := ""
	for i := 0; i < b.Length(); i++ {
		if b.Get(i) {
			if len(res) > 0 {
				res += ", "
			}
			res += fmt.Sprintf("%d", i)
		}
	}
	return "{" + res + "}"
}
