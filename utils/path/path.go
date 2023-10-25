package path

import (
	"cicd_envsubst/utils/logger"
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	path     string
	fullpath string
}

const (
	FILE = iota
	DIRECTORY
)

var _logger = logger.New()

func New(path string) *Path {
	return &Path{
		path:     path,
		fullpath: getFullPath(path),
	}
}

func (p *Path) SetPath(path string) *Path {
	p.path = filepath.FromSlash(path)
	p.fullpath = getFullPath(path)
	return p
}

func (p *Path) GetPath() string {
	return p.path
}

func (p *Path) GetFullPath() string {
	return p.fullpath
}

func (p *Path) SetCompositePath(subPaths ...string) *Path {
	path := strings.Join(subPaths, "/")
	p.path = filepath.FromSlash(path)
	p.fullpath = getFullPath(path)
	return p
}

// Exists returns whether the given file or directory exists
func (p *Path) Exists() bool {
	_, err := os.Stat(p.path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (p *Path) MkdirAll(mode os.FileMode) {
	if err := os.MkdirAll(p.path, mode); err != nil {
		_logger.Errorf("Unable to `mkdir -p` on path { %s }: %v", p.path, err)
	}
}

func (p *Path) IsFileOrDir() int {
	fi, _ := os.Stat(p.path)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return DIRECTORY
	case mode.IsRegular():
		return FILE
	}

	return -1
}

func (p *Path) IsFile() bool {
	fi, _ := os.Stat(p.path)
	return fi.Mode().IsRegular()
}

func (p *Path) IsExistAndFile() bool {
	if p.Exists() {
		return p.IsFile()
	}
	return false
}

func (p *Path) IsExistAndDirectory() bool {
	if p.Exists() {
		return p.IsDirectory()
	}
	return false
}

func (p *Path) IsDirectory() bool {
	fi, _ := os.Stat(p.path)
	return fi.Mode().IsDir()
}

func (p *Path) Remove() *Path {
	if err := os.RemoveAll(p.GetPath()); err != nil {
		_logger.Errorf("[E] Removing { %s } from filesystem: %v", p.GetPath(), err)
	}
	return p
}

func (p *Path) RemoveIfExists() *Path {
	if p.Exists() {
		p.Remove()
	}
	return p
}

func (p *Path) ChdirIfDir() *Path {
	if p.IsDirectory() {
		err := os.Chdir(p.GetFullPath())
		if err != nil {
			_logger.Errorf("[E] switching to { %s } directory: %v", p.GetFullPath(), err)
		}
	}
	return p
}

//func (p *Path) GetFileObject() file.File {
//	var f = file.File{}
//	f.SetPath(p.path)
//	return f
//}

func getFullPath(initialPath string) string {
	fullPath, err := filepath.Abs(initialPath)
	if err != nil {
		_logger.Warnf("[E] getting full path for { %s }: %v", initialPath, err)
		return ""
	}
	return fullPath
}
