package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	socketio "github.com/googollee/go-socket.io"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

// func newServer() {
// 	s, err := NewServer()
// 	if err != nil {
// 		log.Fatalf("new error %s", err)
// 	}
// 	go func() {
// 		err := s.Serve()
// 		if err != nil {
// 			log.Fatalf("serve error %s", err)
// 		}
// 	}()

// 	http.Handle("/socket.io/", s)
// 	log.Println("Serving at localhost:8000...")
// 	ts := httptest.NewServer(nil)
// 	defer ts.Close()
// }

func TestNewServer(t *testing.T) {
	_, err := NewServer()
	if err != nil {
		t.Errorf("new error %s", err)
	}
}

func TestNewClient(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		log.Fatalf("new error %s", err)
	}
	go func() {
		err := s.Serve()
		if err != nil {
			log.Fatalf("serve error %s", err)
		}
	}()

	http.Handle("/socket.io/", s)

	ts := httptest.NewServer(nil)
	log.Println("Serving at...", ts.URL)
	defer ts.Close()

	s.OnConnect("/", func(s socketio.Conn) error {
		return nil
	})

	opts := &socketio_client.Options{
		Transport: "websocket",
	}
	uri := ts.URL + "/socket.io/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})

	err = client.Emit("msg")
	if err != nil {
		log.Fatalln(err)
	}
	client.On("message", func(msg string) {
		log.Printf("on message:%v\n", msg)
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	data, _, _ := reader.ReadLine()
	// 	command := string(data)
	// 	client.Emit("message", command)
	// 	log.Printf("send message:%v\n", command)
	// }
}
