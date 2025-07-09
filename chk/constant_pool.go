package chk

import (
	"errors"

	"github.com/simonwater/gopression/util"
	"github.com/simonwater/gopression/values"
)

// ConstantPool 常量池实现
type ConstantPool struct {
	constants []*values.Value
	indexMap  map[string]int
	tracer    *util.Tracer
}

// NewConstantPool 创建新的常量池
func NewConstantPool(capacity int, tracer *util.Tracer) *ConstantPool {
	initCap := capacity
	if initCap < 10 {
		initCap = 10
	}

	return &ConstantPool{
		constants: make([]*values.Value, 0, initCap),
		indexMap:  make(map[string]int),
		tracer:    tracer,
	}
}

// NewConstantPoolFromBytes 从字节数组创建常量池
func NewConstantPoolFromBytes(data []byte, tracer *util.Tracer) (*ConstantPool, error) {
	cp := &ConstantPool{
		constants: make([]*values.Value, 0),
		indexMap:  make(map[string]int),
		tracer:    tracer,
	}

	if tracer != nil {
		tracer.StartTimer()
		defer tracer.EndTimer("根据字节数组构造常量池。")
	}

	buffer := util.NewBufferFromBytes(data)
	buffer.SetPosition(0) // 重置位置到开始
	for buffer.Remaining() > 0 {
		val, err := values.GetFrom(buffer)
		if err != nil {
			return nil, err
		}
		cp.constants = append(cp.constants, &val)
	}
	return cp, nil
}

// ToBytes 将常量池转换为字节数组
func (cp *ConstantPool) ToBytes() ([]byte, error) {
	if cp.tracer != nil {
		cp.tracer.StartTimer()
		defer cp.tracer.EndTimer("常量池生成字节数组。")
	}

	buffer := util.NewByteBuffer(0)
	for _, val := range cp.constants {
		if err := val.WriteTo(buffer); err != nil {
			return nil, err
		}
	}
	return buffer.ToBytes(), nil
}

// AddConst 添加常量到池中，返回索引
func (cp *ConstantPool) AddConst(value *values.Value) (int, error) {
	if err := cp.checkType(value); err != nil {
		return -1, err
	}

	key := value.String()
	if idx, exists := cp.indexMap[key]; exists {
		return idx, nil
	}

	idx := len(cp.constants)
	cp.constants = append(cp.constants, value)
	cp.indexMap[key] = idx
	return idx, nil
}

// ReadConst 读取指定索引的常量
func (cp *ConstantPool) ReadConst(index int) (*values.Value, error) {
	if index < 0 || index >= len(cp.constants) {
		return nil, errors.New("常量索引越界")
	}
	return cp.constants[index], nil
}

// GetConstIndex 获取常量的索引
func (cp *ConstantPool) GetConstIndex(constant string) (int, bool) {
	idx, exists := cp.indexMap[constant]
	return idx, exists
}

// GetAllConsts 获取所有常量
func (cp *ConstantPool) GetAllConsts() []*values.Value {
	return cp.constants
}

// Clear 清空常量池
func (cp *ConstantPool) Clear() {
	cp.constants = make([]*values.Value, 0)
	cp.indexMap = make(map[string]int)
}

// checkType 检查常量类型是否支持
func (cp *ConstantPool) checkType(value *values.Value) error {
	switch value.GetValueType() {
	case values.Vt_Integer, values.Vt_Long, values.Vt_Float, values.Vt_Double, values.Vt_String, values.Vt_Boolean:
		return nil
	default:
		return errors.New("常量池中暂不支持此类型：" + value.GetValueType().String())
	}
}
