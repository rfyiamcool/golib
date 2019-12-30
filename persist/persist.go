package persist

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	extension = ".gob"
	prefix    = "_sc_"
)

var mutexList = make(map[string]*sync.Mutex)

type Collection struct {
	mutex sync.Mutex
	dir   string
	items []string
}

func New(name string) (*Collection, error) {
	if len(name) <= 0 {
		return &Collection{}, errors.New("Collection name can not be empty!")
	}

	//make file path correct
	dir := prefix + filepath.Clean(name)
	collection := Collection{
		dir: dir,
	}
	return &collection, os.MkdirAll(dir, 0755)
}

func (c *Collection) Put(key string, value interface{}) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty!")
	}

	path := filepath.Join(c.dir, key+extension)
	m := c.getMutex(path)
	m.Lock()
	defer m.Unlock()

	file, err := os.Create(path)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(value)
	}
	return err
}

func (c *Collection) Get(key string, value interface{}) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty!")
	}

	path := filepath.Join(c.dir, key+extension)
	m := c.getMutex(path)
	m.Lock()
	defer m.Unlock()
	if !c.Has(key) {
		return fmt.Errorf("Key %s does not exist!", key)
	}

	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(value)
	}
	return err
}

func (c *Collection) Remove(key string) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty!")
	}

	path := filepath.Join(c.dir, key+extension)
	m := c.getMutex(path)
	m.Lock()
	defer m.Unlock()
	if c.Has(key) {
		return os.Remove(path)
	}
	return fmt.Errorf("Key %s does not exist!", key)
}

func (c *Collection) Flush() error {
	if _, err := os.Stat(c.dir); err == nil {
		os.RemoveAll(c.dir)
		return err
	}
	return nil
}

func (c *Collection) Has(key string) bool {
	if len(key) <= 0 {
		return false
	}

	path := filepath.Join(c.dir, key+extension)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func (c *Collection) List() ([]string, error) {
	var (
		items []string
	)
	files, err := ioutil.ReadDir(c.dir)
	if err != nil {
		return items, err
	}

	for _, f := range files {
		item := f.Name()
		item = strings.Trim(item, extension)
		items = append(items, item)
	}
	return items, err
}

func (c *Collection) TotalItem() int {
	list, _ := c.List()
	return len(list)
}

func (c *Collection) getMutex(path string) *sync.Mutex {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m, ok := mutexList[path]
	if !ok {
		m = &sync.Mutex{}
		mutexList[path] = m
	}
	return m
}
