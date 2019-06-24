package nucleohub

type NucleoList struct {
	Head *NucleoItem
	Tail *NucleoItem
	Size int
}
func newList() * NucleoList{
	query := new(NucleoList)
	query.Size = 0
	return  query
}

func (list *NucleoList ) Add(item *NucleoItem){
	if list.Head==nil {
		list.Head = item
		list.Tail = item
		list.Size++
		return
	}
	list.Tail.Next = item;
	list.Tail = item
	list.Size++
}
func (list *NucleoList) Pop() *NucleoItem{
	if list.Head==nil{
		return nil
	}
	query := list.Head
	list.Head = list.Head.Next
	if list.Head==nil {
		list.Tail=nil
	}
	list.Size--
	return query
}