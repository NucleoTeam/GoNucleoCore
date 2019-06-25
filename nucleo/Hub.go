package nucleohub

import (
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
	consumers []*NucleoConsumer
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
	//hub.producer["broadcast"] = NewProducer("broadcast", brokers, hub)
	hub.consumers = append(hub.consumers, NewConsumer("nucleo.client."+hub.Name, hub.group, hub.brokers, hub))
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
	hub.consumers = append(hub.consumers, NewConsumer(chain, hub.group, hub.brokers, hub))
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

func (hub *NucleoHub) Execute(chain string, data *NucleoData){
	if chain == "nucleo.client."+hub.Name {
		if hub.Response[data.Root] != nil {
			data.Execution.EndStep()
			hub.Response[data.Root](data)
			delete(hub.Response, data.Root)
			return
		}
		return
	}
	if hub.Responders[chain] == nil {
		return
	}
	step := NewStep(chain);
	data.Steps = append(data.Steps, step)
 	hub.Responders[chain](data)
	step.EndStep()

	// Push to elasticsearch using the step count as the version
	hub.Pusher.Push(*data)

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
			hub.Execute(chain, data)
			return
		}
		hub.Queue.Add(NewItem(chain, data))
	} else if response == 1 {
		hub.Queue.Add(NewItem(chain, data))
	} else if response == -1 {
		hub.Queue.Add(NewItem("nucleo.client."+data.Origin, data))
	}
}