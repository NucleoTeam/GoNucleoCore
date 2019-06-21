package main

import (
	"fmt"
	nucleohub "nucleoCore/nucleo"
	"time"
)

func main() {
	//test := nucleohub.NewNucleoData()
	//popcorn := "{\"root\":\"00000000-0000-0000-0000-000000000000\",\"steps\":null,\"chainList\":null,\"origin\":\"\",\"link\":0,\"execution\":{\"step\":\"\",\"start\":1561065550956201000,\"host\":\"Nathaniels-Mac-mini.local\",\"end\":1561065550956203000,\"total\":2000},\"onChain\":0,\"objects\":{\"test\":1},\"chainBreak\":null}"
	//res := nucleohub.NucleoData{}
	hub := nucleohub.NewHub("Go11","nucleoCore-Go", []string{"192.168.1.112:9092"})
	hub.Register("new", func(data * nucleohub.NucleoData) *nucleohub.NucleoData{
		fmt.Println("test")
		data.Objects["test"] = ""
		return data;
	})
	hub.Add("new.taco.bell,taco.one.two", nucleohub.NewNucleoData())
	for{
		time.Sleep(4 * time.Second)
		hub.Add("new.taco.bell,taco.one.two", nucleohub.NewNucleoData())
	}
	//json.Unmarshal([]byte(popcorn), &res)
	//nucleohub.NewStep("player.get.id");
	//test.Execution.EndStep
	/*b, err :=json.Marshal(res)
	if err != nil {
		fmt.Print("error")
	}
 	fmt.Println(string(b));*/
}
