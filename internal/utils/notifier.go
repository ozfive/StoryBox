package utils

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/kataras/iris/v12"
)

var SdNotifyNoSocket = fmt.Errorf("No socket")

// SdNotify sends a message to the init daemon. It is common to ignore the error.
func SdNotify(state string) error {
	socketPath := os.Getenv("NOTIFY_SOCKET")
	if socketPath == "" {
		return SdNotifyNoSocket
	}

	conn, err := net.Dial("unixgram", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	return err
}

func HandleError(errorConst string, err error, ctx iris.Context) {
	_, writeErr := ctx.HTML(errorConst + err.Error() + "</b>")
	if writeErr != nil {
		log.Println(writeErr.Error())
	}
}
