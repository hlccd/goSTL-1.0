package rbTree

//@Title		rbTree
//@Description
//		红黑树-Red Black Tree
//		以二叉树的形式实现
//		二叉树实例保存根节点和比较器以及保存的数量
//		可以在创建时设置节点是否可重复
//		若节点可重复则增加节点中的数值,否则对节点存储元素进行覆盖
//		红黑树在添加和删除时都将对节点进行平衡,从根节点到任意叶子节点经过的黑节点数目相同
//@author     	hlccd		2021-07-23
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//rbTree红黑树结构体
//该实例存储二叉树的根节点
//同时保存该二叉树已经存储了多少个元素
//二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找
//创建时传入是否允许该二叉树出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可
type rbTree struct {
	root    *node
	size    int
	cmp     comparator.Comparator
	isMulti bool
	mutex   sync.Mutex //并发控制锁
}

//rbTree红黑树容器接口
//存放了rbTree红黑树可使用的函数
//对应函数介绍见下方
type rbTreeer interface {
	Iterator() (i *iterator.Iterator) //返回包含该红黑树的所有元素,重复则返回多个
	Size() (num int)                  //返回该红黑树中保存的元素个数
	Clear()                           //清空该红黑树
	Empty() (b bool)                  //判断该v是否为空
	Insert(e interface{})             //向红黑树中插入元素e
	Erase(e interface{})              //从红黑树中删除元素e
	Count(e interface{}) (num int)    //从红黑树中寻找元素e并返回其个数
}

//@title    New
//@description
//		新建一个rbTree红黑树容器并返回
//		初始根节点为nil
//		传入该红黑树是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖
//		若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器
//@author     	hlccd		2021-07-11
//@receiver		nil
//@param    	isMulti		bool						该二叉树是否保存重复值?
//@param    	Cmp			 ...comparator.Comparator	rbTree比较器集
//@return    	rb        	*rbTree						新建的rbTree指针
func New(isMulti bool, cmps ...comparator.Comparator) (rb *rbTree) {
	//判断是否有传入比较器,若有则设为该红黑树默认比较器
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	return &rbTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以rbTree红黑搜索树做接收者
//		将该红黑树中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中
//		若允许重复存储则对于重复元素进行多次放入
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (rb *rbTree) Iterator() (i *iterator.Iterator) {
	if rb == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	rb.mutex.Lock()
	i = iterator.New(rb.root.inOrder())
	rb.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以rbTree红黑搜索树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回-1
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (rb *rbTree) Size() (num int) {
	if rb == nil {
		return -1
	}
	return rb.size
}

//@title    Clear
//@description
//		以rbTree红黑搜索树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	nil
//@return    	nil
func (rb *rbTree) Clear() {
	if rb == nil {
		return
	}
	rb.mutex.Lock()
	rb.root = nil
	rb.size = 0
	rb.mutex.Unlock()
}

//@title    Empty
//@description
//		以rbTree红黑搜索树做接收者
//		判断该红黑树是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (rb *rbTree) Empty() (b bool) {
	if rb == nil {
		return true
	}
	if rb.size > 0 {
		return false
	}
	return true
}

//@title    Insert
//@description
//		以rbTree红黑搜索树做接收者
//		向红黑树插入元素e,若不允许重复则对相等元素进行覆盖
//		如果红黑树为空则之间用根节点承载元素e,否则以根节点开始进行查找
//		不做平衡
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (rb *rbTree) Insert(e interface{}) {
	if rb == nil {
		return
	}
	rb.mutex.Lock()
	if rb.Empty() {
		if rb.cmp == nil {
			rb.cmp = comparator.GetCmp(e)
		}
		if rb.cmp == nil {
			rb.mutex.Unlock()
			return
		}
		//红黑树为空,用根节点承载元素e
		rb.root = newNode(nil, e)
		rb.root.color = BLACK
		rb.size = 1
		rb.mutex.Unlock()
		return
	}
	if rb.root.insert(e, rb.isMulti, rb.cmp) {
		rb.size++
	}
	rb.mutex.Unlock()
}

//@title    Erase
//@description
//		以rbTree红黑搜索树做接收者
//		从搜素红黑树中删除元素e
//		若允许重复记录则对承载元素e的节点中数量记录减一即可
//		若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证二叉树的不发送断裂
//		如果该红黑树仅持有一个元素且根节点等价于待删除元素,则将红黑树根节点置为nil
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	e			interface{}				待删除元素
//@return    	nil
func (rb *rbTree) Erase(e interface{}) {
	if rb == nil {
		return
	}
	if rb.Empty() {
		return
	}
	rb.mutex.Lock()
	if rb.size == 1 && rb.root.value == e {
		//删除跟节点
		rb.root = nil
		rb.size = 0
		rb.mutex.Unlock()
		return
	}
	//删除元素e
	if rb.root.delete(e, rb.cmp) {
		//删除成功
		rb.size--
	}
	rb.mutex.Unlock()
}

//@title    Count
//@description
//		以rbTree红黑搜索树做接收者
//		从红黑树中查找元素e的个数
//		如果找到则返回该二叉树中和元素e相同元素的个数
//		如果不允许重复则最多返回1
//		如果未找到则返回0
//@auth      	hlccd		2021-07-23
//@receiver		rb			*rbTree					接受者rbTree的指针
//@param    	e			interface{}				待查找元素
//@return    	num			int						待查找元素在二叉树中存储的个数
func (rb *rbTree) Count(e interface{}) (num int) {
	if rb == nil {
		return 0
	}
	if rb.Empty() {
		return 0
	}
	rb.mutex.Lock()
	num = rb.root.search(e, rb.cmp)
	rb.mutex.Unlock()
	return num
}
