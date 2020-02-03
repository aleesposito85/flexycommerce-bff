package util

import (
	"fmt"
	"log"
	"runtime"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Println("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Println("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Println("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Println("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
