package utils

import "testing"

func TestSprintf(t *testing.T) {
	// 测试不带千分节的数字
	if Sprintf("%10.2f", 1346.45) != "   1346.45" {
		t.Error("Sprintf 不带千分节测试失败")
	}
	// 测试带千分节的数字
	if Sprintf("%10,.2f", -1346.45) != " -1,346.45" {
		t.Error("Sprintf 带千分节测试失败")
	}
	// 测试带千分节的数字
	if Sprintf("%10,.2f", -346.45) != "   -346.45" {
		t.Error("Sprintf 带千分节测试失败")
	}
	// 测试汉字格式化（右对齐）
	if Sprintf("%10s", "测试数据") != "  测试数据" {
		t.Error("Sprintf 字符串右对齐失败")
	}
	// 测试汉字格式化（左对齐）
	if Sprintf("%10s", "测试数据3A") != "测试数据3A" {
		t.Error("Sprintf 字符串左对齐失败")
	}
	// 测试数字格式化（带千分节）
	if Sprintf("%10,d", 1235) != "     1,235" {
		t.Error("Sprintf 测试正整数失败")
	}
}
