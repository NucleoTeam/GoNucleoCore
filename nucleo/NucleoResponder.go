package nucleohub


type NucleoResponder struct {
	Function func(data *NucleoData) *NucleoData
	Requirements []string
}
func NewResponder(function func(data *NucleoData) *NucleoData, Reqs []string) * NucleoResponder {
	nr := new(NucleoResponder)
	nr.Function = function
	nr.Requirements = Reqs
	return nr;
}