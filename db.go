package main

import (
	"errors"
	"sync"
)

type Database struct {
	mu sync.RWMutex
	m  map[string]string
}

func (d *Database) Add(key string, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, exists := d.m[key]
	if exists {
		return errors.New("ALREADY_EXISTS")
	}

	d.m[key] = value
	return nil
}

func (d *Database) Get(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, exists := d.m[key]
	if !exists {
		return "", errors.New("NOT_FOUND")
	}
	return value, nil
}

func (d *Database) Remove(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, exists := d.m[key]
	if !exists {
		return errors.New("key not found")
	}

	delete(d.m, key)
	return nil
}

func (d *Database) Update(key string, newValue string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, exists := d.m[key]
	if !exists {
		return errors.New("key not found")
	}

	d.m[key] = newValue
	return nil
}
