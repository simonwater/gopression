package gop_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/simonwater/gopression/chk"
	"github.com/simonwater/gopression/gop"
	"github.com/simonwater/gopression/gop/testdata"
	fileutil "github.com/simonwater/gopression/util/files"
	"github.com/stretchr/testify/require"
)

const (
	ser_formulaBatches = 1000
	ser_testDirectory  = "SerializeTest"
)

func TestChunkSerialization(t *testing.T) {
	fmt.Printf("序列化反序列化测试：\n")
	chunkSerializeTest(t)
}

func chunkSerializeTest(t *testing.T) {
	// 创建表达式列表
	lines := testdata.GetExpressions(ser_formulaBatches)
	fmt.Printf("表达式总数：%d\n", len(lines))

	runner := gop.NewGopRunner()

	// 解析和分析表达式
	fmt.Printf("开始解析和分析：\n")
	start := time.Now()
	exprs, err := runner.Parse(lines)
	require.NoError(t, err, "解析失败")

	exprInfos := runner.Analyze(exprs)
	require.NoError(t, err, "分析失败")

	elapsed := time.Since(start)
	fmt.Printf("中间结果生成完成。耗时: %s\n", elapsed)

	// 编译为字节码
	runner.SetTrace(true)
	chunk := runner.CompileIR(exprInfos)

	// 序列化字节码
	chunkSize := chunk.GetByteSize()
	fmt.Printf("开始进行字节码序列化反序列化，字节码大小(KB): %d\n", chunkSize/1024)

	start = time.Now()
	fileName := "Chunks.pb"
	filePath := fileutil.GetTestPath(ser_testDirectory, fileName)

	err = writeChkFile(chunk, filePath)
	require.NoError(t, err, "序列化失败")
	elapsed = time.Since(start)
	fmt.Printf("字节码已序列化到文件：%s 耗时: %s\n", fileName, elapsed)

	// 反序列化字节码
	start = time.Now()
	deserializedChunk, err := readChkFile(filePath)
	require.NoError(t, err, "反序列化失败")
	elapsed = time.Since(start)
	fmt.Printf("完成从文件反序列化字节码。耗时: %s\n", elapsed)

	// 执行反序列化的字节码
	fmt.Printf("开始执行字节码：\n")
	start = time.Now()
	env := testdata.GetEnv(ser_formulaBatches)

	_ = runner.RunChunk(deserializedChunk, env)

	testdata.CheckValues(t, env, ser_formulaBatches)
	elapsed = time.Since(start)
	fmt.Printf("字节码执行完成。耗时: %s\n", elapsed)

	// 执行语法树（IR）作为对比
	fmt.Printf("开始执行语法树\n")
	start = time.Now()
	env = testdata.GetEnv(ser_formulaBatches)

	_ = runner.RunIR(exprInfos, env)

	testdata.CheckValues(t, env, ser_formulaBatches)
	elapsed = time.Since(start)
	fmt.Printf("语法树执行完成。耗时: %s\n", elapsed)

	fmt.Printf("==========\n")
}

func writeChkFile(chunk *chk.Chunk, filePath string) error {
	data := chunk.ToBytes()
	fileutil.CreateParentIfNotExist(filePath)
	return os.WriteFile(filePath, data, 0644)
}

func readChkFile(filePath string) (*chk.Chunk, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	chunk := chk.NewChunkWithBytes(bytes)
	return chunk, nil
}
