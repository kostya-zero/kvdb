package main

import (
	"errors"
	"sync"
)

type Database struct {
	mu   sync.RWMutex
	maps map[string]map[string]string
}

func NewDatabase() *Database {
	return &Database{
		maps: make(map[string]map[string]string),
	}
}

func (d *Database) CreateDb(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exists := d.maps[name]
	if exists {
		return errors.New("ALREADY_EXISTS")
	}

	d.maps[name] = make(map[string]string)

	return nil
}

func (d *Database) DeleteDb(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exists := d.maps[name]
	if !exists {
		return errors.New("DATABASE_NOT_FOUND")
	}

	delete(d.maps, name)

	return nil
}

func (d *Database) Add(db string, key string, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	keysMap, exists := d.maps[db]
	if !exists {
		return errors.New("DATABASE_NOT_FOUND")
	}

	_, exists = keysMap[key]
	if exists {
		return errors.New("ALREADY_EXISTS")
	}

	keysMap[key] = value
	return nil
}

func (d *Database) Get(db string, key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	keysMap, exists := d.maps[db]
	if !exists {
		return "", errors.New("DATABASE_NOT_FOUND")
	}

	value, exists := keysMap[key]
	if !exists {
		return "", errors.New("KEY_NOT_FOUND")
	}

	return value, nil
}

func (d *Database) Remove(db string, key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	keysMap, exists := d.maps[db]
	if !exists {
		return errors.New("DATABASE_NOT_FOUND")
	}

	_, exists = keysMap[key]
	if !exists {
		return errors.New("KEY_NOT_FOUND")
	}

	delete(keysMap, key)
	return nil
}

func (d *Database) Update(db string, key string, newValue string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	keyMaps, exists := d.maps[db]
	if !exists {
		return errors.New("DATABASE_NOT_FOUND")
	}

	_, exists = keyMaps[key]
	if !exists {
		return errors.New("KEY_NOT_FOUND")
	}

	keyMaps[key] = newValue
	return nil
}
