package imageproc

import (
  "unsafe"
)

// Format represents C type v4l2_format
type Format struct {
  StreamType uint32
  fmt        [200]byte // union
}

// PixFormat represents C type v4l2_pix_format
type PixFormat struct {
  Width        uint32
  Height       uint32
  PixelFormat  uint32
  Field        uint32
  BytesPerLine uint32
  SizeImage    uint32
  Colorspace   uint32
  Priv         uint32
  Flags        uint32
  YcbcrEnc     uint32
  Quantization uint32
  XferFunc     uint32
}

// values for pix format v4l2_fields
const (
  FieldAny  uint32 = iota 
  FieldNone               
)

func setFormat(fd uintptr, pixFmt PixFormat) error {
  format := Format{StreamType: BufTypeVideoCapture}

  // a bit of C union type magic with unsafe.Pointer
  *(*PixFormat)(unsafe.Pointer(&format.fmt[0])) = pixFmt

  // encode command VIDIC_S_FMT
  vidiocSetFormat := ioEncRW(
    'V', 5, uintptr(unsafe.Sizeof(Format{})),
  )

  // send command
  return ioctl(fd, vidioSetFormat, uintptr(unsafe.Pointer(&format)))
}
