package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IndexOf(data []string, element int) string {
	for k, v := range data {
		if element == k {
			return v
		}
	}
	return "not found" //not found.
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)s
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func GetDirInfo(selected string, rootDir string, currPath string) ([]string, []string, string, int) {

	var addNav bool
	var cnt int    // total count of slice
	var p string   // selected dir or file
	var d []string // dir slice
	var f []string // file slice

	if currPath == rootDir {
		addNav = false
	} else {
		addNav = true
	}

	if addNav == true {
		d = append(d, "../")
	}

	err := os.Chdir(selected)
	newDir, err := os.Getwd()
	p = newDir
	check(err)

	c, err := ioutil.ReadDir(".")
	check(err)

	for _, entry := range c {
		if entry.IsDir() {
			d = append(d, entry.Name())
		} else {
			if strings.ToLower(filepath.Ext(entry.Name())) == ".ans" || strings.ToLower(filepath.Ext(entry.Name())) == ".asc" {
				f = append(f, entry.Name())
			}
		}
	}
	cnt = len(d) + len(f)
	return d, f, p, cnt
}
