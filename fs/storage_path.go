package fs

import "path/filepath"

type Path struct {
	root     string
	bucket   string
	dir      string
	filename string
}

func NewPath(root, bucket, dir, filename string) *Path {
	return &Path{
		root:     root,
		bucket:   bucket,
		dir:      dir,
		filename: filename,
	}
}

func (p *Path) FullFilePath() string {
	return filepath.Join(p.root, p.bucket, p.dir, p.filename)
}

func (p *Path) FilePath() string {
	return filepath.Join(p.bucket, p.dir, p.filename)
}

func (p *Path) FullDirPath() string {
	return filepath.Join(p.root, p.bucket, p.dir, p.filename)
}

func (p *Path) DirPath() string {
	return filepath.Join(p.bucket, p.dir, p.filename)
}
