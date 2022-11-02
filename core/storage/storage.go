package storage

import "github.com/lunny/log"

const (
	StorageTypeJuiceFS = "juicefs"
)

type StorageDesc struct {
	Uri        string
	Id         string
	Key        string
	MountPoint string
	CacheDir   string
	CacheSize  int64
	SubDir     string
	ReadOnly   bool
}

type Storage interface {
	Init(desc StorageDesc) error
	Deinit() error
	CreateNode(node string) error
	DeleteNode(node string) error
	ListNodes() ([]string, error)
}

var supportedStorageMap = make(map[string]func() Storage, 0)

func register(name string, callFunc func() Storage) {
	supportedStorageMap[name] = callFunc
}

func GetStorage(storageType string) Storage {
	if _, ok := supportedStorageMap[storageType]; !ok {
		log.Fatalf("Storage Type: %s not supported", storageType)
		return nil
	}
	return supportedStorageMap[storageType]()
}
