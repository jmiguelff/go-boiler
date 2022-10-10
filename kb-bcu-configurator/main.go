package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destinantion")
	flag.Parse()

	// Get the IP struct from hostname
	tcpAddr, err := net.ResolveTCPAddr("tcp", *dstPtr)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to the server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Fail to connect to ", *dstPtr)
		log.Fatalln(err)
	}
	// Close socket after return from main
	defer conn.Close()

	// Create Default KCF
	KCF_CHTYPE_CONT := "00000000"
	CHUNK_LEN := "000000c6"
	KCF_CHTYPE_CONT_HDR := "00000001"

	HDR_LENGTH := "0000008a"
	HDR_DATA_ASCII := []byte("Suitable.for:...STN49089.v01.37.KCF:STN49089/KCF.ver.:01.05.....Date:2022-10-10.Author:.J.Ferreira..............TOOL.ver.00.06..HEADERv03$")
	HDR_DATA := hex.EncodeToString(HDR_DATA_ASCII)
	PADDING := "0000"

	KCF := KCF_CHTYPE_CONT + CHUNK_LEN + KCF_CHTYPE_CONT_HDR + HDR_LENGTH + HDR_DATA + PADDING

	KCF_CHTYPE_CAN_CFG_HDR := "0000001f"
	CFG_CAN_DATA_LEN := "0000002a"
	IP_ADDRESS := "0a6e1701"
	PORT := "00003a9d"

	KCF = KCF + KCF_CHTYPE_CAN_CFG_HDR + CFG_CAN_DATA_LEN + IP_ADDRESS + PORT

	NBR_CANIDS := "00000004"
	CYCLE_UDP_PERIOD := "000001f4"

	CANID_1 := "00000021"
	SIZE_1 := "08"
	TYPE_1 := "02"

	CANID_2 := "00000682"
	SIZE_2 := "08"
	TYPE_2 := "02"

	CANID_3 := "00000686"
	SIZE_3 := "08"
	TYPE_3 := "02"

	CANID_4 := "000006a0"
	SIZE_4 := "08"
	TYPE_4 := "02"

	EOF := "ffee"

	KCF = KCF + NBR_CANIDS + CYCLE_UDP_PERIOD + CANID_1 + SIZE_1 + TYPE_1 + CANID_2 + SIZE_2 + TYPE_2 + CANID_3 + SIZE_3 + TYPE_3 + CANID_4 + SIZE_4 + TYPE_4 + EOF
	msg, err := hex.DecodeString(KCF)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Message sent: ", msg)
}
