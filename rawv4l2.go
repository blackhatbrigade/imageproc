package imageproce

// v4l2 command encoding
const (
  iocOpNone  = 0
  iocOpWrite = 1
  iocOpRead  = 2
  iocTypeBits   = 8
  iocNumberBits = 8
  iocSizeBits   = 14
  iocOpBits     = 2
  numberPos = 0
  typePos   = numberPos + iocNumberBits
  sizePos   = typePos + iocTypeBits
  opPos     = sizePos + iocSizeBits
)

// ioctl command encoding
func ioEnc(iocMode, iocType, number, size uintptr) uintptr {
  return (iocMode << opPos) |
    (iocType << typePos) |
    (number << numberPos) |
    (size << sizePos)
}

// convienence funcs
func ioEncR(iocType, number, size uintptr) uintptr {
  return ioEnc(iocOpRead, iocType, number, size)
}

func ioEncW(iocType, number, size uintptr) uintptr {
  return ioEnc(iocOpWrite, iocType, number, size)
}

func ioEncRW(iocType, number, size uintptr) uintptr {
  return ioEnc(iocOpRead|iocOpWrite, iocType, number, size)
}

// fourcc implements C v4l2 fourcc
func fourcc(a, b, c, d uint32) uint32 {
  return (a | b << 8 | c << 16 | d << 24)
}
