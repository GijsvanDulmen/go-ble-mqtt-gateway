package ble

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type CleargrassParser struct {
	UUID []byte
	Data []byte
}

func NewCleargrassParser(data []byte, uuid []byte) CleargrassParser {
	return CleargrassParser{
		UUID: uuid,
		Data: data,
	}
}

func (p CleargrassParser) IsCleargrassAdv() bool {
	return len(p.Data) >= 5 && p.UUIDToString() == "fdcd"
}

func (p CleargrassParser) UUIDToString() string {
	return fmt.Sprintf("%x", Reverse(p.UUID))
}

func (p CleargrassParser) AddToMap(dataMap map[string]string) {
	if p.IsCleargrassAdv() {
		var hex = hex.EncodeToString(p.Data)

		// https://github.com/alexvenom/XiaomiCleargrassInkDislpay/blob/master/XiaomiClearGrassInk.js
		// 080c 78e850342d58 01 04 c100 0502
		// 0807453810342d580104f500da02020145
		if len(hex) >= 24 {

			total := hex[22:24] + hex[20:22]
			parsed, _ := strconv.ParseInt(total, 16, 64)
			temperature := (float64(parsed) / 10)

			dataMap["temp"] = fmt.Sprint(temperature)

			totalHum := hex[26:28] + hex[24:26]
			parsed, _ = strconv.ParseInt(totalHum, 16, 64)
			hum := (float64(parsed) / 10)

			dataMap["hum"] = fmt.Sprint(hum)
		}
	}
}
