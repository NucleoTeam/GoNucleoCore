package nucleohub

import (
	"os"
	"time"
)

type NucleoStep struct {
	Step string `json:"step"`
	Start int64 `json:"start"`
	Host string `json:"host"`
	End int64 `json:"end"`
	Total int64 `json:"total"`
}
func NewStep(chain string) *NucleoStep {
	step := new(NucleoStep)
	step.Step = chain
	name, _ := os.Hostname()
	step.Host = name
	step.Start = time.Now().Unix()
	return step
}
func (s *NucleoStep) EndStep() {
	s.End = time.Now().Unix()
	s.Total = s.End - s.Start
}