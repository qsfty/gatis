package gatis

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"database/sql"
	"strings"
	"fmt"
)

//创建基类与方法

type GatisModel struct {
	Id int64
	GmtCreate time.Time
	GmtUpdate time.Time
}

var db *sql.DB

func init() {

}

func Init(host string,user string,passwd string,dbname string){
	var err error

	fmt.Print(user + ":" + passwd + "@" + host + "/" + dbname)
	db, err = sql.Open("mysql", user + ":" + passwd + "@tcp(" + host + ":3306)/" + dbname)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	//
	if err != nil {
		panic(err)
	}

	Log.Info("connect mysql success")
	//Analysis_sql_file("src/gatis/src/mapper/course.xml")
}


func getRealSql(model string, method string, data interface{}) (tag string,sql string,attrs map[string]string){

	item := FindMethod(model, method)
	if item.Model == "" {
		panic("no such model : " + model + ",method : " + method +" defined,please check your sql xml")
		return "","",nil
	}
	tplsql := item.Tplsql
	tag = item.Tag
	sql = Render(tplsql, data)
	attrs = item.Attrs
	return
}

func Execute(model string, method string, data interface{}) interface{} {

	tag, sql, attrs := getRealSql(model, method,data)

	if tag == "" {
		return nil
	}

	db.Ping()
	Log.Info("执行sql:\n %s", sql)
	switch tag {
	case "insert":
		return insert( sql, attrs)
	case "update":
		return update( sql, attrs)
	case "delete":
		return del( sql, attrs)
	case "select":
		return query( method, sql, attrs)
	default:
		panic("非法标签" + tag)
		return nil
	}
}

func insert( sql string, attrs map[string]string) int64{

	rst,err := db.Exec(sql)

	if err != nil {

		Log.Error("执行错误:%v", err)
		Log.Error("执行sql: %s", sql)

		return -1
	}

	if a,_ := rst.RowsAffected() ; a > 0 {
		if id,_ := rst.LastInsertId(); err == nil{
			return id
		}
	}
	return -1
}

func update( sql string, attrs map[string]string) int64{
	rst,err := db.Exec(sql)
	if err != nil {
		panic(err)
	}

	fmt.Println(rst)

	if a,_ := rst.RowsAffected() ; a > 0 {
		return 1
	}
	return -1
}


func del( sql string, attrs map[string]string) int64{
	rst,err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	if a,_ := rst.RowsAffected(); a>0 {
		return 1
	}
	return -1
}

func query( method string, sql string, attrs map[string]string) interface{}{
	rows,err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	if method == "get" || attrs["one"] == "true" {
		return queryResultWithOne(rows)
	} else if method == "pageCount" || strings.HasPrefix(method, "count") {
		rt := queryResultWithOne(rows)
		//返回一个int64
		return getMapOne(rt)
	} else if method == "pageQuery" || attrs["more"] == "true" {
		return queryResultWithMore(rows)
	} else {
		rts := queryResultWithMore(rows)
		if len(rts) == 1 {
			return rts[0]
		}
		return rts
	}
}

func getMapOne(m map[string]string) interface{}{
	if len(m) != 1 {
		return m
	}
	for _,v := range m{
		return v
	}
	return m
}

func queryResultWithOne(rows *sql.Rows) (record map[string]string){
	defer rows.Close()
	//获取列
	columns,_ := rows.Columns()
	lens := len(columns)
	//scan参数
	scanArgs := make([]interface{}, lens)
	//列值
	columnValues := make([]interface{}, lens)

	//将scan参数与列值一一对应
	for i,_ := range columnValues {
		scanArgs[i] = &columnValues[i]
	}
	record = make(map[string]string)
	if rows.Next() {
		rows.Scan(scanArgs...)
		for i,columnValue := range columnValues {
			record[columns[i]] = string(columnValue.([]byte))
		}
	}
	return
}

func queryResultWithMore(rows *sql.Rows) (records []map[string]string) {
	defer rows.Close()
	//获取列
	columns,_ := rows.Columns()
	lens := len(columns)

	//scan参数
	scanArgs := make([]interface{}, lens)
	//列值
	columnValues := make([]interface{}, lens)

	//将scan参数与列值一一对应
	for i,_ := range columnValues {
		scanArgs[i] = &columnValues[i]
	}

	var initlen, totalcap = 10,10
	records = make([]map[string]string, initlen)

	n := 0

	for rows.Next() {
		rows.Scan(scanArgs...)

		record := make(map[string]string)

		for i,columnValue := range columnValues {
			record[columns[i]] = string(columnValue.([]byte))
		}

		if n < initlen {
			records[n] = record
		} else {
			records = append(records, record)
			totalcap = cap(records)
		}

		n += 1
	}
	if n < totalcap {
		records = records[:n]
	}

	return
}