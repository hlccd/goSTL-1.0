package multiset

//@Title		multiset
//@Description
//		线程不安全,不建议使用
//		multiset可重复集合容器包
//		以切片数组的形式实现
//		该容器可以增减元素后依然保持整体有序
//		若使用默认比较器,则整体序列是升序
//		该集合只能对于相等元素可以存储多个
//		可接纳不同类型的元素,但为了便于比较,建议使用同一个类型
//@author     	hlccd		2021-07-9
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//multiset可重复集合结构体
//包含泛型切片和比较器
//增删节点后会使用比较器保持该切片数组的有序性
type multiset struct {
	data  []interface{}         //泛型切片
	cmp   comparator.Comparator //该可重复集合的比较器
	mutex sync.Mutex            //并发控制锁
}

//multiset可重复集合容器接口
//存放了multiset容器可使用的函数
//对应函数介绍见下方
type multiseter interface {
	Iterator() (i *iterator.Iterator)          //返回一个包含multiset容器中所有使用元素的迭代器
	Size() (num int)                           //返回该可重复集合中存储的元素数量
	Clear()                                    //清空该可重复集合
	Empty() (b bool)                           //判断该可重复集合是否为空
	Insert(e interface{})                      //插入元素e
	Erase(e interface{})                       //删除元素e
	Count(e interface{}) (num int)             //查找元素e并返回该元素个数
	Find(e interface{}) (i *iterator.Iterator) //查找元素e并返回指向该元素的迭代器
}

//@title    New
//@description
//		新建一个multiset可重复集合容器并返回
//		初始multiset的切片数组为空
//		如果有传入比较器,则将传入的第一个比较器设为可重复集合默认比较器
//@author     	hlccd		2021-07-9
//@receiver		nil
//@param    	Cmp			...comparator.Comparator	multiset的比较器集
//@return    	ms        	*multiset					新建的multiset指针
func New(Cmp ...comparator.Comparator) (ms *multiset) {
	var cmp comparator.Comparator
	if len(Cmp) == 0 {
		cmp = nil
	} else {
		cmp = Cmp[0]
	}
	return &multiset{
		data:  make([]interface{}, 0, 0),
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以multiset可重复集合容器做接收者
//		返回一个包含容器中所有使用元素的迭代器
//@author     	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (ms *multiset) Iterator() (i *iterator.Iterator) {
	if ms == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	ms.mutex.Lock()
	i = iterator.New(ms.data)
	ms.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以multiset可重复集合容器做接收者
//		返回该容器当前含有元素的数量
//		当容器不存在时,返回-1
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	nil
//@return    	num        	int						容器中存储元素的个数
func (ms *multiset) Size() (num int) {
	if ms == nil {
		return -1
	}
	return len(ms.data)
}

//@title    Clear
//@description
//		以multiset可重复集合容器做接收者
//		将该容器中所承载的元素清空
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	nil
//@return    	nil
func (ms *multiset) Clear() {
	if ms == nil {
		return
	}
	ms.mutex.Lock()
	ms.data = ms.data[0:0]
	ms.mutex.Unlock()
}

//@title    Empty
//@description
//		以multiset可重复集合容器做接收者
//		判断该multiset容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过含有元素个数进行判断
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (ms *multiset) Empty() (b bool) {
	if ms == nil {
		return true
	}
	return ms.Size() <= 0
}

//@title    Insert
//@description
//		以multiset可重复集合容器做接收者
//		在该集合中插入元素e,通过查找到对应位置进行插入,保证插入后集合仍然处于有序状态
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (ms *multiset) Insert(e interface{}) {
	if ms == nil {
		return
	}
	ms.mutex.Lock()
	if ms.Empty() {
		ms.data = append(ms.data, e)
	} else {
		if ms.cmp == nil {
			ms.cmp = comparator.GetCmp(e)
		}
		if ms.cmp == nil {
			return
		}
		begin := ms.Iterator().Begin()
		end := ms.Iterator().End()
		p := algorithm.LowerBound(begin, end, e, ms.cmp)
		if p == len(ms.data)-1 {
			ms.data = append(ms.data, e)
		} else if p == 0 {
			es := append([]interface{}{}, e)
			ms.data = append(es, ms.data...)
		} else {
			es := append([]interface{}{}, ms.data[p:]...)
			ms.data = append(append(ms.data[:p], e), es...)
		}
	}
	ms.mutex.Unlock()
}

//@title    Erase
//@description
//		以multiset可重复集合容器做接收者
//		删除在集合中的一个元素e
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	e			interface{}				待删除元素
//@return    	nil
func (ms *multiset) Erase(e interface{}) {
	if ms == nil {
		return
	}
	ms.mutex.Lock()
	p := algorithm.Search(ms.Iterator().Begin(), ms.Iterator().End(), e, ms.cmp)
	if p != -1 {
		if len(ms.data) == 1 {
			ms.Clear()
		} else {
			if p == 0 {
				ms.data = ms.data[1:]
			} else {
				es := append([]interface{}{}, ms.data[:p]...)
				ms.data = append(es, ms.data[p+1:]...)
			}
		}
	}
	ms.mutex.Unlock()
}

//@title    Count
//@description
//		以multiset可重复集合容器做接收者
//		返回元素e在集合中的个数
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	e			interface{}				待查找元素
//@return    	num			int						容器的e元素数量
func (ms *multiset) Count(e interface{}) (num int) {
	if ms == nil {
		return 0
	}
	ms.mutex.Lock()
	upper := algorithm.UpperBound(ms.Iterator().Begin(), ms.Iterator().End(), e, ms.cmp)
	if e != ms.data[upper] {
		return 0
	}
	lower := algorithm.LowerBound(ms.Iterator().Begin(), ms.Iterator().Get(upper), e, ms.cmp)
	num = upper - lower + 1
	if num <= 0 {
		num = 0
	}
	ms.mutex.Unlock()
	return num
}

//@title    Find
//@description
//		以multiset可重复集合容器做接收者
//		返回直线元素e的迭代器
//		如果元素e不在集合中存在,返回nil
//@auth      	hlccd		2021-07-9
//@receiver		ms			*multiset				接受者multiset的指针
//@param    	e			interface{}				待查找元素
//@return    	i			*iterator.Iterator		指向e元素的迭代器
func (ms *multiset) Find(e interface{}) (i *iterator.Iterator) {
	if ms == nil {
		return nil
	}
	ms.mutex.Lock()
	p := algorithm.Search(ms.Iterator().Begin(), ms.Iterator().End(), e, ms.cmp)
	if p != -1 {
		ms.mutex.Unlock()
		return ms.Iterator().Get(p)
	}
	ms.mutex.Unlock()
	return nil
}
