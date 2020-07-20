package tcpmax

import (
    "math/bits"
)
const (
    BITS_PER_WORD  uint16 = 6
    WORDS_PER_SIZE uint16 = 64
    UINT64_MASK    uint64 = 0xffffffffffffffff
    BITMAP_SIZE    uint16 = 4096
    INDEX_MASK     uint16 = 0x3f
) 

type BitMap struct {
    flag       uint64
    bitmap [WORDS_PER_SIZE]uint64
}

func getleb(x uint64) int{
    if x == 0 {
        return -1
    }
    pos := bits.TrailingZeros64(x & ((x - 1) ^ UINT64_MASK))
    return int(pos)
 }

 func (bm *BitMap) GetID() int {
    if bm.flag == 0{
       return -1
    }
    index := getleb(bm.flag)
    if index == -1 {
       return -1
    }
    pos := getleb(bm.bitmap[index])
    if pos == -1 {
       return -1
    }
    id := (index << BITS_PER_WORD) + pos
    bm.Set(uint16(id))
    return id
}

func (bm *BitMap) Set(x uint16) {
    if x < BITMAP_SIZE{
        index := (x >> BITS_PER_WORD) & INDEX_MASK
        pos := x & INDEX_MASK
        bm.bitmap[index] = bm.bitmap[index] & ((1 << uint(pos)) ^ UINT64_MASK)
        if bm.bitmap[index] == 0 {
            bm.flag = bm.flag & ((1 << uint(index)) ^ UINT64_MASK)
        }
    }
}

func (bm *BitMap) Clear(x int) {
    id := uint16(x)
    if id < BITMAP_SIZE{
        index := (id >> BITS_PER_WORD) & INDEX_MASK
        pos := id & INDEX_MASK
        bm.bitmap[index] = bm.bitmap[index] | (1 << uint(pos))
        bm.flag = bm.flag | (1 << uint(index))
    }
}

func (bm *BitMap) Init() {
    bm.flag = UINT64_MASK
    for i := uint16(0); i < WORDS_PER_SIZE; i++ {
        bm.bitmap[i] = UINT64_MASK
    }
}
