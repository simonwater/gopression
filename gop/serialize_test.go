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
	ser_formulaBatches = 10000
	ser_testDirectory  = "SerializeTest"
)

func TestChunkSerialization(t *testing.T) {
	fmt.Printf("序列化反序列化测试：")
	chunkSerializeTest(t)
}

func chunkSerializeTest(t *testing.T) {
	// 创建表达式列表
	lines := testdata.GetExpressions(t, ser_formulaBatches)
	fmt.Printf("表达式总数：%d", len(lines))

	runner := gop.NewGopRunner()

	// 解析和分析表达式
	fmt.Printf("开始解析和分析：")
	start := time.Now()
	exprs, err := runner.Parse(lines)
	require.NoError(t, err, "解析失败")

	exprInfos := runner.Analyze(exprs)
	require.NoError(t, err, "分析失败")

	elapsed := time.Since(start)
	fmt.Printf("中间结果生成完成。耗时: %s", elapsed)

	// 编译为字节码
	runner.SetTrace(true)
	chunk := runner.CompileIR(exprInfos)

	// 序列化字节码
	chunkSize := chunk.GetByteSize()
	fmt.Printf("开始进行字节码序列化反序列化，字节码大小(KB): %d", chunkSize/1024)

	start = time.Now()
	fileName := "Chunks.gob"
	filePath := fileutil.GetTestPath(ser_testDirectory, fileName)

	err = fileutil.SerializeObject(chunk, filePath)
	require.NoError(t, err, "序列化失败")
	elapsed = time.Since(start)
	fmt.Printf("字节码已序列化到文件：%s 耗时: %s", fileName, elapsed)

	// 反序列化字节码
	start = time.Now()
	deserializedChunk, err := fileutil.DeserializeObject[chk.Chunk](filePath)
	require.NoError(t, err, "反序列化失败")
	elapsed = time.Since(start)
	fmt.Printf("完成从文件反序列化字节码。耗时: %s", elapsed)

	// 执行反序列化的字节码
	fmt.Printf("开始执行字节码：")
	start = time.Now()
	env := testdata.GetEnv(t, ser_formulaBatches)

	_ = runner.RunChunk(&deserializedChunk, env)

	testdata.CheckValues(t, env, ser_formulaBatches)
	elapsed = time.Since(start)
	fmt.Printf("字节码执行完成。耗时: %s", elapsed)

	// 执行语法树（IR）作为对比
	fmt.Printf("开始执行语法树")
	start = time.Now()
	env = testdata.GetEnv(t, ser_formulaBatches)

	_ = runner.RunIR(exprInfos, env)

	testdata.CheckValues(t, env, ser_formulaBatches)
	elapsed = time.Since(start)
	fmt.Printf("语法树执行完成。耗时: %s", elapsed)

	fmt.Printf("==========")
}
