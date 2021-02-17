package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Open file
	f, err := os.Open("input.bin")
	check(err)

	// Close file descriptor at the end
	defer f.Close()

	// Get size of file
	stats, err := f.Stat()
	check(err)

	// Create slice to hold bytes
	var fileSize = stats.Size()
	bytesArr := make([]byte, fileSize)

	// Read bytes from file
	r := bufio.NewReader(f)
	_, err = r.Read(bytesArr)
	check(err)

	// Create output file
	out, err := os.Create("output.bin")
	check(err)

	// Write first packet
	w := bufio.NewWriter(out)
	_, err = w.WriteString("//Packet 1\n")
	check(err)
	_, err = w.Write(bytesArr[:20])
	check(err)
	_, err = w.WriteString("\n")
	check(err)
	w.Flush()

	fmt.Println(len(bytesArr[20:]))

	var count = 0

	for i := 0; i < len(bytesArr[20:]); i++ {
		err = w.WriteByte(bytesArr[i+20])
		check(err)

		if i%16 == 0 && i != 0 {
			fmt.Println(i)
			_, err = w.WriteString("\n")
			check(err)
		}

		if i%234 == 0 && i != 0 {
			count = count + 1
			fmt.Println("Packet " + strconv.Itoa(count+1) + "\n")
			_, err = w.WriteString("Packet " + strconv.Itoa(count+1) + "\n")
			check(err)
		}

		w.Flush()
	}
}
