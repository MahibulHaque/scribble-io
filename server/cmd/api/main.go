package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	scribbleServer "scribble/internal/server"

	socketio "github.com/googollee/go-socket.io"
)

func main() {

	server := scribbleServer.NewSocketIoServer()
	server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext(conn.Context())
		log.Println("connected:", conn.ID())
		return nil
	})
	scribbleServer.RegisterEventHandlers(server)
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../asset")))

	log.Printf("Serving at localhost:%s...\n", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
