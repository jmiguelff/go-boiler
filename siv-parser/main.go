package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	BlockSize   int = 117
	NbrOfBlocks int = 64
)

type DigitalVars struct {
	index     int
	name      string
	errorCode string
	category  string
	value     bool
}

type sivEvent struct {
	eqpIdx           int
	index            int
	name             string
	errorCode        string
	nbrrOfOccurences int
	occurrences      []sivEventOccurence
}

type sivEventOccurence struct {
	activeVars []DigitalVars
	timeStamp  string
}

var ccuDigitalVars = [64]DigitalVars{
	{index: 1, name: "CAN INV fail", errorCode: "E38001", category: "restrictive", value: false},
	{index: 2, name: "CAN LVPS fail", errorCode: "E38002", category: "restrictive", value: false},
	{index: 3, name: "K901 fail", errorCode: "E38003", category: "restrictive", value: false},
	{index: 4, name: "K902 fail", errorCode: "E38004", category: "restrictive", value: false},
	{index: 5, name: "Vbat max redun", errorCode: "E38005", category: "restrictive", value: false},
	{index: 6, name: "FAN fail", errorCode: "E38006", category: "restrictive", value: false},
	{index: 7, name: "Battery fail", errorCode: "E38007", category: "restrictive", value: false},
	{index: 8, name: "CCU_E_SCR_failed", errorCode: "E10170", category: "restrictive", value: false},
	{index: 9, name: "SIVPulsing detected", errorCode: "E0", category: "restrictive", value: false},
	{index: 10, name: "Req.Reset other HBU", errorCode: "E0", category: "restrictive", value: false},
	{index: 11, name: "Reset this HBU/INV", errorCode: "E0", category: "restrictive", value: false},
	{index: 12, name: "digital event 12", errorCode: "E0"},
	{index: 13, name: "digital event 13", errorCode: "E0"},
	{index: 14, name: "digital event 14", errorCode: "E0"},
	{index: 15, name: "digital event 15", errorCode: "E0"},
	{index: 16, name: "digital event 16", errorCode: "E0"},
	{index: 17, name: "Vin max. INV", errorCode: "E38008", category: "non-restrictive", value: false},
	{index: 18, name: "Vin min. INV", errorCode: "E38009", category: "non-restrictive", value: false},
	{index: 19, name: "digital event 19", errorCode: "E0"},
	{index: 20, name: "digital event 20", errorCode: "E0"},
	{index: 21, name: "digital event 21", errorCode: "E0"},
	{index: 22, name: "digital event 22", errorCode: "E0"},
	{index: 23, name: "digital event 23", errorCode: "E0"},
	{index: 24, name: "digital event 24", errorCode: "E0"},
	{index: 25, name: "digital event 25", errorCode: "E0"},
	{index: 26, name: "digital event 26", errorCode: "E0"},
	{index: 27, name: "digital event 27", errorCode: "E0"},
	{index: 28, name: "digital event 28", errorCode: "E0"},
	{index: 29, name: "digital event 29", errorCode: "E0"},
	{index: 30, name: "digital event 30", errorCode: "E0"},
	{index: 31, name: "digital event 31", errorCode: "E0"},
	{index: 32, name: "digital event 32", errorCode: "E0"},
	{index: 33, name: "inverter enable", errorCode: "E0"},
	{index: 34, name: "lvps enable", errorCode: "E0", category: "other", value: false},
	{index: 35, name: "k901 drive", errorCode: "E0", category: "other", value: false},
	{index: 36, name: "k901 feedback", errorCode: "E0", category: "other", value: false},
	{index: 37, name: "k902 drive", errorCode: "E0", category: "other", value: false},
	{index: 38, name: "k902 feedback", errorCode: "E0", category: "other", value: false},
	{index: 39, name: "k903 drive", errorCode: "E0", category: "other", value: false},
	{index: 40, name: "digital event 40", errorCode: "E0", category: "other", value: false},
	{index: 41, name: "INV OK (K1)", errorCode: "E0", category: "other", value: false},
	{index: 42, name: "LVPS fault (K2)", errorCode: "E0", category: "other", value: false},
	{index: 43, name: "Vout INV OK (K3)", errorCode: "E0", category: "other", value: false},
	{index: 44, name: "HV available (K5)", errorCode: "E0", category: "other", value: false},
	{index: 45, name: "board location", errorCode: "E0", category: "other", value: false},
	{index: 46, name: "V Bat minimum", errorCode: "E0", category: "other", value: false},
	{index: 47, name: "CCU battery low", errorCode: "E37102", category: "other", value: false},
	{index: 48, name: "clock stopped", errorCode: "E38010", category: "other", value: false},
	{index: 49, name: "train stable", errorCode: "E0", category: "no trigger", value: false},
	{index: 50, name: "external stop", errorCode: "E0", category: "no trigger", value: false},
	{index: 51, name: "Fan", errorCode: "E0", category: "no trigger", value: false},
	{index: 52, name: "K902 (cou.) stopped", errorCode: "E0", category: "no trigger", value: false},
	{index: 53, name: "K902 allowed (ext.)", errorCode: "E0", category: "no trigger", value: false},
	{index: 54, name: "K902 request (int.)", errorCode: "E0", category: "no trigger", value: false},
	{index: 55, name: "digital event 55", errorCode: "E0"},
	{index: 56, name: "digital event 56", errorCode: "E0"},
	{index: 57, name: "digital event 57", errorCode: "E0"},
	{index: 58, name: "digital event 58", errorCode: "E0"},
	{index: 59, name: "digital event 59", errorCode: "E0"},
	{index: 60, name: "digital event 60", errorCode: "E0"},
	{index: 61, name: "digital event 61", errorCode: "E0"},
	{index: 62, name: "digital event 62", errorCode: "E0"},
	{index: 63, name: "digital event 63", errorCode: "E0"},
	{index: 64, name: "digital event 64", errorCode: "E0"},
}

var invDigitalVars = [64]DigitalVars{
	{index: 1, name: "Vce locked", errorCode: "E33115", category: "restrictive", value: false},
	{index: 2, name: "Idyn locked", errorCode: "E33121", category: "restrictive", value: false},
	{index: 3, name: "Vacmax locked", errorCode: "E33125", category: "restrictive", value: false},
	{index: 4, name: "Vphase locked", errorCode: "E33151", category: "restrictive", value: false},
	{index: 5, name: "Ilim locked", errorCode: "E33116", category: "restrictive", value: false},
	{index: 6, name: "com_err locked", errorCode: "E37130", category: "restrictive", value: false},
	{index: 7, name: "aux supply error", errorCode: "E37131", category: "restrictive", value: false},
	{index: 8, name: "Temp max A301", errorCode: "E33127", category: "restrictive", value: false},
	{index: 9, name: "Temp max A302", errorCode: "E33129", category: "restrictive", value: false},
	{index: 10, name: "Temp max T301", errorCode: "E33128", category: "restrictive", value: false},
	{index: 11, name: "Temp max T302", errorCode: "E33130", category: "restrictive", value: false},
	{index: 12, name: "Twarn", errorCode: "E33162", category: "restrictive", value: false},
	{index: 13, name: "Temp locked", errorCode: "E33161", category: "restrictive", value: false},
	{index: 14, name: "digital event 14", errorCode: "E0"},
	{index: 15, name: "digital event 15", errorCode: "E0"},
	{index: 16, name: "digital event 16", errorCode: "E0"},
	{index: 17, name: "VCE-Mod1 Branch1", errorCode: "E33103", category: "non-restrictive", value: false},
	{index: 18, name: "VCE-Mod1 Branch2", errorCode: "E33104", category: "non-restrictive", value: false},
	{index: 19, name: "VCE-Mod1 Branch3", errorCode: "E33105", category: "non-restrictive", value: false},
	{index: 20, name: "VCE-Mod2 Branch1", errorCode: "E33106", category: "non-restrictive", value: false},
	{index: 21, name: "VCE-Mod2 Branch2", errorCode: "E33107", category: "non-restrictive", value: false},
	{index: 22, name: "VCE-Mod2 Branch3", errorCode: "E33108", category: "non-restrictive", value: false},
	{index: 23, name: "Idyn HW", errorCode: "E33160", category: "non-restrictive", value: false},
	{index: 24, name: "Idyn SW 1", errorCode: "E33117", category: "non-restrictive", value: false},
	{index: 25, name: "Idyn SW 2", errorCode: "E33118", category: "non-restrictive", value: false},
	{index: 26, name: "Vac max HW", errorCode: "E33124", category: "non-restrictive", value: false},
	{index: 27, name: "Vac max SW", errorCode: "E33124", category: "non-restrictive", value: false},
	{index: 28, name: "Vdc max HW", errorCode: "E33122", category: "non-restrictive", value: false},
	{index: 29, name: "Vdc max SW", errorCode: "E33122", category: "non-restrictive", value: false},
	{index: 30, name: "Vphase", errorCode: "E33151", category: "non-restrictive", value: false},
	{index: 31, name: "Startup error", errorCode: "E33152", category: "non-restrictive", value: false},
	{index: 32, name: "Vdc unsym", errorCode: "E33163", category: "non-restrictive", value: false},
	{index: 33, name: "Vdc2 max SW", errorCode: "E33122", category: "non-restrictive", value: false},
	{index: 34, name: "digital event 34", errorCode: "E0"},
	{index: 35, name: "digital event 35", errorCode: "E0"},
	{index: 36, name: "digital event 36", errorCode: "E0"},
	{index: 37, name: "digital event 37", errorCode: "E0"},
	{index: 38, name: "digital event 38", errorCode: "E0"},
	{index: 39, name: "digital event 39", errorCode: "E0"},
	{index: 40, name: "digital event 40", errorCode: "E0"},
	{index: 41, name: "digital event 41", errorCode: "E0"},
	{index: 42, name: "digital event 42", errorCode: "E0"},
	{index: 43, name: "digital event 43", errorCode: "E0"},
	{index: 44, name: "digital event 44", errorCode: "E0"},
	{index: 45, name: "digital event 45", errorCode: "E0"},
	{index: 46, name: "digital event 46", errorCode: "E0"},
	{index: 47, name: "digital event 47", errorCode: "E0"},
	{index: 48, name: "digital event 48", errorCode: "E0"},
	{index: 49, name: "Vac min", errorCode: "E33100", category: "other", value: false},
	{index: 50, name: "Vdc min", errorCode: "E33100", category: "other", value: false},
	{index: 51, name: "Ilim", errorCode: "E33100", category: "other", value: false},
	{index: 52, name: "com error", errorCode: "E37130", category: "other", value: false},
	{index: 53, name: "Bat nvdiag low", errorCode: "E37132", category: "other", value: false},
	{index: 54, name: "Vdc min2", errorCode: "E33100", category: "other", value: false},
	{index: 55, name: "Run", errorCode: "E0", category: "other", value: false},
	{index: 56, name: "Ready", errorCode: "E0", category: "other", value: false},
	{index: 57, name: "Error", errorCode: "E0", category: "other", value: false},
	{index: 58, name: "Enable", errorCode: "E0", category: "other", value: false},
	{index: 59, name: "digital event 59", errorCode: "E0"},
	{index: 60, name: "digital event 60", errorCode: "E0"},
	{index: 61, name: "digital event 61", errorCode: "E0"},
	{index: 62, name: "digital event 62", errorCode: "E0"},
	{index: 63, name: "digital event 63", errorCode: "E0"},
	{index: 64, name: "digital event 64", errorCode: "E0"},
}

var lvpsDigitalVars = [64]DigitalVars{
	{index: 1, name: "digital event 1", errorCode: "E0"},
	{index: 2, name: "V-BC-max locked", errorCode: "E34125", category: "restrictive", value: false},
	{index: 3, name: "Idyn BC-err locked", errorCode: "E34121", category: "restrictive", value: false},
	{index: 4, name: "Therm1", errorCode: "E0", category: "restrictive", value: false},
	{index: 5, name: "Temp BC", errorCode: "E34127", category: "restrictive", value: false},
	{index: 6, name: "Therm3", errorCode: "E0", category: "restrictive", value: false},
	{index: 7, name: "Therm4", errorCode: "E0", category: "restrictive", value: false},
	{index: 8, name: "Temp BC transf", errorCode: "E34134", category: "restrictive", value: false},
	{index: 9, name: "DR-err-BC locked", errorCode: "E34115", category: "restrictive", value: false},
	{index: 10, name: "15V OK error", errorCode: "E37141", category: "restrictive", value: false},
	{index: 11, name: "CAN LVPS fail", errorCode: "E37145", category: "restrictive", value: false},
	{index: 12, name: "Selftest-BC error", errorCode: "E34135", category: "restrictive", value: false},
	{index: 13, name: "Shutdown V-compare", errorCode: "E34125", category: "restrictive", value: false},
	{index: 14, name: "BC CU PLD fail", errorCode: "E37146", category: "restrictive", value: false},
	{index: 15, name: "Shutdown Charge err", errorCode: "E60000", category: "restrictive", value: false},
	{index: 16, name: "V-BC range err.", errorCode: "E34153", category: "restrictive", value: false},
	{index: 17, name: "Voltage charge err.", errorCode: "E34154", category: "restrictive", value: false},
	{index: 18, name: "Low batt. ch. curr.", errorCode: "E34155", category: "restrictive", value: false},
	{index: 19, name: "BC CU batt low", errorCode: "E37142", category: "restrictive", value: false},
	{index: 20, name: "digital event 20", errorCode: "E0"},
	{index: 21, name: "digital event 21", errorCode: "E0"},
	{index: 22, name: "digital event 22", errorCode: "E0"},
	{index: 23, name: "digital event 23", errorCode: "E0"},
	{index: 24, name: "digital event 24", errorCode: "E0"},
	{index: 25, name: "digital event 25", errorCode: "E0"},
	{index: 26, name: "digital event 26", errorCode: "E0"},
	{index: 27, name: "digital event 27", errorCode: "E0"},
	{index: 28, name: "digital event 28", errorCode: "E0"},
	{index: 29, name: "digital event 29", errorCode: "E0"},
	{index: 30, name: "digital event 30", errorCode: "E0"},
	{index: 31, name: "digital event 31", errorCode: "E0"},
	{index: 32, name: "digital event 32", errorCode: "E0"},
	{index: 33, name: "DR1-error BC", errorCode: "E34103", category: "non-restrictive", value: false},
	{index: 34, name: "DR2-error BC", errorCode: "E34106", category: "non-restrictive", value: false},
	{index: 35, name: "DR3-error BC", errorCode: "E34109", category: "non-restrictive", value: false},
	{index: 36, name: "Idyn BC error", errorCode: "E34117", category: "non-restrictive", value: false},
	{index: 37, name: "Current high SW", errorCode: "E34167", category: "non-restrictive", value: false},
	{index: 38, name: "V-BC-max", errorCode: "E34124", category: "non-restrictive", value: false},
	{index: 39, name: "V_BC too high SW", errorCode: "E34168", category: "non-restrictive", value: false},
	{index: 40, name: "Vin LVPS max", errorCode: "E34123", category: "non-restrictive", value: false},
	{index: 41, name: "Vin too high SW", errorCode: "E34169", category: "non-restrictive", value: false},
	{index: 42, name: "digital event 42", errorCode: "E0"},
	{index: 43, name: "digital event 43", errorCode: "E0"},
	{index: 44, name: "digital event 44", errorCode: "E0"},
	{index: 45, name: "digital event 45", errorCode: "E0"},
	{index: 46, name: "digital event 46", errorCode: "E0"},
	{index: 47, name: "digital event 47", errorCode: "E0"},
	{index: 48, name: "digital event 48", errorCode: "E0"},
	{index: 49, name: "VBC/LVPS low", errorCode: "E34100", category: "no trigger", value: false},
	{index: 50, name: "BC enable", errorCode: "E0", category: "no trigger", value: false},
	{index: 51, name: "Const Curr. Charge", errorCode: "E0", category: "no trigger", value: false},
	{index: 52, name: "Const. Volt. Charge", errorCode: "E0", category: "no trigger", value: false},
	{index: 53, name: "Float Charge", errorCode: "E0", category: "no trigger", value: false},
	{index: 54, name: "BC Error TBat.", errorCode: "E0", category: "no trigger", value: false},
	{index: 55, name: "Emergency Charge", errorCode: "E0", category: "no trigger", value: false},
	{index: 56, name: "Udc min", errorCode: "E0", category: "no trigger", value: false},
	{index: 57, name: "digital event 57", errorCode: "E0"},
	{index: 58, name: "digital event 58", errorCode: "E0"},
	{index: 59, name: "digital event 59", errorCode: "E0"},
	{index: 60, name: "digital event 60", errorCode: "E0"},
	{index: 61, name: "digital event 61", errorCode: "E0"},
	{index: 62, name: "digital event 62", errorCode: "E0"},
	{index: 63, name: "digital event 63", errorCode: "E0"},
	{index: 64, name: "digital event 64", errorCode: "E0"},
}

func main() {
	// Define a string flag for the input filename
	inputFile := flag.String("input", "", "Input filename")

	// Parse the command-line flags
	flag.Parse()

	// Check if the input filename is provided
	if *inputFile == "" {
		fmt.Println("Input filename must be provided using the -input flag.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("MTM Simens fleet events parser, file: %s\n\n", *inputFile)

	// Move the name of the file to an application variable
	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("Error getting file info:", err)
	}

	// Allocate a byte slice to hold the file content
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	// Read the file into the buffer
	bytesRead, err := file.Read(buffer)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	if bytesRead != BlockSize*NbrOfBlocks {
		log.Fatal("Invalid file size:", bytesRead)
	}

	for i := 0; i < NbrOfBlocks; i++ {
		// Get the block
		block := buffer[i*BlockSize : (i+1)*BlockSize]

		// Parse the block
		parseBlock(block)
	}

	fmt.Println("End of program")
	os.Exit(0)
}

func parseBlock(block []byte) {
	fmt.Println("Parsing block")

	// Add all the index values to consts
	occurrences := (int)(uint16(block[15]) + (uint16(block[16]) << 8))
	fmt.Printf("Block Idx: %d has %d occurrences \n", block[9], occurrences)

	// Pretty print the contents of the block
	fmt.Println("Block Contents:")
	for i := 0; i < len(block); i++ {
		fmt.Printf("%02X ", block[i])
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	// If there is any occurence create a new event
	if occurrences == 0 {
		fmt.Printf("\n\n")
		return
	}

	eqpIdx := (int)(uint16(block[1]) + (uint16(block[0]) << 8))
	index := (int)(uint16(block[9]) + (uint16(block[10]) << 8))
	var evtName string
	var errorCode string
	if eqpIdx == 1 {
		evtName = ccuDigitalVars[index-1].name
		errorCode = ccuDigitalVars[index-1].errorCode
	} else if eqpIdx == 2 {
		evtName = invDigitalVars[index-1].name
		errorCode = invDigitalVars[index-1].errorCode
	} else if eqpIdx == 3 {
		evtName = lvpsDigitalVars[index-1].name
		errorCode = lvpsDigitalVars[index-1].errorCode
	} else {
		fmt.Println("Invalid equipment index")
		return
	}

	// Create an event
	event := sivEvent{
		eqpIdx:           eqpIdx,
		index:            index,
		name:             evtName,
		errorCode:        errorCode,
		nbrrOfOccurences: occurrences,
	}

	// Parse the occurrences, max 6 or the number of occurrences (each occurernce is 16 bytes)
	for i := 0; i < occurrences && i < 6; i++ {

		var occurrence sivEventOccurence
		occurrence.timeStamp = fmt.Sprintf("%02X-%02X-20%02X %02X:%02X:%02X", block[25+i*16], block[26+i*16], block[27+i*16], block[24+i*16], block[23+i*16], block[22+i*16])

		// Parse the active vars
		for j := 0; j < 8; j++ {
			for varIndex := 0; varIndex < 8; varIndex++ {
				varValue := (int)(0x01 & (block[28+j+i*16] >> byte(varIndex)))
				if varValue == 1 && eqpIdx == 1 {
					ccuDigitalVars[j*8+varIndex].value = true
					occurrence.activeVars = append(occurrence.activeVars, ccuDigitalVars[j*8+varIndex])
				} else if varValue == 1 && eqpIdx == 2 {
					invDigitalVars[j*8+varIndex].value = true
					occurrence.activeVars = append(occurrence.activeVars, invDigitalVars[j*8+varIndex])
				} else if varValue == 1 && eqpIdx == 3 {
					lvpsDigitalVars[j*8+varIndex].value = true
					occurrence.activeVars = append(occurrence.activeVars, lvpsDigitalVars[j*8+varIndex])
				} else {
					fmt.Println("Invalid equipment index")
					return
				}
			}
		}

		// Append occurence to the event
		event.occurrences = append(event.occurrences, occurrence)
	}

	// Print the event
	fmt.Printf("Event: %d, Name: %s, ErrorCode: %s, Number of occurrences: %d\n", event.index, event.name, event.errorCode, event.nbrrOfOccurences)

	// Print the active vars
	for _, v := range event.occurrences {
		fmt.Println("Occurrence: ", v.timeStamp)
		for _, activeVar := range v.activeVars {
			fmt.Printf("Active Signal: Index: %d, Name: %s, ErrorCode: %s, Category: %s\n", activeVar.index, activeVar.name, activeVar.errorCode, activeVar.category)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n\n")
}
