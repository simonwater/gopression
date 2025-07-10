package exec_test

import (
	"strings"
	"testing"

	"github.com/simonwater/gopression/exec"
	"github.com/simonwater/gopression/gop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var AssemblyContents = []string{
	"POSITION  CODE                PARAMETER           ORDER\n",
	"0         OP_BEGIN            6                   6\n",
	"5         OP_GET_GLOBAL       a                   6\n",
	"10        OP_GET_GLOBAL       b                   6\n",
	"15        OP_GET_GLOBAL       c                   6\n",
	"20        OP_MULTIPLY                             6\n",
	"21        OP_ADD                                  6\n",
	"22        OP_SET_GLOBAL       y                   6\n",
	"27        OP_SET_GLOBAL       x                   6\n",
	"32        OP_END                                  6\n",
	"33        OP_BEGIN            0                   0\n",
	"38        OP_GET_GLOBAL       a                   0\n",
	"43        OP_GET_GLOBAL       b                   0\n",
	"48        OP_GET_GLOBAL       c                   0\n",
	"53        OP_MULTIPLY                             0\n",
	"54        OP_ADD                                  0\n",
	"55        OP_CONSTANT         100                 0\n",
	"60        OP_CONSTANT         5                   0\n",
	"65        OP_CONSTANT         2                   0\n",
	"70        OP_CONSTANT         1                   0\n",
	"75        OP_POWER                                0\n",
	"76        OP_POWER                                0\n",
	"77        OP_DIVIDE                               0\n",
	"78        OP_SUBTRACT                             0\n",
	"79        OP_END                                  0\n",
	"80        OP_BEGIN            1                   1\n",
	"85        OP_GET_GLOBAL       a                   1\n",
	"90        OP_GET_GLOBAL       b                   1\n",
	"95        OP_GET_GLOBAL       c                   1\n",
	"100       OP_MULTIPLY                             1\n",
	"101       OP_ADD                                  1\n",
	"102       OP_CONSTANT         6                   1\n",
	"107       OP_GREATER_EQUAL                        1\n",
	"108       OP_END                                  1\n",
	"109       OP_BEGIN            2                   2\n",
	"114       OP_CONSTANT         1                   2\n",
	"119       OP_CONSTANT         2                   2\n",
	"124       OP_ADD                                  2\n",
	"125       OP_CONSTANT         3                   2\n",
	"130       OP_SUBTRACT                             2\n",
	"131       OP_END                                  2\n",
	"132       OP_BEGIN            3                   3\n",
	"137       OP_CONSTANT         3                   3\n",
	"142       OP_CONSTANT         2                   3\n",
	"147       OP_CONSTANT         1                   3\n",
	"152       OP_ADD                                  3\n",
	"153       OP_MULTIPLY                             3\n",
	"154       OP_END                                  3\n",
	"155       OP_BEGIN            4                   4\n",
	"160       OP_GET_GLOBAL       a                   4\n",
	"165       OP_GET_GLOBAL       b                   4\n",
	"170       OP_GET_GLOBAL       c                   4\n",
	"175       OP_SUBTRACT                             4\n",
	"176       OP_ADD                                  4\n",
	"177       OP_END                                  4\n",
	"178       OP_BEGIN            5                   5\n",
	"183       OP_GET_GLOBAL       a                   5\n",
	"188       OP_CONSTANT         2                   5\n",
	"193       OP_MULTIPLY                             5\n",
	"194       OP_GET_GLOBAL       b                   5\n",
	"199       OP_GET_GLOBAL       c                   5\n",
	"204       OP_SUBTRACT                             5\n",
	"205       OP_ADD                                  5\n",
	"206       OP_END                                  5\n",
	"207       OP_BEGIN            7                   7\n",
	"212       OP_GET_GLOBAL       a                   7\n",
	"217       OP_CONSTANT         1                   7\n",
	"222       OP_GREATER                              7\n",
	"223       OP_JUMP_IF_FALSE    :5->to:233          7\n",
	"228       OP_JUMP             :12->to:245         7\n",
	"233       OP_POP                                  7\n",
	"234       OP_GET_GLOBAL       b                   7\n",
	"239       OP_CONSTANT         1                   7\n",
	"244       OP_GREATER                              7\n",
	"245       OP_JUMP_IF_FALSE    :5->to:255          7\n",
	"250       OP_JUMP             :12->to:267         7\n",
	"255       OP_POP                                  7\n",
	"256       OP_GET_GLOBAL       c                   7\n",
	"261       OP_CONSTANT         1                   7\n",
	"266       OP_GREATER                              7\n",
	"267       OP_JUMP_IF_FALSE    :5->to:277          7\n",
	"272       OP_JUMP             :12->to:289         7\n",
	"277       OP_POP                                  7\n",
	"278       OP_GET_GLOBAL       d                   7\n",
	"283       OP_CONSTANT         1                   7\n",
	"288       OP_GREATER                              7\n",
	"289       OP_END                                  7\n",
	"290       OP_BEGIN            8                   8\n",
	"295       OP_GET_GLOBAL       aa                  8\n",
	"300       OP_CONSTANT         11                  8\n",
	"305       OP_GREATER                              8\n",
	"306       OP_JUMP_IF_FALSE    :12->to:323         8\n",
	"311       OP_POP                                  8\n",
	"312       OP_GET_GLOBAL       bb                  8\n",
	"317       OP_CONSTANT         11                  8\n",
	"322       OP_GREATER                              8\n",
	"323       OP_JUMP_IF_FALSE    :12->to:340         8\n",
	"328       OP_POP                                  8\n",
	"329       OP_GET_GLOBAL       cc                  8\n",
	"334       OP_CONSTANT         11                  8\n",
	"339       OP_GREATER                              8\n",
	"340       OP_JUMP_IF_FALSE    :12->to:357         8\n",
	"345       OP_POP                                  8\n",
	"346       OP_GET_GLOBAL       dd                  8\n",
	"351       OP_CONSTANT         11                  8\n",
	"356       OP_GREATER                              8\n",
	"357       OP_END                                  8\n",
	"358       OP_EXIT                                 8\n",
}

func TestDisassembler(t *testing.T) {
	lines := []string{
		"a + b * c - 100 / 5 ** 2 ** 1",
		"a + b * c >= 6",
		"1 + 2 - 3",
		"3 * (2 + 1)",
		"a + (b - c)",
		"a * 2 + (b - c)",
		"x = y = a + b * c",
		"a > 1 || b > 1 || c > 1 || d > 1",
		"aa > 11 && bb > 11 && cc > 11 && dd > 11",
	}

	runner := gop.NewGopRunner()
	chunk, err := runner.CompileSource(lines)
	require.NoError(t, err, "编译失败")

	// 收集反汇编输出
	var output []string
	printer := func(msg string) {
		output = append(output, msg)
	}

	disassembler := exec.NewDisassembler(printer)
	disassembler.Execute(chunk)

	// 验证字节码大小(单位：字节)
	assert.Equal(t, 441, chunk.GetByteSize(), "字节码总大小不匹配")
	assert.Equal(t, 359, chunk.GetCodesSize(), "指令部分大小不匹配")
	assert.Equal(t, 79, chunk.GetConstsSize(), "常量池大小不匹配")
	assert.Equal(t, 3, chunk.GetVarsSize(), "变量信息大小不匹配")

	// 验证反汇编输出行数
	assert.Equal(t, len(AssemblyContents), len(output), "输出行数不匹配")

	// 验证每行输出内容
	for i, line := range output {
		if i < len(AssemblyContents) {
			expected := strings.Join(strings.Fields(AssemblyContents[i]), " ")
			actual := strings.Join(strings.Fields(line), " ")
			assert.Equal(t, expected, actual, "第 %d 行输出不匹配", i+1)
		}
	}
}
