package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	path := os.Args[1]

	host := "127.0.0.1"
	port := "1979";

	conn, err := net.Dial("tcp", net.JoinHostPort(host, port));
	if err != nil {
		panic(err);
	}
	defer conn.Close();

	req := fmt.Sprintf("%v:%v %v\n", host, port, path);
	_, err = conn.Write([]byte(req));
	if err != nil {
		panic(err);
	}

	buf := make([]byte, 4096);
	n, err := conn.Read(buf);
	if err != nil {
		fmt.Println(err);
	}
	fmt.Printf("%v", string(buf[:n]));
}
