package nucleohub



type NucleoHub struct {
	Methods  map[string]MethodRegistration
	Query []*NucleoQuery
	Name string

}