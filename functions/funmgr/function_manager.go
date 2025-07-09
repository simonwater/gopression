package funmgr

import (
	"sync"

	"github.com/simonwater/gopression/functions"
	"github.com/simonwater/gopression/functions/impl/numfunc"
	"github.com/simonwater/gopression/functions/impl/sysfunc"
)

type FunctionManager struct {
	functions map[string]functions.CallableFunction
	mu        sync.RWMutex
}

var (
	instance     *FunctionManager
	onceInstance sync.Once
)

func GetFunctionManager() *FunctionManager {
	onceInstance.Do(func() {
		instance = &FunctionManager{
			functions: make(map[string]functions.CallableFunction),
		}
		// 注册内置函数
		instance.RegistFunction(numfunc.NewAbs())
		instance.RegistFunction(sysfunc.NewClock())
	})
	return instance
}

func (fm *FunctionManager) GetFunction(name string) functions.CallableFunction {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.functions[name]
}

func (fm *FunctionManager) RegistFunction(fn functions.CallableFunction) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.functions[fn.GetName()] = fn
}

func (fm *FunctionManager) RemoveFunction(name string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	delete(fm.functions, name)
}
