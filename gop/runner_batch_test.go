package gop_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/simonwater/gopression/gop"
	"github.com/simonwater/gopression/gop/testdata"
	fileutil "github.com/simonwater/gopression/util/files"
	"github.com/stretchr/testify/require"
)

const (
	gop_formulaBatches = 1000
	gop_testDirectory  = "BatchRunnerTest"
)

// 测试解析执行
func TestBatchRunner_IR(t *testing.T) {
	fmt.Printf("批量运算测试(解析执行)")
	lines := testdata.GetExpressions(gop_formulaBatches)
	runner := gop.NewGopRunner()
	runner.SetExecuteMode(gop.SyntaxTree)
	runner.SetTrace(true)
	env := testdata.GetEnv(gop_formulaBatches)

	_, err := runner.ExecuteBatch(lines, env)
	require.NoError(t, err, "批量执行失败")

	testdata.CheckValues(t, env, gop_formulaBatches)
	fmt.Printf("==========")
}

// 测试编译+字节码执行
func TestBatchRunner_CompileChunk(t *testing.T) {
	fmt.Printf("批量运算测试(编译+字节码执行)")
	start := time.Now()

	lines := testdata.GetExpressions(gop_formulaBatches)
	runner := gop.NewGopRunner()
	runner.SetTrace(true)

	chunk, err := runner.CompileSource(lines)
	require.NoError(t, err, "编译失败")

	env := testdata.GetEnv(gop_formulaBatches)
	_ = runner.RunChunk(chunk, env)

	testdata.CheckValues(t, env, gop_formulaBatches)

	elapsed := time.Since(start)
	fmt.Printf("总耗时: %s", elapsed)
	fmt.Printf("==========")

	// 序列化字节码
	fileName := "Chunks.xp"
	filePath := fileutil.GetTestPath(gop_testDirectory, fileName)
	err = writeChkFile(chunk, filePath)
	require.NoError(t, err, "序列化字节码失败")
}

// 测试字节码直接执行
func TestBatchRunner_Chunk(t *testing.T) {
	fmt.Printf("批量运算测试(字节码直接执行)\n")
	start := time.Now()

	// 反序列化字节码
	fileName := "Chunks.xp"
	filePath := fileutil.GetTestPath(gop_testDirectory, fileName)
	chunk, err := readChkFile(filePath)
	require.NoError(t, err, "反序列化字节码失败")

	elapsed := time.Since(start)
	fmt.Printf("完成从文件反序列化字节码。耗时: %s\n", elapsed)

	runner := gop.NewGopRunner()
	runner.SetTrace(true)
	env := testdata.GetEnv(gop_formulaBatches)

	_ = runner.RunChunk(chunk, env)
	require.NoError(t, err, "执行字节码失败")

	testdata.CheckValues(t, env, gop_formulaBatches)
	fmt.Printf("==========\n")
}
