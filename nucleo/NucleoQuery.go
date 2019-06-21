package nucleohub

type NucleoItem struct{
	Next *NucleoItem
	Data *NucleoData
	Chain string
}
func NewItem(chain string, data * NucleoData) * NucleoItem{
	q := new(NucleoItem)
	q.Data = data
	q.Chain = chain
	return q
}