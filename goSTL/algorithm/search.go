package algorithm

//@Title		algorithm
//@Description
//		算法包
//		该部分为对迭代器中元素进行二分查找
//		warning:仅对有序元素集合有效
//@author     	hlccd		2021-07-2
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
)

//@title    Search
//@description
//		通过比较器对传入两个迭代器中的元素集合进行二分查找
//		找到后返回该元素的下标
//		若该元素不在该部分内存在,则返回-1
//@author     	hlccd		2021-07-2
//@receiver		nil
//@param    	begin		*iterator.Iterator			待排序的起始迭代器
//@param    	end			*iterator.Iterator			待排序的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...comparator.Comparator	比较器
//@return    	idx			int							待查找元素下标
func Search(begin, end *iterator.Iterator, e interface{}, Cmp ...comparator.Comparator) (idx int) {
	//判断比较器是否有效,若无效则寻找默认比较器
	var cmp comparator.Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = comparator.GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		//若并非默认类型且未传入比较器则直接结束
		return -1
	}
	//查找开始
	return search(begin, end, e, cmp)
}

//@title    search
//@description
//		通过比较器对传入两个迭代器中的元素集合进行二分查找
//		找到后返回该元素的下标
//		若该元素不在该部分内存在,则返回-1
//@author     	hlccd		2021-07-2
//@receiver		nil
//@param    	begin		*iterator.Iterator			待排序的起始迭代器
//@param    	end			*iterator.Iterator			待排序的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	cmp			comparator.Comparator	比较器
//@return    	idx			int							待查找元素下标
func search(begin, end *iterator.Iterator, e interface{}, cmp comparator.Comparator) (idx int) {
	//通过二分查找的方式寻找该元素
	m, l, r := 0, begin.Index(), end.Index()
	for l < r {
		m = (l + r) / 2
		if cmp(begin.Get(m).Value(), e) < 0 {
			l = m + 1
		} else {
			r = m
		}
	}
	//查找结束
	if begin.Get(l).Value() == e {
		//该元素存在,返回下标
		return l
	}
	//该元素不存在,返回-1
	return -1
}
