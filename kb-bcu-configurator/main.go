package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"gopkg.in/yaml.v2"
)

type Can_id struct {
	Id   uint32 `yaml:"id"`
	Size uint8  `yaml:"size"`
	Nv   uint8  `yaml:"nv"`
}

type Kcf_opt struct {
	Hdr_data string `yaml:"header"`
	Udp_net  struct {
		Ip     string `yaml:"ip"`
		Port   uint32 `yaml:"port"`
		Period uint32 `yaml:"period"`
	} `yaml:"udp"`
	Can_ids []Can_id `yaml:"can"`
}

type Can_msg struct {
	CANID [4]byte
	SIZE  byte
	TYPE  byte
}

type Kcf_msg struct {
	KCF_CHTYPE_CONT        [4]byte
	CHUNK_LEN              [4]byte
	KCF_CHTYPE_CONT_HDR    [4]byte
	HDR_LENGTH             [4]byte
	HDR_DATA               []byte
	PADDING                [2]byte
	KCF_CHTYPE_CAN_CFG_HDR [4]byte
	CFG_CAN_DATA_LEN       [4]byte
	IP_ADDRESS             [4]byte
	PORT                   [4]byte
	NBR_CANIDS             [4]byte
	CYCLE_UDP_PERIOD       [4]byte
	CAN                    []Can_msg
	KCF_MSG                []byte
}

func ip_to_bin(ip string) uint32 {
	var ret uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &ret)
	return ret
}

func main() {
	dstPtr := flag.String("dst", "localhost:8000", "data destinantion")
	flag.Parse()

	// Open yaml file
	fd, err := ioutil.ReadFile("kcf-parameters.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal yaml file
	var kcf_opts Kcf_opt
	err = yaml.Unmarshal(fd, &kcf_opts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("KCF CONFIGURATIONS FROM FILE")
	fmt.Println(kcf_opts.Hdr_data)
	fmt.Println(kcf_opts.Udp_net.Ip)
	fmt.Println(kcf_opts.Udp_net.Port)
	fmt.Println(kcf_opts.Udp_net.Period)
	for _, v := range kcf_opts.Can_ids {
		fmt.Println(v.Id)
		fmt.Println(v.Size)
		fmt.Println(v.Nv)
	}

	// Create Default KCF TODO: Automatic calculation of length values
	var kcfFrame Kcf_msg
	binary.BigEndian.PutUint32(kcfFrame.KCF_CHTYPE_CONT[:], 0)
	binary.BigEndian.PutUint32(kcfFrame.CHUNK_LEN[:], 0)
	binary.BigEndian.PutUint32(kcfFrame.KCF_CHTYPE_CONT_HDR[:], 0x01)

	binary.BigEndian.PutUint32(kcfFrame.HDR_LENGTH[:], uint32(len(kcf_opts.Hdr_data)))
	kcfFrame.HDR_DATA = []byte(kcf_opts.Hdr_data)
	binary.BigEndian.PutUint16(kcfFrame.PADDING[:], 0)
	binary.BigEndian.PutUint32(kcfFrame.KCF_CHTYPE_CAN_CFG_HDR[:], 0x1f)

	binary.BigEndian.PutUint32(kcfFrame.IP_ADDRESS[:], ip_to_bin(kcf_opts.Udp_net.Ip))
	binary.BigEndian.PutUint32(kcfFrame.PORT[:], uint32(kcf_opts.Udp_net.Port))

	binary.BigEndian.PutUint32(kcfFrame.NBR_CANIDS[:], uint32(len(kcf_opts.Can_ids)))
	binary.BigEndian.PutUint32(kcfFrame.CYCLE_UDP_PERIOD[:], uint32(kcf_opts.Udp_net.Period))

	var aux Can_msg
	for _, c := range kcf_opts.Can_ids {
		binary.BigEndian.PutUint32(aux.CANID[:], c.Id)
		aux.SIZE = byte(c.Size)
		aux.TYPE = byte(c.Nv)
		kcfFrame.CAN = append(kcfFrame.CAN, aux)
	}
	chunkLength := len(kcfFrame.KCF_CHTYPE_CONT) + len(kcfFrame.CHUNK_LEN) +
		len(kcfFrame.KCF_CHTYPE_CONT_HDR) + len(kcfFrame.HDR_LENGTH) +
		len(kcf_opts.Hdr_data) + len(kcfFrame.PADDING) +
		len(kcfFrame.KCF_CHTYPE_CAN_CFG_HDR) + len(kcfFrame.CFG_CAN_DATA_LEN) +
		len(kcfFrame.IP_ADDRESS) + len(kcfFrame.PORT) + len(kcfFrame.NBR_CANIDS) +
		len(kcfFrame.CYCLE_UDP_PERIOD) + (6 * len(kcf_opts.Can_ids))
	binary.BigEndian.PutUint32(kcfFrame.CHUNK_LEN[:], uint32(chunkLength))

	cfgLength := len(kcfFrame.KCF_CHTYPE_CAN_CFG_HDR) + len(kcfFrame.CFG_CAN_DATA_LEN) +
		len(kcfFrame.IP_ADDRESS) + len(kcfFrame.PORT) + len(kcfFrame.NBR_CANIDS) +
		len(kcfFrame.CYCLE_UDP_PERIOD) + (6 * len(kcf_opts.Can_ids))
	binary.BigEndian.PutUint32(kcfFrame.CFG_CAN_DATA_LEN[:], uint32(cfgLength))

	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.KCF_CHTYPE_CONT[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.CHUNK_LEN[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.KCF_CHTYPE_CONT_HDR[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.HDR_LENGTH[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.HDR_DATA...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.PADDING[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.KCF_CHTYPE_CAN_CFG_HDR[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.CFG_CAN_DATA_LEN[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.IP_ADDRESS[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.PORT[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.NBR_CANIDS[:]...)
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, kcfFrame.CYCLE_UDP_PERIOD[:]...)
	// CAN ID ports from settings file
	for _, v := range kcfFrame.CAN {
		kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, v.CANID[:]...)
		kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, v.SIZE, v.TYPE)
	}
	// End of frame bytes
	kcfFrame.KCF_MSG = append(kcfFrame.KCF_MSG, 0xff, 0xee)

	fmt.Println(kcfFrame.KCF_MSG, len(kcfFrame.KCF_MSG)-2)

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

	_, err = conn.Write(kcfFrame.KCF_MSG)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Message sent: ", kcfFrame.KCF_MSG)
}
