package ring

//@Title		ring
//@Description
//		ring环容器包
//		以切片数组的形式实现
//		该容器可以在当前节点实现线性增减元素
//		可接纳不同类型的元素
//@author     	hlccd		2021-07-8
import "github.com/hlccd/goSTL/utils/iterator"

//ring环结构体
//包含泛型切片和该切片的当前位置指针
//增删节点都会重新规划空间
type ring struct {
	data  []interface{} //泛型切片
	index int           //当前节点指针
}

//ring环容器接口
//存放了ring容器可使用的函数
//对应函数介绍见下方
type ringer interface {
	Iterator() (i *iterator.Iterator) //返回一个包含ring容器中所有使用元素的迭代器
	Size() (num int)                  //返回该ring容器中所含有的元素个数
	Clear()                           //清空该ring容器
	Empty() (b bool)                  //判断该ring容器是否为空
	Insert(e interface{})             //在当前节点元素后面添加一个元素
	Erase()                           //删除该节点元素
	Next()                            //将该ring容器节点后移
	Pre()                             //将该ring容器节点前移
	Value() (e interface{})           //返回该ring容器当前节点元素
}

//@title    New
//@description
//		新建一个ring环容器并返回
//		初始ring的切片数组为空
//		初始ring的节点指针置零
//@author     	hlccd		2021-07-8
//@receiver		nil
//@param    	nil
//@return    	r        	*ring					新建的ring指针
func New() (r *ring) {
	return &ring{
		data:  make([]interface{}, 0, 0),
		index: 0,
	}
}

//@title    Iterator
//@description
//		以ring环容器做接收者
//		返回一个包含容器中所有使用元素的迭代器
//		以当前节点元素作为迭代器首元素
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (r *ring) Iterator() (i *iterator.Iterator) {
	if r == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	if r.Empty() {
		return iterator.New(make([]interface{}, 0, 0))
	}
	return iterator.New(append(r.data[r.index:], r.data[:r.index]...))
}

//@title    Size
//@description
//		以ring环容器做接收者
//		返回该容器当前含有元素的数量
//		当容器不存在时,返回-1
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	num        	int						容器中存储元素的个数
func (r *ring) Size() (num int) {
	if r == nil {
		return -1
	}
	return len(r.data)
}

//@title    Clear
//@description
//		以ring环容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的节点指针置0
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	nil
func (r *ring) Clear() {
	if r == nil {
		return
	}
	r.data = r.data[0:0]
	r.index = -1
}

//@title    Empty
//@description
//		以ring环容器做接收者
//		判断该ring环容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过含有元素个数进行判断
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (r *ring) Empty() (b bool) {
	if r == nil {
		return true
	}
	return r.Size() <= 0
}

//@title    Insert
//@description
//		以ring环容器做接收者
//		在容器当前节点下一位置插入元素,并将原节点以后元素后移一位
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (r *ring) Insert(e interface{}) {
	if r == nil {
		return
	}
	if r.index < r.Size()-1 {
		es := append([]interface{}{}, r.data[r.index+1:]...)
		r.data = append(append(r.data[:r.index+1], e), es...)
	} else {
		r.data = append(r.data, e)
	}
}

//@title    Erase
//@description
//		以ring环容器做接收者
//		如果元素集合为空则直接结束
//		否则删除当前节点元素并将该节点后面元素前移一位,该节点指向原节点的下一位
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	nil
func (r *ring) Erase() {
	if r == nil {
		return
	}
	if r.Empty() {
		return
	}
	if r.index == 0 {
		r.data = r.data[1:]
	} else if r.index == r.Size()-1 {
		r.data = r.data[:r.Size()-1]
		r.index = 0
	} else {
		es := append([]interface{}{}, r.data[:r.index]...)
		r.data = append(es, r.data[r.index+1:]...)
	}
}

//@title    Next
//@description
//		以ring环容器做接收者
//		如果元素集合为空则直接结束
//		否则将该节点前移,令原节点位置指向下一位元素
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	nil
func (r *ring) Next() {
	if r == nil {
		return
	}
	if r.Empty() {
		return
	}
	r.index = (r.index + 1) % r.Size()
}

//@title    Pre
//@description
//		以ring环容器做接收者
//		如果元素集合为空则直接结束
//		否则将该节点后移,令原节点位置指向前一位元素
//@auth      	hlccd		2021-07-8
//@receiver		r			*ring					接受者ring的指针
//@param    	nil
//@return    	nil
func (r *ring) Pre() {
	if r == nil {
		return
	}
	if r.Empty() {
		return
	}
	r.index = (r.index - 1 + r.Size()) % r.Size()
}

//@title    Value
//@description
//		以ring环容器做接收者
//		返回当前节点指向元素
//		若该容器当前为空,则返回nil
//		若容器为nil则返回nil
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (r *ring) Value() (e interface{}) {
	if r == nil {
		return nil
	}
	if r.Empty() {
		return nil
	}
	return r.data[r.index]
}
