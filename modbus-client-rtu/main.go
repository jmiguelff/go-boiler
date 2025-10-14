package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
)

// ---- PT constants (from your meter config) ----
const (
	PT1 = 30000.0 // primary volts
	PT2 = 180.0   // secondary volts
)

func bytesToFloat32(b []byte) float32 {
	bits := binary.BigEndian.Uint32(b)
	return math.Float32frombits(bits)
}

// secondaryToPrimary applies your fixed PT ratio.
func secondaryToPrimary(vSecondary float32) float64 {
	return float64(vSecondary) * (PT1 / PT2)
}

func main() {
	// Define flags
	floatFlag := flag.Bool("float", false, "Convert response to float32 values")
	primaryFlag := flag.Bool("primary", false, "Read from address 16394 (2 registers) and convert to primary voltage using PT ratio")
	flag.Parse()

	// Handle --primary flag separately
	if *primaryFlag {
		args := flag.Args()
		if len(args) < 2 {
			fmt.Println("Usage: go run main.go --primary <serial_port> <baud_rate>")
			fmt.Println("Example: go run main.go --primary /dev/ttyUSB0 9600")
			os.Exit(1)
		}

		serialPort := args[0]
		baudRate, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Invalid baud rate:", err)
		}

		// Configure RTU handler
		handler := modbus.NewRTUClientHandler(serialPort)
		handler.BaudRate = baudRate
		handler.DataBits = 8
		handler.Parity = "N"
		handler.StopBits = 1
		handler.SlaveId = 1
		handler.Timeout = 5 * time.Second

		log.Printf("Connecting to %s @ %d baud...", serialPort, baudRate)

		// Connect
		err = handler.Connect()
		if err != nil {
			log.Fatalln("Error connecting:", err)
		}
		defer handler.Close()

		client := modbus.NewClient(handler)

		// Read from address 16394 (2 registers = 4 bytes for float32)
		const address = 16394
		log.Printf("Reading 2 registers from address %d (logical: 4%04d) for primary voltage calculation", address, address+1)
		res, err := client.ReadHoldingRegisters(uint16(address), 2)
		if err != nil {
			log.Fatalln("Error reading:", err)
		}

		// Convert to float32
		secondaryVoltage := bytesToFloat32(res)
		primaryVoltage := secondaryToPrimary(secondaryVoltage)

		fmt.Printf("\nSecondary Voltage: %.6f V\n", secondaryVoltage)
		fmt.Printf("Primary Voltage: %.6f V\n", primaryVoltage)
		fmt.Printf("PT Ratio: %.1f:%.1f\n", PT1, PT2)

		return
	}

	args := flag.Args()
	if len(args) < 4 {
		fmt.Println("Usage: go run main.go [--float] <serial_port> <baud_rate> <start_address> <quantity>")
		fmt.Println("       go run main.go --primary <serial_port> <baud_rate>")
		fmt.Println("Example: go run main.go /dev/ttyUSB0 9600 0 10")
		fmt.Println("Example: go run main.go --float /dev/ttyUSB0 9600 40001 5")
		fmt.Println("Example: go run main.go --primary /dev/ttyUSB0 9600")
		fmt.Println("\nNote: start_address can be protocol address (0-based) or logical address (40001+)")
		os.Exit(1)
	}

	serialPort := args[0]
	baudRate, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln("Invalid baud rate:", err)
	}

	startAddress, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalln("Invalid start address:", err)
	}

	// Convert logical address to protocol address if needed
	if startAddress >= 40001 {
		startAddress = startAddress - 40001
		log.Printf("Converted logical address to protocol address: %d", startAddress)
	}

	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalln("Invalid quantity:", err)
	}

	// Configure RTU handler
	handler := modbus.NewRTUClientHandler(serialPort)
	handler.BaudRate = baudRate
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	log.Printf("Connecting to %s @ %d baud...", serialPort, baudRate)

	// Connect
	err = handler.Connect()
	if err != nil {
		log.Fatalln("Error connecting:", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	// Read holding registers
	log.Printf("Reading %d registers starting at address %d (logical: 4%04d)", quantity, startAddress, startAddress+1)
	res, err := client.ReadHoldingRegisters(uint16(startAddress), uint16(quantity))
	if err != nil {
		log.Fatalln("Error reading:", err)
	}

	fmt.Printf("\nResponse: %v\n", res)
	fmt.Printf("Response (hex): %X\n", res)

	// Convert to float if flag is set
	if *floatFlag {
		fmt.Println("\nFloat32 values:")
		// Each float32 requires 2 registers (4 bytes)
		if len(res)%4 != 0 {
			log.Println("Warning: Response length is not a multiple of 4. Some data may not convert properly.")
		}

		for i := 0; i+3 < len(res); i += 4 {
			floatVal := bytesToFloat32(res[i : i+4])
			registerPair := (i / 4) * 2
			fmt.Printf("Registers %d-%d (4%04d-4%04d): %.6f\n",
				startAddress+registerPair,
				startAddress+registerPair+1,
				startAddress+registerPair+1,
				startAddress+registerPair+2,
				floatVal)
		}
	} else {
		// Display as 16-bit integers
		fmt.Println("\nUint16 values:")
		for i := 0; i+1 < len(res); i += 2 {
			val := binary.BigEndian.Uint16(res[i : i+2])
			register := i / 2
			fmt.Printf("Register %d (4%04d): %d (0x%04X)\n",
				startAddress+register,
				startAddress+register+1,
				val, val)
		}
	}
}
