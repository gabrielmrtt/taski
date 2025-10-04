package storagedatabase

import (
	"os"
	"path/filepath"

	"github.com/gabrielmrtt/taski/config"
)

type LocalStorageRepository struct {
	basePath string
}

func NewLocalStorageRepository() *LocalStorageRepository {
	basePath := config.GetInstance().StorageLocalBasePath

	return &LocalStorageRepository{
		basePath: basePath,
	}
}

func (r *LocalStorageRepository) GetFile(dir string, filename string) ([]byte, error) {
	path := filepath.Join(r.basePath, dir, filename)

	return os.ReadFile(path)
}

func (r *LocalStorageRepository) StoreFile(dir string, filename string, file []byte) error {
	path := filepath.Join(r.basePath, dir, filename)

	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}

func (r *LocalStorageRepository) DeleteFile(dir string, filename string) error {
	path := filepath.Join(r.basePath, dir, filename)

	return os.Remove(path)
}
