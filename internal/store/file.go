package store

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var ErrNotFound = errors.New("store: key not found")

type FileStore struct {
	root string
	mu   sync.Mutex
}

func New(root string) *FileStore {
	return &FileStore{root: root}
}

func (s *FileStore) Save(key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.pathFor(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *FileStore) Load(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.pathFor(key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

func (s *FileStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.pathFor(key))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *FileStore) Exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, err := os.Stat(s.pathFor(key))
	return err == nil && !info.IsDir()
}

func (s *FileStore) Root() string {
	return s.root
}

func (s *FileStore) pathFor(key string) string {
	name := strings.ReplaceAll(key, "/", string(os.PathSeparator))
	return filepath.Join(s.root, name+".json")
}
