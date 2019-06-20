package nucleohub

type NucleoQuery struct{
	Next *NucleoQuery
	Data *NucleoData
}