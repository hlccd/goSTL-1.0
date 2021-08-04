package set

//@Title		set
//@Description
//		set集合容器包
//		以切片数组的形式实现
//		该容器可以增减元素后依然保持整体有序
//		若使用默认比较器,则整体序列是升序
//		该集合只能对于相等元素仅能存储一个
//		可接纳不同类型的元素,但为了便于比较,建议使用同一个类型
//@author     	hlccd		2021-07-9
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//set集合结构体
//包含泛型切片和比较器
//增删节点后会使用比较器保持该切片数组的有序性
type set struct {
	data  []interface{}         //泛型切片
	cmp   comparator.Comparator //该集合的比较器
	mutex sync.Mutex            //并发控制锁
}

//set集合容器接口
//存放了set容器可使用的函数
//对应函数介绍见下方
type seter interface {
	Iterator() (i *iterator.Iterator)          //返回一个包含set容器中所有使用元素的迭代器
	Size() (num int)                           //返回该集合中存储的元素数量
	Clear()                                    //清空该集合
	Empty() (b bool)                           //判断该集合是否为空
	Insert(e interface{})                      //插入元素e
	Erase(e interface{})                       //删除元素e
	Count(e interface{}) (num int)             //查找元素e并返回该元素个数
	Find(e interface{}) (i *iterator.Iterator) //查找元素e并返回指向该元素的迭代器
}

//@title    New
//@description
//		新建一个set集合容器并返回
//		初始set的切片数组为空
//		如果有传入比较器,则将传入的第一个比较器设为集合默认比较器
//@author     	hlccd		2021-07-9
//@receiver		nil
//@param    	Cmp			...comparator.Comparator	set的比较器集
//@return    	s        	*set						新建的set指针
func New(Cmp ...comparator.Comparator) (s *set) {
	var cmp comparator.Comparator
	if len(Cmp) == 0 {
		cmp = nil
	} else {
		cmp = Cmp[0]
	}
	return &set{
		data:  make([]interface{}, 0, 0),
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以set集合容器做接收者
//		返回一个包含容器中所有使用元素的迭代器
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (s *set) Iterator() (i *iterator.Iterator) {
	if s == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	s.mutex.Lock()
	i = iterator.New(s.data)
	s.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以set集合容器做接收者
//		返回该容器当前含有元素的数量
//		当容器不存在时,返回-1
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	nil
//@return    	num        	int						容器中存储元素的个数
func (s *set) Size() (num int) {
	if s == nil {
		return -1
	}
	return len(s.data)
}

//@title    Clear
//@description
//		以set集合容器做接收者
//		将该容器中所承载的元素清空
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	nil
//@return    	nil
func (s *set) Clear() {
	if s == nil {
		return
	}
	s.mutex.Lock()
	s.data = s.data[0:0]
	s.mutex.Unlock()
}

//@title    Empty
//@description
//		以set集合容器做接收者
//		判断该set集合容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过含有元素个数进行判断
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (s *set) Empty() (b bool) {
	if s == nil {
		return true
	}
	return s.Size() <= 0
}

//@title    Insert
//@description
//		以set集合容器做接收者
//		在该集合中插入元素e,通过查找到对应位置进行插入,保证插入后集合仍然处于有序状态
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (s *set) Insert(e interface{}) {
	if s == nil {
		return
	}
	s.mutex.Lock()
	if s.Empty() {
		s.data = append(s.data, e)
	} else {
		if s.cmp == nil {
			s.cmp = comparator.GetCmp(e)
		}
		if s.cmp == nil {
			return
		}
		begin := s.Iterator().Begin()
		end := s.Iterator().End()
		p := algorithm.LowerBound(begin, end, e, s.cmp)
		if s.data[p] == e {
			return
		} else {
			if p == len(s.data)-1 {
				s.data = append(s.data, e)
			} else if p == 0 {
				es := append([]interface{}{}, e)
				s.data = append(es, s.data...)
			} else {
				es := append([]interface{}{}, s.data[p:]...)
				s.data = append(append(s.data[:p], e), es...)
			}
		}
	}
	s.mutex.Unlock()
}

//@title    Erase
//@description
//		以set集合容器做接收者
//		删除在集合中的元素e
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	e			interface{}				待删除元素
//@return    	nil
func (s *set) Erase(e interface{}) {
	if s == nil {
		return
	}
	s.mutex.Lock()
	p := algorithm.Search(s.Iterator().Begin(), s.Iterator().End(), e, s.cmp)
	if p != -1 {
		if len(s.data) == 1 {
			s.Clear()
		} else {
			if p == 0 {
				s.data = s.data[1:]
			} else {
				es := append([]interface{}{}, s.data[:p]...)
				s.data = append(es, s.data[p+1:]...)
			}
		}
	}
	s.mutex.Unlock()
}

//@title    Count
//@description
//		以set集合容器做接收者
//		返回元素e在集合中的个数
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	e			interface{}				待查找元素
//@return    	num			int						容器的e元素数量
func (s *set) Count(e interface{}) (num int) {
	if s == nil {
		return 0
	}
	s.mutex.Lock()
	p := algorithm.Search(s.Iterator().Begin(), s.Iterator().End(), e, s.cmp)
	if p != -1 {
		s.mutex.Unlock()
		return 1
	}
	s.mutex.Unlock()
	return 0
}

//@title    Find
//@description
//		以set集合容器做接收者
//		返回直线元素e的迭代器
//		如果元素e不在集合中存在,返回nil
//@auth      	hlccd		2021-07-9
//@receiver		s			*set					接受者set的指针
//@param    	e			interface{}				待查找元素
//@return    	i			*iterator.Iterator		指向e元素的迭代器
func (s *set) Find(e interface{}) (i *iterator.Iterator) {
	if s == nil {
		return nil
	}
	s.mutex.Lock()
	p := algorithm.Search(s.Iterator().Begin(), s.Iterator().End(), e, s.cmp)
	if p != -1 {
		s.mutex.Unlock()
		return s.Iterator().Get(p)
	}
	s.mutex.Unlock()
	return nil
}
