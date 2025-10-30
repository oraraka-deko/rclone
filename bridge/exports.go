package bridge

/*
#include "stdint.h"
#include "dart_api/dart_api_dl.h"
#include "dart_api/dart_api_dl.c"
#include "dart_api/dart_native_api.h"
bool GoDart_PostCObject(Dart_Port_DL port, Dart_CObject* obj) {
  return Dart_PostCObject_DL(port, obj);
}
*/
import "C"
import (
	"encoding/json"
)

// ==========================================
// EXPORT PATTERN EXAMPLES FOR DART BRIDGE
// ==========================================
//
// Pattern: //export functionName
// This comment directs cgo to export the function as a C symbol
// that can be called from Dart via FFI (Foreign Function Interface)
//
// Function Parameters:
// - Always start with: port C.longlong (for Dart response port)
// - Additional parameters: *C.char for strings (convert with C.GoString)
// - Other C types as needed
//
// Response Pattern:
// - Always call: SendResponseToPort(int64(port), &DartResponse{...})
// - Never use return values; communicate via the response port
//
// Error Handling:
// - Return early with DartResponse{Success: false, Message: err.Error()}
// - Always send a response, even on error

// ==========================================
// EXAMPLE 1: Simple Get Function
// ==========================================

// ConfigGeneralSettings represents the general configuration
type ConfigGeneralSettings struct {
	AllowLan    *bool   `json:"allowLan,omitempty"`
	BindAddress *string `json:"bindAddress,omitempty"`
	LogLevel    *string `json:"logLevel,omitempty"`
	Port        *int    `json:"port,omitempty"`
	SocksPort   *int    `json:"socksPort,omitempty"`
	MixedPort   *int    `json:"mixedPort,omitempty"`
	RedirPort   *int    `json:"redirPort,omitempty"`
	TProxyPort  *int    `json:"tproxyPort,omitempty"`
}

// Export getGeneralConfig - retrieves general configuration settings
// Called from Dart with just a port number
//
//export getGeneralConfig
func getGeneralConfig(port C.longlong) {
	// TODO: Replace with actual config retrieval logic
	// Example from rclone:
	// general := executor.GetGeneral()

	general := &ConfigGeneralSettings{
		AllowLan:    ptrBool(true),
		BindAddress: ptrString("0.0.0.0"),
		Port:        ptrInt(8080),
	}

	data, err := json.Marshal(general)
	if err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
		Data:    string(data),
	})
}

// ==========================================
// EXAMPLE 2: Update Function with JSON Patch
// ==========================================

// Export patchGeneralConfig - updates general configuration settings
// Accepts JSON patch data from Dart and applies changes
//
//export patchGeneralConfig
func patchGeneralConfig(port C.longlong, patchStr *C.char) {
	patch := C.GoString(patchStr)

	// Unmarshal the patch JSON
	general := &ConfigGeneralSettings{}
	if err := json.Unmarshal([]byte(patch), general); err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Apply each field that was provided
	// Pattern: check if pointer is not nil, then apply
	if general.AllowLan != nil {
		// P.SetAllowLan(*general.AllowLan)  // Replace with actual setter
	}

	if general.BindAddress != nil {
		// P.SetBindAddress(*general.BindAddress)
	}

	if general.Port != nil {
		// P.ReCreateHTTP(*general.Port, ...)
	}

	if general.LogLevel != nil {
		// log.SetLevel(*general.LogLevel)
	}

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ==========================================
// EXAMPLE 3: Function with String Parameter
// ==========================================

// Export getRemoteConfig - retrieves configuration for a specific remote
//
//export getRemoteConfig
func getRemoteConfig(port C.longlong, remoteName *C.char) {
	name := C.GoString(remoteName)

	// TODO: Replace with actual config retrieval
	// config, err := getConfigForRemote(name)

	remoteConfig := map[string]interface{}{
		"name": name,
		"type": "s3",
		"url":  "https://example.com",
	}

	data, err := json.Marshal(remoteConfig)
	if err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
		Data:    string(data),
	})
}

// ==========================================
// EXAMPLE 4: Function with Multiple Parameters
// ==========================================

// Export updateRemoteConfig - updates configuration for a specific remote
//
//export updateRemoteConfig
func updateRemoteConfig(port C.longlong, remoteName *C.char, configData *C.char) {
	name := C.GoString(remoteName)
	data := C.GoString(configData)

	// Unmarshal the configuration data
	var remoteConfig map[string]interface{}
	if err := json.Unmarshal([]byte(data), &remoteConfig); err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Apply configuration changes
	for key, value := range remoteConfig {
		// TODO: Replace with actual config setter
		// if err := setConfigValue(name, key, value); err != nil {
		//     SendResponseToPort(int64(port), &DartResponse{
		//         Success: false,
		//         Message: err.Error(),
		//     })
		//     return
		// }
		_, _, _ = name, key, value // Suppress unused warnings
	}

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ==========================================
// EXAMPLE 5: Delete Function
// ==========================================

// Export deleteRemote - deletes a remote configuration
//
//export deleteRemote
func deleteRemote(port C.longlong, remoteName *C.char) {
	name := C.GoString(remoteName)

	// TODO: Replace with actual deletion logic
	// if err := deleteConfigRemote(name); err != nil {
	//     SendResponseToPort(int64(port), &DartResponse{
	//         Success: false,
	//         Message: err.Error(),
	//     })
	//     return
	// }
	_ = name // Suppress unused warning

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ==========================================
// UTILITY FUNCTIONS (for pointer creation)
// ==========================================

func ptrBool(v bool) *bool {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func ptrInt(v int) *int {
	return &v
}
