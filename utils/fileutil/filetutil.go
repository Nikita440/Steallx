package fileutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

// IsFileExists checks if the file exists in the provided path
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDirExists checks if the folder exists
func IsDirExists(folder string) bool {
	info, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile reads the file from the provided path
func ReadFile(filename string) (string, error) {
	s, err := os.ReadFile(filename)
	return string(s), err
}

// CopyDir copies the directory from the source to the destination
// skip the file if you don't want to copy
func CopyDir(src, dst, skip string) error {
	s := cp.Options{Skip: func(info os.FileInfo, src, dst string) (bool, error) {
		return strings.HasSuffix(strings.ToLower(src), skip), nil
	}}
	return cp.Copy(src, dst, s)
}

// CopyFile copies the file from the source to the destination
func CopyFile(src, dst string) error {
	s, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, s, 0o600)
	if err != nil {
		return err
	}
	return nil
}

// Filename returns the filename from the provided path
func Filename(browser, dataType, ext string) string {
	replace := strings.NewReplacer(" ", "_", ".", "_", "-", "_")
	return strings.ToLower(fmt.Sprintf("%s_%s.%s", replace.Replace(browser), dataType, ext))
}

func BrowserName(browser, user string) string {
	replace := strings.NewReplacer(" ", "_", ".", "_", "-", "_", "Profile", "user")
	return strings.ToLower(fmt.Sprintf("%s_%s", replace.Replace(browser), replace.Replace(user)))
}

// ParentDir returns the parent directory of the provided path
func ParentDir(p string) string {
	return filepath.Dir(filepath.Clean(p))
}

// BaseDir returns the base directory of the provided path
func BaseDir(p string) string {
	return filepath.Base(p)
}

// ParentBaseDir returns the parent base directory of the provided path
func ParentBaseDir(p string) string {
	return BaseDir(ParentDir(p))
}

// CompressDir compresses the directory into a zip file
func CompressDir(dir string) error {
	b := new(bytes.Buffer)
	zw := zip.NewWriter(b)

	// Recursive closure function to process each file and directory
	var walkFn func(string, string) error
	walkFn = func(path string, baseInZip string) error {
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, f := range files {
			fullPath := filepath.Join(path, f.Name())
			relPath := filepath.Join(baseInZip, f.Name())

			if f.IsDir() {
				// Create a directory entry in the zip file
				_, err := zw.Create(relPath + "/")
				if err != nil {
					return err
				}
				// Recursively walk into the directory
				if err := walkFn(fullPath, relPath); err != nil {
					return err
				}
			} else {
				// Create a file entry in the zip file
				fw, err := zw.Create(relPath)
				if err != nil {
					return err
				}
				fileContent, err := os.ReadFile(fullPath)
				if err != nil {
					return err
				}
				_, err = fw.Write(fileContent)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Start walking from the root directory
	if err := walkFn(dir, ""); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	filename := filepath.Join(dir, fmt.Sprintf("%s.zip", filepath.Base(dir)))
	outFile, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = b.WriteTo(outFile)
	if err != nil {
		return err
	}

	return nil
}
