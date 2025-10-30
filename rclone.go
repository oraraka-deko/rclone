// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

/*
#include "stdint.h"
*/
import "C"
import (
	"unsafe"

	_ "github.com/rclone/rclone/backend/all" // import all backends
	"github.com/rclone/rclone/cmd"
	_ "github.com/rclone/rclone/cmd/all"    // import all commands
	_ "github.com/rclone/rclone/lib/plugin" // import plugins
	bridge "hiddify.com/hiddify/bridge"
)

//export initNativeDartBridge
func initNativeDartBridge(api unsafe.Pointer) {
	bridge.InitDartApi(api)
}
func main() {
	cmd.Main()
}
