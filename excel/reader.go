package excel

import (
	"fmt"
	"os"
	"peach/data"
	"peach/utils"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	dateFormat1 *regexp.Regexp = regexp.MustCompile(`^\d{2}-\d{2}-\d{2}$`)
	dateFormat2 *regexp.Regexp = regexp.MustCompile(`^\d+$`)
)

// FormatDate 格式化 Excel 文件中的日期
// 将日期统格式化成 YYYY-MM-DD 格式，
// 无法格式化的，沿用原来的值
func FormatDate(s string) string {
	if dateFormat1.MatchString(s) {
		h := strings.Split(s, "-")
		return fmt.Sprintf("20%s-%s-%s", h[2], h[0], h[1])
	}
	if dateFormat2.MatchString(s) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			d, _ := Date(f)
			return d[:10]
		}
	}
	return s
}

// ExcelFile 定义 Excel 文件的读取接口
type ExcelFile interface {
	GetSheetList() []string
	GetValues(sheet int) ([][]string, error)
	ReadSheet(sheet int, skipRow int, ch chan<- []any, cvfns ...data.ConvertFunc)
}

// ExcelReader 定义读取 Excel 数据的类型
type ExcelReader struct {
	fp *os.File
	ExcelFile
}

// NewExcelReader ExcelReader 构造函数
func NewExcelReader(file string) (reader *ExcelReader, err error) {
	var excelFile ExcelFile
	path := utils.NewPath(file)
	// 文件类型检查
	if !path.HasExt(".xls", ".xlsx", ".xlsm") {
		err = fmt.Errorf("文件类型错误:%s", file)
		return
	}
	//打开文件
	fp, err := os.Open(path.String())
	if err != nil {
		return
	}
	if path.HasExt(".xls") {
		excelFile, err = NewXlsFile(fp)
	} else {
		excelFile, err = NewXlsxFile(fp)
	}
	if err == nil {
		reader = &ExcelReader{fp: fp, ExcelFile: excelFile}
	}
	return
}

// Close 关闭 ExcelReader
func (r *ExcelReader) Close() error {
	return r.fp.Close()
}

// GetSheets 获取需要读取的 Worksheet 清单
func (r *ExcelReader) GetSheets(sheets any) (result []int, err error) {
	ws := r.GetSheetList()
	if sheets == nil {
		result = make([]int, len(ws))
		for i := range ws {
			result[i] = i
		}
	} else if names, ok := sheets.([]string); ok {
		result = make([]int, 0)
		fmt.Println(result, names)
		for _, name := range names {
			idx := slices.Index(ws, name)
			if idx < 0 {
				err = fmt.Errorf("sheet %s doesn't exist", name)
				return
			} else {
				result = append(result, idx)
			}
		}
	} else if sts, ok := sheets.([]int); ok {
		result = sts
	} else if name, ok := sheets.(string); ok {
		result = make([]int, 1)
		idx := slices.Index(ws, name)
		if idx < 0 {
			err = fmt.Errorf("sheet %s doesn't exist", name)
			return
		} else {
			result[0] = idx
		}
	} else if idx, ok := sheets.(int); ok {
		result = make([]int, 1)
		result[0] = idx
	}
	return
}

/*
// 读取Excel 数据，通过 data 的数据通道发送
func ReadExcel(file string, sheets any, usecols string, skiprows int, data *utils.Data) {
	var r ExcelReader
	defer close(data.Data)
	path := utils.NewPath(file)
	f, err := os.Open(path.String())
	if err != nil {
		data.Cancel(err)
		return
	}
	defer f.Close()
	if path.HasExt(".xls") {
		r, err = NewXlsReader(f)
	} else if path.HasExt(".xlsx", ".xlsm") {
		r, err = NewXlsxReader(f)
	} else {
		err = fmt.Errorf("%s 非 excel 文件", path.FileInfo().Name())
	}
	if err != nil {
		data.Cancel(err)
		return
	}
	wss, err := GetSheets(r, sheets)
	if err != nil {
		data.Cancel(err)
		return
	}
	for _, idx := range wss {
		d, err := r.GetValues(idx)
		if err != nil {
			data.Cancel(err)
			return
		}
		if skiprows > 0 {
			d = d[skiprows:]
		}
		for _, row := range d {
			select {
			case <-data.Done():
				return
			default:
				data.Data <- utils.Slice(row)
			}
		}
	}
}
*/
