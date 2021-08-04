## 介绍

**支持泛型**

goSTL 是 Go 的数据结构和算法库，旨在提供类似于 C++ STL的功能，但功能更强大。结合go语言的特点，所有数据结构都实现了goroutine-safe。

**set和multiset由于高并发时线程不安全,不建议使用,实例中已经删除该两项实例,建议用红黑树替代**

[TOC]

## 示例

### 工具

#### 迭代器

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/utils/iterator"
)

func main() {
	arr:=make([]interface{},0,0)
	arr=append(arr,"测试数据1")
	arr=append(arr,"测试数据2")
	arr=append(arr,"测试数据3")
	i:=iterator.New(arr)
	i.Begin()//返回指向迭代器首迭代的迭代器
	i.End()//返回指向迭代器尾节点的迭代器
	i.Get(1)//返回指向位置1的迭代器
	i.Set("测试")//在迭代器当前位置设置值为"测试",不建议使用
	i.Index()//返回迭代器当前下标值
	i.Value()//返回迭代器当前值
	i.HasPre()//判断迭代器是否可以后移
	i.Next()//迭代器后移
	i.HasPre()//判断迭代器是否可以前移
	i.Pre()//迭代器前移
	//迭代器遍历
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	for i=i.End();i.HasPre();i.Pre(){
		fmt.Println(i.Value())
	}
}
```

#### 比较器

```go
//比较器本身用于配合其他数据结构完成相应类型的支持

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}

//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
```

### 数据结构

#### 向量-vector

向量,可从尾部插入或删除,也可在任意位置插入或删除

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
func main() {
	//实例化
	v := vector.New()
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			v.PushBack(e) //尾部添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	v.Insert(0, pair{-1, -1}) //自定义位置添加元素
	i := v.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", v.Size())
	//逆转vector
	fmt.Println("逆转后的结果")
	v.Reverse()
	i = v.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", v.Size())
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			v.PopBack() //尾部弹出元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = v.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", v.Size())
	v.Erase(0) //自定义位置删除元素
	v.Clear() //清空
	v.Empty() //判断该实例是否为空
	fmt.Println("is empty?",v.Empty())
}

```

#### 队列-queue

队列,遵守先进先出原则

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/queue"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}

func main() {
	//实例化
	q := queue.New()
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			q.Push(e) //尾部添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	fmt.Println("队首元素",q.Front())//获取队首元素
	fmt.Println("队尾元素",q.Back())//获取队尾元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			q.Pop() //尾部弹出元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	q.Clear()  //清空
	q.Empty()  //判断该实例是否为空
	fmt.Println("is empty?", q.Empty())
}

```

#### 双向队列-deque

队列的双向版本,支持头尾插入和退出

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/deque"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}

func main() {
	//实例化
	q := deque.New()
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			q.PushBack(e) //尾部添加元素
			q.PushFront(e)//队首添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	fmt.Println("队首元素",q.Front())//获取队首元素
	fmt.Println("队尾元素",q.Back())//获取队尾元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			q.PopBack() //队尾弹出元素
			q.PopFront()//队首弹出元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	q.Clear()  //清空
	q.Empty()  //判断该实例是否为空
	fmt.Println("is empty?", q.Empty())
}
```

#### 栈-stack

遵守先进后出原则

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/stack"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
func main() {
	//实例化
	q := stack.New()
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			q.Push(e) //向栈顶添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	fmt.Println("队首元素",q.Top())//获取栈顶元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			q.Pop() //弹出栈顶元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	q.Clear()  //清空
	q.Empty()  //判断该实例是否为空
	fmt.Println("is empty?", q.Empty())
}

```

#### 环-ring

环形结构体

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/ring"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
func main() {
	//实例化
	q := ring.New()
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			q.Insert(e) //向环中添加元素,添加到该元素的下一位
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	fmt.Println("当前指向元素",q.Value())//获取当前指向元素
	q.Next()//后移一位
	fmt.Println("当前指向元素",q.Value())//获取当前指向元素
	q.Pre()//前移一位
	fmt.Println("当前指向元素",q.Value())//获取当前指向元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			q.Erase() //弹出当前元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = q.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", q.Size())
	q.Clear()  //清空
	q.Empty()  //判断该实例是否为空
	fmt.Println("is empty?", q.Empty())
}
```

#### 堆-heap

可自行选择传入比较器即可(若不传入比较器,则自动寻找默认比较器,仅对常规类型元素有效)

极值元素处于堆顶

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/heap"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	h := heap.New(cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			h.Push(e) //向堆添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := h.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", h.Size())
	fmt.Println("堆顶元素", h.Top()) //获取堆顶
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			h.Pop() //弹出堆顶元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = h.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", h.Size())
	h.Clear() //清空
	h.Empty() //判断该实例是否为空
	fmt.Println("is empty?", h.Empty())
}

```

#### 二叉搜索树-bsTree

二叉树搜索是由二叉树有序集合,构造时需要传入"isMulti",即需要传入是否允许重复

随后可自行选择传入比较器即可(若不传入比较器,则自动寻找默认比较器,仅对常规类型元素有效)

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/bsTree"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//二叉搜索树需要先传入是否允许保存重复值
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	s := bsTree.New(true,cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			s.Insert(e) //向集合添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := s.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", s.Size())
	fmt.Println("查找元素数量", s.Count(pair{-1,-1})) //查找元素并返回该元素在树中的数量
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			s.Erase(e) //删除元素e
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = s.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", s.Size())
	s.Clear() //清空
	s.Empty() //判断该实例是否为空
	fmt.Println("is empty?", s.Empty())
}

```

#### 完全二叉树-cbTree

完全二叉树是由二叉树实现的堆,构造时不需要像其他二叉树一样传入"isMulti",即不需要传入是否允许重复

也可自行选择传入比较器即可(若不传入比较器,则自动寻找默认比较器,仅对常规类型元素有效)

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/cbTree"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	t := cbTree.New(cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			t.Push(e) //向集合添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	fmt.Println("查看树顶元素", t.Top()) //返回树顶元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			t.Pop() //删除元素e
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	t.Clear() //清空
	t.Empty() //判断该实例是否为空
	fmt.Println("is empty?", t.Empty())
}

```

#### 树堆-treap

符合二叉搜索树的大小关系,同时利用随机的优先级使得整个二叉树达到随机平衡,即依概率平衡

使用时基本等同于是二叉搜索树

仅通过随机的优先级实现依概率平衡

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/treap"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//树堆需要先传入是否允许保存重复值
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	t := treap.New(true,cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			t.Insert(e) //向集合添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	fmt.Println("查看元素在树中的数量", t.Count(pair{-1,-1})) //返回树顶元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			t.Erase(e) //删除元素e
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	t.Clear() //清空
	t.Empty() //判断该实例是否为空
	fmt.Println("is empty?", t.Empty())
}

```

#### 平衡二叉树-avlTree

符合二叉搜索树的大小关系,同时利用节点深度对节点进行动态调整,使得该二叉树满足任意节点的左右子节点深度差不超过1的平衡

使用时基本等同于是二叉搜索树

仅通过深度调节满足平衡

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/avlTree"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//平衡二叉树需要先传入是否允许保存重复值
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	t := avlTree.New(true,cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			t.Insert(e) //向集合添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	fmt.Println("查看元素在树中的数量", t.Count(pair{-1,-1})) //返回树顶元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			t.Erase(e) //删除元素e
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	t.Clear() //清空
	t.Empty() //判断该实例是否为空
	fmt.Println("is empty?", t.Empty())
}

```

#### 红黑树-rbTree

符合二叉搜索树的大小关系,同时满足从根节点到任意叶子节点经过的黑节点数量相等,且无连续的红节点

使用时基本等同于是二叉搜索树

仅通过节点颜色进行调整

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/rbTree"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	//实例化
	//平衡二叉树需要先传入是否允许保存重复值
	//由于需要大小关系,故需要传入比较器,若为自带类型,可不传入比较器,直接使用默认比较器
	t := rbTree.New(true,cmp)
	var wg sync.WaitGroup
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e := pair{m, m}
			t.Insert(e) //向集合添加元素
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("添加后结果")
	i := t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	fmt.Println("查看元素在树中的数量", t.Count(pair{-1,-1})) //返回树顶元素
	//删除元素
	for j := 0; j < 10000; j++ {
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			t.Erase(e) //删除元素e
			wg.Done()
		}(j)
	}
	wg.Wait()
	//获取其迭代器
	fmt.Println("删除后结果")
	i = t.Iterator()
	for i = i.Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("size=", t.Size())
	t.Clear() //清空
	t.Empty() //判断该实例是否为空
	fmt.Println("is empty?", t.Empty())
}

```

### 算法

由于需要使用算法的的只有vector,故以下皆以vector为例

#### 排序-sort

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	var wg sync.WaitGroup
	v:=vector.New()
	for x:=0;x<100;x++{
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			v.PushBack(e)
			wg.Done()
		}(x)
	}
	wg.Wait()
	fmt.Println("排序前")
	i:=v.Iterator()
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println("排序后")
	algorithm.Sort(i.Begin(),i.End(),cmp)
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
}

```

#### 查找-search

查找只针对有序序列有效

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	var wg sync.WaitGroup
	v:=vector.New()
	for x:=0;x<100;x++{
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			v.PushBack(e)
			wg.Done()
		}(x)
	}
	wg.Wait()
	i:=v.Iterator()
	algorithm.Sort(i.Begin(),i.End(),cmp)
	//查找返回其下标,-1为不存在
	fmt.Println(algorithm.Search(i.Begin(),i.End(),pair{-1,-1},cmp))
	fmt.Println(algorithm.Search(i.Begin(),i.End(),pair{7,7},cmp))
}

```

#### 第n大的值-nth_element

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	var wg sync.WaitGroup
	v:=vector.New()
	for x:=0;x<100;x++{
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			v.PushBack(e)
			wg.Done()
		}(x)
	}
	wg.Wait()
	i:=v.Iterator()
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println()
	fmt.Println(i.End().Value())
	algorithm.NthElement(i.Begin(),i.End(),99,cmp)
	fmt.Println()
	fmt.Println(i.End().Value())
}

```

#### 上界-upperBound

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	var wg sync.WaitGroup
	v:=vector.New()
	for x:=0;x<100;x++{
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			v.PushBack(e)
			v.PushBack(e)
			v.PushBack(e)
			wg.Done()
		}(x)
	}
	wg.Wait()
	i:=v.Iterator()
	algorithm.Sort(i.Begin(),i.End(),cmp)
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println("LowerBound:",algorithm.LowerBound(i.Begin(),i.End(),pair{77,77},cmp))
}

```

#### 下界-lowerBound

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/vector"
	"sync"
)

//自定义结构体,需配合对应的比较器
type pair struct {
	value1 int
	value2 int
}
//对应结构体比较器,若为系统自带类型,则可不提供比较器
func cmp(a, b interface{}) int {
	if a == b {
		return 0
	}
	if a.(pair).value1 > b.(pair).value1 {
		return 1
	} else if a.(pair).value1 < b.(pair).value1 {
		return -1
	} else {
		if a.(pair).value2 > b.(pair).value2 {
			return 1
		} else if a.(pair).value2 < b.(pair).value2 {
			return -1
		}
	}
	return 0
}
func main() {
	var wg sync.WaitGroup
	v:=vector.New()
	for x:=0;x<100;x++{
		wg.Add(1)
		go func(m int) {
			e:=pair{m,m}
			v.PushBack(e)
			v.PushBack(e)
			v.PushBack(e)
			wg.Done()
		}(x)
	}
	wg.Wait()
	i:=v.Iterator()
	algorithm.Sort(i.Begin(),i.End(),cmp)
	for i=i.Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	fmt.Println("UpperBound:",algorithm.UpperBound(i.Begin(),i.End(),pair{77,77},cmp))
}

```