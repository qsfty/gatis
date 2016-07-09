package gatis

import (
	"github.com/astaxie/beego/logs"
	"fmt"
)

var Log *logs.BeeLogger

func init() {
	Log = logs.NewLogger(1 << 10)
	Log.SetLogger("console", "")
}


func Ps(args ...interface{}) {
	for _, v := range args {
		fmt.Print(v, " ")
	}
	fmt.Println("")
}
