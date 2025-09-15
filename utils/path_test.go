package utils

import (
	"testing"
)

func TestExpand(t *testing.T) {
	/*
			p := NewPath("~hunter/abc")
			if p.String() != `C:\Users\hunter\abc` {
				t.Errorf("Test Failed!")
			}

		path := Expand("$programdata/abc")
		if path != `C:\ProgramData\abc` {
			t.Errorf("Test Failed!")
		}
		path = Expand("%programdata%/abc")

		if path != `C:\ProgramData\abc` {
			t.Errorf("Test Expand Failed!")
		}
	*/
}

func TestPath(t *testing.T) {
	tmp := NewPath(TempDir)
	tmp = tmp.Join("Abc.txt")
	if tmp.Exists() {
		t.Error("测试检查文件是否存在失败")
	}
	if tmp.Ext() != ".txt" {
		t.Error("测试文件扩展名失败")
	}
	if !tmp.HasExt(".txt", ".csv") {
		t.Error("测试文件的扩展名是否在指定列表中扩失败")
	}
	if tmp.HasExt(".xls", ".xlsx") {
		t.Error("测试文件的扩展名是否在指定列表中扩失败")
	}
	t2 := NewPath("abc.tar.gz")
	if !t2.HasExt(".tar.gz") {
		t.Error("测试文件扩展名检查失败")
	}
}
