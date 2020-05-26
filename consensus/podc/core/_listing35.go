

// Sample program to show how a bytes.Buffer can also be used
// with the io.Copy function.
package main

import (
"bytes"
"fmt"
"io"
"os"
)

// main is the entry point for the application.
func main() {
	var b bytes.Buffer
	//type Buffer struct {
	//	buf      []byte // contents are the bytes buf[off : len(buf)]
	//	off      int    // read at &buf[off], write at &buf[len(buf)]
	//	lastRead readOp // last read operation, so that Unread* can work correctly.
	//}


	// Write a string to the buffer.
	b.Write([]byte("Hello"))  //buf에 bytes.Buffer의 메소드 Write를 이용하여, buf에 "Hello"를 저장

	// Use Fprintf to concatenate a string to the Buffer.
	fmt.Fprintf(&b, "World!")  // buf의 메모리번지를 넘겨줘서, World를 추가함.

	// Write the content of the Buffer to stdout.
	io.Copy(os.Stdout, &b)  // buf의 메모리번지를 넘겨주고, 표준출력으로 내보낸다.
}
