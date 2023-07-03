package imageproc

import (
  "fmt"

  "golang.org/x/sys/unix"
)

func Test() {
  fmt.Println("Image processing imported successfully")
  fmt.Println("Capturing image from webcam")

  devName := "/dev/video0"

  // open the device and get a file descriptor
  devFile, _ := os.OpenFile(devName, sys.O_RDWR|sys.O_NONBLOCK, 0)

}

func SetupCamera() {
  devFile, err := os.OpenFile(devName, sys.O_RDWR|sys.O_NONBLOCK, 0)
  if err != nil {
    fmt.Println("Error opening device!")
    return
  }

  fd := devFile.Fd()
}
