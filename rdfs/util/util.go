package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	RDFS_VER        string = "0.1"
	RDFS_HOME_DIR   string = "/pss/rdfs/"
	RDFS_CONFIG_DIR string = RDFS_HOME_DIR + ".config/"
	RDFS_UP_DIR     string = RDFS_HOME_DIR + "up/"
	RDFS_DOWN_DIR   string = RDFS_HOME_DIR + "down/"

	GETH_DATA_DIR string = "/pss/geth_data/"
)

const (
	CMD_STORE = iota
	CMD_PURCHASE
	CMD_LIST
	CMD_TEST
	CMD_HELP
	CMD_EXIT
	CMD_NETVERSION
	CMD_COINBASE
	CMD_ISMINING
	CMD_BLOCKNUMBER
	CMD_BALANCE
	CMD_ACCOUNTS
	CMD_GCTEST
)

var CMD = map[int]string{
	CMD_EXIT: "exit" +
		":close RDFS",
	CMD_PURCHASE: "purchase <hash> (dir path)" +
		":purchase an item with hash. (option) dir path",
	CMD_STORE: "store <path>" +
		":store an item with file(dir) path",
	CMD_LIST: "ls <option> <hash>" +
		":show all of the files stored within the given hash recursively. " +
		"(option) -f for local files / -h for hash files",
	CMD_TEST: "test" +
		":test codes for debugging",
	CMD_NETVERSION: "netversion" +
		":show current geth netversion" +
		"(1 Mainnet, 2 Morde, 3 Ropsten, 208518 rdfs)",
	CMD_COINBASE: "coinbase" +
		":show current the client coinbase address",
	CMD_ISMINING: "ismining" +
		":return true if client is actively mining.",
	CMD_BALANCE: "balance <address> " +
		":return the balance of the account of given <address> in wei. " +
		"Based on latest block.",
	CMD_ACCOUNTS: "accounts" +
		":return a list of addresses owned by client.",
	CMD_BLOCKNUMBER: "blocknumber" +
		":return the number of most recent block.",
	CMD_HELP: "help" +
		":list available commands",
}

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

	for _, fi := range fis {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}
		fsize := du(path+"/"+fi.Name(), fi)
		fmt.Printf("%s (%d bytes)\n", path+"/"+fi.Name(), fsize)
		size += fsize
	}

	//fmt.Printf("%s (%d)\n", path, size/10)

	return size

}

func List(path string) bool {
	info, err := os.Lstat(path)

	if err != nil {
		fmt.Printf("[-] Could not walk the given path.\n")
		return false
	}

	du(path, info)

	return true
}

func Shell(cmdln string) (int, []string) {
	args := strings.Fields(cmdln)
	l := len(args)

	switch l {
	case 0:
		return -1, nil
	case 1:
		if strings.Compare(args[0], "exit") == 0 {
			return CMD_EXIT, nil
		} else if strings.Compare(args[0], "help") == 0 {
			return CMD_HELP, nil
		} else if strings.Compare(args[0], "netversion") == 0 {
			return CMD_NETVERSION, nil
		} else if strings.Compare(args[0], "coinbase") == 0 {
			return CMD_COINBASE, nil
		} else if strings.Compare(args[0], "ismining") == 0 {
			return CMD_ISMINING, nil
		} else if strings.Compare(args[0], "blocknumber") == 0 {
			return CMD_BLOCKNUMBER, nil
		} else if strings.Compare(args[0], "accounts") == 0 {
			return CMD_ACCOUNTS, nil
		}
	case 2:
		if strings.Compare(args[0], "store") == 0 {
			return CMD_STORE, args[1:]
		} else if strings.Compare(args[0], "purchase") == 0 {
			return CMD_PURCHASE, args[1:]
		} else if strings.Compare(args[0], "balance") == 0 {
			return CMD_BALANCE, args[1:]
		} else if strings.Compare(args[0], "test") == 0 {
			return CMD_TEST, args[1:]
		}
	case 3:
		if strings.Compare(args[0], "purchase") == 0 {
			return CMD_PURCHASE, args[1:]
		} else if strings.Compare(args[0], "ls") == 0 {
			return CMD_LIST, args[1:]
		}
	case 5:
		if strings.Compare(args[0], "test") == 0 {
			return CMD_TEST, args[1:]
		}
	}

	return -1, nil
}

/*
// Exec the given command, return pid
func BashS(cmdln string) int {
	args := strings.Fields(cmdln)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	//out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("[-] Something went wrong with '%s'\n", err)
		return -1
	}

	fmt.Printf("[+] Done\n%d\n", cmd.Process.Pid)

	return cmd.Process.Pid
}

func BashL(cmdln string) int {
	args := strings.Fields(cmdln)
	cmd := exec.Command(args[0], args[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()

	if err != nil {
		fmt.Printf("[-] Something went wrong with '%s'\n", err)
		return -1
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("[-] Something went wrong with '%s'\n", err)
		return -1
	}

	if errStdout != nil || errStderr != nil {
		fmt.Println("[-] Failed to capture stdout or stderr ")
		return -1
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("%s\n%s", outStr, errStr)

	return cmd.Process.Pid
}
*/
