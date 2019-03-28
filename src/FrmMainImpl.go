// 由res2go自动生成。
// 在这里写你的事件。

package main

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"strings"
	"unicode"
)

// 一些缩写单词
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

// Example:
// 	fmtFieldName("foo_id")
// Output: FooID
func FmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	for len(name) > 0 && name[0] == '_' {
		name = name[1:]
	}

	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

//::private::
type TFrmMainFields struct {
}

type SQLDD struct {
	TableName    string
	FmtTableName string
	Cols         []SQLDDCol
}

type SQLDDCol struct {
	Name    string
	Type    string
	FmtName string
}

var sqldd SQLDD

func (f *TFrmMain) OnBtnGenClick(sender vcl.IObject) {
	if f.MmoSQL.Text() == "" {
		vcl.ShowMessage("请粘贴要进行转换SQL的创建表语句到文本框")
		f.MmoSQL.SetFocus()
		return
	}
	tree, err := sqlparser.ParseStrictDDL(f.MmoSQL.Text())
	if err != nil {
		vcl.ShowMessage("确保SQL语句正确")
		f.MmoSQL.SetFocus()
		return
	}
	f.ClbFields.Clear()
	ddl := tree.(*sqlparser.DDL)

	tableName := ddl.NewName.Name.String()

	GoStruct := ""

	GoStruct += fmt.Sprintln("type " + FmtFieldName(tableName) + " struct {")

	sqldd.TableName = tableName

	sqldd.FmtTableName = FmtFieldName(tableName)

	for i, col := range ddl.TableSpec.Columns {
		name := col.Name.String()
		mysqlType := strings.ToLower(ddl.TableSpec.Columns[i].Type.Type)
		if mysqlType == "int" || mysqlType == "tinyint" || mysqlType == "mediumint" || mysqlType == "smallint" {
			mysqlType = "int"
		}
		if mysqlType == "varchar" || mysqlType == "char" || mysqlType == "text" || mysqlType == "mediumtext" {
			mysqlType = "string"
		}
		if mysqlType == "double" {
			mysqlType = "float32"
		}
		if mysqlType == "set" {
			mysqlType = "[]string"
		}
		if mysqlType == "timestamp" {
			mysqlType = "time.Time"
		}
		GoStruct += fmt.Sprintln("    ", FmtFieldName(name), mysqlType, "`sql:\""+name+"\"`")
		var sqlddcol SQLDDCol
		sqlddcol.Name = name
		sqlddcol.Type = mysqlType
		sqlddcol.FmtName = FmtFieldName(name)
		sqldd.Cols = append(sqldd.Cols, sqlddcol)
		f.ClbFields.Items().Add(name)
	}
	GoStruct += fmt.Sprintln("}")
	f.MmoStruct.SetText(GoStruct)
}

func (f *TFrmMain) OnFormCreate(sender vcl.IObject) {

}

func (f *TFrmMain) OnBtnGetOneClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	var i int32
	var ValSQL string // 变量
	var SelectFields []string
	var ValSQLFields []string
	var ValSQLScanFileds []string
	var StSQL string // 结构体
	var StSQLScanFileds []string

	var checkedCount int

	for i = 0; i < f.ClbFields.Items().Count(); i++ {
		item := f.ClbFields.Items().ValueFromIndex(i)

		for j := 0; j < len(sqldd.Cols); j++ {
			col := sqldd.Cols[j]
			if item == col.Name && f.ClbFields.Checked(i) {
				SelectFields = append(SelectFields, col.Name)
				ValSQLFields = append(ValSQLFields, "var F"+col.FmtName+" "+col.Type)
				ValSQLScanFileds = append(ValSQLScanFileds, "&F"+col.FmtName)
				StSQLScanFileds = append(StSQLScanFileds, "&F"+sqldd.FmtTableName+"."+col.FmtName)
				checkedCount++
			}
		}
	}
	if checkedCount == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	ValSQL += "//查询到变量\r\n"
	ValSQL += fmt.Sprintln(strings.Join(ValSQLFields, "\r\n"))
	ValSQL += "DB.QueryRow(\"SELECT " + strings.Join(SelectFields, ",") + " FROM " + sqldd.TableName + " WHERE ? LIMIT 1\", 1).Scan(" + strings.Join(ValSQLScanFileds, ",") + ")"

	StSQL += "//查询到结构\r\n"
	StSQL += "var F" + sqldd.FmtTableName + " " + sqldd.FmtTableName + "\r\n"
	StSQL += "DB.QueryRow(\"SELECT " + strings.Join(SelectFields, ",") + " FROM " + sqldd.TableName + " WHERE ? LIMIT 1\", 1).Scan(" + strings.Join(StSQLScanFileds, ",") + ")"

	f.MmoResult.SetText(fmt.Sprintln(ValSQL + "\r\n\r\n" + StSQL))

}

func (f *TFrmMain) OnBtnGetAllClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	var i int32
	var ValSQL string // 变量
	var SelectFields []string
	var ValSQLFields []string
	var ValSQLScanFileds []string
	var StSQL string // 结构体
	var StSQLScanFileds []string

	var checkedCount int

	for i = 0; i < f.ClbFields.Items().Count(); i++ {
		item := f.ClbFields.Items().ValueFromIndex(i)

		for j := 0; j < len(sqldd.Cols); j++ {
			col := sqldd.Cols[j]
			if item == col.Name && f.ClbFields.Checked(i) {
				SelectFields = append(SelectFields, col.Name)
				ValSQLFields = append(ValSQLFields, "    var F"+col.FmtName+" "+col.Type)
				ValSQLScanFileds = append(ValSQLScanFileds, "&F"+col.FmtName)
				StSQLScanFileds = append(StSQLScanFileds, "&F"+sqldd.FmtTableName+"."+col.FmtName)
				checkedCount++
			}
		}
	}
	if checkedCount == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	ValSQL += "//查询到变量\r\n"
	ValSQL += "rows, err := DB.Query(\"SELECT " + strings.Join(SelectFields, ",") + " FROM " + sqldd.TableName + " WHERE ?\", 1)"
	ValSQL += `
if err != nil {
    log.Println(err)
}
defer rows.Close()
for rows.Next() {
`
	ValSQL += fmt.Sprintln(strings.Join(ValSQLFields, "\r\n"))
	ValSQL += "    err := row.Scan(" + strings.Join(ValSQLScanFileds, ",") + ")"
	ValSQL += `
    if err != nil {
        log.Println(err)
    }`
	ValSQL += "\r\n    // Your code..."
	ValSQL += "\r\n}"

	StSQL += "//查询到结构\r\n"
	StSQL += "rows, err := DB.Query(\"SELECT " + strings.Join(SelectFields, ",") + " FROM " + sqldd.TableName + " WHERE ?\", 1)"
	StSQL += `
if err != nil {
    log.Println(err)
}
defer rows.Close()
for rows.Next() {
`
	StSQL += "    var F" + sqldd.FmtTableName + " " + sqldd.FmtTableName + "\r\n"
	StSQL += "    err := rows.Scan(" + strings.Join(StSQLScanFileds, ",") + ")"
	StSQL += `
    if err != nil {
        log.Println(err)
    }`
	StSQL += "\r\n    // Your code..."
	StSQL += "\r\n}"

	f.MmoResult.SetText(fmt.Sprintln(ValSQL + "\r\n\r\n" + StSQL))
}

func (f *TFrmMain) OnBtnUpdateClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	var i int32
	var UpdateSQL string // 变量
	var UpdateFields []string
	var UpdateDfValue []string

	var checkedCount int

	for i = 0; i < f.ClbFields.Items().Count(); i++ {
		item := f.ClbFields.Items().ValueFromIndex(i)

		for j := 0; j < len(sqldd.Cols); j++ {
			col := sqldd.Cols[j]
			if item == col.Name && f.ClbFields.Checked(i) {
				UpdateFields = append(UpdateFields, col.Name+"=?")
				var dfval string
				if col.Type == "int" {
					dfval = "0"
				}
				if col.Type == "string" {
					dfval = "\"\""
				}
				if col.Type == "float" {
					dfval = "0.0"
				}
				if col.Type == "time.Time" {
					dfval = "time.Now()"
				}
				UpdateDfValue = append(UpdateDfValue, "    "+dfval+", // "+col.Name)
				checkedCount++
			}
		}
	}
	if checkedCount == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}

	UpdateSQL = `_, err := DB.Exec("UPDATE ` + sqldd.TableName + ` SET ` + strings.Join(UpdateFields, ",") + ` WHERE ?",
` + strings.Join(UpdateDfValue, "\r\n") + `
    1,  // where start
)
if err != nil {
    log.Println(err)
}`
	f.MmoResult.SetText(UpdateSQL)
}

func (f *TFrmMain) OnBtnDeleteClick(sender vcl.IObject) {
	if sqldd.TableName == "" {
		vcl.ShowMessage("没有数据表")
		return
	}
	DeleteSQL := `_, err := DB.Exec("DELETE FROM ` + sqldd.TableName + ` WHERE ?", 
    1,  // where start
)
if err != nil {
    log.Println(err)
}`
	f.MmoResult.SetText(DeleteSQL)
}

func (f *TFrmMain) OnBtnInsertClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}
	var i int32
	var InsertSQL string // 变量
	var InsertFields []string
	var InsertValue []string
	var InsertDfValue []string

	var checkedCount int

	for i = 0; i < f.ClbFields.Items().Count(); i++ {
		item := f.ClbFields.Items().ValueFromIndex(i)

		for j := 0; j < len(sqldd.Cols); j++ {
			col := sqldd.Cols[j]
			if item == col.Name && f.ClbFields.Checked(i) {
				InsertFields = append(InsertFields, col.Name)
				InsertValue = append(InsertValue, "?")
				var dfval string
				if col.Type == "int" {
					dfval = "0"
				}
				if col.Type == "string" {
					dfval = "\"\""
				}
				if col.Type == "float" {
					dfval = "0.0"
				}
				if col.Type == "time.Time" {
					dfval = "time.Now()"
				}
				InsertDfValue = append(InsertDfValue, "    "+dfval+", // "+col.Name)
				checkedCount++
			}
		}
	}
	if checkedCount == 0 {
		vcl.ShowMessage("请选择需要操作的字段")
		return
	}

	InsertSQL = `_, err := DB.Exec("INSERT INTO ` + sqldd.TableName + `(` + strings.Join(InsertFields, ",") + `) VALUES (` + strings.Join(InsertValue, ",") + `)",
` + strings.Join(InsertDfValue, "\r\n") + `
)
if err != nil {
    log.Println(err)
}`
	f.MmoResult.SetText(InsertSQL)
}

func (f *TFrmMain) OnMmiSelectAllClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() > 0 {
		f.ClbFields.CheckAll(types.CbChecked, true, true)
	}
}

func (f *TFrmMain) OnMmiUnSelectAllClick(sender vcl.IObject) {
	if f.ClbFields.Items().Count() > 0 {
		f.ClbFields.CheckAll(types.CbUnchecked, true, true)
	}
}
