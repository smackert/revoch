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
var chunkSize int64 = 512
var readBufferSize int
var writeBufferSize int

func main() {

	inputFile := flag.String("f", "", "Path to the file you wish to analyze.")
	chunkSize = *flag.Int64("s", 512, "Define sector size in bytes.")
	readBufferSize = *flag.Int("rb", 40, "Read buffer size in MB")
	writeBufferSize = *flag.Int("wb", 40, "Write buffer size in MB")
	flag.Parse()
	if *inputFile == "" {
		fmt.Println("A file name is required. Exiting..")
		os.Exit(0)
	}

	startTime := time.Now()
	ReadFile(*inputFile)
	fmt.Printf("Completed in %s", time.Since((startTime).Truncate(time.Second)))

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

	filestat, err := file.Stat()
	errCheck(err)

	totalChunks := filestat.Size() / chunkSize
	fmt.Println("File:", filestat.Name(), "\nSize:", filestat.Size(), "\nChunks:", totalChunks)

	var offset int64
	chunk := 0
	mapFile := makeMapFile(filePath)
	br := bufio.NewReaderSize(file, readBufferSize*1048576)     //  1048576 = 1mb
	bw := bufio.NewWriterSize(mapFile, writeBufferSize*1048576) //  1048576 = 1mb
	progress := 0
	for offset = 0; (offset + chunkSize) < filestat.Size(); offset += chunkSize { // this might cut the tail off. check that l8r
		ent := calcEntropy(br, chunkSize)
		writeMap(offset, ent, bw)
		chunk += 1
		curProgress := int((float32(chunk) / float32(totalChunks)) * float32(100))
		if curProgress >= progress+1 {
			fmt.Printf("\r%v%% complete...", curProgress)
			progress = curProgress
		}
	}

	bw.Flush()
	mapFile.Close()
	fmt.Println("Finished!             ")
}

func makeMapFile(fp string) *os.File {

	f, err := os.OpenFile(fp+".map.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	errCheck(err)
	//defer f.Close()
	return f

}

func calcEntropy(br *bufio.Reader, chunkSize int64) int {
	buf := make([]byte, chunkSize)
	_, err := br.Read(buf)
	errCheck(err)
	return Shannon(buf[:])
}

func writeMap(offset int64, ent int, bw *bufio.Writer) {

	_, err := bw.WriteString(strconv.FormatInt(offset, 10) + "," + strconv.Itoa(ent) + "\n")
	errCheck(err)

}
