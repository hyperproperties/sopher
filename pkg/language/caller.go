package language

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"runtime"
)

func Caller() uint64 {
	callers := make([]uintptr, 10)
	runtime.Callers(3, callers)
	frames := runtime.CallersFrames(callers)

	frame, _ := frames.Next()

	// Combine PC, function name, and file/line into a single identifier
	uniqueData := fmt.Sprintf("%s:%s:%d:%x", frame.Function, frame.File, frame.Line, frame.PC)

	hash := sha256.Sum256([]byte(uniqueData))
	uniqueID := binary.LittleEndian.Uint64(hash[:8])

	return uint64(uniqueID)
}