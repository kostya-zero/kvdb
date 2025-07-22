package main

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var db *Database

func StartServer(port int, file string, saveInterval int) error {
	LogInfo("server", "KVDB "+version+" is starting...")
	db = NewDatabase(file)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	if file != "" {
		err := db.LoadFromFile()
		if err != nil {
			LogError("server", "failed to load database: "+err.Error())
			LogWarn("server", "Falling back to in-memory database")
		} else {
			LogInfo("server", "Using file database.")
			wg.Add(1)
			go func() {
				defer wg.Done()
				BackupService(ctx, time.Millisecond*time.Duration(saveInterval), db)
			}()
		}
	} else {
		LogInfo("server", "Using in-memory database.")
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listener.Close()

	LogInfo("server", "Starting TCP server on port "+strconv.Itoa(port))

	go func() {
		<-ctx.Done()
		LogInfo("server", "Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			LogError("server", "failed to accept connection: "+err.Error())
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer conn.Close()
			handleConn(conn)
		}(conn)
	}

	// I'll left if for a while.
	// LogInfo("Waiting for all workers to finish.")
	// wg.Wait()
	// LogInfo("Goodbye.")
	return nil
}

func BackupService(ctx context.Context, interval time.Duration, db *Database) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			PerformBackup(db)
		case <-ctx.Done():
			LogInfo("backup_service", "backing up database...")
			PerformBackup(db)
			return
		}
	}
}

func PerformBackup(db *Database) {
	if !db.Dirty {
		return
	}

	if err := db.SaveToFile(); err != nil {
		LogError("backup_service", "backup failed: "+err.Error())
	} else {
		LogInfo("backup_service", "database backup completed.")
		db.Dirty = false
	}
}

func sendResponse(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		LogError("server", "failed to send response: "+err.Error())
	}
}

var bufPool = sync.Pool{
	New: func() any {
		return make([]byte, 2048)
	},
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(&buf)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			LogError("server", "failed to read from connection: "+err.Error())
			return
		}
		receiveQuery := string(buf[:n])
		query, err := parseQuery(receiveQuery)
		if err != nil {
			sendResponse(conn, ResponseBadQuery)
			continue
		}

		switch {
		case query.CreateDB != nil:
			handleCreateDB(query.CreateDB, conn)
		case query.Get != nil:
			handleGet(query.Get, conn)
		case query.Set != nil:
			handleSet(query.Set, conn)
		case query.Remove != nil:
			switch query.Remove.Which {
			case "DB":
				handleRemoveDB(query.Remove, conn)
			case "KEY":
				handleRemoveKey(query.Remove, conn)
			default:
				sendResponse(conn, ResponseBadQuery)
			}
		case query.Update != nil:
			handleUpdate(query.Update, conn)
		default:
			sendResponse(conn, ResponseBadQuery)
		}
	}
}

func IsValid(data string) bool {
	return strings.Contains(data, ":")
}

func handleCreateDB(query *CreateDBQuery, conn net.Conn) {
	name := query.Name

	err := db.CreateDB(name)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database", "map '"+name+"' has beed created.")

	sendResponse(conn, ResponseOk)
}

func handleGet(query *GetQuery, conn net.Conn) {
	targetDB := query.Location.DB
	key := query.Location.Key

	value, err := db.Get(targetDB, key)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	sendResponse(conn, value)
}

func handleSet(query *SetQuery, conn net.Conn) {
	value := strings.Trim(query.Value, "\"")
	targetDB := query.Location.DB
	key := query.Location.Key

	if IsValid(value) {
		sendResponse(conn, ResponseIllegalChars)
		return
	}

	err := db.Add(targetDB, key, value)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database", "created key '"+key+"' on map '"+targetDB+"' with value '"+value+"'")
	sendResponse(conn, ResponseOk)
}

func handleRemoveDB(query *RemoveQuery, conn net.Conn) {
	targetDB := query.DB

	err := db.DeleteDB(targetDB)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database", "map '"+targetDB+"' has been removed.")
	sendResponse(conn, ResponseOk)
}

func handleRemoveKey(query *RemoveQuery, conn net.Conn) {
	targetDB := query.DB
	key := query.Key

	if key == nil {
		sendResponse(conn, ResponseKeyNotProvided)
		return
	}

	err := db.Remove(targetDB, *key)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database", "key '"+*key+"' from map '"+targetDB+"' has been removed")
	sendResponse(conn, ResponseOk)
}

func handleUpdate(query *UpdateQuery, conn net.Conn) {
	targetDB := query.Location.DB
	key := query.Location.Key
	value := strings.Trim(query.Value, "\n")

	if IsValid(value) {
		sendResponse(conn, ResponseIllegalChars)
		return
	}

	err := db.Update(targetDB, key, value)
	if err != nil {
		sendResponse(conn, err.Error())
		return
	}

	LogInfo("database", "key '"+key+"' from map '"+targetDB+"' has been updated to value '"+value+"'")
	sendResponse(conn, ResponseOk)
}
