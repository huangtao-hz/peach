package utils

import (
	"testing"
)

func TestNewDate(t *testing.T) {
	d := NewDate("2020-3-1")
	if d.Year() != 2020 || d.Month() != 3 || d.Day() != 1 {
		t.Error("日期不正确")
	}

	d = NewDate("20200301")
	if d.Year() != 2020 || d.Month() != 3 || d.Day() != 1 {
		t.Error("日期不正确")
	}
}
func TestToday(t *testing.T) {
	d := NewDate("2020-02-29")
	if d.String() != "2020-02-29" {
		t.Error("测试失败")
	}
	if d.Format("%F") != "2020-02-29" {
		t.Error("Format  %Q Failed ")
	}
	if d.Format("%W") != "星期六" {
		t.Error("Format  %Q Failed ")
	}
	if d.Quater() != 1 {
		t.Error("季度测试失败")
	}
	if d.Weekday() != 6 {
		t.Error(d.Weekday())
	}
	s := d.Add(0, 0, 1)
	if s.String() != "2020-03-01" {
		t.Error("测试失败")
	}
	if s.Weekday() != 0 {
		t.Error(s.Weekday())
	}
	k := s.Format("%y")
	if k != "20" {
		t.Error("Format Failed")
	}
}
