package queue

//@Title		queue
//@Description
//		queue队列容器包
//		以切片数组的形式实现
//		该容器可以在尾部实现线性增加元素,在首部实现先学减少元素
//		该容器满足FIFO的先进先出模式
//		可接纳不同类型的元素
//@author     	hlccd		2021-07-5
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//queue队列结构体
//包含泛型切片和该切片的首尾指针
//当删除节点时仅仅需要后移首指针一位即可
//当剩余长度小于实际占用空间长度的一半时会重新规划以释放掉多余占用的空间
//当添加节点时若未占满全部已分配空间则尾指针后移一位同时进行覆盖存放
//当添加节点时尾指针大于已分配空间长度,则新增空间
//首节点指针始终不能超过尾节点指针
type queue struct {
	data  []interface{} //泛型切片
	begin int           //首节点指针
	end   int           //尾节点指针
	mutex sync.Mutex    //并发控制锁
}

//queue队列容器接口
//存放了queue容器可使用的函数
//对应函数介绍见下方

type queuer interface {
	Iterator() *iterator.Iterator //返回一个包含queue中所有使用元素的迭代器
	Size() (num int)              //返回该队列中元素的使用空间大小
	Clear()                       //清空该队列
	Empty() (b bool)              //判断该队列是否为空
	Push(e interface{})           //将元素e添加到该队列末尾
	Pop() (e interface{})         //将该队列首元素弹出并返回
	Front() (e interface{})       //获取该队列首元素
	Back() (e interface{})        //获取该队列尾元素
}

//@title    New
//@description
//		新建一个queue队列容器并返回
//		初始queue的切片数组为空
//		初始queue的首尾指针均置零
//@author     	hlccd		2021-07-5
//@receiver		nil
//@param    	nil
//@return    	q        	*queue					新建的queue指针
func New() (q *queue) {
	return &queue{
		data:  make([]interface{}, 0, 0),
		begin: 0,
		end:   0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以queue队列容器做接收者
//		将queue队列容器中不使用空间释放掉
//		返回一个包含容器中所有使用元素的迭代器
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (q *queue) Iterator() (i *iterator.Iterator) {
	if q == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	q.mutex.Lock()
	q.data = q.data[q.begin:q.end]
	q.begin = 0
	q.end = len(q.data)
	i = iterator.New(q.data)
	q.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以queue队列容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		若容器为空则返回-1
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (q *queue) Size() (num int) {
	if q == nil {
		return -1
	}
	return q.end - q.begin
}

//@title    Clear
//@description
//		以queue队列容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的首尾指针均置0
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	nil
func (q *queue) Clear() {
	if q == nil {
		return
	}
	q.mutex.Lock()
	q.data = q.data[0:0]
	q.begin = 0
	q.end = 0
	q.mutex.Unlock()
}

//@title    Empty
//@description
//		以queue队列容器做接收者
//		判断该queue队列容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过首尾指针数值进行判断
//		当尾指针数值等于首指针时说明不含有元素
//		当尾指针数值大于首指针时说明含有元素
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (q *queue) Empty() (b bool) {
	if q == nil {
		return true
	}
	return q.Size() <= 0
}

//@title    Push
//@description
//		以queue队列向量容器做接收者
//		在容器尾部插入元素
//		若尾指针小于切片实际使用长度,则对当前指针位置进行覆盖,同时尾指针后移一位
//		若尾指针等于切片实际使用长度,则新增切片长度同时使尾指针后移一位
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (q *queue) Push(e interface{}) {
	if q == nil {
		return
	}
	q.mutex.Lock()
	if q.end < len(q.data) {
		q.data[q.end] = e
	} else {
		q.data = append(q.data, e)
	}
	q.end++
	q.mutex.Unlock()
}

//@title    Pop
//@description
//		以queue队列容器做接收者
//		弹出容器第一个元素,同时首指针后移一位
//		当剩余元素数量小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//		同时返回队首元素
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e 			interface{}				队首元素
func (q *queue) Pop() (e interface{}) {
	if q == nil {
		return nil
	}
	if q.Empty() {
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.begin]
	q.begin++
	if q.begin*2 >= q.end {
		q.data = q.data[q.begin:q.end]
		q.begin = 0
		q.end = len(q.data)
	}
	q.mutex.Unlock()
	return e
}

//@title    Front
//@description
//		以queue队列容器做接收者
//		返回该容器的第一个元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (q *queue) Front() (e interface{}) {
	if q == nil {
		return nil
	}
	if q.Empty() {
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.begin]
	q.mutex.Unlock()
	return e
}

//@title    Back
//@description
//		以queue队列容器做接收者
//		返回该容器的最后一个元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-5
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e			interface{}				容器的最后一个元素
func (q *queue) Back() (e interface{}) {
	if q == nil {
		return nil
	}
	if q.Empty() {
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.end-1]
	q.mutex.Unlock()
	return e
}
