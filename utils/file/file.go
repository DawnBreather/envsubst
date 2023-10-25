package file

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/logger"
	"cicd_envsubst/utils/placeholder"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _logger = logger.New()

type File struct {
	path    string
	content []byte
	base64  string
	md5     string
}

func (f File) GetFileName() string {
	return filepath.Base(f.path)
}
func (f File) GetBaseDir() string {
	dir, _ := filepath.Split(f.path)
	return dir
}
func (f File) GetPath() string {
	return f.path
}
func (f File) GetContent() []byte {
	return f.content
}
func (f File) GetBase64() string {
	return f.base64
}
func (f File) GetMd5() string {
	return f.md5
}

func (f *File) SetPath(path string) *File {
	f.path = path
	return f
}
func (f *File) SetCompositePath(subPaths ...string) *File {
	path := strings.Join(subPaths, "/")
	f.path = filepath.FromSlash(path)
	return f
}
func (f *File) SetContent(content []byte) *File {
	f.content = content
	return f
}
func (f *File) SetBase64(b64string string) *File {
	f.base64 = b64string
	return f
}

func (f *File) ParseContentToBase64() *File {
	f.base64 = base64.StdEncoding.EncodeToString(f.content)
	return f
}

func (f *File) ParseContentToMd5() *File {
	f.md5 = fmt.Sprintf("%x", md5.Sum([]byte(f.content)))
	return f
}

func (f *File) ParseBase64ToContent() *File {
	var err error
	f.content, err = base64.StdEncoding.DecodeString(f.base64)
	if err != nil {
		_logger.Errorf("Unable to decode base64 encoded string into file content: %v", err)
	}
	return f
}

func (f *File) ReadContent() *File {
	var err error
	f.content, err = ioutil.ReadFile(f.path)
	if err != nil {
		_logger.Errorf("Error reading file { %s }: %v", f.path, err)
	}
	return f
}

func (f *File) SaveTo(dstPath string) *File {
	var err error

	dstPath = filepath.FromSlash(dstPath)

	if dstPath != "" {
		err = ioutil.WriteFile(dstPath, f.content, 0644)
		if err != nil {
			_logger.Errorf("Unable to save file content into { %s }: %v", dstPath, err)
		}
		if f.path == "" {
			f.path = dstPath
		}
	} else {
		_logger.Errorf("Unable to save file content: %v", "no destination path provided")
	}

	return f
}

func (f *File) FillContentFromEnvironmentVariable(envVarName string) *File {
	e := env_var.EnvVar{}
	e.SetName(envVarName)
	if e.IsExist() {
		f.SetContent([]byte(e.Value()))
	} else {
		_logger.Errorf("Unable to fill file content from environment variable { %s }: not found", envVarName)
	}

	return f
}

func (f *File) Save() *File {
	return f.SaveTo(f.path)
}
func (f *File) Delete() *File {
	err := os.Remove(f.path)
	if err != nil {
		_logger.Errorf("Unable to remove file { %s }: %v", f.path, err)
	}

	return f
}

func (f *File) ReplaceEnvVarsPlaceholder(prefix, suffix string) *File {
	var p = placeholder.New(prefix, suffix, env_var.REGEX_MASK)
	f.content = []byte(p.ReplacePlaceholdersWithEnvVars(string(f.content)))

	return f
}

func (f *File) ReplaceVarsSetPlaceholder(varsSet map[string]string, prefix, suffix string) *File {
	var p = placeholder.
		New(prefix, suffix, env_var.REGEX_MASK).
		SetVarsSet(varsSet)
	f.content = []byte(p.ReplacePlaceholdersWithVarsSet(string(f.content)))

	return f
}

func (f *File) ReplaceEnvVarsPlaceholderByExplicitRegex(prefix, suffix string, regexMask string) *File {
	var p = placeholder.New(prefix, suffix, regexMask)
	f.content = []byte(p.ReplacePlaceholdersWithEnvVars(string(f.content)))

	return f
}

func (f *File) ReplaceStringAll(old string, new string) *File {
	f.content = []byte(strings.ReplaceAll(string(f.content), old, new))

	return f
}

func (f *File) IsExists() (res bool) {
	if _, err := os.Stat(f.path); !errors.Is(err, os.ErrNotExist) {
		res = true
	}
	return
}
