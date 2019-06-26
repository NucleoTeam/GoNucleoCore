package nucleohub

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/estransport"
	"os"
	"strings"
	"time"
)


type ElasticSearchPusher struct {
	Hosts []string
	List NucleoList
}

func (es * ElasticSearchPusher) Push(data * NucleoData){
	data.Version++
	es.List.Add( NewItem("", *data))
}

func NewESPusher(hosts []string) *ElasticSearchPusher{
	es := new(ElasticSearchPusher)
	es.Hosts = hosts
	es.List = * newList()
	go es.process()
	return es
}

func (es * ElasticSearchPusher) process(){
	cfg := elasticsearch.Config{
		Addresses: es.Hosts,
		Logger: &estransport.ColorLogger{Output: os.Stdout},
	}
	eClient, err := elasticsearch.NewClient(cfg)
	if err !=nil {
		fmt.Print(err)
	}
	for {
		if es.List.Size > 0 {
			e := es.List.Pop()
			if e!=nil {
				item := e.Data.(NucleoData)
				dataJson, _ := json.Marshal(e.Data)
				eClient.Index(
					"nucleo",
					strings.NewReader(string(dataJson)),
					eClient.Index.WithRefresh("true"),
					eClient.Index.WithVersion(item.Version),
					eClient.Index.WithVersionType("external"),
					eClient.Index.WithPretty(),
					eClient.Index.WithDocumentID(item.Origin+"-"+item.Root.String()),
				)
			}
		}
		time.Sleep(time.Millisecond)
	}
}