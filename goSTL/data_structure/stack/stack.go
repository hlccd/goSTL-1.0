package stack

//@Title		stack
//@Description
//		stack栈容器包
//		以切片数组的形式实现
//		该容器可以在顶部实现线性增减元素
//		通过interface实现泛型
//		可接纳不同类型的元素
//@author     	hlccd		2021-07-7
import "github.com/hlccd/goSTL/utils/iterator"

//vector向量结构体
//包含泛型切片和该切片的顶部指针
//当删除节点时仅仅需要前移顶部指针一位即可
//当剩余长度小于实际占用空间长度的一半时会重新规划以释放掉多余占用的空间
//当添加节点时若未占满全部已分配空间则顶部指针后移一位同时进行覆盖存放
//当添加节点时顶部指针大于已分配空间长度,则新增空间
type stack struct {
	data []interface{} //泛型切片
	top  int           //顶部指针
}

//stack栈容器接口
//存放了stack容器可使用的函数
//对应函数介绍见下方
type stacker interface {
	Iterator() (i *iterator.Iterator) //返回一个包含栈中所有元素的迭代器
	Size() (num int)                  //返回该栈中元素的使用空间大小
	Clear()                           //清空该栈容器
	Empty() (b bool)                  //判断该栈容器是否为空
	Push(e interface{})               //将元素e添加到栈顶
	Pop()                             //弹出栈顶元素
	Top() (e interface{})             //返回栈顶元素
}

//@title    New
//@description
//		新建一个stack栈容器并返回
//		初始stack的切片数组为空
//		初始stack的顶部指针置0
//@author     	hlccd		2021-07-7
//@receiver		nil
//@param    	nil
//@return    	s        	*stack					新建的stack指针
func New() (s *stack) {
	return &stack{
		data: make([]interface{}, 0, 0),
		top:  0,
	}
}

//@title    Iterator
//@description
//		以stack栈容器做接收者
//		将stack栈容器中不使用空间释放掉
//		返回一个包含容器中所有使用元素的迭代器
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (s *stack) Iterator() (i *iterator.Iterator) {
	if s == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	s.data = s.data[:s.top]
	return iterator.New(s.data)
}

//@title    Size
//@description
//		以stack栈容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		如果容器为nil返回-1
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (s *stack) Size() (num int) {
	if s == nil {
		return -1
	}
	return s.top
}

//@title    Clear
//@description
//		以stack栈容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的尾指针置0
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	nil
func (s *stack) Clear() {
	if s == nil {
		return
	}
	s.data = s.data[0:0]
	s.top = 0
}

//@title    Empty
//@description
//		以stack栈容器做接收者
//		判断该stack栈容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过顶部指针数值进行判断
//		当顶部指针数值为0时说明不含有元素
//		当顶部指针数值大于0时说明含有元素
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (s *stack) Empty() (b bool) {
	if s == nil {
		return true
	}
	return s.Size() <= 0
}

//@title    Push
//@description
//		以stack栈容器做接收者
//		在容器顶部插入元素
//		若顶部指针小于切片实际使用长度,则对当前指针位置进行覆盖,同时顶部指针上移一位
//		若顶部指针等于切片实际使用长度,则新增切片长度同时使尾指针上移一位
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	e			interface{}				待插入顶部的元素
//@return    	nil
func (s *stack) Push(e interface{}) {
	if s == nil {
		return
	}
	if s.top < len(s.data) {
		s.data[s.top] = e
	} else {
		s.data = append(s.data, e)
	}
	s.top++
}

//@title    Pop
//@description
//		以stack栈容器做接收者
//		弹出容器顶部元素,同时顶部指针下移一位
//		当顶部指针小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	nil
func (s *stack) Pop() {
	if s == nil {
		return
	}
	if s.Empty() {
		return
	}
	s.top--
	if s.top*2 <= len(s.data) {
		s.data = s.data[0:s.top]
	}
}

//@title    Top
//@description
//		以stack栈容器做接收者
//		返回该容器的顶部元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-7
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	e			interface{}				容器的顶部元素
func (s *stack) Top() (e interface{}) {
	if s == nil {
		return nil
	}
	if s.Empty() {
		return nil
	}
	return s.data[s.top-1]
}
