package nucleohub

import (
	"encoding/json"
	"github.com/google/uuid"
	"strings"
)

type NucleoHub struct {
	Name string
	Responders  map[string]*NucleoResponder
	Response map[uuid.UUID] func(data *NucleoData)
	group string
	Pusher *ElasticSearchPusher
	producer *NucleoProducer
	origin uuid.UUID
	brokers []string
	consumers *NucleoConsumer
}

func NewHub(name string, group string, brokers []string, elasticServers []string) *NucleoHub{
	hub := new(NucleoHub)
	hub.group = group
	hub.Name = name
	hub.Pusher = NewESPusher(elasticServers)
	hub.origin, _ = uuid.NewRandom()
	hub.producer = NewProducer(brokers, hub)
	hub.Responders = map[string] * NucleoResponder{}
	hub.Response = map[uuid.UUID]func(data *NucleoData) {}
	hub.brokers = brokers
	hub.consumers = NewConsumer(hub.group, hub.brokers, hub)
	hub.consumers.AddTopic("nucleo.client."+hub.Name)
	return hub
}

func (hub * NucleoHub) Start(){
	hub.consumers.Start()
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
	hub.producer.Queue.Add(NewItem(data.GetCurrentChain(), data))
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

	hub.consumers.AddTopic(chain)
	hub.Responders[chain] = NewResponder(function, chainReqs)
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

func (hub *NucleoHub) Execute(chain string, data *NucleoData){
	if chain == "nucleo.client."+hub.Name {
		if hub.Response[data.Root] != nil {
			data.Execution.EndStep()
			hub.Response[data.Root](data)
			hub.Pusher.Push(data)
			delete(hub.Response, data.Root)
			return
		}
		return
	}
	if hub.Responders[chain] == nil {
		return
	}
	responder := hub.Responders[chain]
	if len(responder.Requirements)>0 {
		missingRequirements := containsAll(data, responder.Requirements)
		//fmt.Println(missingRequirements)
		if len(missingRequirements) > 0 {
			data.ChainBreak.BreakChain = true
			reqStr, _ := json.Marshal(missingRequirements)
			data.ChainBreak.BreakReasons = append(data.ChainBreak.BreakReasons, "Missing required chains "+ string(reqStr))
			hub.Pusher.Push(data)
			hub.producer.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
			return
		}
	}
	step := NewStep(chain);
	data.Steps = append(data.Steps, step)
 	hub.Responders[chain].Function(data)
	step.EndStep()

	// Push to elasticsearch using the step count as the version
	hub.Pusher.Push(data)

	if data.ChainBreak!=nil {
		if data.ChainBreak.BreakChain {
			hub.producer.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
			return
		}
	}
	response := data.Increment()
	chain = data.GetCurrentChain()
	if response == 0 {
		if responder!=nil {
			hub.Execute(chain, data)
			return
		}
		hub.producer.Queue.Add(NewItem(chain, data))
	} else if response == 1 {
		hub.producer.Queue.Add(NewItem(chain, data))
	} else if response == -1 {
		hub.producer.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
	}
}