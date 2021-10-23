package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	inputFile := flag.String("f", "", "Path to the file you wish to analyze.")
	inputBufferSize := flag.Int64("s", 512, "Define sector size in bytes.")
	flag.Parse()
	if *inputFile == "" {
		fmt.Println("A file name is required. Exiting..")
		os.Exit(0)
	}

	//const binaryExe = "test_data/binary.exe"
	//const text = "test_data/sample1.txt"
	//const volImg = "test_data/test_volume.img" // 599mb (628,097,024 bytes) contains 100mb and 50mb encrypted volumes
	//const movie = "test_data/movie.mkv"        // 4.37 GB (4,695,146,496 bytes)
	//const encVol = "test_data/encrypted_volume" // 600mb (629,145,600 bytes)
	//const loremipsum = "test_data/loremipsum.txt"
	//const bufferSize int64 = 837463 // 1mb (1000000)

	startTime := time.Now()
	ReadFile(*inputFile, *inputBufferSize)
	fmt.Printf("Completed in %s", time.Since((startTime).Truncate(time.Second)))

	//ReadFile(volImg, bufferSize)
	//eadFile(movie, bufferSize)
	//ReadFile(encVol, bufferSize)
	//ReadFile(binaryExe, bufferSize)
	//ReadFile(loremipsum, bufferSize)

	//fmt.Println("Done")
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(filePath string, bufferSize int64) {

	file, err := os.Open(filePath)
	errCheck(err)
	defer file.Close()

	f, err := file.Stat()
	errCheck(err)

	fmt.Println("File:", f.Name(), "\nSize:", f.Size(), "\nChunks:", (f.Size() / bufferSize))

	buffer := make([]byte, bufferSize)

	var offset int64

	chunk := 0
	mapFile := makeMapFile(filePath)
	for offset = 0; (offset + bufferSize) < f.Size(); offset += bufferSize { // this cuts the tail off. fix that l8r
		_, err := file.ReadAt(buffer, offset)
		errCheck(err)
		sectorEntropy := make(map[int64]int)
		sectorEntropy[offset] = Shannon(buffer[:])
		writeMapBuf(mapFile, sectorEntropy)
		chunk += 1
		fmt.Printf("\r%f%% Complete..", (float32(float32(chunk)/(float32(f.Size())/float32(512))) * float32(100)))
	}

	fmt.Println("\nFinished!")
}

func makeMapFile(fp string) string {

	f, err := os.Create(fp + ".map.txt")
	errCheck(err)
	defer f.Close()
	return f.Name()

}

func writeMapBuf(mf string, data map[int64]int) {
	/*
		    var writeBuffer = make(map[int64]int)
			for k, v := range data {
				writeBuffer[k] = v
			}
	*/
	//if len(writeBuffer) > 1000 {
	f, err := os.OpenFile(mf, os.O_APPEND|os.O_WRONLY, 0600)
	errCheck(err)
	for k, v := range data {
		_, err = f.WriteString(strconv.FormatInt(k, 10) + "," + strconv.Itoa(v) + "\n")
		errCheck(err)
	}
	defer f.Close()
	//}
}
