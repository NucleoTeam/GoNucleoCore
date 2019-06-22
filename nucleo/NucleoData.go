package nucleohub

import (
	"github.com/google/uuid"
)
type NucleoData struct {
	Root uuid.UUID `json:"root"`
	Steps []*NucleoStep `json:"steps"`
	ChainList [][]string `json:"chainList"`
	Origin string `json:"origin"`
	Link int `json:"link"`
	Execution *NucleoStep `json:"execution"`
	OnChain int `json:"onChain"`
	Objects map[string]interface{} `json:"objects"`
	ChainBreak *NucleoChainStatus `json:"chainBreak"`
}
func NewNucleoData() *NucleoData {
	data := new(NucleoData);
	step := NewStep("")
	data.ChainBreak = NewChainStatus()
	o, _ := uuid.NewRandom()
	data.Root = o
	data.Objects = map[string]interface{}{}
	data.Execution = step
	return data
}

func (d * NucleoData) GetCurrentChain() string{
	if d.ChainList == nil {
		return ""
	}
	if d.ChainList[d.OnChain] == nil {
		return ""
	}
	chainText := ""
	for x:=0;x<=d.Link;x++ {
		if chainText != "" {
			chainText+="."+d.ChainList[d.OnChain][x]
		} else {
			chainText=d.ChainList[d.OnChain][x]
		}

	}
	return chainText
}
func (d * NucleoData) Increment() int {
	if d.Link+1 < len(d.ChainList[d.OnChain]) {
		d.Link++
		return 0
	}
	if d.OnChain+1 < len(d.ChainList) {
		d.OnChain++
		d.Link = 0
		return 1
	}
	return -1
}