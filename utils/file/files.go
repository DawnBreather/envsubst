package file

import (
	path2 "cicd_envsubst/utils/path"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func FindFilesByRegexpRecursively(basePath, pattern string) []*File {
	libRegEx, e := regexp.Compile(pattern)
	if e != nil {
		log.Fatal(e)
	}

	var resFiles []*File

	e = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			pth := path2.Path{}
			if pth.SetPath(path).IsFile() {
				resFiles = append(resFiles, getFileFromPath(pth))
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	return resFiles
}

func FindFilesRecursively(basePath string) []*File {
	var e error
	var res []*File

	e = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			var pth = path2.Path{}
			if pth.SetPath(path).IsFile() {
				f := getFileFromPath(pth)
				res = append(res, f)
			}
		}
		return nil
	})
	if e != nil {
		log.Fatal(e)
	}

	return res
}

func RemoveFilesRecursively(basePath string) []*File {
	var res []*File

	res = FindFilesRecursively(basePath)
	for _, f := range res {
		f.Delete()
	}

	return res
}

func getFileFromPath(p path2.Path) *File {
	var f = File{}
	f.SetPath(p.GetPath())
	return &f
}
