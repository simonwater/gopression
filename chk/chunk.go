package chk

import (
	"github.com/simonwater/gopression/util"
)

// Chunk 表示执行块的数据结构
type Chunk struct {
	Codes     []byte // 字节码
	Constants []byte // 常量池
	Vars      []byte // 变量信息
}

// NewChunk 创建新的空 Chunk
func NewChunk() *Chunk {
	return &Chunk{
		Codes:     []byte{},
		Constants: []byte{},
		Vars:      []byte{},
	}
}

// NewChunkWithData 使用给定数据创建 Chunk
func NewChunkWithData(codes, constants, vars []byte) *Chunk {
	return &Chunk{
		Codes:     codes,
		Constants: constants,
		Vars:      vars,
	}
}

func NewChunkWithBytes(bytes []byte) *Chunk {
	buf := util.NewBufferFromBytes(bytes)
	buf.SetPosition(0)
	codeSz, _ := buf.GetInt()
	codes, _ := buf.GetBytes(int(codeSz))
	constSz, _ := buf.GetInt()
	constants, _ := buf.GetBytes(int(constSz))
	varSz, _ := buf.GetInt()
	vars, _ := buf.GetBytes(int(varSz))
	return NewChunkWithData(codes, constants, vars)
}

func (c *Chunk) ToBytes() []byte {
	sz := c.GetByteSize() + 3*4
	buf := util.NewByteBuffer(sz)
	buf.PutInt(int32(c.GetCodesSize()))
	buf.PutBytes(c.Codes)
	buf.PutInt(int32(c.GetConstsSize()))
	buf.PutBytes(c.Constants)
	buf.PutInt(int32(c.GetVarsSize()))
	buf.PutBytes(c.Vars)
	return buf.ToBytes()
}

// GetByteSize 获取总字节大小
func (c *Chunk) GetByteSize() int {
	return len(c.Codes) + len(c.Constants) + len(c.Vars)
}

// GetCodesSize 获取字节码大小
func (c *Chunk) GetCodesSize() int {
	return len(c.Codes)
}

// GetConstsSize 获取常量池大小
func (c *Chunk) GetConstsSize() int {
	return len(c.Constants)
}

// GetVarsSize 获取变量信息大小
func (c *Chunk) GetVarsSize() int {
	return len(c.Vars)
}
