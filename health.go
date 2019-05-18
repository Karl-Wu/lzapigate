package main

func encodeHealthReport() []byte {

	timestamp := uint32(NowInMS())
	// freq_khz := uint32(0)

	msg := []byte{}
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

	/**
	 * 状态 0 故障，1 正常
	 **/
	//unsigned char status;
	msg[pos] = 1 //unsigned char m_ucDisplay;
	pos = pos + 1

	pos = pos + 5 // char bReserve[5];

	byteOrder.PutUint16(msg[2:4], uint16(pos)) //set length

	return msg
}
