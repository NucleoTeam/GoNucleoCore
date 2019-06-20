package nucleohub


type MethodRegistration struct {
	Name string
	Run func(NucleoData);
}