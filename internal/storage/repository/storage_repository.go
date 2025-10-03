package storagerepo

type StorageRepository interface {
	GetFile(dir string, filename string) ([]byte, error)
	StoreFile(dir string, filename string, file []byte) error
	DeleteFile(dir string, filename string) error
}
