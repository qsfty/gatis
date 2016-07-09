package gatis

import (
	"text/template"
	"strconv"
	"time"
	"strings"
)

var Tplfuncs template.FuncMap

func init() {

	Tplfuncs = map[string]interface{}{
		"q" : addSingleQuote,
		"eq" : addEqual,
		"w" : addWhere,
		"s" : addSet,
		"now" : addNow,
	}
}

// val
func addSingleQuote(v...interface{}) string {
	if len(v) == 0 {
		return "''"
	}
	return "'" + escape(itos(v[0])) + "'"
}

func escape(v string) string {
	if v == "" {
		return ""
	}
	return strings.Replace(v, "'", "\\'", -1)
}

//col val
func addEqual(v ...interface{}) string{
	if len(v) == 1{
		return ""
	}

	return itos(v[0]) + " = " + addSingleQuote(v[1])
}

//col val
func addWhere(v ...interface{}) string{
	if len(v) == 1{
		return ""
	}

	if v[1] != "" {
		return "and " + addEqual(v[0], v[1])
	}
	return ""
}

func addSet(v...interface{}) string{
	if len(v) == 1{
		return ""
	}

	if v[1] != "" {
		return addEqual(v[0], v[1]) + ","
	}

	return ""
}

func addNow(col string,format ...interface{}) string{
	formatLayout := "2006-01-02 03:04:05"
	if len(format) == 1{
		switch  itos(format[0]) {
		case "date":
			formatLayout = "2006-01-02"
		case "timestamp":
			formatLayout = ""
		}
	}
	if formatLayout != "" {
		return addEqual(col, time.Now().Format(formatLayout))
	}
	return addEqual(col, time.Now().Unix())
}


func itos(v interface{}) string{
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	default:
		return ""
	}
}