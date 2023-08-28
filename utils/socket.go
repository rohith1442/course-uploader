package utils

import (
	"fmt"
	"log"

	// "net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func socket() {
	router := gin.New()

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Socket id connected:", s.ID(), s.RemoteAddr(), s.LocalAddr())
		return nil
	})
	server.OnEvent("/", "test", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "Calm down bro, we are still connected. Stop sending this `"+msg+"` message")
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		err := s.Close()
		if err != nil {
			return ""
		}
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("Socket id disconnected", s.ID(), msg)
	})
	go func() {
		if err := server.Serve(); err != nil {
			fmt.Printf("socketio listen error: %s\n", err)
		}
	}()
	defer func(server *socketio.Server) {
		err := server.Close()
		if err != nil {
			fmt.Println("Socket serveris close", err)
		}
	}(server)

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

}
