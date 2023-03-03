package GoBitMap

import (
	"sync"
)

type BitMap struct {
	Bits []uint8
	lock sync.Mutex
}

func (_self *BitMap) autoAppend(num uint) {
	sectionIndex := int(num >> 3)
	currLen := len(_self.Bits)
	if sectionIndex >= len(_self.Bits) {
		_self.Bits = append(_self.Bits, make([]uint8, sectionIndex+1-currLen)...)
	}
}

func (_self *BitMap) Set(num uint) {
	defer _self.lock.Unlock()
	_self.lock.Lock()
	if _self.Exist(num) {
		return
	}
	_self.autoAppend(num)
	_self.Bits[num>>3] |= 1 << (num % 8)
}

func (_self *BitMap) UnSet(num uint) {
	defer _self.lock.Unlock()
	_self.lock.Lock()
	if !_self.Exist(num) {
		return
	}
	_self.autoAppend(num)
	_self.Bits[num>>3] &^= 1 << (num % 8)
}

func (_self *BitMap) Exist(num uint) bool {
	if int(num/8) >= len(_self.Bits) {
		return false
	}
	return 0 != (_self.Bits[num>>3] & (1 << (num % 8)))
}

func (_self *BitMap) Dump() []uint {
	result := make([]uint, 0, 256)
	_self.Each(func(num uint) {
		result = append(result, num)
	})
	return result
}

func (_self *BitMap) Each(f func(num uint)) {
	for sectionIndex, section := range _self.Bits {
		if section == 0 {
			continue
		}
		//10101010
		//i是从右到左索引
		for i := 7; i >= 0; i-- {
			if ((section << i) >> 7) == 1 {
				f(uint(sectionIndex*8 + 7 - i))
			}
		}
	}
}

func (_self *BitMap) Merge(bm *BitMap) {
	var largestBits []uint8
	var minBits []uint8
	if len(_self.Bits) >= len(bm.Bits) {
		largestBits = _self.Bits
		minBits = bm.Bits
	} else {
		largestBits = bm.Bits
		minBits = _self.Bits
	}
	for sectionIndex, section := range minBits {
		largestBits[sectionIndex] |= section
	}
	_self.Bits = largestBits
}

func (_self *BitMap) Count() uint {
	var count uint
	for _, section := range _self.Bits {
		if section == 0 {
			continue
		}
		//10101010
		//i是从右到左索引
		for i := 7; i >= 0; i-- {
			if ((section << i) >> 7) == 1 {
				count += 1
			}
		}
	}
	return count
}
