package chk

import (
	"fmt"

	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

// ChunkReader 用于读取字节码块
type ChunkReader struct {
	codeBuffer *util.ByteBuffer
	constPool  *ConstantPool
	isVarConst *util.BitSet
	tracer     *util.Tracer
}

// NewChunkReader 创建新的块读取器
func NewChunkReader(chunk *Chunk, tracer *util.Tracer) *ChunkReader {
	codeBuffer := util.NewBufferFromBytes(chunk.Codes)
	codeBuffer.SetPosition(0) // 重置位置到开始
	constPool, err := NewConstantPoolFromBytes(chunk.Constants, tracer)
	if err != nil {
		panic(fmt.Sprintf("Failed to create constant pool from chunk data: %v", err))
	}
	return &ChunkReader{
		codeBuffer: codeBuffer,
		constPool:  constPool,
		isVarConst: util.NewBitSetFromBytes(chunk.Vars),
		tracer:     tracer,
	}
}

// ReadByte 读取一个字节
func (cr *ChunkReader) ReadByte() (byte, error) {
	b, err := cr.codeBuffer.Get()
	return b, err
}

// ReadShort 读取一个短整型 (16位)
func (cr *ChunkReader) ReadShort() (int16, error) {
	s, err := cr.codeBuffer.GetShort()
	return s, err
}

// ReadInt 读取一个整型 (32位)
func (cr *ChunkReader) ReadInt() (int32, error) {
	i, err := cr.codeBuffer.GetInt()
	return i, err
}

// ReadOpCode 读取操作码
func (cr *ChunkReader) ReadOpCode() (OpCode, error) {
	code, err := cr.ReadByte()
	if err != nil {
		return OP_EXIT, err
	}
	return OpCodeFromValue(code)
}

// ReadConst 读取指定索引的常量
func (cr *ChunkReader) ReadConst(index int) (*values.Value, error) {
	return cr.constPool.ReadConst(index)
}

// GetVariables 获取所有变量名
func (cr *ChunkReader) GetVariables() []string {
	if cr.tracer != nil {
		cr.tracer.StartTimer()
		defer cr.tracer.EndTimer("构造变量列表")
	}
	allConsts := cr.constPool.GetAllConsts()
	result := make([]string, 10)

	for i, value := range allConsts {
		if cr.isVarConst.Get(i) {
			result = append(result, value.String())
		}
	}

	return result
}

// Position 获取当前读取位置
func (cr *ChunkReader) Position() int {
	return cr.codeBuffer.Position()
}

// NewPosition 设置新的读取位置
func (cr *ChunkReader) NewPosition(pos int) error {
	return cr.codeBuffer.SetPosition(pos)
}
