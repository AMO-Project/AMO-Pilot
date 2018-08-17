package util

import (
	"fmt"
	"io"
	"os"
)

func du(path string, info os.FileInfo) int64 {
	size := info.Size()

	if !info.IsDir() {
		return size
	}

	dir, err := os.Open(path)
	if err != nil {
		return size
	}

	defer dir.Close()

	fis, err := dir.Readdir(-1)

	for i, fi := range fis {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}
		fsize := du(path+fi.Name(), fi)
		fmt.Printf("[%d] %s (%d bytes)\n", i+1, fi.Name(), fsize)
		size += fsize
	}

	return size
}

func List(path string) bool {
	info, err := os.Lstat(path)

	if err != nil {
		fmt.Printf("[-] Could not walk the given path.\n")
		return false
	}

	fmt.Printf("[+] Files located in %s (RDFS_UP_DIR):\n", path)
	du(path, info)

	return true
}

func Copy(from string, to string) bool {
	from_file, err := os.Open(from)
	if err != nil {
		return false
	}
	defer from_file.Close()

	to_file, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	defer to_file.Close()

	_, err = io.Copy(to_file, from_file)
	if err != nil {
		return false
	}

	return true
}
