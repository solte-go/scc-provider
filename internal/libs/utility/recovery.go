package utility

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	src := "path/to/source"
	dst := "path/to/destination"
	excludeFile := "excluded.txt" // file you want to exclude

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error in accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() || filepath.Base(path) == excludeFile {
			return nil
		}

		// The file isn't a directory and isn't the excluded file, so move it
		relativePath, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dst, relativePath)

		// Create the directories for the destination path, if they don't exist
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		// Open source and destination file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Copy the contents to the destination file
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		// Delete the source file
		if err := os.Remove(path); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", src, err)
		return
	}
}
