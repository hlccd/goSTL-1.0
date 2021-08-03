package algorithm

//@Title		algorithm
//@Description
//		算法包
//		该包内通过传入迭代器和比较器查找位于第n个的节点
//		对二分排序的变形,当只对该节点位置存在的某一局部进行查找即可
//@author     	hlccd		2021-07-3
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
)

//@title    NthElement
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找
//		以传入的比较器进行比较
//		该部分主要进行比较器的判断以及n的位置处理
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待排序的起始迭代器
//@param    	end			*iterator.Iterator			待排序的末尾迭代器
//@param    	n			int							待查找的是第n位,从0计数
//@param    	Cmp			...comparator.Comparator	比较器
//@return    	nil
func NthElement(begin, end *iterator.Iterator, n int, Cmp ...comparator.Comparator) {
	//判断末尾迭代器是否在起始迭代器前方
	gap := end.Index() - begin.Index()
	if gap <= 0 {
		return
	}
	//判断比较器是否有效
	var cmp comparator.Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = comparator.GetCmp(begin.Value())
	}
	if cmp == nil {
		return
	}
	//判断待确认的第n位是否在该集合范围内
	if n <= begin.Index() {
		n = begin.Index()
	}
	if n >= end.Index() {
		n = end.Index()
	}
	//进行查找
	nthElement(begin, end, n, cmp)
}

//@title    nthElement
//@description
//		对传入的开启和结尾的两个比较器中的值进行查找
//		以传入的比较器进行比较
//		通过局部二分的方式进行查找并将第n小的元素放到第n位置(大小按比较器进行确认,默认未小)
//@author     	hlccd		2021-07-3
//@receiver		nil
//@param    	begin		*iterator.Iterator			待排序的起始迭代器
//@param    	end			*iterator.Iterator			待排序的末尾迭代器
//@param    	n			int							待查找的是第n位,从0计数
//@param    	cmp			comparator.Comparator		比较器
//@return    	nil
func nthElement(begin, end *iterator.Iterator, n int, cmp comparator.Comparator) {
	//二分该区域并对此进行预排徐
	l, r := begin.Index(), end.Index()
	if l >= r {
		return
	}
	m := begin.Get((r + l) / 2).Value()
	i, j := l-1, r+1
	for i < j {
		i++
		for cmp(begin.Get(i).Value(), m) < 0 {
			i++
		}
		j--
		for cmp(end.Get(j).Value(), m) > 0 {
			j--
		}
		if i < j {
			ti := begin.Get(i).Value()
			tj := begin.Get(j).Value()
			begin.Get(i).Set(tj)
			begin.Get(j).Set(ti)
		}
	}
	//确认第n位的范围进行局部二分
	if n-1 >= i {
		nthElement(begin.Get(j+1), end.Get(r), n, cmp)
	} else {
		nthElement(begin.Get(l), end.Get(j), n, cmp)
	}
}
