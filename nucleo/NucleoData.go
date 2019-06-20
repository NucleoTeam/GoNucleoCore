package nucleohub
import (
	"github.com/google/uuid"
)
type NucleoData struct {
	Root uuid.UUID `json:"root"`
	Steps []*NucleoStep `json:"steps"`
	ChainList []string `json:"chainList"`
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
	data.Execution = step
	return data;
}