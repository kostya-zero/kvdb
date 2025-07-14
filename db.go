package main

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

func Save(db *Database) {
	err := db.SaveToFile()
	if err != nil {
		LogError("failed to save database: " + err.Error())
	}
}

type Database struct {
	Path string
	Mu   sync.RWMutex
	Maps map[string]map[string]string
}

type DatabaseMirror struct {
	Path string
	Maps map[string]map[string]string
}

func (d *Database) LoadFromFile() error {
	f, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if err.Error() == "EOF" {
			LogWarn("Database file not found. Creating new one.")
			return nil
		}
		return err
	}

	defer f.Close()

	var dbMirror DatabaseMirror
	decoder := gob.NewDecoder(f)
	if err = decoder.Decode(&dbMirror); err != nil {
		if err.Error() == "EOF" {
			LogWarn("Database file not found. Creating new one.")
			return nil
		}
		return err
	}

	d.Maps = dbMirror.Maps

	return nil
}

func (d *Database) SaveToFile() error {
	f, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	dbMirror := DatabaseMirror{
		Maps: d.Maps,
	}

	encoder := gob.NewEncoder(f)
	if err = encoder.Encode(dbMirror); err != nil {
		return err
	}

	return nil
}

func NewDatabase(path string) *Database {
	return &Database{
		Path: path,
		Maps: make(map[string]map[string]string),
	}
}

func (d *Database) CreateDB(name string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	_, exists := d.Maps[name]
	if exists {
		return errors.New(ResponseAlreadyExists)
	}

	if d.Path != "" {
		defer Save(d)
	}

	d.Maps[name] = make(map[string]string)

	return nil
}

func (d *Database) DeleteDB(name string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	_, exists := d.Maps[name]
	if !exists {
		return errors.New(ResponseDatabaseNotFound)
	}

	if d.Path != "" {
		defer Save(d)
	}

	delete(d.Maps, name)

	return nil
}

func (d *Database) Add(db string, key string, value string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	keysMap, exists := d.Maps[db]
	if !exists {
		return errors.New(ResponseDatabaseNotFound)
	}

	_, exists = keysMap[key]
	if exists {
		return errors.New(ResponseAlreadyExists)
	}

	if d.Path != "" {
		defer Save(d)
	}

	keysMap[key] = value
	return nil
}

func (d *Database) Get(db string, key string) (string, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	keysMap, exists := d.Maps[db]
	if !exists {
		return "", errors.New(ResponseDatabaseNotFound)
	}

	value, exists := keysMap[key]
	if !exists {
		return "", errors.New(ResponseKeyNotFound)
	}

	return value, nil
}

func (d *Database) Remove(db string, key string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	keysMap, exists := d.Maps[db]
	if !exists {
		return errors.New(ResponseDatabaseNotFound)
	}

	_, exists = keysMap[key]
	if !exists {
		return errors.New(ResponseKeyNotFound)
	}

	if d.Path != "" {
		defer Save(d)
	}

	delete(keysMap, key)
	return nil
}

func (d *Database) Update(db string, key string, newValue string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	keyMaps, exists := d.Maps[db]
	if !exists {
		return errors.New(ResponseDatabaseNotFound)
	}

	_, exists = keyMaps[key]
	if !exists {
		return errors.New(ResponseKeyNotFound)
	}

	if d.Path != "" {
		defer Save(d)
	}

	keyMaps[key] = newValue
	return nil
}
