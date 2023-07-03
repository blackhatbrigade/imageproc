package imageproc

import (
  "golang.org/x/sys/unix"
)

// wrapper for ioctl system call
func ioctl(fd, req, arg uintptr) (err error) {
  errno := sys.Syscall(sys.SYS_IOCTL, fd, req, arg)
  if errno != 0 {
    err = errno
    return
  }
  return nil
}
