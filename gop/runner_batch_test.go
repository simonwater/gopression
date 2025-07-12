package gop_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/simonwater/gopression/chk"
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
	fmt.Println("===批量运算测试(解析语法树执行)")
	lines := testdata.GetExpressions(gop_formulaBatches)
	runner := gop.NewGopRunner()
	runner.SetExecuteMode(gop.SyntaxTree)
	runner.SetTrace(true)
	env := testdata.GetEnv(gop_formulaBatches)

	_, err := runner.ExecuteBatch(lines, env)
	require.NoError(t, err, "批量执行失败")

	testdata.CheckValues(t, env, gop_formulaBatches)
	fmt.Printf("==========\n")
}

// 测试编译+字节码执行
func TestBatchRunner_CompileChunk(t *testing.T) {
	fmt.Println("===批量运算测试(编译+字节码执行)")
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
	fmt.Printf("==========\n")
}

// 测试编译到文件再从文件读取执行
func TestBatchRunner_MultiChunks(t *testing.T) {
	fmt.Println("===字节码编译到文件再从文件读取执行")
	filePath := fileutil.GetTestPath(gop_testDirectory, "Chunks.pb")
	chunk, err := createAndGetChunk(t, filePath)
	require.NoError(t, err, "字节码创建失败")
	fmt.Printf("字节码大小%vKB\n", chunk.GetByteSize()/1024)

	start, cnt := time.Now(), 1
	for i := 0; i < cnt; i++ {
		runner := gop.NewGopRunner()
		env := testdata.GetEnv(gop_formulaBatches)
		_ = runner.RunChunk(chunk, env)
		require.NoError(t, err, "执行字节码失败")
		testdata.CheckValues(t, env, gop_formulaBatches)
	}
	fmt.Printf("执行完成。执行次数：%v。耗时: %s\n", cnt, time.Since(start))
	fmt.Printf("==========\n")
}

func createAndGetChunk(t *testing.T, filePath string) (*chk.Chunk, error) {
	lines := testdata.GetExpressions(gop_formulaBatches)
	runner := gop.NewGopRunner()
	start := time.Now()
	chunk, err := runner.CompileSource(lines)
	require.NoError(t, err, "编译失败")
	fmt.Printf("编译完成。耗时: %s\n", time.Since(start))

	// 序列化字节码
	start = time.Now()
	err = writeChkFile(chunk, filePath)
	require.NoError(t, err, "序列化字节码失败")
	fmt.Printf("序列化到文件完成。耗时: %s\n", time.Since(start))

	// 反序列化
	start = time.Now()
	chunk, err = readChkFile(filePath)
	fmt.Printf("从文件反序列化完成。耗时: %s\n", time.Since(start))
	return chunk, err
}
