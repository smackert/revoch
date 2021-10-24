package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

//var inputFile string = ""
var bufferSize int64 = 512

//var mapFile string = ""

func main() {

	inputFile := flag.String("f", "", "Path to the file you wish to analyze.")
	bufferSize = *flag.Int64("s", 512, "Define sector size in bytes.")
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
	ReadFile(*inputFile)
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

func ReadFile(filePath string) {

	file, err := os.Open(filePath)
	errCheck(err)
	defer file.Close()

	f, err := file.Stat()
	errCheck(err)

	fmt.Println("File:", f.Name(), "\nSize:", f.Size(), "\nChunks:", (f.Size() / bufferSize))

	var offset int64
	//writeChan := make(chan map[int64]int, 5000)
	chunk := 0
	mapFile := makeMapFile(filePath)
	bw := bufio.NewWriterSize(mapFile, 1048576)
	/*
		    for offset
			inputdata = read up to 1mb
			for len(inputdata) / buffersize
				go ent(inputdata, bufoffset)
				bufoffset += buffesize
			offset += len(inputdata)
	*/

	for offset = 0; (offset + bufferSize) < f.Size(); offset += bufferSize { // this cuts the tail off. fix that l8r
		calcEntropy(*file, offset, bw)
		chunk += 1
		fmt.Printf("\r%f%% Complete..", (float32(float32(chunk)/(float32(f.Size())/float32(512))) * float32(100)))
	}
	//bw.Flush()
	mapFile.Close()
	fmt.Println("\nFinished!")
}

func makeMapFile(fp string) *os.File {

	f, err := os.OpenFile(fp+".map.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	errCheck(err)
	//defer f.Close()
	return f

}

func calcEntropy(file os.File, offset int64, bw *bufio.Writer) {
	buf := make([]byte, bufferSize)
	_, err := file.ReadAt(buf, offset)
	errCheck(err)
	se := make(map[int64]int)
	se[offset] = Shannon(buf[:])
	writeMapBuf(se, bw)
}

func writeMapBuf(se map[int64]int, bw *bufio.Writer) {
	/*
		    var writeBuffer = make(map[int64]int)
			for k, v := range data {
				writeBuffer[k] = v
			}
	*/
	for k, v := range se {
		_, err := bw.WriteString(strconv.FormatInt(k, 10) + "," + strconv.Itoa(v) + "\n")
		errCheck(err)
	}

}
