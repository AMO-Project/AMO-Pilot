package util

import "strings"

const (
	CMD_PURCHASE = iota
	CMD_STORE
	CMD_TOKENBAL
	CMD_LIST
	CMD_PEERFILE
	CMD_TEST
	CMD_HELP
	CMD_EXIT
	CMD_NETVERSION
	CMD_COINBASE
	CMD_ISMINING
	CMD_BLOCKNUMBER
	CMD_ETHBALANCE
	CMD_ACCOUNTS
	CMD_SENDTX
	CMD_UNLOCK
)

var CMD = map[int]string{
	CMD_PURCHASE: "purchase <hash>" +
		":purchase an item with hash, saved in " + RDFS_DOWN_DIR,
	CMD_STORE: "store <path>" +
		":store an item with file(dir) path, copied into " + RDFS_UP_DIR,
	CMD_TOKENBAL: "tokenbal <address>" +
		":return the token balance of the given <address>",
	CMD_LIST: "ls <option> <hash>" +
		":show all of the files stored within the given hash recursively. " +
		"(option) -f for local files / -h for hash files",
	CMD_PEERFILE: "peerfile" +
		":show all of files belonging to client's peers.",
	CMD_TEST: "test" +
		":test codes for debugging",
	CMD_HELP: "help" +
		":list available commands",
	CMD_EXIT: "exit" +
		":close RDFS",
	CMD_NETVERSION: "netversion" +
		":show current geth netversion" +
		"(1 Mainnet, 2 Morde, 3 Ropsten, 208518 rdfs)",
	CMD_COINBASE: "coinbase" +
		":show current the client coinbase address",
	CMD_ISMINING: "ismining" +
		":return true if client is actively mining.",
	CMD_BLOCKNUMBER: "blocknumber" +
		":return the number of most recent block.",
	CMD_ETHBALANCE: "ethbal <address> " +
		":return the ethereum balance of the given <address> in wei. " +
		"Based on latest block.",
	CMD_ACCOUNTS: "accounts" +
		":return a list of addresses owned by client.",
	CMD_SENDTX: "sendtx" +
		":send transaction.(unsupported yet)",
	CMD_UNLOCK: "unlock <passphrase> <time>" +
		":unlock current client's account",
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
		} else if strings.Compare(args[0], "peerfile") == 0 {
			return CMD_PEERFILE, nil
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
		} else if strings.Compare(args[0], "sendtx") == 0 {
			return CMD_SENDTX, nil
		}
	case 2:
		if strings.Compare(args[0], "store") == 0 {
			return CMD_STORE, args[1:]
		} else if strings.Compare(args[0], "purchase") == 0 {
			return CMD_PURCHASE, args[1:]
		} else if strings.Compare(args[0], "tokenbal") == 0 {
			return CMD_TOKENBAL, args[1:]
		} else if strings.Compare(args[0], "ls") == 0 {
			return CMD_LIST, args[1:]
		} else if strings.Compare(args[0], "ethbal") == 0 {
			return CMD_ETHBALANCE, args[1:]
		} else if strings.Compare(args[0], "test") == 0 {
			return CMD_TEST, args[1:]
		}
	case 3:
		if strings.Compare(args[0], "ls") == 0 {
			return CMD_LIST, args[1:]
		} else if strings.Compare(args[0], "purchase") == 0 {
			return CMD_PURCHASE, args[1:]
		}
	case 4:
		if strings.Compare(args[0], "unlock") == 0 {
			return CMD_UNLOCK, args[1:]
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
