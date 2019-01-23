package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tvm-sdk/proto/tm"
	t_server "tvm-sdk/trans"
	t_conf "tvm-sdk/config"
)

func executeContract(writer http.ResponseWriter, request *http.Request) {
	//body, _ := ioutil.ReadAll(request.Body)
	//body_str := string(body)
	//fmt.Println(body_str)
	var contractRequest *tm.ExecuteContractRequest;
	if err:= json.NewDecoder(request.Body).Decode(&contractRequest); err != nil {
		fmt.Println(err)
		request.Body.Close()
	}
	server := t_server.NewMWService();
	response := server.ExecuteContract(nil, contractRequest)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/executeContract", executeContract)
	fmt.Println(t_conf.GetPort())
	err := http.ListenAndServe(":"+t_conf.GetPort(), nil)
	if err != nil {
		fmt.Println(err)
	}
}