// goURLs
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	//"sync"
	"runtime"
	"time"
	//"os"

	//"os/signal"
)

type MonitorSite struct {
	Url      string
	Data     string
	Errcount int
	Checkcount int
	Last     time.Duration
}

var UrlMap map[string]*MonitorSite

//var locker sync.Mutex
func checkUrl(url string, data string) bool {

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return strings.Contains(string(body), data)
}

func ReportStatus() {
    memStats := &runtime.MemStats{}
    runtime.ReadMemStats(memStats)
    nsInMs := float64(time.Millisecond)
    prefix := "MT"
    
	for {
	  time.Sleep(10e8)
	  //fmt.Println("=========")
	   for _, u := range UrlMap {
		 fmt.Println("=========")
		 fmt.Println("rpt",u)
	  }
	  
	  
       if false {
        
        fmt.Println(fmt.Sprintf("%s.goroutines", prefix),
            float64(runtime.NumGoroutine()))
        fmt.Println(fmt.Sprintf("%s.memory.allocated", prefix),
            float64(memStats.Alloc))
        fmt.Println(fmt.Sprintf("%s.memory.mallocs", prefix),
            float64(memStats.Mallocs))
        fmt.Println(fmt.Sprintf("%s.memory.frees", prefix),
            float64(memStats.Frees))
        fmt.Println(fmt.Sprintf("%s.memory.gc.total_pause", prefix),
            float64(memStats.PauseTotalNs)/nsInMs)
        fmt.Println(fmt.Sprintf("%s.memory.heap", prefix),
            float64(memStats.HeapAlloc))
        fmt.Println(fmt.Sprintf("%s.memory.stack", prefix),
            float64(memStats.StackInuse))
        
      }
	  
	}
}

func  CheckSite(u *MonitorSite) {
			for {
				t := time.Now()
				
				u.Checkcount=u.Checkcount+1
				if checkUrl(u.Url, u.Data) {
					u.Last = time.Since(t)
					u.Errcount = 0
					//fmt.Println("OK",u)
					
				} else {
					u.Errcount = u.Errcount + 1
					fmt.Println("*ERR*",u)
					
					todaylog = todaylog+"<br />"+fmt.Sprintf("*ERR*%v",u)
				}
				time.Sleep(10e8*5)
			}
} 
		
		
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, todaylog)
}

var todaylog string
func runchecksite() {


	UrlMap := make(map[string]*MonitorSite)
	
	UrlMap["http://www.chat4support.com"]=&MonitorSite{Url: "http://www.chat4support.com", Data: "chat4support"}
	UrlMap["http://qchat.cn"]=&MonitorSite{Url: "http://qchat.cn", Data: "chat4support"}
	UrlMap["http://www.chatonwebsite.com/weboperator"]=&MonitorSite{Url: "http://www.chatonwebsite.com/weboperator", Data: "Account"}
	 



	//println(checkUrl("http://qchat.cn", "chat4support"))
	
	//go ReportStatus()

	//locker

	for _, u := range UrlMap {
	    
		go CheckSite(u)
	}
	
	for  {
		for _, u := range UrlMap {
			if u.Errcount>=3 {
				fmt.Println("*********************ALARM*********************************")
				fmt.Println("****",u,"****")
				fmt.Println("***********************************************************")
				todaylog = todaylog+"<br />*****ALARM******"+fmt.Sprintf("%v",u)
			
			}
			
		}
		time.Sleep(10e8*5)
	}

}
func main() {

	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt)
	//	go func() {
	//	  for sig := range c {
	//		 fmt.Printf("captured %v, stopping profiler and exiting..", sig)
	//		 os.Exit(1)
	//	  }
	//	}()
	
	go runchecksite()

     port:=":8888"
	 
	http.HandleFunc("/", handler)
	fmt.Println("Server listen on",port)
	todaylog="Server Started"
	http.ListenAndServe(port, nil)

	//    done := make(chan string)
//	<-done
}
