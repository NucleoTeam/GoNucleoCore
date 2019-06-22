package main

import (
	"encoding/json"
	"fmt"
	nucleohub "nucleoCore/nucleo"
	"time"
)

func main() {
	hub := nucleohub.NewHub("Go-Client","nucleoCore-Go", []string{"192.169.1.1:9092"})
	hub.Register("pop", func(data * nucleohub.NucleoData) *nucleohub.NucleoData{
		data.Objects["test"] = "wow"
		return data;
	})
	hub.Register("pop.corn", func(data * nucleohub.NucleoData) *nucleohub.NucleoData{
		data.Objects["Corn"] = "GOGOGO"
		return data;
	})
	for{
		fmt.Println("Created new request")
		hub.Add("pop.corn", nucleohub.NewNucleoData(), func(data *nucleohub.NucleoData) {
			d, _ := json.Marshal(data)
			fmt.Println(string(d))
		})
		time.Sleep(10 * time.Second)
	}
}
