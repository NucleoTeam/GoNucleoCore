package nucleohub

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type NucleoHub struct {
	Name string
	Responders  map[string] func(data *NucleoData) *NucleoData
	Response map[uuid.UUID] func(data *NucleoData)
	group string
	Pusher *ElasticSearchPusher
	Queue *NucleoList
	producer map[string]*NucleoProducer
	origin uuid.UUID
	brokers []string
	consumers map[string]*NucleoConsumer
}

func NewHub(name string, group string, brokers []string, elasticServers []string) *NucleoHub{
	hub := new(NucleoHub)
	hub.group = group
	hub.Name = name
	hub.Pusher = NewESPusher(elasticServers)
	hub.Queue = newList()
	hub.origin, _ = uuid.NewRandom()
	hub.producer = map[string]*NucleoProducer{}
	hub.Responders = map[string]func(data *NucleoData) *NucleoData{}
	hub.Response = map[uuid.UUID]func(data *NucleoData) {}
	hub.brokers = brokers
	hub.consumers = map[string]*NucleoConsumer{}
	//hub.producer["broadcast"] = NewProducer("broadcast", brokers, hub)
	hub.consumers["nucleo.client."+hub.Name] = NewConsumer("nucleo.client."+hub.Name, []string{}, hub.group, hub.brokers, hub)
	go hub.PollQueue()

	return hub
}

func (hub * NucleoHub) Add(chains string, data *NucleoData, function func(data *NucleoData)){
	data.Origin = hub.Name
	data.Link = 0
	data.OnChain = 0
	chainList := strings.Split(strings.Replace(chains, " ", "", -1), ",")
	chainListArr := [][]string{}
	for _, chain :=  range chainList {
		chainListArr = append(chainListArr, strings.Split(chain, "."))
	}
	data.ChainList = chainListArr
	data.GetCurrentChain()
	hub.Response[data.Root] = function
	hub.Queue.Add(NewItem(data.GetCurrentChain(), data))
}

func (hub *NucleoHub) Register(chain string, function func(data *NucleoData) *NucleoData ) {
	if chain == "" {
		return
	}
	chain = strings.Replace(chain, " ", "", -1)
	chainReqs := strings.Split(chain, ">")
	if len(chainReqs) == 1 {
		chain = chainReqs[0]
		chainReqs = []string{}
		//fmt.Println(chainReqs)
		//fmt.Println(chain)
	} else {
		chain = chainReqs[len(chainReqs)-1]
		chainReqs = chainReqs[0:len(chainReqs)-1]
		//fmt.Println(chainReqs)
		//fmt.Println(chain)
	}
	hub.consumers[chain] = NewConsumer(chain, chainReqs, hub.group, hub.brokers, hub)
	hub.Responders[chain] = function
}

func (hub *NucleoHub) PollQueue(){
	for {
		item := hub.Queue.Pop()
		if item != nil {
			if hub.producer[item.Chain] == nil {
				hub.producer[item.Chain] = NewProducer(item.Chain, hub.brokers, hub)
			}
			hub.producer[item.Chain].Queue.Add(item)
		}
		time.Sleep(1 * time.Microsecond)
	}
}

func containsAll(data *NucleoData, search []string) []string {
	notFound := []string{}
	for _, req := range search {
		var found = false
		for _, step := range data.Steps {
			if step.Step == req {
				found = true
			}
		}
		if found == false {
			notFound = append(notFound, req)
		}
	}
	return notFound
}

func (hub *NucleoHub) Execute(chain string, data *NucleoData, reqs []string){
	if chain == "nucleo.client."+hub.Name {
		fmt.Println(chain + "=nucleo.client."+hub.Name)
		if hub.Response[data.Root] != nil {
			data.Execution.EndStep()
			hub.Response[data.Root](data)
			hub.Pusher.Push(data)
			delete(hub.Response, data.Root)
			return
		}
		return
	} else {
		fmt.Println(chain + "=nucleo.client."+hub.Name)
	}
	if len(reqs)>0 {
		missingRequirements := containsAll(data, reqs)
		//fmt.Println(missingRequirements)
		if len(missingRequirements) > 0 {
			data.ChainBreak.BreakChain = true
			reqStr, _ := json.Marshal(missingRequirements)
			data.ChainBreak.BreakReasons = append(data.ChainBreak.BreakReasons, "Missing required chains "+ string(reqStr))
			hub.Pusher.Push(data)
			hub.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
			return
		}
	}
	if hub.Responders[chain] == nil {
		return
	}
	step := NewStep(chain);
	data.Steps = append(data.Steps, step)
 	hub.Responders[chain](data)
	step.EndStep()

	// Push to elasticsearch using the step count as the version
	hub.Pusher.Push(data)

	if data.ChainBreak!=nil {
		if data.ChainBreak.BreakChain {
			hub.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
			return
		}
	}
	response := data.Increment()
	chain = data.GetCurrentChain()
	if response == 0 {
		if hub.Responders[chain]!=nil {
			hub.consumers[chain].Exec(data)
			hub.Execute(chain, data, []string{})
			return
		}
		hub.Queue.Add(NewItem(chain, data))
	} else if response == 1 {
		hub.Queue.Add(NewItem(chain, data))
	} else if response == -1 {
		hub.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
	}
}