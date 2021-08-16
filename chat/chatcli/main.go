package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")

	var nick string
	fmt.Print("Enter nick: ")
	fmt.Scanln(os.Stdin, &nick)

	writer := bufio.NewWriter(conn)
	writer.WriteString(nick)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z

	fmt.Printf("%s: exit", conn.LocalAddr())
}
