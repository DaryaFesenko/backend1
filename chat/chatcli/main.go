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
	if err != nil {
		log.Fatal(err)
	}

	var nick string
	fmt.Print("Enter nick: ")
	fmt.Scanln(os.Stdin, &nick)

	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(nick)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err = io.Copy(conn, os.Stdin) // until you send ^Z
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: exit", conn.LocalAddr())
}
