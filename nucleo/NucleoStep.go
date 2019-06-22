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
	Total float64 `json:"total"`
}
func NewStep(chain string) *NucleoStep {
	step := new(NucleoStep)
	step.Step = chain
	name, _ := os.Hostname()
	step.Host = name
	step.Start = time.Now().UnixNano()
	return step
}
func (s *NucleoStep) EndStep() {
	s.End = time.Now().UnixNano()
	s.Total = float64(float64( s.End - s.Start ) / float64(time.Millisecond))
}