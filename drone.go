package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	seqGenerator = uint16(0)
	msgGenerator = uint16(0)
	droneStore   = make(map[string]*droneStruct)

	byteOrder = binary.LittleEndian
)

type droneStruct struct {
	seqNum    uint16
	lzID      string
	staled    bool
	msgFromLz interface{}
}

func (drn *droneStruct) encode(remove bool) ([]byte, int) {

	timestamp := uint32(0)
	// freq_khz := uint32(0)

	drone := drn.msgFromLz.(map[string]interface{})
	// fmt.Println("drone: ", drone["name"])
	// fmt.Println("sensor: ", drone["seen_sensor"])
	// fmt.Println(reflect.TypeOf(drone["seen_sensor"]))
	seenSensors := (drone["seen_sensor"]).([]interface{})
	for _, s := range seenSensors {
		sensor := s.(map[string]interface{})
		fmt.Println("sensor: ", sensor["detected_freq_khz"])
		v, _ := sensor["detected_time"].(float64)
		timestamp = uint32(v)
		// v, _ = sensor["detected_freq_khz"].(float64)
		// freq_khz = uint32(v)
		break
	}

	if remove {
		//TO BE DONE
	}

	msg := make([]byte, 1024)
	pos := 0
	byteOrder.PutUint16(msg[pos:pos+2], genMsgNum()) //unsigned short wMsgNo;
	pos = pos + 2
	byteOrder.PutUint16(msg[pos:pos+2], 0) //unsigned short wMsgLenth;
	pos = pos + 2
	byteOrder.PutUint32(msg[pos:pos+4], ip2Long(myIP)) //unsigned long dwSrcIP;
	pos = pos + 4
	byteOrder.PutUint16(msg[pos:pos+2], myPort) //unsigned short wSrcPort;
	pos = pos + 2
	//padding here ?
	//pos = pos + 2
	byteOrder.PutUint32(msg[pos:pos+4], ip2Long(commanderIP)) //unsigned long dwDestIP;
	pos = pos + 4

	byteOrder.PutUint16(msg[pos:pos+2], commanderPort) //unsigned short wDestPort;
	pos = pos + 2

	//padding here ?
	//pos = pos + 2

	byteOrder.PutUint16(msg[pos:pos+2], myID) //unsigned short wDevID;
	pos = pos + 2

	byteOrder.PutUint32(msg[pos:pos+4], timestamp) //unsigned long dTimeStamp;
	pos = pos + 4

	msg[pos] = 2 ////unsigned char cIsFixedFreq;
	pos = pos + 1
	//pos = pos + 1

	byteOrder.PutUint16(msg[pos:pos+2], drn.seqNum) //unsigned short seqNum;
	pos = pos + 2

	byteOrder.PutUint64(msg[pos:pos+8], 0) //long long WorkMode;
	pos = pos + 8

	byteOrder.PutUint64(msg[pos:pos+8], 0) //double dMinFrequence; size of double is 8 ?
	pos = pos + 8

	byteOrder.PutUint64(msg[pos:pos+8], 0) //double dMaxFrequence;
	pos = pos + 8

	byteOrder.PutUint64(msg[pos:pos+8], 0) //double fBearing;
	pos = pos + 8

	byteOrder.PutUint64(msg[pos:pos+8], 0) //double m_fAngle;
	pos = pos + 8

	msg[pos] = 0 //unsigned char m_ucQuality;
	pos = pos + 1
	//pos = pos + 3

	byteOrder.PutUint64(msg[pos:pos+8], 0) //long long m_Time;
	pos = pos + 8

	msg[pos] = 0 //unsigned char m_ucDisplay;
	pos = pos + 1
	//pos = pos + 1

	byteOrder.PutUint16(msg[pos:pos+2], 0) // unsigned short m_sDisplayLevel;
	pos = pos + 2

	byteOrder.PutUint16(msg[pos:pos+2], 0) // unsigned short m_usFreqNum;
	pos = pos + 2

	byteOrder.PutUint16(msg[pos:pos+2], 2) // short sVisalLimitLevel;
	pos = pos + 2

	byteOrder.PutUint16(msg[pos:pos+2], 0) // short sLimitLevel;
	pos = pos + 2

	msg[pos] = 0 //bool bParametersAsk;
	pos = pos + 1

	msg[pos] = 0 //unsigned char type;
	pos = pos + 1

	return msg, pos
}

func genSeqNum() uint16 {
	seqGenerator = seqGenerator + 1
	return seqGenerator
}

func genMsgNum() uint16 {
	msgGenerator = msgGenerator + 1
	return msgGenerator
}

func upsertDroneList(droneList []interface{}) (bool, []*droneStruct) {
	var (
		ok     bool
		rmList = []*droneStruct{}
	)
	changed := false
	for _, droneItem := range droneStore {
		droneItem.staled = true
	}

	for _, v := range droneList {
		drone := v.(map[string]interface{})
		// fmt.Println("drone: ", drone["name"])
		// fmt.Println(reflect.TypeOf(drone["seen_sensor"]))
		lzID := drone["id"].(string)
		var cur *droneStruct
		if cur, ok = droneStore[lzID]; !ok {
			cur = &droneStruct{
				seqNum:    genSeqNum(),
				lzID:      lzID,
				staled:    false,
				msgFromLz: v,
			}
			changed = true
		} else {
			cur.staled = false
			cur.msgFromLz = v
		}
		droneStore[lzID] = cur
	}

	for k, droneItem := range droneStore {
		if droneItem.staled == true {
			delete(droneStore, k)
			rmList = append(rmList, droneItem)
			changed = true
		}
	}

	return changed, rmList
}

const droneQueryStr = `{
	drone {
		id,
		image,
		name,
		description,
		state,
		can_attack,
		can_takeover,
		direction,
		created_time,
		deleted_time,
		whitelisted,
		attacking,
		seen_sensor{
			sensor_id,
			detected_freq_khz,
			detected_time
		}
	}
}`

func droneListFetch() {
	var (
		rmList []*droneStruct
	)
	updated := false
	data, err := postGraphQL(graphqlURL, droneQueryStr)
	if err == nil {
		val := data["drone"]
		droneList := val.([]interface{})

		updated, rmList = upsertDroneList(droneList)

		// fmt.Println(reflect.TypeOf(droneList))
		// for _, v := range droneList {
		// 	drone := v.(map[string]interface{})
		// 	fmt.Println("drone: ", drone["name"])
		// 	// fmt.Println("sensor: ", drone["seen_sensor"])
		// 	fmt.Println(reflect.TypeOf(drone["seen_sensor"]))
		// 	seenSensors := (drone["seen_sensor"]).([]interface{})
		// 	for _, s := range seenSensors {
		// 		fmt.Println("sensor: ", s)
		// 		sensor := s.(map[string]interface{})
		// 		fmt.Println("sensor: ", sensor["detected_freq_khz"])
		// 	}

		// }

	}

	if updated {
		for _, droneItem := range droneStore {
			droneItem.sendDetectMsg(false)
		}
	}

	for _, droneItem := range rmList {
		droneItem.sendDetectMsg(true)
	}
}

func (drn *droneStruct) sendDetectMsg(rmmove bool) {

	buf, len := drn.encode(rmmove)
	dstAddr := &net.UDPAddr{
		Port: int(myID),
		IP:   net.ParseIP(commanderIP),
	}

	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	DEBUG.Println(DRONE, "drone message ", dstAddr, len, buf[0:len])
	n, err := conn.WriteTo(buf[0:len], dstAddr)
	DEBUG.Println(DRONE, "send detect message: ", n, err)
}

func serverDetect() {
	for appStop == false {
		droneListFetch()
		time.Sleep(2 * time.Second)
	}
}
