package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	//格式化字符串
	s := "a" + "a"
	fmt.Println(fmt.Sprintf("%s = %d", s, 12))

	//字符串转整型
	i, _ := strconv.Atoi("123")
	fmt.Println(fmt.Sprintf("%s = %d", "asd", i))

	//日期
	timeNow := time.Now()

	fmt.Println(timeNow.String())
	//2013-11-23 23:21:09.9267656 +0800 +0800

	fmt.Println(timeNow.Unix())
	//1385220106

	y, m, d := timeNow.Date()
	fmt.Println(y, m, d, timeNow.Hour(), timeNow.Minute(), timeNow.Second(), timeNow.Nanosecond())

	t := time.Date(2009, time.November, 10, 23, 1, 2, 3, time.UTC)
	fmt.Println(fmt.Sprintf("%s\n", t.Local()))

	//字符串转日期

	const longForm = "2006-1-2 15:04:05"
	t, _ = time.Parse(longForm, "2013-2-1 18:30:50")
	fmt.Println(t)

	const shortForm = "2006-1-2"
	t, _ = time.Parse(shortForm, "2013-10-03")
	fmt.Println(t)

}
