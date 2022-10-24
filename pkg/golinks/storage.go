package golinks

import (
	"fmt"
	"net/url"
	"sync"
)

type StorageFunc func() Storage

var storageFuncs map[string]StorageFunc
var storageFuncsLock sync.Mutex

func init() {
	storageFuncs = make(map[string]StorageFunc)
	storageFuncsLock = sync.Mutex{}
}

func RegisterStorage(storageType string, storageFunc StorageFunc) {
	storageFuncsLock.Lock()
	defer storageFuncsLock.Unlock()

	if _, exists := storageFuncs[storageType]; exists {
		errMsg := fmt.Sprintf("[ERROR] storage type '%s' already exists", storageType)
		panic(errMsg)
	}
	storageFuncs[storageType] = storageFunc
}

type Storage interface {
	GetLink(name string) (*url.URL, error)
	SetLink(name string, url url.URL) error

	Migrate() error
}

func NewStorage(storageType string) Storage {
	f, ok := storageFuncs[storageType]
	if !ok {
		errMsg := fmt.Sprintf("[ERROR] storage type '%s' does not exist", storageType)
		panic(errMsg)
	}
	return f()
}
