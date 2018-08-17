package util

type Config struct {
	IPFS struct {
		Enode string `json:"enode"`
		Port  int    `json:"port"`
	} `json:"ipfs"`
}

/*
  Implement config json file control, if needed near in the future.
*/
