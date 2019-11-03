package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
)

var (
	flagHost, flagUser, flagPasswd, flagDatabase, flagTable, flagCharset string
	flagPort, flagTimes int
)
func init()  {
	flag.StringVar(&flagHost, "h", "", "db host")
	flag.StringVar(&flagUser, "u", "", "db user")
	flag.StringVar(&flagPasswd, "p", "", "db passwd")
	flag.StringVar(&flagDatabase, "B", "", "database")
	flag.StringVar(&flagTable, "T", "", "table")
	flag.IntVar(&flagPort, "P",0,"db port")
	flag.IntVar(&flagTimes, "n", 1, "times")
	flag.StringVar(&flagCharset, "t", "utf8mb4", "charset")
}

func main()  {
	flag.Parse()
	if len(os.Args) == 1 {
		flag.PrintDefaults()
	}
	// 初始化结构体
	db := dbinfo{}
	db.user = flagUser
	db.passwd = flagPasswd
	db.host = flagHost
	db.port = flagPort
	db.database = flagDatabase
	db.table = flagTable
	db.charset = flagCharset
	//fmt.Println(engine)
	//fmt.Println("engine is", engine)
	//定义一个channel
	channel := make(chan string, flagTimes)
	start(db, db.table, flagTimes, channel)
	fenxi(channel)
}

func fenxi(ch chan string)  {
	resultMap := map[string]int{}
	for key := range ch {
		if _, ok := resultMap[key]; ok {
			resultMap[key] +=1
		} else {
			resultMap[key] = 1
		}

	}
	fmt.Println("节点", "core90", "访问了", resultMap["core90"], "次,占比为",float64(resultMap["core90"])/float64(flagTimes)*100)
	fmt.Println("节点", "core91", "访问了", resultMap["core91"], "次,占比为",float64(resultMap["core91"])/float64(flagTimes)*100)
	fmt.Println("节点", "core92", "访问了", resultMap["core92"], "次,占比为",float64(resultMap["core92"])/float64(flagTimes)*100)
	//for key, value := range resultMap {
	//	fmt.Println("节点", key, "访问了", value, "次,占比为",float64(value)/float64(flagTimes)*100)
	//}
}


func start(db dbinfo, table string, times int, ch chan string)  {
	sql := fmt.Sprintf("select name from %s;", table)
	for i:=1;i<=times;i++ {
		conn := db.getEngine()
		mapSlice, err := conn.QueryString(sql)
		if err != nil {
			fmt.Println("run sql is bad, sql is", sql, "err is", err)
			os.Exit(0)
		}
		//fmt.Println("run sql is ok")
		for _, dict := range mapSlice {
			value := dict["name"]
			fmt.Println("i is", i, "value is", value)
			ch <- value
		}
	}
	close(ch)
}


type dbinfo struct {
	user string
	passwd string
	host string
	port int
	database string
	table string
	charset string
}

func (db *dbinfo)getEngine() *xorm.Engine {
	//fmt.Println("init db", time.Now())
	endPoint := fmt.Sprintf("%s:%d", db.host, db.port)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true", db.user, db.passwd, endPoint, db.database, db.charset)
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		fmt.Println("init db connection", endPoint, "is bad")
		os.Exit(1)
	}
	//fmt.Println("init db is ok")
	return engine
}








