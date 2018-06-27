package main

import (
	"encoding/json"
	"fmt"
	"mvs_api"
	"strings"
)

func main() {
	r := mvs_api.NewRPCClient("http://127.0.0.1:8820/rpc/v2", "1s")
	rpcResp, err := r.Getdid("BIAM")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Getdid("")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Setminingaccount("Alice", "A123456", "MLasJFxZQnA49XEvhTHmRKi2qstkj9ppjo")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply string
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Getwork("", "")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply []string
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Getnewmultisig("Alice", "A123456", 2, 3, "0380990a7312b87abda80e5857ee6ebf798a2bf62041b07111287d19926c429d11",
		[]string{"02578ad340083e85c739f379bbe6c6937c5da2ced52e09ac1eec43dc4c64846573", "03af3a99f1c3279dbe1c22fd767fb98b5dbd138f6e0511c2fc11128e44c0373cad"},
		"test mvs api")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply["address"])
	}

	rpcResp, err = r.Listmultisig("Alice", "A123456")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Deletemultisig("Alice", "A123456", "359mjCL3V8PaxLUzU9mJSNtLSEXHFJmzfA")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Getaccountasset("Alice", "A123456", "", true)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Getaccountasset("Alice", "A123456", "", false)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Listtxs("Alice", "A123456", "MLasJFxZQnA49XEvhTHmRKi2qstkj9ppjo", [2]uint64{1000, 1001}, "", 0, 0)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Importaccount(strings.Split("notice judge certain company novel quality plunge list blind library ride uncover fold wink biology original aim whale stand coach hire clinic fame robot", " "), "", "robot", "robot123456", 10)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Deleteaccount("robot", "robot123456", "robot")
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Dumpkeyfile("Alice", "A123456", "robot", "", true)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply map[string]interface{}
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}

	rpcResp, err = r.Dumpkeyfile("Alice", "A123456", "robot", "", false)
	if err != nil {
		fmt.Println(err)
	} else {
		var reply string
		json.Unmarshal(*rpcResp.Result, &reply)
		fmt.Println(reply)
	}
}
