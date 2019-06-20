package nucleohub

type NucleoList struct {
	Head *NucleoQuery
	Tail *NucleoQuery
	Size int
}

func (list *NucleoList ) Add(query *NucleoQuery){
	if list.Head==nil {
		list.Head = query
		list.Tail = query
		return
	}
	list.Tail.Next = query;
	list.Tail = query
}
func (list *NucleoList) Pop() *NucleoQuery{
	if list.Head==nil{
		return nil
	}
	query := list.Head
	list.Head = list.Head.Next
	if list.Head==nil {
		list.Tail=nil
	}
	return query
}