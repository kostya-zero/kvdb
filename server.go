package main

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var db Database

func StartServer() error {
	db = Database{
		m: make(map[string]string),
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return err
	}
	defer listener.Close()
	println("Starting TCP server on port 3000")

	go func() {
		<-ctx.Done()
		println("Received CTRL+C. Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			println("failed to accept connection: " + err.Error())
			continue
		}

		go handleConn(conn)
	}
	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			println("failed to read from connection: " + err.Error())
			return
		}
		receiveQuery := string(buf[:n])
		println("Received query: " + receiveQuery)
		query, err := parseQuery(receiveQuery)
		if err != nil {
			conn.Write([]byte("BAD_QUERY"))
			continue
		}

		switch {
		case query.Get != nil:
			value, err := db.Get(query.Get.Key)
			if err != nil {
				conn.Write([]byte(err.Error()))
				continue
			}
			conn.Write([]byte(value))
		case query.Set != nil:
			value := strings.Trim(query.Set.Value, "\"")
			err = db.Add(query.Set.Key, value)
			if err != nil {
				conn.Write([]byte(err.Error()))
				continue
			}
			conn.Write([]byte("OK"))
		case query.Delete != nil:
			err = db.Remove(query.Delete.Key)
			if err != nil {
				conn.Write([]byte(err.Error()))
				continue
			}
			conn.Write([]byte("OK"))
		default:
			conn.Write([]byte("BAD_QUERY"))
		}

	}
}
