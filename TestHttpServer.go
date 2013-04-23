package main
 
import (
    "net/http"
    "log"
    "os"
    "fmt"
    "html"
)
 
func handler(c http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(c,"%q",html.EscapeString(r.URL.Path))
}
 
func main() {
    http.HandleFunc("/", handler)
    log.Println("Start serving on port 7777")
    http.ListenAndServe(":7777", nil)
    os.Exit(0)
}
