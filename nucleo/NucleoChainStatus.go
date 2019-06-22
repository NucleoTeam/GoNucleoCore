package nucleohub


type NucleoChainStatus struct {
	BreakChain bool
	BreakReasons []string
}
func NewChainStatus() * NucleoChainStatus{
	c := new(NucleoChainStatus)
	c.BreakReasons = []string{}
	return c
}