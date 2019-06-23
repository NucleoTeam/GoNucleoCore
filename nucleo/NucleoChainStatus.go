package nucleohub


type NucleoChainStatus struct {
	BreakChain bool `json:"breakChain"`
	BreakReasons []string `json:"breakReasons"`
}
func NewChainStatus() * NucleoChainStatus{
	c := new(NucleoChainStatus)
	c.BreakReasons = []string{}
	return c
}