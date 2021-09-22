package trie

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

type trie struct {
	root  *node      //根节点指针
	mutex sync.Mutex //并发控制锁
}

type trieer interface {
	Iterator() (i *iterator.Iterator) //返回包含该树堆的所有元素,重复则返回多个
	Size() (num int)                  //返回该树堆中保存的元素个数
	Clear()                           //清空该树堆
	Empty() (b bool)                  //判断该树堆是否为空
	Insert(s string, e interface{})   //向树堆中插入元素e
	Erase(s string)                   //从树堆中删除元素e
	Count(s string) (num int)         //从树堆中寻找元素e并返回其个数
	Find(s string) (e interface{})
}

func New() (t *trie) {
	return &trie{
		root:  newNode(nil),
		mutex: sync.Mutex{},
	}
}
func (t *trie) Iterator() (i *iterator.Iterator) {
	if t == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	t.mutex.Lock()
	i = iterator.New(t.root.inOrder(""))
	t.mutex.Unlock()
	return i
}
func (t *trie) Size() (num int) {
	if t == nil {
		return -1
	}
	if t.root == nil {
		return -1
	}
	return t.root.num
}
func (t *trie) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = newNode(nil)
	t.mutex.Unlock()
}
func (t *trie) Empty() (b bool) {
	if t.Size() > 0 {
		return false
	}
	return true
}
func (t *trie) Insert(s string, e interface{}) {
	//判断容器是否存在
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root.num++
	now := t.root
	for i := 0; i < len(s); i++ {
		if now.son[s[i]-'a'] == nil {
			now.son[s[i]-'a'] = newNode(nil)
		}
		now = now.son[s[i]-'a']
		now.num++
	}
	if e == nil {
		now.value = s
	} else {
		now.value = e
	}
	t.mutex.Unlock()
}
func (t *trie) Erase(s string) {
	if t == nil {
		return
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	if s == "" {
		t.root = newNode(nil)
	} else {
		num := 0
		for i, now := 0, t.root; i < len(s); i++ {
			if now.son[s[i]-'a'] == nil {
				num = 0
				break
			}
			now = now.son[s[i]-'a']
			num = now.num
		}
		if num > 0 {
			for i, now := 0, t.root; i < len(s); i++ {
				now.num -= num
				if i == len(s)-1 {
					for j := 0; j < 26; j++ {
						now.son[j] = nil
					}
				} else {
					now = now.son[s[i]-'a']
				}
			}
		}
	}
	t.mutex.Unlock()
}
func (t *trie) Count(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	now := t.root
	for i := 0; i < len(s); i++ {
		if now.son[s[i]-'a'] == nil {
			num = 0
			break
		}
		now = now.son[s[i]-'a']
		num = now.num
	}
	t.mutex.Unlock()
	//树堆存在,从根节点开始查找该元素
	return num
}
func (t *trie) Find(s string) (e interface{}) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	now := t.root
	for i := 0; i < len(s); i++ {
		if now.son[s[i]-'a'] == nil {
			e = nil
			break
		}
		now = now.son[s[i]-'a']
		e = now.value
	}
	t.mutex.Unlock()
	//树堆存在,从根节点开始查找该元素
	return e
}
