package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case cmdNick:
			s.nick(cmd.client, cmd.args)
		case cmdJoin:
			s.join(cmd.client, cmd.args)
		case cmdRoom:
			s.listRooms(cmd.client)
		case cmdMsg:
			s.msg(cmd.client, cmd.args)
		case cmdQuit:
			s.quit(cmd.client)

		case cmdHelp:
			s.help(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	log.Printf(">Yeni kullanıcı katıldı : %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "name",
		commands: s.commands,
	}

	//fmt.Println("Hoş Geldiniz. /help komutu ile komutları öğrenebilirsiniz.")

}

func (s *server) nick(c *client, args []string) {
	if len(args) == 0 {
		c.msg(">Lütfen adınızı girin. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("Hoş Geldin %s", c.nick))
}
func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Sohbet Odaları: %s", strings.Join(rooms, ", ")))
}
func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("Lütfen oda adini yaziniz. usage: /join [OdaAdı]")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s odaya katildi", c.nick))

	c.msg(fmt.Sprintf(">Hosgeldin  %s odasina", roomName))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("mesaj cok kisa yada yanlis kod kullandiniz, usage: /msg [Mesaj]")
		fmt.Println("")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf(">%s adlı Kullanici cikti ", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg(">Tekrar bekleriz :)")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s odadan çıktı", c.nick))
	}
}

func (s *server) help(c *client) {

	c.msg(">Komutlar: \n\r *nick : Adini yaz.\r\n* join : odaya katil.\r\n* msg  : mesaj gonder.\r\n* rooms : sohbet odalarini görmek için.\r\n* help : help list.\r\n* quit : cikis yap.\r\n")

}
