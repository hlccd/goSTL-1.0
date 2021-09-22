package radix

type node struct {
	name  string
	num   int
	value interface{}
	son   []*node
}

func newNode(name string, e interface{}) (n *node) {
	return &node{
		name:  name,
		num:   0,
		value: e,
		son:   make([]*node, 0, 0),
	}
}
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es = append(es, s+n.name)
	}
	for i := 0; i < len(n.son); i++ {
		es = append(es, n.son[i].inOrder(s+n.name+"/")...)
	}
	return es
}
