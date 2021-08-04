package heap

//@Title		heap
//@Description
//		heap堆集合容器包
//		以切片数组的形式实现
//		该容器可以增减元素使最值元素处于顶端
//		若使用默认比较器,顶端元素是最小元素
//		该集合只能对于相等元素可以存储多个
//		可接纳不同类型的元素,但为了便于比较,建议使用同一个类型
//@author     	hlccd		2021-07-10
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//heap堆集合结构体
//包含泛型切片和比较器
//增删节点后会使用比较器保持该切片数组的有序性
type heap struct {
	data []interface{}         //泛型切片
	cmp  comparator.Comparator //该堆的比较器
	mutex   sync.Mutex            //并发控制锁
}

//heap堆容器接口
//存放了heap容器可使用的函数
//对应函数介绍见下方
type heaper interface {
	Iterator() (i *iterator.Iterator) //返回一个包含heap容器中所有使用元素的迭代器
	Size() (num int)                  //返回该容器存储的元素数量
	Clear()                           //清空该容器
	Empty() (b bool)                  //判断该容器是否为空
	Push(e interface{})               //将元素e插入该容器
	Pop()                             //弹出顶部元素
	Top() (e interface{})             //返回顶部元素
}

//@title    New
//@description
//		新建一个heap堆容器并返回
//		初始heap的切片数组为空
//		如果有传入比较器,则将传入的第一个比较器设为可重复集合默认比较器
//@author     	hlccd		2021-07-10
//@receiver		nil
//@param    	Cmp			...comparator.Comparator	heap的比较器集
//@return    	h        	*heap						新建的heap指针
func New(cmps ...comparator.Comparator) (h *heap) {
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	return &heap{
		data: make([]interface{}, 0, 0),
		cmp:  cmp,
		mutex:   sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以heap容器做接收者
//		返回一个包含容器中所有使用元素的迭代器
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (h *heap) Iterator() (i *iterator.Iterator) {
	if h == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	h.mutex.Lock()
	i=iterator.New(h.data)
	h.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以heap容器做接收者
//		返回该容器当前含有元素的数量
//		当容器不存在时,返回-1
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	num        	int						容器中存储元素的个数
func (h *heap) Size() (num int) {
	if h == nil {
		return -1
	}
	return len(h.data)
}

//@title    Clear
//@description
//		以heap容器做接收者
//		将该容器中所承载的元素清空
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	nil
func (h *heap) Clear() {
	if h == nil {
		return
	}
	h.mutex.Lock()
	h.data = h.data[0:0]
	h.mutex.Unlock()
}

//@title    Empty
//@description
//		以heap容器做接收者
//		判断该heap容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过含有元素个数进行判断
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (h *heap) Empty() bool {
	if h == nil {
		return true
	}
	return h.Size() <= 0
}

//@title    Push
//@description
//		以heap容器做接收者
//		在该堆中插入元素e,利用比较器和交换使得堆保持相对有序状态
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (h *heap) Push(e interface{}) {
	if h == nil {
		return
	}
	h.mutex.Lock()
	if h.cmp == nil {
		h.cmp = comparator.GetCmp(e)
	}
	if h.cmp == nil {
		h.mutex.Unlock()
		return
	}
	h.data = append(h.data, e)
	h.up(len(h.data) - 1)
	h.mutex.Unlock()
}

//@title    up
//@description
//		以heap容器做接收者
//		判断待上升节点与其父节点的大小关系以确定是否进行递归上升
//		从而保证父节点必然都大于或都小于子节点
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	p			int						待上升节点的位置
//@return    	nil
func (h *heap) up(p int) {
	if p == 0 {
		return
	}
	if h.cmp(h.data[(p-1)/2], h.data[p]) > 0 {
		h.data[p], h.data[(p-1)/2] = h.data[(p-1)/2], h.data[p]
		h.up((p - 1) / 2)
	}
}

//@title    Pop
//@description
//		以heap容器做接收者
//		在该堆中删除顶部元素,利用比较器和交换使得堆保持相对有序状态
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	nil
func (h *heap) Pop() {
	if h == nil {
		return
	}
	h.mutex.Lock()
	if h.Empty() {
		h.mutex.Unlock()
		return
	}
	h.data[0] = h.data[h.Size()-1]
	h.data = h.data[:h.Size()-1]
	if h.Empty() {
		h.mutex.Unlock()
		return
	}
	h.down(0)
	h.mutex.Unlock()
}

//@title    down
//@description
//		以heap容器做接收者
//		判断待下沉节点与其左右子节点的大小关系以确定是否进行递归上升
//		从而保证父节点必然都大于或都小于子节点
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	p			int						待下沉节点的位置
//@return    	nil
func (h *heap) down(p int) {
	q := p
	if 2*p+1 <= h.Size()-1 && h.cmp(h.data[p], h.data[2*p+1]) > 0 {
		q = 2*p + 1
	}
	if 2*p+2 < h.Size()-1 && h.cmp(h.data[q], h.data[2*p+2]) > 0 {
		q = 2*p + 2
	}
	if p != q {
		h.data[p], h.data[q] = h.data[q], h.data[p]
		h.down(q)
	}
}

//@title    Top
//@description
//		以heap容器做接收者
//		返回该堆容器的顶部元素
//		如果容器不存在或容器为空,返回nil
//@author     	hlccd		2021-07-10
//@receiver		h			*heap					接受者heap的指针
//@param    	nil
//@return    	e			interface{}				堆顶元素
func (h *heap) Top() (e interface{}) {
	if h == nil {
		return nil
	}
	if h.Empty() {
		return nil
	}
	h.mutex.Lock()
	e=h.data[0]
	h.mutex.Unlock()
	return e
}
