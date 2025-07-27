package fs

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type SubDir string

const (
	Input  SubDir = "input"
	Output SubDir = "output"
)

type FileStorage struct {
	Root      string
	SubDir    SubDir
	Threshold time.Duration
	DirLayout string
}

func NewFileStorage(root string, dir SubDir, threshold int) *FileStorage {
	return &FileStorage{
		Root:      filepath.Join(root, string(dir)),
		SubDir:    dir,
		Threshold: time.Duration(threshold) * time.Minute,
		DirLayout: "20060102-1504",
	}
}

func (s *FileStorage) getBucketName() string {
	t := time.Now().UTC()
	rounded := t.Truncate(s.Threshold)
	return rounded.Format(s.DirLayout)
}

func (s *FileStorage) Save(reader io.Reader, filename string, subDirs ...string) (string, error) {
	var (
		dirPath []string
	)
	bucket := s.getBucketName()
	dirPath = append(dirPath, s.Root, bucket)
	dirPath = append(dirPath, subDirs...)
	path := filepath.Join(dirPath...)

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	outFile, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, reader); err != nil {
		return "", err
	}

	visiblePath := filepath.Join(dirPath[1:]...)
	return filepath.Join(visiblePath, filename), nil
}

func (s *FileStorage) Dirs() ([]string, error) {
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

func (s *FileStorage) RemoveDir(dir string) error {
	return os.RemoveAll(filepath.Join(s.Root, dir))
}

func (s *FileStorage) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}
