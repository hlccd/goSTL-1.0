package rbTree

//@Title		rbTree
//@Description
//		红黑树的节点
//		可通过节点实现二叉搜索树的添加删除
//		也可通过节点返回整个二叉搜索树的所有元素
//		同时可实现增删节点后的平衡调整
//		仅平衡黑色节点即可,同时不允许两个相邻的红色节点存在
//@author     	hlccd		2021-07-13
import "github.com/hlccd/goSTL/utils/comparator"

//node树节点结构体
//该节点是红黑树的树节点
//若该红黑树允许重复则对节点num+1即可,否则对value进行覆盖
//红黑树通过旋转进行调整
type node struct {
	value  interface{} //节点承载的元素
	num    int         //承载的元素数量
	parent *node       //父节点指针
	left   *node       //左节点指针
	right  *node       //右节点指针
	color  bool        //颜色
}

//树节点颜色
const (
	RED   bool = false //树节点颜色,红色
	BLACK bool = true  //树节点颜色,黑色
)

//@title    newNode
//@description
//		新建一个红黑树节点并返回
//		传入的父节点作为新建节点的父节点
//		将传入的元素e作为该节点的承载元素
//		该节点的num默认为1,左右子节点设为nil
//		颜色默认为红色,若默认黑色则无需调整,会引起不平衡
//@auth      	hlccd		2021-07-23
//@receiver		nil
//@param    	parent		*node					新建节点的父节点
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的红黑树节点的指针
func newNode(parent *node, e interface{}) (n *node) {
	return &node{
		value:  e,
		num:    1,
		parent: parent,
		left:   nil,
		right:  nil,
		color:  RED,
	}
}

//@title    inOrder
//@description
//		以node红黑树节点做接收者
//		以中缀序列返回节点集合
//		若允许重复存储则对于重复元素进行多次放入
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	es        	[]interface{}			以该节点为起点的中缀序列
func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := 0; i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}

//@title    getParent
//@description
//		以node红黑树节点做接收者
//		返回该节点的父节点
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m        	*node					该节点的父节点
func (n *node) getParent() (m *node) {
	if n == nil {
		return nil
	}
	return n.parent
}

//@title    getGrandParent
//@description
//		以node红黑树节点做接收者
//		返回该节点的祖父节点
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m        	*node					该节点的祖父节点
func (n *node) getGrandParent() (m *node) {
	parent := n.getParent()
	if parent == nil {
		return nil
	}

	return parent.getParent()
}

//@title    getUncle
//@description
//		以node红黑树节点做接收者
//		返回该节点的叔叔节点,即父节点的兄弟节点
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m        	*node					该节点的叔叔节点
func (n *node) getUncle() (m *node) {
	parent := n.getParent()
	if parent == nil {
		return nil
	}
	grandParent := n.getGrandParent()
	if grandParent == nil {
		return nil
	}
	if grandParent.left == parent {
		return grandParent.right
	} else if grandParent.right == parent {
		return grandParent.left
	}
	return nil
}

//@title    leftRotate
//@description
//		以node红黑树树节点做接收者
//		将该节点向左节点方向转动,使右节点作为原来节点
//		同时将右节点的左节点设为原节点的右节点
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) leftRotate() {
	if n == nil {
		//节点不存在,旋转失败
		return
	}
	if n.right == nil {
		//右子树不存在,旋转失败
		return
	}
	m := &node{
		value:  n.value,
		num:    n.num,
		parent: n,
		left:   n.left,
		right:  n.right.left,
		color:  n.color,
	}
	n.left = m
	if m.left != nil {
		m.left.parent = m
	}
	if m.right != nil {
		m.right.parent = m
	}
	//节点替换
	n.value = n.right.value
	n.num = n.right.num
	n.color = n.right.color
	n.right = n.right.right
	if n.right != nil {
		n.right.parent = n
	}
}

//@title    rightRotate
//@description
//		以node红黑树树节点做接收者
//		将该节点向右节点方向转动,使左节点作为原来节点
//		同时将左节点的右节点设为原节点的左节点
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) rightRotate() {
	if n == nil {
		//节点不存在,旋转失败
		return
	}
	if n.left == nil {
		//左子树不存在,旋转失败
		return
	}
	m := &node{
		value:  n.value,
		num:    n.num,
		parent: n,
		left:   n.left.right,
		right:  n.right,
		color:  n.color,
	}
	n.right = m
	if m.left != nil {
		m.left.parent = m
	}
	if m.right != nil {
		m.right.parent = m
	}
	//节点替换
	n.value = n.left.value
	n.num = n.left.num
	n.color = n.left.color
	n.left = n.left.left
	if n.left != nil {
		n.left.parent = n
	}
}

//@title    insert
//@description
//		以node红黑树树节点做接收者
//		从n节点中插入元素e
//		如果n节点中承载元素与e不同则根据大小从左右子树插入该元素
//		如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖
//		插入成功返回true,插入失败或不允许重复插入返回false
//		插入成功后对该节点即祖辈节点进行调整
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待插入元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	b        	bool					是否插入成功?
func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		//节点不存在,插入失败
		return false
	}
	if cmp(n.value, e) > 0 {
		//待插入元素应插入左子树
		if n.left == nil {
			//左子树为空,直接插入
			n.left = newNode(n, e)
			n = n.left
			b = true
		} else {
			//递归插入
			b = n.left.insert(e, isMulti, cmp)
		}
		if b {
			//插入成功,对该节点进行调整
			n.insertAdjust()
		}
		return b
	}
	if cmp(n.value, e) < 0 {
		//待插入元素应插入右子树
		if n.right == nil {
			//右子树为空,直接插入
			n.right = newNode(n, e)
			n = n.right
			b = true
		} else {
			//递归插入
			b = n.right.insert(e, isMulti, cmp)
		}
		if b {
			//插入成功,对该节点进行调整
			n.insertAdjust()
		}
		return b
	}
	//在该节点找到等价元素
	if isMulti {
		//允许重复,数值+1即可
		n.num++
		return true
	}
	//不允许重复,覆盖原数值
	n.value = e
	return false
}

//@title    insertAdjust
//@description
//		以node红黑树树节点做接收者
//		对该节点进行调整以实现黑节点平衡
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) insertAdjust() {
	if n == nil {
		return
	}
	var uncle *node
	for n.color == RED && n.parent.color == RED {
		//当自己和父节点都是红色时
		//即存在相邻的红色节点时候,需要进行调整

		//获取新插入的节点的叔叔节点
		uncle = n.getUncle()
		if uncle != nil && uncle.color == RED {
			//叔叔节点也是红色,只需要将父节点和叔叔节点都设为黑色,祖父节点设为红色即可
			//这样可以保证从祖父节点到其后代节点经过的黑色节点数目相同
			//同时对祖父节点即新的红色节点进行平衡检测
			uncle.color, n.parent.color = BLACK, BLACK
			n = n.parent.parent
			if n.getParent() == nil {
				return
			}
			n.color = RED

		} else {
			//当叔叔节点不存在或为黑色时候
			//只需要将祖父节点向叔叔节点方向旋转即可
			//即将父节点放在祖父节点的位置,祖父节点和自己作为父节点的子节点
			//同时父节点设为黑色,自己和祖父节点设为红色即可
			//由于祖父节点本身是黑色,故不影响上面的平衡
			if n.parent == n.parent.parent.left {
				if n == n.parent.right {
					//当自己是右节点的时候,需要先左转否则无法平衡
					//当自己是右节点的时候,祖父节点右转会让该节点跑到祖父节点那边去
					//所以需要先左转让父节点成为左节点
					n = n.parent
					n.leftRotate()
				}
				n = n.parent
				n.color = BLACK
				n = n.parent
				if n!=nil{
					n.color = RED
					n.rightRotate()
				}
			} else {
				//该部分原因同上
				if n == n.parent.left {
					n = n.parent
					n.rightRotate()
				}
				n = n.parent
				n.color = BLACK
				n = n.parent
				if n!=nil{
					n.color = RED
					n.leftRotate()
				}
			}
			return
		}
	}
}

//@title    delete
//@description
//		以node红黑树树节点做接收者
//		从n节点中删除元素e
//		如果n节点中承载元素与e不同则根据大小从左右子树删除该元素
//		如果n节点与该元素相等,且允许重复值,则将num-1否则直接删除该元素
//		删除时先寻找该元素的前缀节点,若不存在则寻找其后继节点进行替换
//		替换后删除该节点
//		删除后对其进行调整以保证祖辈节点的平衡性
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待删除元素
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	b        	bool					是否删除成功?
func (n *node) delete(e interface{}, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//找到待删除节点
	m := n
	for m != nil {
		if cmp(e, m.value) < 0 {
			m = m.left
		} else if cmp(e, m.value) > 0 {
			m = m.right
		} else {
			break
		}
	}
	if m == nil {
		//没找到和元素e相同的节点,删除失败
		return false
	}
	//存在和元素e相同的节点,可以删除
	if m.num > 1 {
		//存在相同元素,减一即可完成删除
		m.num--
		return true
	}
	//找到该节点的前缀节点或者后继节点,以保证被删除的节点不是根节点
	c := m
	if m.right != nil {
		c = m.right
		for c.left != nil {
			c = c.left
		}
	} else if m.left != nil {
		c = m.left
		for c.right != nil {
			c = c.right
		}
	}
	//交换节点存储元素,随后进行删除
	m.value = c.value
	m.num = c.num
	e = c.value
	m = c
	parent := m.parent
	//删除该节点,同时可以确认,该节点必然只有左节点或只有右节点
	//该节点不可能同时拥有左右节点,但该节点可能是叶子节点
	if m.left == nil && m.right == nil {
		//如果要删除的节点没有孩子，直接删掉它就可以
		if m.parent.left == m {
			m.parent.left = nil
		} else {
			m.parent.right = nil
		}
	} else if m.left != nil && m.right == nil {
		//如果要删除的节点只有左孩子或右孩子
		//让这个节点的父节点指向它的指针指向它的孩子即可
		m.left.parent = m.parent
		if m.parent.left == m {
			m.parent.left = m.left
		} else {
			m.parent.right = m.left
		}
	} else if m.left == nil && m.right != nil {
		m.right.parent = m.parent
		if m.parent.left == m {
			m.parent.left = m.right
		} else {
			m.parent.right = m.right
		}
	}
	//节点调整
	if m.color == BLACK {
		//待删除节点不可能是根节点,根节点在删除前已经被处理
		if m == parent.left {
			parent.left.deleteAdjust()
		} else if m == parent.right {
			parent.right.deleteAdjust()
		}
	}
	return true
}

//@title    deleteAdjust
//@description
//		以node红黑树树节点做接收者
//		对该节点进行调整以实现黑节点平衡
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) deleteAdjust() {
	if n == nil {
		return
	}
	var brother *node

	for n.parent != nil && n.color == BLACK {
		if n.parent.left == n && n.parent.right != nil {
			//自己是左节点,兄弟节点是右节点
			//同时,自己是黑节点
			brother = n.parent.right
			if brother.color == RED {
				//到兄弟节点经过的黑节点少1
				brother.color = BLACK
				n.parent.color = RED
				n.parent.leftRotate()
			} else if brother.color == BLACK && brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == BLACK {
				//到兄弟节点的子节点比到自己节点的子节点经过的黑色节点数多1
				brother.color = RED
				n = n.parent
			} else if brother.color == BLACK && brother.left != nil && brother.left.color == RED && brother.right != nil && brother.right.color == BLACK {
				//同上,但是是只有一个多1
				brother.color = RED
				brother.left.color = BLACK
				brother.rightRotate()
			} else if brother.color == BLACK && brother.right != nil && brother.right.color == RED {
				brother.color = RED
				brother.right.color = BLACK
				brother.parent.color = BLACK
				brother.parent.leftRotate()
				break
			} else {
				return
			}

		} else if n.parent.right == n && n.parent.left != nil {
			//原因同上
			brother = n.parent.left
			if brother.color == RED {
				brother.color = BLACK
				n.parent.color = RED
				n.parent.rightRotate()
			} else if brother.color == BLACK && brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == BLACK {
				brother.color = RED
				n = n.parent
			} else if brother.color == BLACK && brother.left != nil && brother.left.color == BLACK && brother.right != nil && brother.right.color == RED {
				brother.color = RED
				brother.right.color = BLACK
				brother.leftRotate()
			} else if brother.color == BLACK && brother.left != nil && brother.left.color == RED {
				brother.color = RED
				brother.left.color = BLACK
				brother.parent.color = BLACK
				brother.parent.rightRotate()
				break
			} else {
				return
			}
		} else {
			return
		}
	}
}

//@title    search
//@description
//		以node红黑树树节点做接收者
//		从n节点中查找元素e并返回存储的个数
//		如果n节点中承载元素与e不同则根据大小从左右子树查找该元素
//		如果n节点与该元素相等,则直接返回其个数
//@auth      	hlccd		2021-07-23
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待查找元素
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	num        	int						待查找元素在红黑树中存储的数量
func (n *node) search(e interface{}, cmp comparator.Comparator) int {
	if n == nil {
		//节点不存在,返回0
		return 0
	}
	if cmp(n.value, e) > 0 {
		//从左子树继续查找
		return n.left.search(e, cmp)
	} else if cmp(n.value, e) < 0 {
		//从右子树继续查找
		return n.right.search(e, cmp)
	}
	//返回查找结果
	return n.num
}
