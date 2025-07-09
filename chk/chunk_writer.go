package chk

import (
	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

// ChunkWriter 用于写入字节码块
type ChunkWriter struct {
	codeBuffer *util.ByteBuffer // 字节码缓冲区
	constPool  *ConstantPool
	isVarConst *util.BitSet
	tracer     *util.Tracer
}

// NewChunkWriter 创建新的块写入器
func NewChunkWriter(capacity int, tracer *util.Tracer) *ChunkWriter {
	initCap := capacity
	if initCap < 128 {
		initCap = 128
	}

	return &ChunkWriter{
		codeBuffer: util.NewByteBuffer(capacity),
		constPool:  NewConstantPool(capacity, tracer),
		isVarConst: util.NewBitSet(0),
		tracer:     tracer,
	}
}

// Flush 将数据刷新到块
func (cw *ChunkWriter) Flush() *Chunk {
	codeBytes := cw.codeBuffer.ToBytes()
	constBytes, _ := cw.constPool.ToBytes()
	varBytes := cw.isVarConst.ToBytes()

	return &Chunk{
		Codes:     codeBytes,
		Constants: constBytes,
		Vars:      varBytes,
	}
}

// Clear 清空写入器
func (cw *ChunkWriter) Clear() {
	cw.codeBuffer.Clear()
	cw.constPool.Clear()
	cw.isVarConst = util.NewBitSet(0)
}

// WriteByte 写入一个字节
func (cw *ChunkWriter) WriteByte(value byte) {
	cw.codeBuffer.Put(value)
}

// WriteShort 写入一个短整型 (16位)
func (cw *ChunkWriter) WriteShort(value int16) {
	cw.codeBuffer.PutShort(value)
}

// WriteInt 写入一个整型 (32位)
func (cw *ChunkWriter) WriteInt(value int32) {
	cw.codeBuffer.PutInt(value)
}

// UpdateInt 更新指定位置的整数值
func (cw *ChunkWriter) UpdateInt(index int, value int32) {
	cw.codeBuffer.PutIntAt(index, value)
}

// WriteCode 写入操作码
func (cw *ChunkWriter) WriteCode(opCode OpCode) {
	cw.WriteByte(byte(opCode))
}

// AddConstant 添加常量到池中
func (cw *ChunkWriter) AddConstant(value *values.Value) (int, error) {
	return cw.constPool.AddConst(value)
}

// SetVariables 设置变量集合
func (cw *ChunkWriter) SetVariables(vars []string) {
	n := len(vars)
	cw.isVarConst = util.NewBitSet(n)
	for _, varName := range vars {
		index, exists := cw.constPool.GetConstIndex(varName)
		if !exists {
			value := values.NewStringValue(varName)
			index, _ = cw.constPool.AddConst(&value)
		}
		cw.isVarConst.Set(index)
	}
}

// Position 获取当前写入位置
func (cw *ChunkWriter) Position() int {
	return cw.codeBuffer.Position()
}
