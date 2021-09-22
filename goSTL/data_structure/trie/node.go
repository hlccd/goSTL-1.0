package trie


type node struct {
	num   int
	son   [26]*node
	value interface{}
}

func newNode(e interface{}) (n *node) {
	return &node{
		num:   0,
		value: e,
	}
}
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es=append(es,s)
	}
	for i := 0; i < 26; i++ {
		if n.son[i] != nil {
			es = append(es, n.son[i].inOrder(s + string(i+'a'))...)
		}
	}
	//if n.value != nil {
	//	fmt.Println(s)
	//}else if b {
	//	fmt.Println(s)
	//}

	//if n.right != nil {
	//	es = append(es, n.right.inOrder()...)
	//}
	return es
}
