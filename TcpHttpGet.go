// You can edit this code!
// Click here and start typing.
package main

import ("net"
"bytes"
"fmt"
)

func main() {

conn, err:=net.Dial("tcp","www.qchat.cn:80")
_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
result,err := readFully(conn)
fmt.Println(String(result))

}