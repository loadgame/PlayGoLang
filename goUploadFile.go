package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	//"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"
	"time"
)

//<html>
//<head>
//    <title>上传文件</title>
//</head>
//<body>
//<form enctype="multipart/form-data" action="/upload" method="post">
//  <input type="file" name="uploadfile" />
//  <input type="hidden" name="token" value="{{.}}"/>
//  <input type="submit" value="upload" />
//</form>
//</body>
//</html>

var uploadTemplate = template.Must(template.ParseFiles("upload.gtpl"))

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := uploadTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}

}

func stop(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

var token string

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()

		//生成 token ,防止被乱 post

		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token = fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		//获取 token ,防止被乱 post
		tokenRcv := r.FormValue("token")
		if (tokenRcv == "") || (token != tokenRcv) {
			fmt.Println(token + "!=" + tokenRcv)
			fmt.Fprintf(w, "wrong toke")
			return
		}
		token = ""

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {

	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/stop", stop)
	http.ListenAndServe(":8888", nil)
}
