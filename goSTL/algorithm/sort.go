package algorithm

//@Title		algorithm
//@Description
//		算法包
//		该包内通过传入迭代器和比较器进行排序
//		当前支持二分排序和归并排序
//		warning:当传入两个迭代器并非同一个迭代器的不同位置时会发送故障(但应该没啥人能干出这事吧)
//@author     	hlccd		2021-07-2
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
)

//@title    Sort
//@description
//		对传入的开启和结尾的两个比较器中的值进行排序
//		以传入的比较器进行比较
//		若未传入比较器则寻找默认比较器,默认比较器排序结果为升序
//		若该泛型类型并非系统默认类型之一,则不进行排序
//		当待排序元素个数超过2^16个时使用归并排序
//		当待排序元素个数少于2^16个时使用二分排序
//@author     	hlccd		2021-07-2
//@receiver		nil
//@param    	begin		*iterator.Iterator			待排序的起始迭代器
//@param    	end			*iterator.Iterator			待排序的末尾迭代器
//@param    	Cmp			...comparator.Comparator	比较器
//@return    	nil
func Sort(begin, end *iterator.Iterator, Cmp ...comparator.Comparator) {
	//获取两个迭代器之间的差值,若末尾迭代器不在起始迭代器后方则终止
	gap := end.Index() - begin.Index()
	if gap <= 0 {
		return
	}
	//判断迭代器是否存在,如有传入比较器则按照传入执行,否则去寻找默认比较器
	var cmp comparator.Comparator
	cmp = nil
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	} else {
		cmp = comparator.GetCmp(begin.Value())
	}
	if cmp == nil {
		//未传入比较器且并非默认类型导致未找到默认比较器则直接终止排序
		return
	}
	//
	if gap >= 65536 {
		merge(begin, end, cmp)
	} else {
		quick(begin, end, cmp)
	}
}

//@title    merge
//@description
//		归并排序
//		对传入的两个迭代器中的内容使用比较器进行归并排序
//@author     	hlccd		2021-07-2
//@receiver		nil
//@param    	begin		*iterator.Iterator		待排序的起始迭代器
//@param    	end			*iterator.Iterator		待排序的末尾迭代器
//@param    	Cmp			comparator.Comparator	比较器
//@return    	nil
func merge(begin, end *iterator.Iterator, cmp comparator.Comparator) {
	l, r := begin.Index(), end.Index()
	if l >= r {
		return
	}
	m := (r + l) / 2
	//对待排序内容进行二分
	merge(begin.Get(l), begin.Get(m), cmp)
	merge(end.Get(m+1), end.Get(r), cmp)
	//二分结束后依次比较进行归并
	i, j := l, m+1
	var tmp []interface{}
	for i <= m && j <= r {
		if cmp(begin.Get(i).Value(), begin.Get(j).Value()) <= 0 {
			tmp = append(tmp, begin.Get(i).Value())
			i++
		} else {
			tmp = append(tmp, begin.Get(j).Value())
			j++
		}
	}
	//当一方比较到头时将另一方剩余内容全部加入进去
	for ; i <= m; i++ {
		tmp = append(tmp, begin.Get(i).Value())
	}
	for ; j <= r; j++ {
		tmp = append(tmp, begin.Get(j).Value())
	}
	//将局部排序结果放入迭代器中
	for i, j = l, 0; i <= r; i, j = i+1, j+1 {
		begin.Get(i).Set(tmp[j])
	}
}

//@title    quick
//@description
//		二分排序
//		对传入的两个迭代器中的内容使用比较器进行二分排序
//@author     	hlccd		2021-07-2
//@receiver		nil
//@param    	begin		*iterator.Iterator		待排序的起始迭代器
//@param    	end			*iterator.Iterator		待排序的末尾迭代器
//@param    	Cmp			comparator.Comparator	比较器
//@return    	nil
func quick(begin, end *iterator.Iterator, cmp comparator.Comparator) {
	//对当前部分进行预排序,使得两侧都大于或小于中间值
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
	//对分好的两侧进行迭代二分排序
	quick(begin.Get(l), end.Get(j), cmp)
	quick(begin.Get(j+1), end.Get(r), cmp)
}
