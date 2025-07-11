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
	LogInfo("KVDB " + version + " is starting...")
	db = NewDatabase(file)

	if file != "" {
		err := db.LoadFromFile()
		if err != nil {
			LogError("failed to load database: " + err.Error())
		} else {
			LogInfo("Using file database.")
		}
	} else {
		LogInfo("Using in-memory database.")
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

func sendResponse(conn *net.Conn, msg string) {
	_, err := (*conn).Write([]byte(msg))
	if err != nil {
		LogError("failed to send response: " + err.Error())
	}
}

func IsValid(data string) bool {
	return strings.ContainsAny(data, ":")
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
			sendResponse(&conn, "BAD_QUERY")
			continue
		}

		switch {
		case query.CreateDb != nil:
			handleCreateDb(query.CreateDb, &conn)
		case query.Get != nil:
			handleGet(query.Get, &conn)
		case query.Set != nil:
			handleSet(query.Set, &conn)
		case query.Remove != nil:
			switch query.Remove.Which {
			case "DB":
				handleRemoveDb(query.Remove, &conn)
			case "KEY":
				handleRemoveKey(query.Remove, &conn)
			}
			sendResponse(&conn, "BAD_QUERY")
		case query.Update != nil:
			handleUpdate(query.Update, &conn)
		default:
			sendResponse(&conn, "BAD_QUERY")
		}
	}
}

func handleCreateDb(query *CreateDbQuery, conn *net.Conn) {
	name := query.Name
	if !IsValid(name) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	err := db.CreateDb(name)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database '" + name + "' has beed created.")

	sendResponse(conn, "OK")
	return
}

func handleGet(query *GetQuery, conn *net.Conn) {
	targetDb := query.Location.Db
	key := query.Location.Key

	if !IsValid(targetDb) || !IsValid(key) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	value, err := db.Get(targetDb, key)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	sendResponse(conn, value)
	return
}

func handleSet(query *SetQuery, conn *net.Conn) {
	value := strings.Trim(query.Value, "\"")
	targetDb := query.Location.Db
	key := query.Location.Key

	if !IsValid(value) || !IsValid(targetDb) || !IsValid(key) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	err := db.Add(targetDb, key, value)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("created key '" + key + "' on database '" + targetDb + "' with value '" + value + "'")
	sendResponse(conn, "OK")
}

func handleRemoveDb(query *RemoveQuery, conn *net.Conn) {
	targetDb := query.DB

	if !IsValid(targetDb) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	err := db.DeleteDb(targetDb)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database '" + targetDb + "' has been removed.")
	sendResponse(conn, "OK")
}

func handleRemoveKey(query *RemoveQuery, conn *net.Conn) {
	targetDb := query.DB
	key := query.Key

	if key == nil {
		sendResponse(conn, "KEY_NOT_PROVIDED")
		return
	}

	if !IsValid(targetDb) || !IsValid(*key) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	err := db.Remove(targetDb, *key)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("key '" + *key + "' from database '" + targetDb + "' has been removed")
	sendResponse(conn, "OK")
}

func handleUpdate(query *UpdateQuery, conn *net.Conn) {
	targetDb := query.Location.Db
	key := query.Location.Key
	value := query.Value

	if !IsValid(value) || !IsValid(targetDb) || !IsValid(key) {
		sendResponse(conn, "ILLEGAL_CHARACTERS")
		return
	}

	err := db.Update(targetDb, key, value)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("key '" + key + "' from database '" + targetDb + "' has been updated to value '" + value + "'")
	sendResponse(conn, "OK")
}
