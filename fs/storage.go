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

func (s *FileStorage) DirPath(subDirs ...string) string {
	var (
		path []string
	)
	bucket := s.getBucketName()
	path = append(path, bucket)
	path = append(path, subDirs...)
	return filepath.Join(path...)
}

func (s *FileStorage) RooDirPath(subDirs ...string) string {
	return filepath.Join(s.Root, s.DirPath(subDirs...))
}

func (s *FileStorage) FilePath(filename string, subDirs ...string) string {
	return filepath.Join(s.DirPath(subDirs...), filename)
}

func (s *FileStorage) RootFilePath(filename string, subDirs ...string) string {
	return filepath.Join(s.Root, s.FilePath(filename, subDirs...))
}

func (s *FileStorage) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *FileStorage) Save(reader io.Reader, filename string, subDirs ...string) (string, error) {
	dirPath := s.DirPath(subDirs...)
	rootDirPath := filepath.Join(s.Root, dirPath)

	if err := s.CreateDir(rootDirPath); err != nil {
		return "", err
	}

	outFile, err := os.Create(filepath.Join(rootDirPath, filename))
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, reader); err != nil {
		return "", err
	}

	return filepath.Join(dirPath, filename), nil
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
