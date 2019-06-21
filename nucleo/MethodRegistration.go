package nucleohub


type MethodRegistration struct {
	Chain string
	Function *func(*NucleoData);
}

func NewMethod(chain string, function *func(*NucleoData)) *MethodRegistration {
	reg := new(MethodRegistration)
	reg.Chain = chain
	reg.Function = function
	return reg
}