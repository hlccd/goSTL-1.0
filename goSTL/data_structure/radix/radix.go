package radix

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"strings"
	"sync"
)

type radix struct {
	root  *node
	mutex sync.Mutex
}

type radixer interface {
	Iterator() (i *iterator.Iterator) //返回包含该树堆的所有元素,重复则返回多个
	Size() (num int)                  //返回该树堆中保存的元素个数
	Clear()                           //清空该树堆
	Empty() (b bool)                  //判断该树堆是否为空
	Insert(s string, e interface{})   //向树堆中插入元素e
	Erase(s string)                   //从树堆中删除元素e
	Count(s string) (num int)         //从树堆中寻找元素e并返回其个数
	Find(s string) (e interface{})
}

func New() (t *radix) {
	return &radix{
		root:  newNode("", nil),
		mutex: sync.Mutex{},
	}
}
func (t *radix) Iterator() (i *iterator.Iterator) {
	if t == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	t.mutex.Lock()
	i = iterator.New(t.root.inOrder(""))
	t.mutex.Unlock()
	return i
}
func (t *radix) Size() (num int) {
	if t == nil {
		return -1
	}
	if t.root == nil {
		return -1
	}
	return t.root.num
}
func (t *radix) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = newNode("", nil)
	t.mutex.Unlock()
}
func (t *radix) Empty() (b bool) {
	if t.Size() > 0 {
		return false
	}
	return true
}
func (t *radix) Insert(s string, e interface{}) {
	//判断容器是否存在
	if t == nil {
		return
	}
	t.mutex.Lock()
	now := t.root
	b := true
	ss := strings.Split(s, "/")
	for i := 0; i < len(ss); i++ {
		if ss[i] != "" {
			b = true
			now.num++
			for j := 0; j < len(now.son); j++ {
				if now.son[j].name == ss[i] {
					b = false
					now = now.son[j]
				}
			}
			if b {
				now.son = append(now.son, newNode(ss[i], nil))
				now = now.son[len(now.son)-1]
				now.name = ss[i]
			}
			if i == len(ss)-1 {
				now.num++
				if e == nil {
					now.value = ss[len(ss)-1]
				} else {
					now.value = e
				}
			}
		}
	}
	t.mutex.Unlock()
}
func (t *radix) Erase(s string) {
	if t == nil {
		return
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	if s == "" {
		t.root = newNode("", nil)
	} else {
		num := 0
		ss := strings.Split(s, "/")
		for i := 0; i < len(ss); i++ {
			if ss[i] == "" {
				ss = append(ss[:i], ss[i+1:]...)
				i--
			}
		}
		b := true
		for i, now := 0, t.root; i < len(ss); i++ {
			b = true
			for j := 0; j < len(now.son); j++ {
				if now.son[j].name == ss[i] {
					now = now.son[j]
					num = now.num
					b = false
					break
				}
			}
			if b {
				num = 0
				break
			}
		}
		if num > 0 {
			for i, now := 0, t.root; i < len(ss); i++ {
				now.num -= num
				if i == len(ss)-1 {
					for j := 0; j < len(now.son); j++ {
						if now.son[j].name == ss[i] {
							now.son = append(now.son[:j], now.son[j+1:]...)
						}
					}
					//now.son = make([]*node, 0, 0)
				} else {
					for j := 0; j < len(now.son); j++ {
						if now.son[j].name == ss[i] {
							now = now.son[j]
						}
					}
				}
			}
		}
	}

	t.mutex.Unlock()
}
func (t *radix) Count(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	ss := strings.Split(s, "/")
	b := true
	for i, now := 0, t.root; i < len(ss); i++ {
		if ss[i] != "" {
			b = true
			for j := 0; j < len(now.son); j++ {
				if now.son[j].name == ss[i] {
					now = now.son[j]
					num = now.num
					b = false
					break
				}
			}
			if b {
				num = 0
				break
			}
		}
	}
	t.mutex.Unlock()
	return num
}
func (t *radix) Find(s string) (e interface{}) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	ss := strings.Split(s, "/")
	b := true
	for i, now := 0, t.root; i < len(ss); i++ {
		if ss[i] != "" {
			b = true
			for j := 0; j < len(now.son); j++ {
				if now.son[j].name == ss[i] {
					now = now.son[j]
					e = now.value
					b = false
					break
				}
			}
			if b {
				e = nil
				break
			}
		}
	}
	t.mutex.Unlock()
	return e
}
