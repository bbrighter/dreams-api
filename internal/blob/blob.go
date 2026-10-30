package blob

import (
	"os"

	"path/filepath"
	"strconv"
)

type BlobStorage struct {
	folder string
}

func (b *BlobStorage) init() error {
	return os.MkdirAll(b.folder, 0700) // Only owner can do everything
}

func (b *BlobStorage) filePath(id uint) string {
	return filepath.Join(b.folder, "recording_"+strconv.FormatUint(uint64(id), 10))
}

func NewBlobStorage(folderName string) *BlobStorage {
	storage := &BlobStorage{folder: folderName}
	if err := storage.init(); err != nil {
		panic(err)
	}
	return storage
}

func (b *BlobStorage) Get(id uint) ([]byte, error) {
	file, err := os.ReadFile(b.filePath(id))
	return file, err
}

func (b *BlobStorage) Delete(id uint) error {
	return os.Remove(b.filePath(id))
}

func (b *BlobStorage) Create(id uint, file []byte) error {
	return os.WriteFile(b.filePath(id), file, 0600) // Only owner can read and write
}
