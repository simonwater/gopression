package chk

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
