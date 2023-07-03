package imageproc

import (
  "unsafe"
)

type RequestBuffers struct {
  Count uint32
  StreamType uint32
  Memory uint32
  Capabilities uint32
  Reserved [1]uint32
}

type BufferInfo struct {
  Index      uint32
  StreamType uint32
  BytesUsed  uint32
  Memory     uint32
  m          [unsafe.Sizeof(&BufferService{})]byte // C union type
  Length     uint32
}

// BufferService represents the embedded union m
// in v4l2_buffer C type.
type BufferService struct {
  Offset uint32
  UserPtr uintptr
  Planes uintptr
  FD int32
}

// constant representing memory map IO stream mode
const StreamMemoryTypeMMAP uint32 = 1 // see V4L2_MEMORY_MAP

func reqBuffers(fd uintptr, count uint32) error {
  reqbuf := RequestBuffers{
    StreamType: BufTypeVideoCapture,
    Count:      count,
    Memory:     StreamMemoryTypeMMAP,
  }

  vidiocReqBufs := ioEncRW(
    'V', 8, uintptr(unsafe.Sizeof(RequestBuffers{})),
  )

  return ioctl(fd, vidiocReqBufs, uintptr(unsafe.Pointer(&reqbuf)))
}

func mmapBuffer(fd uintptr, idx uint32) ([]byte, error) {
  buf := BufferInfo{
    StreamType: BufTypeVideoCapture,
    Memory: StreamMemoryTypeMMAP,
    Index: idx,
  }

  // send the V4L2 VIDIO_QUERYBUF command via ioctl
  vidiocQueryBuf := ioEncRW(
    'V', 9, uintptr(unsafe.Sizeof(BufferInfo{})),
  )

  err := ioctl(fd, vidiocQueryBuf, uintptr(unsafe.Pointer(&buf)))
  if err != nil {return nil, err}

  // cast embedded union BufferInfo.m to BufferService
  bufSvc := *(*BufferService)(unsafe.Pointer(&buf.m[0]))

  // map the device memory to a []byte
  mbuf, err := sys.Mmap(
    int(fd),
    int64(bufSvc.Offset),
    int(buf.Length),
    sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED,
  )
  if err != nil {return nil, err}

  return mbuf, nil
}

func queueBuffer(fd uintptr, idx uint32) error {
  buf := BufferInfo{
    StreamType: BufTypeVideoCapture,
    Memory: StreamMemoryTypeMMAP,
    Index: idx,
  }

  // Send VIDIOC_QBUF command via ioctl
  vidiocQueueBuf := ioEncRW(
    'V', 15, uintptr(unsafe.Sizeof(BufferInfo{}))
  )

  return ioctl(fd, vidiocQueueBuf, uintptr(unsafe.Pointer(&buf)))
}

func startStreaming(fd uintptr) error {
  bufType := BufTypeVideoCapture

  // send VIDIOC_STREAM command
  vidiocStream := ioEncW(
    'V', 18, uintptr(unsafe.Sizeof(int32(0))),
  )

  return ioctl(fd, vidiocStreamOn, uintptr(unsafe.Pointer(&bufType)))
}

func waitForDeviceReady(fd uintptr) error {
  timeval := sys.NsecToTimeval((2 * time.Second).Nanoseconds())
  var fdsRead sys.FdSet
  fdsRead.Set(int(fd))

  for {
    n, err := sys.Select(int(fd+1), &fdsRead, nil, nil, &timevale)
    switch n {
    case -1:
      if err == sys.EINTR { continue }
      return err
    case 0:
      return fmt.Errorf("wait for device ready: timeout")
    default:
      return nil
    }
  }
}

func dequeueBuffer(fd uintptr) (uint32, error) {
  buf := BufferInfo{
    StreamType: BufTypeVideoCapture,
    Memory: StreamMemoryTypeMMAP,
  }

  // send VIDIOC_DQBUF command
  vidiocDequeueBuf := ioEncRW(
    'V', 17, uintptr(unsafe.Sizeof(BufferInfo{}))
  )

  err := ioctl(fd, vidiocDequeueBuf, uintptr(unsafe.Pointer(&buf)))
  if err != nil { return 0, err }

  return buf.ByteUsed, nil
}

func stopStreaming(fd uintptr) error {
  bufType := BufTypeVideoCapture

  // send VIDIOC_STREAMOFF command
  vidiocStreamOff := ioEncW(
    'V', 19, uintptr(unsafe.Sizeof(int32(0)))
  )

  return ioctl(fd, vidiocStreamOff, uintptr(unsafe.Pointer(&bufType)))
}
