package storage

import (
	"fmt"
	jfscmd "github.com/juicedata/juicefs/cmd"
	"net/url"
	"os"
	"strings"
	"time"
)

func init() {
	register(StorageTypeJuiceFS, getJuiceFSStorage)
}

func getJuiceFSStorage() Storage {
	return new(JuiceFSStorage)
}

type JuiceFSStorage struct {
	StorageDesc
}

func (fs *JuiceFSStorage) Init(desc StorageDesc) error {
	fs.StorageDesc = desc
	return fs.mount(desc)
}

func (fs *JuiceFSStorage) Deinit() error {
	args := []string{
		"juicefs",
		"umount",
		fs.MountPoint,
	}
	var err error = nil
	for i := 0; i < 10; i++ {
		if err = jfscmd.Main(args); err == nil {
			return nil
		}
		time.Sleep(time.Second * 10)
	}
	return err
}

func (fs *JuiceFSStorage) CreateNode(node string) error {
	if _, err := os.Stat(node); err == nil {
		return fmt.Errorf("node already exists: %s", node)
	}
	return os.Mkdir(node, os.ModePerm)
}

func (fs *JuiceFSStorage) DeleteNode(node string) error {
	if s, err := os.Stat(node); err != nil {
		return nil
	}
	if err := os.Remove(node); err != nil {
		return err
	}
	return nil
}

func (fs *JuiceFSStorage) ListNodes() ([]string, error) {
	entries, err := os.ReadDir(fs.MountPoint)
	if err != nil {
		return []string{}, err
	}
	res := make([]string, len(entries))
	for i, entry := range entries {
		res[i] = entry.Name()
	}
	return res, nil
}

func (fs *JuiceFSStorage) constructUri(uri, id, key string) (string, error) {
	if id == "" && key == "" {
		return uri, nil
	}
	parts := strings.Split(uri, "://")
	if len(parts) != 2 {
		return uri, fmt.Errorf("invalid uri: %s", uri)
	}
	return fmt.Sprintf("%s://%s:%s@%s", parts[0], id, url.QueryEscape(key), parts[1]), nil
}
func (fs *JuiceFSStorage) mount(desc StorageDesc) error {
	args := []string{
		"juicefs",
		"mount",
		"--writeback",
		"--cache-dir", fs.CacheDir,
		"--cache-size", fmt.Sprintf("%d", fs.CacheSize),
	}
	if fs.SubDir != "" {
		args = append(args, fs.SubDir)
	}
	if fs.ReadOnly {
		args = append(args, "--read-only")
	}
	if uri, err := fs.constructUri(fs.Uri, fs.Id, fs.Key); err == nil {
		args = append(args, uri)
	} else {
		return err
	}
	args = append(args, fs.MountPoint)
	return jfscmd.Main(args)
}
