package algorithm

//@Title		algorithm
//@Description
//		算法包
//		该包内算法可查找有序序列中某一元素的上界和下届
//		当该元素存在时,上界和下届返回的下标位置指向该元素
//		当该元素不存在时,上界和下届指向下标错位且指向位置并非该元素
//@author     	hlccd		2021-07-3
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
)

//@title    UpperBound
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找待查找元素的上界
//		以传入的比较器进行比较
//		如果该元素存在,则上界指向元素为该元素
//		如果该元素不存在,上界指向元素为该元素的前一个元素
//		该部分主要对迭代器和比较器进行判断
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待查找的起始迭代器
//@param    	end			*iterator.Iterator			待查找的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...comparator.Comparator	比较器
//@return    	idx 		int							待查找元素的上界
func UpperBound(begin, end *iterator.Iterator, e interface{}, Cmp ...comparator.Comparator) (idx int) {
	//判断末尾迭代器是否在起始迭代器前方
	gap := end.Index() - begin.Index()
	if gap <= 0 {
		return
	}
	//判断比较器是否有效
	var cmp comparator.Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = comparator.GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的上界
	return upperBound(begin, end, e, cmp)
}

//@title    upperBound
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找待查找元素的上界
//		以传入的比较器进行比较
//		如果该元素存在,则上界指向元素为该元素
//		如果该元素不存在,上界指向元素为该元素的前一个元素
//		以二分查找的方式寻找该元素的上界
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待查找的起始迭代器
//@param    	end			*iterator.Iterator			待查找的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	cmp			comparator.Comparator		比较器
//@return    	idx 		int							待查找元素的上界
func upperBound(begin, end *iterator.Iterator, e interface{}, cmp comparator.Comparator) (idx int) {
	m, l, r := 0, begin.Index(), end.Index()
	for l < r {
		m = (l + r + 1) / 2
		if cmp(begin.Get(m).Value(), e) <= 0 {
			l = m
		} else {
			r = m - 1
		}
	}
	return l
}

//@title    LowerBound
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找待查找元素的下界
//		以传入的比较器进行比较
//		如果该元素存在,则下界指向元素为该元素
//		如果该元素不存在,下界指向元素为该元素的下一个元素
//		该部分主要对迭代器和比较器进行判断
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待查找的起始迭代器
//@param    	end			*iterator.Iterator			待查找的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	Cmp			...comparator.Comparator	比较器
//@return    	idx 		int							待查找元素的下界
func LowerBound(begin, end *iterator.Iterator, e interface{}, Cmp ...comparator.Comparator) (idx int) {
	//判断末尾迭代器是否在起始迭代器前方
	gap := end.Index() - begin.Index()
	if gap <= 0 {
		return
	}
	//判断比较器是否有效
	var cmp comparator.Comparator
	cmp = nil
	if len(Cmp) == 0 {
		cmp = comparator.GetCmp(e)
	} else {
		cmp = Cmp[0]
	}
	if cmp == nil {
		return -1
	}
	//寻找该元素的下界
	return lowerBound(begin, end, e, cmp)
}

//@title    lowerBound
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找待查找元素的下界
//		以传入的比较器进行比较
//		如果该元素存在,则下界指向元素为该元素
//		如果该元素不存在,下界指向元素为该元素的下一个元素
//		以二分查找的方式寻找该元素的下界
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待查找的起始迭代器
//@param    	end			*iterator.Iterator			待查找的末尾迭代器
//@param    	e			interface{}					待查找元素
//@param    	cmp			comparator.Comparator		比较器
//@return    	idx 		int							待查找元素的下界
func lowerBound(begin, end *iterator.Iterator, e interface{}, cmp comparator.Comparator) (idx int) {
	m, l, r := 0, begin.Index(), end.Index()
	for l < r {
		m = (l + r) / 2
		if cmp(begin.Get(m).Value(), e) >= 0 {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}
