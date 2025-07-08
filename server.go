package main

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var db *Database

func StartServer(port int, file string) error {
	db = NewDatabase(file)

	if file != "" {
		err := db.LoadFromFile()
		if err != nil {
			LogError("failed to load database: " + err.Error())
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listener.Close()
	LogInfo("Starting TCP server on port " + strconv.Itoa(port))

	go func() {
		<-ctx.Done()
		LogInfo("Received CTRL+C. Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			LogError("failed to accept connection: " + err.Error())
			continue
		}

		go handleConn(conn)
	}
}

func sendResponse(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		LogError("failed to send response: " + err.Error())
	}
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
			LogError("failed to read from connection: " + err.Error())
			return
		}
		receiveQuery := string(buf[:n])
		query, err := parseQuery(receiveQuery)
		if err != nil {
			sendResponse(conn, "BAD_QUERY")
			continue
		}

		switch {
		case query.CreateDb != nil:
			err = db.CreateDb(query.CreateDb.Name)
			if err != nil {
				sendResponse(conn, err.Error())
				continue
			}
			LogInfo("database '" + query.CreateDb.Name + "' has beed created.")
			sendResponse(conn, "OK")
		case query.Get != nil:
			value, err := db.Get(query.Get.Location.Db, query.Get.Location.Key)
			if err != nil {
				sendResponse(conn, err.Error())
				continue
			}
			sendResponse(conn, value)
		case query.Set != nil:
			value := strings.Trim(query.Set.Value, "\"")
			err = db.Add(query.Set.Location.Db, query.Set.Location.Key, value)
			if err != nil {
				sendResponse(conn, err.Error())
				continue
			}
			LogInfo("created key '" + query.Set.Location.Key + "' on database '" + query.Set.Location.Db + "' with value '" + query.Set.Value + "'")
			sendResponse(conn, "OK")
		case query.Remove != nil:
			switch query.Remove.Which {
			case "DB":
				err = db.DeleteDb(query.Remove.DB)
				if err != nil {
					sendResponse(conn, err.Error())
					continue
				}
				LogInfo("database '" + query.Remove.DB + "' has been removed.")
				sendResponse(conn, "OK")
			case "KEY":
				err = db.Remove(query.Remove.DB, *query.Remove.Key)
				if err != nil {
					sendResponse(conn, err.Error())
					continue
				}
				LogInfo("key '" + *query.Remove.Key + "' from database '" + query.Remove.DB + "' has been removed")
				sendResponse(conn, "OK")
			}
			sendResponse(conn, "BAD_QUERY")
		case query.Update != nil:
			err = db.Update(query.Update.Location.Db, query.Update.Location.Key, query.Update.Value)
			if err != nil {
				sendResponse(conn, err.Error())
				continue
			}
			LogInfo("key '" + query.Update.Location.Key + "' from database '" + query.Update.Location.Db + "' has been updated to value '" + query.Update.Value + "'")
			conn.Write([]byte("OK"))
		default:
			sendResponse(conn, "BAD_QUERY")
		}
	}
}
