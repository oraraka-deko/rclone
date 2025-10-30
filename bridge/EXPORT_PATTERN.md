# Dart Bridge Export Functions - Documentation

## Overview

This document explains the pattern used to expose Go functions to Dart via the cgo/FFI bridge in the rclone project. The pattern allows Dart code to call Go functions and receive JSON-formatted responses through a port-based IPC mechanism.

## Key Pattern

### Basic Structure

```go
//export functionName
func functionName(port C.longlong, param1 *C.char, param2 *C.char) {
    // 1. Convert C parameters to Go
    p1 := C.GoString(param1)
    p2 := C.GoString(param2)
    
    // 2. Perform business logic
    result, err := doSomething(p1, p2)
    
    // 3. Handle errors early
    if err != nil {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }
    
    // 4. Marshal result to JSON
    data, err := json.Marshal(result)
    if err != nil {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }
    
    // 5. Send success response
    SendResponseToPort(int64(port), &DartResponse{
        Success: true,
        Data:    string(data),
    })
}
```

## Essential Components

### 1. Export Comment
```go
//export functionName
```
- **Must be immediately before function declaration**
- Tells cgo to export this function as a C symbol
- Required for Dart FFI to discover the function

### 2. Function Signature Requirements

#### Port Parameter (Required)
```go
func myFunc(port C.longlong, ...) {
    // port is always the first parameter
    // Convert to int64 when sending response
}
```
- **Always first parameter**
- Type: `C.longlong`
- Used to identify which Dart Isolate/Port should receive the response
- Convert with: `int64(port)`

#### String Parameters
```go
func myFunc(port C.longlong, name *C.char, path *C.char) {
    nameStr := C.GoString(name)
    pathStr := C.GoString(path)
}
```
- Type: `*C.char`
- Convert with: `C.GoString(cParam)`
- Required for string input from Dart

#### Other C Types
- `C.int` - integer values
- `C.double` - floating point values
- `C.bool` - boolean values
- For complex types, use JSON strings via `*C.char`

### 3. Response Mechanism

All functions must communicate responses through `SendResponseToPort`:

```go
// Success response with data
SendResponseToPort(int64(port), &DartResponse{
    Success: true,
    Data:    jsonString,
})

// Error response
SendResponseToPort(int64(port), &DartResponse{
    Success: false,
    Message: "error description",
})

// Success response without data
SendResponseToPort(int64(port), &DartResponse{
    Success: true,
})
```

**DartResponse Structure:**
```go
type DartResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    string `json:"data"`
}
```

### 4. JSON Marshaling Pattern

For input parameters with complex types:
```go
// Receive JSON string from Dart
var config map[string]interface{}
if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
    SendResponseToPort(int64(port), &DartResponse{
        Success: false,
        Message: err.Error(),
    })
    return
}
```

For output responses:
```go
data, err := json.Marshal(result)
if err != nil {
    // Handle error
}
SendResponseToPort(int64(port), &DartResponse{
    Success: true,
    Data:    string(data),
})
```

## Complete Examples

### Example 1: Simple Read (No Parameters)

```go
// Export getStatus - retrieves current status
//export getStatus
func getStatus(port C.longlong) {
    status := map[string]interface{}{
        "running": true,
        "version": "1.0",
    }
    
    data, err := json.Marshal(status)
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
```

### Example 2: Single String Parameter

```go
// Export getFileInfo - retrieves information about a file
//export getFileInfo
func getFileInfo(port C.longlong, filePath *C.char) {
    path := C.GoString(filePath)
    
    info, err := os.Stat(path)
    if err != nil {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }
    
    fileInfo := map[string]interface{}{
        "name":  info.Name(),
        "size":  info.Size(),
        "isDir": info.IsDir(),
    }
    
    data, err := json.Marshal(fileInfo)
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
```

### Example 3: JSON Input/Output

```go
// Export updateSettings - updates application settings from JSON
//export updateSettings
func updateSettings(port C.longlong, settingsJson *C.char) {
    jsonStr := C.GoString(settingsJson)
    
    // Parse incoming JSON
    var settings map[string]interface{}
    if err := json.Unmarshal([]byte(jsonStr), &settings); err != nil {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: "invalid JSON: " + err.Error(),
        })
        return
    }
    
    // Apply settings
    for key, value := range settings {
        if err := applySetting(key, value); err != nil {
            SendResponseToPort(int64(port), &DartResponse{
                Success: false,
                Message: fmt.Sprintf("failed to apply %s: %v", key, err),
            })
            return
        }
    }
    
    SendResponseToPort(int64(port), &DartResponse{
        Success: true,
    })
}
```

### Example 4: Multiple Parameters

```go
// Export copyFile - copies a file from source to destination
//export copyFile
func copyFile(port C.longlong, src *C.char, dst *C.char) {
    srcPath := C.GoString(src)
    dstPath := C.GoString(dst)
    
    // Validate paths
    if srcPath == "" || dstPath == "" {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: "source and destination paths cannot be empty",
        })
        return
    }
    
    // Perform copy operation
    if err := copyFileImpl(srcPath, dstPath); err != nil {
        SendResponseToPort(int64(port), &DartResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }
    
    SendResponseToPort(int64(port), &DartResponse{
        Success: true,
    })
}
```

## Best Practices

### 1. Error Handling
- **Always check errors** and return early
- **Never panic** in exported functions
- **Always send a response**, even on error
- Use descriptive error messages

### 2. Resource Management
- **Keep functions short** - avoid long-running operations in exported functions
- **Use goroutines** for heavy work if needed
- **Clean up resources** properly

### 3. Type Safety
- Use **struct types** for complex JSON objects instead of `map[string]interface{}`
- Define **JSON tags** on struct fields
- Use **pointer fields** when values are optional

Example:
```go
type Settings struct {
    Port    int     `json:"port"`
    Host    string  `json:"host"`
    Timeout *int    `json:"timeout,omitempty"` // Optional field
}
```

### 4. Naming Convention
- Use **camelCase** for function names
- Use **verb + noun** pattern: `getConfig`, `updateSettings`, `deleteRemote`
- Keep names **short but descriptive**

### 5. Documentation
- Add **godoc comments** above exported functions
- Explain **parameters and return types** (via JSON)
- Document **error conditions**

Example:
```go
// getConfig retrieves the application configuration.
//
// Returns JSON object with current config:
//   {
//     "port": 8080,
//     "host": "localhost",
//     "debug": false
//   }
//
// On error, returns JSON:
//   {"success": false, "message": "error description"}
//
//export getConfig
func getConfig(port C.longlong) {
    // implementation
}
```

## Integration with Dart

When calling from Dart:

```dart
// Import the FFI library
import 'dart:ffi' as ffi;

// Define the C function signature
typedef GetConfigNative = ffi.Void Function(ffi.Int64);
typedef GetConfigDart = void Function(int);

// Call the function
dylib.lookupFunction<GetConfigNative, GetConfigDart>('getConfig')(port);
```

## File Location

Exported functions should be placed in:
- **Location**: `/bridge/exports.go` (or similar)
- **Package**: `bridge`
- **Must include** the cgo preamble with Dart API headers

## Common Patterns Summary

| Pattern | Use Case | Example |
|---------|----------|---------|
| `getXXX(port)` | Read-only operations | `getStatus`, `getRemoteList` |
| `patchXXX(port, json)` | Update operations | `patchConfig`, `updateSettings` |
| `deleteXXX(port, name)` | Delete operations | `deleteRemote`, `deleteFile` |
| `setXXX(port, name, value)` | Single value set | `setOption`, `setValue` |
| `createXXX(port, json)` | Create operations | `createRemote`, `createUser` |

## Troubleshooting

### Function not exported
- **Check**: Is `//export functionName` immediately above the function?
- **Check**: Is function name capitalized? (exported functions must be)
- **Check**: Is the package built with cgo enabled?

### C types undefined
- **Check**: Are cgo headers included in the preamble?
- **Check**: Is `import "C"` present after the preamble?

### Response not received in Dart
- **Check**: Is port correctly converted to `int64(port)`?
- **Check**: Is `SendResponseToPort` being called?
- **Check**: Is response JSON valid?

## See Also

- Example implementation: `bridge/exports.go`
- Bridge support functions: `bridge/lib.go`
- Dart FFI documentation: https://dart.dev/guides/libraries/c-interop
