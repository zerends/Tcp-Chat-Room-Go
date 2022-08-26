package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		
		msg = strings.Trim(msg, "\n\r")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     cmdNick,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     cmdJoin,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     cmdRoom,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     cmdMsg,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     cmdQuit,
				client: c,
				args:   args,
			}
		case "/help":
			c.commands <- command{
				id:     cmdHelp,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("yanlis komut: %s \n\r /help komutu ile komutları öğrenebilirsiniz.\n\r", cmd))
		}
	}
	
}

func (c *client) err(err error) {
	c.conn.Write([]byte("Eror: " + err.Error() + "\n"))

}
func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n\r"))

}
