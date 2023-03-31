package fsops

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/Krauze67/flib/stringtools"
)

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) error {
	return copyFileInt(src, dst, false)
}

func copyFileInt(src, dst string, deletesrc bool) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return err
		}
	}
	// if err = os.Link(src, dst); err == nil {
	// 	return err
	// }
	//TODO: in case "delete original" == true and src and dest are on the same volume - use os.Rename()
	res := copyFileContents(src, dst)
	if deletesrc {
		os.Remove(src)
	}
	return res
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}

// CopyFolder copies a folder's content to a (new one) another one
func CopyFolder(source string, dest string) error {
	return copyFolderInt(source, dest, -1, false)
}

// CopyFolderLevels copies a folder's content to a (new one) another one
func CopyFolderLevels(source string, dest string, levels int) error {
	return copyFolderInt(source, dest, levels, false)
}

func copyFolderInt(source string, dest string, levels int, deletesource bool) error {

	srcstat, err := os.Stat(source)
	if err != nil {
		return err
	}
	srcName := srcstat.Name()

	// build dest:
	fullDest := filepath.Join(dest, srcName)
	if _, err = os.Stat(fullDest); os.IsNotExist(err) {
		err = os.MkdirAll(fullDest, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}

	copyFolderContentInt(source, fullDest, levels, deletesource)

	if deletesource {
		os.RemoveAll(source)
	}

	return err
}

// CopyFolderContent copies a folder's content only
func CopyFolderContent(source string, dest string) error {
	return copyFolderContentInt(source, dest, -1, false)
}

// CopyFolderContentLevels copies a folder's content only
func CopyFolderContentLevels(source string, dest string, levelsDeeper int) error {
	return copyFolderContentInt(source, dest, levelsDeeper, false)
}

func copyFolderContentInt(source string, dest string, levelsDeeper int, deletesource bool) error {

	// check dirs:
	var err error
	if _, err = os.Stat(source); os.IsNotExist(err) {
		return err
	}
	if _, err = os.Stat(dest); os.IsNotExist(err) {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			if 0 > levelsDeeper { // recursively all
				err = copyFolderInt(sourcefilepointer, destinationfilepointer, levelsDeeper, deletesource)
			} else if 1 <= levelsDeeper { // some levels only - continue
				err = copyFolderInt(sourcefilepointer, destinationfilepointer, levelsDeeper-1, deletesource)
			} else {
				// 0 == levelsDeeper - means "do not copy next levels"
				// Do nothing
			}
			if deletesource {
				os.RemoveAll(source)
			}
		} else {
			err = copyFileInt(sourcefilepointer, destinationfilepointer, deletesource)
		}

		if err != nil {
			return err
		}

	}
	return nil
}

func IsDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func IsFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

// MoveFolder moves a folder's content only
func MoveFolder(source string, dest string) error {
	//Check is Dir:
	if !IsDir(source) {
		return fmt.Errorf("%s source is not a directory or inaccessible", source)
	}

	if filepath.VolumeName(source) == filepath.VolumeName(dest) { //on the same disk // TODO switch those if and else.
		if strings.HasPrefix(dest, source) {
			return fmt.Errorf("directory %s cannot be moved to a nested directory %s", source, dest)
		}

		if !IsDir(dest) {
			errmake := os.MkdirAll(dest, os.ModeDir|os.ModePerm)
			if errmake != nil {
				return errmake
			}
		}

		srcstat, err := os.Stat(source)
		if err != nil {
			return err
		}
		srcName := srcstat.Name()
		// build dest:
		fullDest := filepath.Join(dest, srcName)
		return os.Rename(source, fullDest)

	}

	// Different disks:
	return copyFolderInt(source, dest, -1, true)
}

// MoveFile moves a folder's content only
func MoveFile(source string, dest string) error {
	//Check is File:
	if !IsFile(source) {
		return fmt.Errorf("%s source is not a file or inaccessible", source)
	}

	if filepath.VolumeName(source) == filepath.VolumeName(dest) { //on the same disk // TODO switch those if and else.
		if strings.HasPrefix(dest, source) {
			return fmt.Errorf("file %s cannot be moved to directory %s - destination directory can not be created", source, dest)
		}
		return os.Rename(source, dest)
	}

	// Different disks:
	return copyFileInt(source, dest, true)
}

// GetFilename returns a filename
func GetFilename(fullfilename string, withExt bool) string {
	dir := filepath.Dir(fullfilename)
	filenameWithExt, err := stringtools.GetSubstring(fullfilename, dir, "")
	if nil != err {
		return ""
	}
	if withExt || !strings.Contains(filenameWithExt, ".") {
		return filenameWithExt[1 : len(filenameWithExt)-1]
	}

	filenameWithoutExt := filenameWithExt[1:strings.LastIndex(filenameWithExt, ".")]
	return filenameWithoutExt
}

func IsFileContains(fpath string, test []byte) bool {
	in, err := os.Open(fpath)
	if err != nil {
		return false
	}
	in.Close()

	if 0 == len(test) {
		return false
	}

	fcontent, err := ioutil.ReadFile(fpath)
	if nil != err || nil == fcontent {
		return false
	}

	bfound := false
	for i := 0; i < len(fcontent) && !bfound; i++ {
		if test[0] == fcontent[i] {
			testmatches := true
			for j := 1; j < len(test) && testmatches; j++ {
				if test[j] != fcontent[i+j] {
					testmatches = false
				}
			}
			if testmatches {
				return true
			}
		} // if first byte of test slice found
	} // loop by whole file content
	return false
}
