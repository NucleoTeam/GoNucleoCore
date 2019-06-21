package nucleohub

import (
	"github.com/google/uuid"
	"strings"
)

type NucleoHub struct {
	Name string
	Responders  map[string] func(data *NucleoData) *NucleoData
	group string
	Queue *NucleoList
	producer map[string]*NucleoProducer
	origin uuid.UUID
	brokers []string
	consumers []*NucleoConsumer
}

func NewHub(name string, group string, brokers []string) *NucleoHub{
	hub := new(NucleoHub)
	hub.group = group
	hub.Name = name
	hub.Queue = newList()
	hub.origin, _ = uuid.NewRandom()
	hub.producer = map[string]*NucleoProducer{}
	hub.Responders = map[string]func(data *NucleoData) *NucleoData{}
	hub.brokers = brokers
	hub.producer["broadcast"] = newProducer("broadcast", brokers, hub)
	hub.consumers = append(hub.consumers, NewConsumer("nucleo.client."+hub.Name, hub.group, hub.brokers, hub))
	go hub.PollQueue()

	return hub
}

func (hub * NucleoHub) Add(chains string, data *NucleoData){
	//
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
	hub.Queue.Add(NewItem("", data))
}

func (hub *NucleoHub) Register(chain string, function func(data *NucleoData) *NucleoData ) {
	hub.consumers = append(hub.consumers, NewConsumer(chain, hub.group, hub.brokers, hub))
	hub.Responders[chain] = function
}

func (hub *NucleoHub) PollQueue(){
	for {
		item := hub.Queue.Pop()
		chain := ""
		if item != nil {
			chain = item.Data.GetCurrentChain()
			if hub.producer[chain] == nil {
				hub.producer[chain] = newProducer(chain, hub.brokers, hub)
			}
			hub.producer[chain].Queue.Add(item)
		}
	}
}

func (hub *NucleoHub) Execute(data *NucleoData){

}
