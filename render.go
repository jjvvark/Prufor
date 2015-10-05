package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type walker func(path string, f os.FileInfo) error

var setter walker

func initRender() {

	var copyer walker

	copyer = func(path string, f os.FileInfo) error {

		var err error

		s := fmt.Sprintf("%s/%s", path, f.Name())
		d := fmt.Sprintf("%s%s/%s", destDir, path[len(sourceDir):], f.Name())

		if f.IsDir() {

			err = os.MkdirAll(d, 0777)
			if err != nil {
				return err
			}

			err = readDir(s, copyer)
			if err != nil {
				return err
			}

		} else {

			if !checkFileType(f.Name(), "jst") {
				err = CopyFile(s, d)
				if err != nil {
					return err
				}
			}

		}

		return nil
	}

	setter = func(path string, f os.FileInfo) error {

		n := f.Name()
		s := fmt.Sprintf("%s/%s", path, n)

		if f.IsDir() {

			readDir(s, setter)

		} else {

			if checkFileType(n, "jst") {

				// render with data!
				d := fmt.Sprintf("%s%s/%s", destDir, path[len(sourceDir):], n[:len(n)-4])

				t := template.New("template")
				t, err := t.ParseFiles(s)
				if err != nil {
					return err
				}

				o, err := os.Create(d)
				if err != nil {
					return err
				}

				defer o.Close()

				writer := bufio.NewWriter(o)
				defer writer.Flush()

				t.Execute(writer, map[string]interface{}{
					"p": GetData(),
				})

			}

		}

		return nil

	}

	var err error

	err = os.RemoveAll(destDir)
	if err != nil {
		log.Panic(err)
	}

	err = os.MkdirAll(destDir, 0777)
	if err != nil {
		log.Panic(err)
	}

	err = readDir(sourceDir, copyer)
	if err != nil {
		log.Panic(err)
	}

	err = readDir(sourceDir, setter)
	if err != nil {
		log.Panic(err)
	}

}

func Render() error {

	var err error

	err = readDir(sourceDir, setter)
	if err != nil {
		return err
	}

	return nil

}

func checkFileType(name, ft string) bool {

	s := strings.Split(name, ".")
	l := len(s)

	if l == 0 {
		return false
	}

	return s[l-1] == ft

}

func readDir(dir string, fn walker) error {

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range fileInfos {
		err = fn(dir, f)
		if err != nil {
			return err
		}

	}

	return nil
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
