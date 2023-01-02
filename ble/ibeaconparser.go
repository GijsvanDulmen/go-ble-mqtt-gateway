package ble

import (
	"fmt"
	"strings"
)

type IBeaconParser struct {
	Data []byte
}

func NewIBeaconParser(data []byte) IBeaconParser {
	return IBeaconParser{
		Data: data,
	}
}

func (p IBeaconParser) IsIbeaconAdv() bool {
	return len(p.Data) >= 4 && ReadUInt32BE(p.Data, 0) == 0x4c000215
}

func (p IBeaconParser) ParserLogMessage(message string) {
	fmt.Println("ibeacon - ", message)
}

func (p IBeaconParser) ParserLog(key string, value string) {
	fmt.Println("ibeacon - ", key, " = ", value)
}

func (p IBeaconParser) parseUUID() string {
	return strings.Join([]string{
		fmt.Sprintf("%x", p.Data[4:8]),
		fmt.Sprintf("%x", p.Data[8:10]),
		fmt.Sprintf("%x", p.Data[10:12]),
		fmt.Sprintf("%x", p.Data[12:14]),
		fmt.Sprintf("%x", p.Data[14:20]),
	}, "-")
}

func (p IBeaconParser) parseMajor() uint16 {
	return ReadUint16BE(p.Data[20:22], 0)
}

func (p IBeaconParser) parseMinor() uint16 {
	return ReadUint16BE(p.Data[22:24], 0)
}

func (p IBeaconParser) parseTxPower() int8 {
	return ReadInt8(p.Data[24:25], 0)
}

func (p IBeaconParser) AddToMap(dataMap map[string]string) {
	if p.IsIbeaconAdv() {
		if len(p.Data) >= 25 {
			dataMap["uuid"] = p.parseUUID()
			dataMap["major"] = fmt.Sprint(p.parseMajor())
			dataMap["minor"] = fmt.Sprint(p.parseMinor())
			dataMap["txPower"] = fmt.Sprint(p.parseTxPower())
		}
	}
}
