package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	conn net.Conn
	nick string
	ch   chan string
}

var (
	currentExpression *Expression
	clients           []client
	questions         chan string
	winners           = make(chan string)
	leaving           = make(chan client)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	questions = make(chan string, 1)

	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	nick, _ := bufio.NewReader(conn).ReadString('\n')
	nick = strings.Trim(nick, "\n")

	c := client{
		conn: conn,
		nick: nick,
		ch:   make(chan string),
	}
	clients = append(clients, c)

	go clientWriter(c)

	c.ch <- "Welcome " + nick

	if currentExpression != nil {
		c.ch <- currentExpression.expression
	}

	go generateExpression()
	go writeWinner()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if currentExpression.answers == input.Text() {
			winners <- nick
		}
	}

	leaving <- c
	conn.Close()
}

func clientWriter(c client) {
	for msg := range c.ch {
		fmt.Fprintln(c.conn, msg)
	}
}

func deleteClient(c client) {
	for idx := range clients {
		if clients[idx].conn.RemoteAddr().String() == c.conn.RemoteAddr().String() {
			copy(clients[idx:], clients[idx+1:])
			return
		}
	}
}

func broadcaster() {
	for {
		cli := <-leaving
		deleteClient(cli)
	}
}

func generateExpression() {
	for {
		question := NewExpression()
		questions <- question.expression
		currentExpression = question
		for _, val := range clients {
			val.ch <- question.expression
		}
	}
}

func writeWinner() {
	for val := range winners {
		for _, cli := range clients {
			cli.ch <- val + " win\n"
		}
		<-questions
	}
}
