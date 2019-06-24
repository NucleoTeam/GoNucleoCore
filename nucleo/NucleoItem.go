package nucleohub

type NucleoItem struct{
	Next *NucleoItem
	Data interface{}
	Chain string
}
func NewItem(chain string, data interface{}) * NucleoItem{
	q := new(NucleoItem)
	q.Data = data
	q.Chain = chain
	return q
}