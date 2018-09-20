package ipfs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"rdfs/util"

	"github.com/ipfs/go-ipfs-api"
)

func Open() int {
	cmd := exec.Command("ipfs", "daemon")
	err := cmd.Start()

	if err != nil {
		fmt.Printf("[-] Failed to initialize 'IPFS' with '%s'\n", err)
		cmd.Process.Kill()
		os.Exit(1)
	}

	fmt.Printf("[+] Successfully initialized 'IPFS' with pid %d\n",
		cmd.Process.Pid)

	return cmd.Process.Pid
}

func Close(pid int) {
	if pid == -1 {
		return
	}
	syscall.Kill(pid, syscall.SIGTERM)
	fmt.Println("[+] Successfully Closed 'IPFS'")
}

func List(shell *shell.Shell, hash string, upper string) bool {
	flist, err := shell.List(hash)

	if err != nil {
		fmt.Println("[-] Error occured while checking the given hash")
		return false
	}

	for _, f := range flist {
		if f.Type == 1 {
			fmt.Printf("%s - %s/%s (%d)\n", f.Hash, upper, f.Name, f.Size)
			List(shell, f.Hash, upper+"/"+f.Name)
			continue
		}

		fmt.Printf("%s - %s/%s (%d)\n", f.Hash, upper, f.Name, f.Size)
	}

	return true
}

func Add(shell *shell.Shell, file *[]byte) (hash string, err error) {
	r := bytes.NewReader(*file)
	hash, err = shell.Add(r)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", err
	}

	return hash, nil
}

func Get(shell *shell.Shell, hash string, path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Println("[-] Could not check the file info. " +
			"Check the file path, please")
		return false
	}

	if !fileInfo.IsDir() {
		fmt.Println("[-] Wrong dir path. Fix the dir path with config file.")

		return false
	}

	err = shell.Get(hash, path)

	if err != nil {
		fmt.Println("[-] Could not get the file with the given hash. " +
			"Check the hash once again, please.")

		return false
	}

	fmt.Printf("[+] Got %s on %s\n", hash, path)

	return true
}

func PublishDefaultDir(shell *shell.Shell) bool {
	hash, err := shell.AddDir(util.RDFS_UP_DIR)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	ok := Publish(shell, hash)
	if ok == false {
		fmt.Printf("[-] Couldn't publish '%s'\n", hash)
		return false
	}

	return true
}

func Publish(shell *shell.Shell, hash string) bool {
	nodeID, err := shell.ID()
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	err = shell.Publish("/ipfs/"+hash, "/ipns/"+nodeID.ID)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully published %s(%s) to '%s'\n", hash, util.RDFS_UP_DIR, nodeID.ID)

	return true
}

func Resolve(shell *shell.Shell, nodeID string) string {
	resp, err := shell.Resolve(nodeID)
	if err != nil {
		//fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	return resp
}

/*
func Set(shell *shell.Shell, path string, hash *string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Println("[-] Could not check the file info. " +
			"Check the file path, please")
		return false
	}

	if fileInfo.IsDir() == true {
		*hash, err = shell.AddDir(path)

		if err != nil {
			fmt.Println("[-] Error occured while adding the directory.")
			return false
		}

	} else {
		file, err := os.Open(path)

		if err != nil {
			fmt.Println("[-] Could not open the file. " +
				"Check the file path, please.")
			return false
		}

		defer file.Close()

		r := bufio.NewReader(file)
		*hash, err = shell.Add(r)

		if err != nil {
			fmt.Println("[-] Error occured while adding the file")
			return false
		}
	}

	fmt.Printf("[+] Added '%s' : '%s'\n", path, *hash)

	return true
}
*/
