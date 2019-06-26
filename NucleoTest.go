package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	nucleohub "nucleoCore/nucleo"
	"time"
)

func main() {
	nameUnique, _ := uuid.NewRandom()
	hub := nucleohub.NewHub(
		"Go-Client-"+nameUnique.String(),
		"nucleoCore-Go",
		[]string{"192.168.1.112:9092"},
		[]string{"http://192.168.1.112:9200"},
	)
	hub.Register("pop", func(data * nucleohub.NucleoData) *nucleohub.NucleoData{
		data.Objects["test"] = "wow"
		return data;
	})
	hub.Register("taco > pop.corn", func(data * nucleohub.NucleoData) *nucleohub.NucleoData{
		data.Objects["Corn"] = "GOGOGO"
		return data;
	})

	hub.Start()

	for{
		fmt.Println("Created new request")
		hub.Add("pop.corn", nucleohub.NewNucleoData(), func(data *nucleohub.NucleoData) {
			d, _ := json.Marshal(data)
			fmt.Println(string(d))
		})
		time.Sleep(2 * time.Second)
	}
}
