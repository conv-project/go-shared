package fs

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type DataType string

const (
	Input  DataType = "input"
	Output DataType = "output"
)

type FileStorage struct {
	Root      string
	Type      DataType
	Threshold time.Duration
	DirLayout string
}

func NewFileStorage(root string, t DataType, threshold int) *FileStorage {
	return &FileStorage{
		Root:      filepath.Join(root, string(t)),
		Type:      t,
		Threshold: time.Duration(threshold) * time.Minute,
		DirLayout: "20060102-1504",
	}
}

func (s *FileStorage) getBucketName() string {
	t := time.Now().UTC()
	rounded := t.Truncate(s.Threshold)
	return rounded.Format(s.DirLayout)
}

func (s *FileStorage) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *FileStorage) WriteFile(reader io.Reader, filename string, subDirs ...string) (*Path, error) {
	path := s.GetPath(filename, subDirs...)
	if err := s.CreateDir(path.FullDirPath()); err != nil {
		return nil, err
	}

	outFile, err := os.Create(path.FullFilePath())
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, reader); err != nil {
		return nil, err
	}

	return path, nil
}

func (s *FileStorage) GetPath(filename string, subDirs ...string) *Path {
	return NewPath(s.Root, s.getBucketName(), filepath.Join(subDirs...), filename)
}

func (s *FileStorage) GetBuckets() ([]string, error) {
	var (
		dirs []string
	)

	entries, err := os.ReadDir(s.Root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

func (s *FileStorage) RemoveBucket(bucketName string) error {
	return os.RemoveAll(filepath.Join(s.Root, bucketName))
}
