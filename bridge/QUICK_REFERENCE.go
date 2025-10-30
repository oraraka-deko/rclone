// QUICK REFERENCE: Adding Exported Functions to the Dart Bridge
// ================================================================
// This file shows how to quickly add new //export functions following the established pattern

// FILE LOCATION: /bridge/exports.go
// PACKAGE: bridge

// Step 1: Add to cgo preamble (already done in exports.go)
/*
#include "stdint.h"
#include "dart_api/dart_api_dl.h"
#include "dart_api/dart_api_dl.c"
#include "dart_api/dart_native_api.h"
*/

package bridge

import (
	"C"
	"encoding/json"
)

// ================================================================
// TEMPLATE 1: Simple Getter (no parameters, returns data)
// ================================================================

// Export getSimpleData - template for simple read operations
// Simply call this template to create your own
//
//export getSimpleData
func getSimpleData(port C.longlong) {
	// YOUR BUSINESS LOGIC HERE
	result := map[string]interface{}{
		"key":   "value",
		"count": 42,
	}

	// Marshal to JSON
	data, err := json.Marshal(result)
	if err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Send success response
	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
		Data:    string(data),
	})
}

// ================================================================
// TEMPLATE 2: Simple Update (JSON input, no return data)
// ================================================================

//export updateSimpleData
func updateSimpleData(port C.longlong, jsonInput *C.char) {
	// Convert C string to Go string
	input := C.GoString(jsonInput)

	// Parse JSON
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(input), &config); err != nil {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// YOUR BUSINESS LOGIC HERE
	// Apply changes from config
	for key, value := range config {
		// Handle key/value pairs
		_ = key
		_ = value
	}

	// Send success response (no data needed for updates)
	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ================================================================
// TEMPLATE 3: String Parameter Read (get by identifier)
// ================================================================

//export getDataByName
func getDataByName(port C.longlong, name *C.char) {
	// Convert C string to Go string
	nameStr := C.GoString(name)

	// Validate input
	if nameStr == "" {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: "name parameter cannot be empty",
		})
		return
	}

	// YOUR BUSINESS LOGIC HERE
	result := map[string]interface{}{
		"name": nameStr,
		"data": "result for " + nameStr,
	}

	// Marshal and send
	data, err := json.Marshal(result)
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

// ================================================================
// TEMPLATE 4: Multiple String Parameters
// ================================================================

//export processData
func processData(port C.longlong, source *C.char, destination *C.char) {
	// Convert C strings to Go strings
	src := C.GoString(source)
	dst := C.GoString(destination)

	// Validate inputs
	if src == "" || dst == "" {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: "source and destination cannot be empty",
		})
		return
	}

	// YOUR BUSINESS LOGIC HERE
	// Process src -> dst

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ================================================================
// TEMPLATE 5: Delete Operation
// ================================================================

//export deleteItem
func deleteItem(port C.longlong, itemID *C.char) {
	// Convert C string to Go string
	id := C.GoString(itemID)

	// Validate input
	if id == "" {
		SendResponseToPort(int64(port), &DartResponse{
			Success: false,
			Message: "itemID cannot be empty",
		})
		return
	}

	// YOUR BUSINESS LOGIC HERE
	// Delete item with given id

	SendResponseToPort(int64(port), &DartResponse{
		Success: true,
	})
}

// ================================================================
// ADDING A NEW EXPORTED FUNCTION - CHECKLIST
// ================================================================
//
// 1. [ ] Add //export functionName comment IMMEDIATELY before function
// 2. [ ] First parameter is ALWAYS: port C.longlong
// 3. [ ] Additional string parameters are: *C.char (convert with C.GoString)
// 4. [ ] Function body:
//    [ ] Convert all C parameters to Go types
//    [ ] Validate input parameters
//    [ ] Execute business logic
//    [ ] Handle errors by returning early with DartResponse{Success: false, Message: err.Error()}
//    [ ] For success: marshal result to JSON and send DartResponse{Success: true, Data: jsonString}
// 5. [ ] Add documentation comment above the //export line
// 6. [ ] Test calling from Dart code
//
// REMEMBER:
// - NEVER use return statements with values (communication happens through SendResponseToPort)
// - ALWAYS send a response through SendResponseToPort (even on error)
// - ALWAYS convert port to int64(port) when sending response
// - ALWAYS handle errors and send them as DartResponse messages
// - Functions can't be private (must start with uppercase letter in Go)

// ================================================================
// COMMON ERRORS TO AVOID
// ================================================================

// ❌ WRONG: Attempting to return a value
// func getConfig(port C.longlong) string {
//     return "{}"  // DON'T DO THIS!
// }

// ✅ CORRECT: Send response through port
// func getConfig(port C.longlong) {
//     SendResponseToPort(int64(port), &DartResponse{
//         Success: true,
//         Data:    "{}",
//     })
// }

// ❌ WRONG: Not handling errors
// func updateConfig(port C.longlong, cfg *C.char) {
//     data := json.Marshal(C.GoString(cfg))  // What if this errors?
//     SendResponseToPort(int64(port), &DartResponse{...})
// }

// ✅ CORRECT: Check errors immediately
// func updateConfig(port C.longlong, cfg *C.char) {
//     var config map[string]interface{}
//     if err := json.Unmarshal([]byte(C.GoString(cfg)), &config); err != nil {
//         SendResponseToPort(int64(port), &DartResponse{
//             Success: false,
//             Message: err.Error(),
//         })
//         return  // Return early on error
//     }
//     // ... rest of logic
// }

// ❌ WRONG: Forgetting to convert port to int64
// SendResponseToPort(port, ...)  // WRONG - port is C.longlong

// ✅ CORRECT: Convert port properly
// SendResponseToPort(int64(port), ...)

// ================================================================
// CALLING FROM DART
// ================================================================
//
// Future<void> callExportedFunction() async {
//   // Create a ReceivePort to listen for responses
//   final receivePort = ReceivePort();
//
//   // Call the exported Go function with the port's sendPort
//   nativeLib.callYourFunction(receivePort.sendPort.nativePort);
//
//   // Wait for response
//   final response = await receivePort.first;
//   receivePort.close();
//
//   // Parse JSON response
//   final jsonResponse = jsonDecode(response);
//   if (jsonResponse['success']) {
//     final data = jsonDecode(jsonResponse['data']);
//     print('Success: $data');
//   } else {
//     print('Error: ${jsonResponse['message']}');
//   }
// }
