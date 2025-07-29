package main

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

type Database struct {
	Path  string
	Mu    sync.RWMutex
	Maps  map[string]map[string]string
	Dirty bool
}

func (d *Database) LoadFromFile() error {
	f, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if err.Error() == "EOF" {
			LogWarn("database", "Database file not found. Creating new one.")
			return nil
		}
		return err
	}

	defer f.Close()

	decoder := gob.NewDecoder(f)
	if err = decoder.Decode(&d.Maps); err != nil {
		if err.Error() == "EOF" {
			return nil
		}
		return err
	}

	return nil
}

func (d *Database) SaveToFile() error {
	f, err := os.OpenFile(d.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := gob.NewEncoder(f)
	if err = encoder.Encode(d.Maps); err != nil {
		return err
	}

	return nil
}

func NewDatabase(path string) *Database {
	return &Database{
		Path:  path,
		Maps:  make(map[string]map[string]string),
		Dirty: false,
	}
}

func (d *Database) CreateDB(name string) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	_, exists := d.Maps[name]
	if exists {
		return errors.New(ResponseAlreadyExists)
	}

	d.Dirty = true

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

	d.Dirty = true

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

	d.Dirty = true

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

	d.Dirty = true

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

	d.Dirty = true

	keyMaps[key] = newValue
	return nil
}

func (d *Database) List() (*[]string, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	if len(d.Maps) == 0 {
		return nil, errors.New(ResponseDatabaseIsEmpty)
	}

	var databases []string
	for k := range d.Maps {
		databases = append(databases, k)
	}

	return &databases, nil
}

func (d *Database) ListKeys(db string) ([]string, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	if len(d.Maps) == 0 {
		return nil, errors.New(ResponseDatabaseIsEmpty)
	}

	keysMap, exists := d.Maps[db]
	if !exists {
		return nil, errors.New(ResponseDatabaseNotFound)
	}

	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}

	return keys, nil
}
