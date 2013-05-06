package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net"
	//"runtime/debug"
	//"runtime"
	//"strconv"
	/*
				"errors"
		"bufio"
			"bytes"
				"io"
				"log"
				"net/url"
				"path"
				"strings"
				"sync"
				"time"
	*/
)
import (
	"os"
	"syscall"
	"time"
)

func main() {
	Log("====================")
	Log("Hello SSL World")
	Log("====================")

	Log("ReadOption")
	readOptions("sslSocket.json")
	Log("Start")
	go StartListen()
	go StartListenSSL()
//StartListenSSL()
	timeout := make(chan bool, 1)
	<-timeout
	/*
		for {
			runtime.Gosched()
		}
	*/
	Log("END")
	return
}

const noLimit int64 = (1 << 63) - 1

func newConn(srcConn net.Conn) (c *conn, err error) {
	c = new(conn)

	c.srcConn = srcConn

	//	c.lr = io.LimitReader(srcConn, noLimit).(*io.LimitedReader)
	//	br := bufio.NewReader(c.lr)
	//	bw := bufio.NewWriter(srcConn)
	//	c.buf = bufio.NewReadWriter(br, bw)

	return c, nil
}

//main -> StartListenSSL -> go c.serve() -> ClientHandler -> go ClientReader() go ClientReader

func StartListenSSLOld() {

	if myopt.ListenSSL == "" {
		return
	}

	var err error
	certFile := myopt.CrtFile
	keyFile := myopt.PemFile
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		Log("LoadX509KeyPair:", myopt.PemFile, myopt.CrtFile, err)
		return
	}

	addr := myopt.ListenSSL
	lstn, err := net.Listen("tcp", addr)
	if err != nil {
		Log("Listen:", err)
		return
	}
	tlsListener := tls.NewListener(lstn, config)

	Log("SSL Listening on", addr)

	l := tlsListener
	for {
		srcConn, err := l.Accept()
		if err != nil {
			Log("SSL Accept:", err)
			break
		}
		c, err := newConn(srcConn)
		if err != nil {
			Log("SSL newConn:", err)
			continue
		}
		c.Name = "SSL " + srcConn.RemoteAddr().String()
		c.isSSL = false
		go c.serve()
	}
	return

}
func StartListenSSL() {

	if myopt.ListenSSL == "" {
		return
	}

	var err error
	certFile := myopt.CrtFile
	keyFile := myopt.PemFile

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		Log("LoadX509KeyPair:", myopt.PemFile, myopt.CrtFile, err)
		return
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Time = time.Now
	config.Rand = rand.Reader

	addr := myopt.ListenSSL
	lstn, err := tls.Listen("tcp", addr, &config)
	if err != nil {
		Log("Listen:", err)
		return
	}

	Log("SSL Listening on", addr)

	l := lstn
	for {
		srcConn, err := l.Accept()
		if err != nil {
			Log("SSL Accept:", err)
			break
		}
		c, err := newConn(srcConn)
		if err != nil {
			Log("SSL newConn:", err)
			continue
		}
		c.Name = "SSL " + srcConn.RemoteAddr().String()
		c.isSSL = true
		go c.serve()
	}
	return

}

func StartListen() {

	if myopt.Listen == "" {
		return
	}

	var err error

	addr := myopt.Listen
	lstn, err := net.Listen("tcp", addr)
	if err != nil {
		Log("Listen:", err)
		return
	}
	Log("Listening on", addr)

	l := lstn
	for {
		srcConn, err := l.Accept()
		if err != nil {
			Log("Accept:", err)
			break
		}
		c, err := newConn(srcConn)
		if err != nil {
			Log("newConn:", err)
			continue
		}
		c.Name = srcConn.RemoteAddr().String()
		go c.serve()
	}
	return

}

//=======================

type conn struct {
	remoteAddr string // network address of remote side
	//server     *Server              // the Server on which the connection arrived
	srcConn net.Conn // i/o connection
	//lr      *io.LimitedReader // io.LimitReader(rwc)
	//buf     *bufio.ReadWriter // buffered(lr,rwc), reading from bufio->limitReader->rwc
	//hijacked bool                 // connection has been hijacked by handler
	//tlsState *tls.ConnectionState // or nil when not using TLS
	//body     []byte
	desConn net.Conn
	isSSL bool
	Name    string
	NameDes string
}

/*
func (c *conn) close() {
	if c.buf != nil {
		c.buf.Flush()
		c.buf = nil
	}
	if c.srcConn != nil {
		c.srcConn.Close()
		c.srcConn = nil
	}
}*/

func (c *conn) ClientReader() {
	buffer := make([]byte, 2048)
	for {
		n1, err1 := c.srcConn.Read(buffer)
		if err1 != nil {
			c.srcConn.Close()
			c.desConn.Close()
			break
		}
		 if n1 !=0 {
		 
		n2, err2 := c.desConn.Write(buffer[:n1])
		if err2 != nil {
			c.srcConn.Close()
			c.desConn.Close()
			break
		}
		if c.isSSL && (n1>1) {
			if (string(buffer[:8])=="Operator") && (buffer[n1-2]==13) && (buffer[n1-1]==10) {
		 	_, err3 := c.desConn.Write(buffer[n1-2:n1])
                	if err3 != nil {
                        	c.srcConn.Close()
                        	c.desConn.Close()
                       	 break
                	}
			Log(c.Name, "ADD 0d 0a")	
			}
		}


//			Log("["+string(buffer[:n1])+"]")
//			Log(c.Name, "---> ",n1,"->", n2)
//			Log(buffer[:n1])
			n1=n2 //nomean avoid complie error		
		}
	}
	Log(c.Name, "-X->")
}

func (c *conn) ClientSender() {
	buffer := make([]byte, 2048)
	for {
		n1, err1 := c.desConn.Read(buffer)
		if err1 != nil {
			c.desConn.Close()
			c.srcConn.Close()
			break
		}

		if n1 !=0 {
			n2, err2 := c.srcConn.Write(buffer[:n1])
			if err2 != nil {
				c.srcConn.Close()
				c.desConn.Close()
				break
			} 
//			Log("["+string(buffer[:n1])+"]")
//			Log(c.Name, "<--- ",n2,"<-", n1)
			n1=n2 //nomean
		}

	}
	Log(c.Name, "<-X-")
}

func (c *conn) StartRW() {

	destination := myopt.Proxy

	Log(c.Name, "-->-", destination)
	des, err := net.Dial("tcp", destination)

	if err != nil {
		Log(err)
		c.srcConn.Close()
		return
	}

	c.desConn = des
	c.NameDes = des.RemoteAddr().String()

	Log(c.Name, "<-->", c.NameDes)

	go c.ClientSender()
	go c.ClientReader()

	//clientList.PushBack(* c)

}

func (c *conn) serve() {
	/*
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			var buf bytes.Buffer
			fmt.Fprintf(&buf, "http: panic serving %v: %v\n", c.remoteAddr, err)
			buf.Write(debug.Stack())
			fmt.Println(buf.String())

			if c.srcConn != nil { // may be nil if connection hijacked
				c.srcConn.Close()
			}
		}()

		if tlsConn, ok := c.srcConn.(*tls.Conn); ok {
			if err := tlsConn.Handshake(); err != nil {
				c.close()
				return
			}
			c.tlsState = new(tls.ConnectionState)
			*c.tlsState = tlsConn.ConnectionState()
		}*/

	c.StartRW()

	/*
		buffer := make([]byte, 2048)

		for {
			n, status := c.rwc.Read(buffer)
			if status != nil {
				c.close()
				break
			}
			fmt.Println(n)

			msg := "HTTP/1.1 200 OK\nContent-Length: " + strconv.Itoa(n) + "\n\n" + string(buffer[:n])
			fmt.Println(msg)
			c.rwc.Write([]byte(msg))

		} 
		c.close()
	*/
}

type myOptionsSockets struct {
	ListenSSL string
	Listen    string
	Proxy     string
	CrtFile   string
	PemFile   string
	Nomsg     int
}

var myopt myOptionsSockets

func readOptions(filename string) {
	file, e := ioutil.ReadFile(filename)

	if e != nil {
		Log(e)
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(file))

	json.Unmarshal(file, &myopt)
	Log(myopt)

}

func filelog(path string, s string) {

	t := time.Now()
	yy, mm, dd := t.Date()
	fn := fmt.Sprintf("%s%d%s%d.log.txt", path, yy, mm, dd)

	f, err := os.OpenFile(fn, syscall.O_RDWR|syscall.O_APPEND|syscall.O_CREAT, 0666) //f,err := os.Create("goLog.txt")

	if err != nil {
		fmt.Print(err)
	} else {
		defer f.Close()

		str := fmt.Sprintf("%02d:%02d:%02d.%07d : %s\n",
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond()/100, s)

		_, err2 := f.WriteString(str)
		if err2 != nil {
			fmt.Print(err2)
		}
	}
}

func Log(v ...interface{}) {
	filelog("SSL.", fmt.Sprint(v...))

	if myopt.Nomsg == 0 {
		fmt.Println(v...)
	}
}
