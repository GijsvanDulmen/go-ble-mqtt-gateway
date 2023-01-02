package ble

import "fmt"

type XiaomiParser struct {
	UUID []byte
	Data []byte
}

type FrameControl struct {
	IsFactoryNew    bool
	IsConnected     bool
	IsCentral       bool
	IsEncrypted     bool
	HasMacAddress   bool
	HasCapabilities bool
	HasEvent        bool
	HasCustomData   bool
	HasSubtitle     bool
	HasBinding      bool
}

const (
	EventTypeTemperature            uint16 = 4100
	EventTypeHumidity               uint16 = 4102
	EventTypeLux                    uint16 = 4103
	EventTypeMoisture               uint16 = 4104
	EventTypeFertility              uint16 = 4105
	EventTypeBattery                uint16 = 4106
	EventTypeTemperatureAndHumidity uint16 = 4109
)

func NewXiaomiParser(data []byte, uuid []byte) XiaomiParser {
	return XiaomiParser{
		UUID: uuid,
		Data: data,
	}
}

func (p XiaomiParser) parseFrameControl() FrameControl {
	frameControl := ReadUint16LE(p.Data, 0)

	return FrameControl{
		IsFactoryNew:    (frameControl & (1 << 0)) != 0,
		IsConnected:     (frameControl & (1 << 1)) != 0,
		IsCentral:       (frameControl & (1 << 2)) != 0,
		IsEncrypted:     (frameControl & (1 << 3)) != 0,
		HasMacAddress:   (frameControl & (1 << 4)) != 0,
		HasCapabilities: (frameControl & (1 << 5)) != 0,
		HasEvent:        (frameControl & (1 << 6)) != 0,
		HasCustomData:   (frameControl & (1 << 7)) != 0,
		HasSubtitle:     (frameControl & (1 << 8)) != 0,
		HasBinding:      (frameControl & (1 << 9)) != 0,
	}
}

func (p XiaomiParser) parseEventOffset() int {
	offset := 5
	frameControl := p.parseFrameControl()
	if frameControl.HasMacAddress {
		offset = 11
	}
	if frameControl.HasCapabilities {
		offset++
	}
	return offset
}

func (p XiaomiParser) parseEventType() uint16 {
	offset := p.parseEventOffset()
	return ReadUint16LE(p.Data, offset)
}

func (p XiaomiParser) parseEventLength() uint8 {
	return ReadUint8(p.Data, p.parseEventOffset()+2)
}

func (p XiaomiParser) parseTemperatureEvent() int16 {
	return ReadInt16LE(p.Data, p.parseEventOffset()+3) / 10
}

func (p XiaomiParser) parseHumidityEvent() uint16 {
	return ReadUint16LE(p.Data, p.parseEventOffset()+3) / 10
}

func (p XiaomiParser) parseHumidityEventWhenTemperatureAsWell() uint16 {
	return ReadUint16LE(p.Data, p.parseEventOffset()+5) / 10
}

func (p XiaomiParser) parseBatteryEvent() uint8 {
	return ReadUint8(p.Data, p.parseEventOffset()+3)
}

func (p XiaomiParser) parseVersion() uint8 {
	return ReadUint8(p.Data, 1) >> 4
}

func (p XiaomiParser) parseProductId() uint16 {
	return ReadUint16LE(p.Data, 2)
}

func (p XiaomiParser) parseFrameCounter() uint8 {
	return ReadUint8(p.Data, 4)
}

func (p XiaomiParser) parseMacAddress() string {
	macBuffer := p.Data[5:(5 + 6)]
	return fmt.Sprintf("%x", Reverse(macBuffer))
}

func (p XiaomiParser) ParserLogMessage(message string) {
	fmt.Println(p.UUIDToString(), " - ", message)
}

func (p XiaomiParser) ParserLog(key string, value string) {
	fmt.Println(p.UUIDToString(), " - ", key, " = ", value)
}

func (p XiaomiParser) UUIDToString() string {
	return fmt.Sprintf("%x", Reverse(p.UUID))
}

func (p XiaomiParser) IsXiaomiAdv() bool {
	return len(p.Data) >= 5 && p.UUIDToString() == "fe95"
}

func (p XiaomiParser) Format() {
	if p.IsXiaomiAdv() {
		frameControl := p.parseFrameControl()

		p.ParserLog("version", fmt.Sprint(p.parseVersion()))
		p.ParserLog("productId", fmt.Sprint(p.parseProductId()))
		p.ParserLog("frameCounter", fmt.Sprint(p.parseFrameCounter()))

		if frameControl.HasMacAddress {
			p.ParserLog("mac", p.parseMacAddress())
		}

		if frameControl.HasEvent {
			p.ParserLogMessage("has event")

			eventType := p.parseEventType()
			if eventType == EventTypeBattery {
				p.ParserLogMessage("battery event")
				p.ParserLog("bat", fmt.Sprint(p.parseBatteryEvent()))
			} else if eventType == EventTypeTemperature {
				p.ParserLogMessage("temp event")
				p.ParserLog("temp", fmt.Sprint(p.parseTemperatureEvent()))
			} else if eventType == EventTypeHumidity {
				p.ParserLogMessage("hum event")
				p.ParserLog("hum", fmt.Sprint(p.parseHumidityEvent()))
			} else if eventType == EventTypeTemperatureAndHumidity {
				p.ParserLogMessage("temp + hum event")
				p.ParserLog("temp", fmt.Sprint(p.parseTemperatureEvent()))
				p.ParserLog("hum", fmt.Sprint(p.parseHumidityEventWhenTemperatureAsWell()))
			} else {
				p.ParserLogMessage("type not implemented yet")
			}
		} else {
			p.ParserLogMessage("has no event")
			fmt.Println("Has no event")
		}
	} else {
		p.ParserLogMessage("has no capable frame control")
	}
}

func (p XiaomiParser) AddToMap(dataMap map[string]string) {
	if p.IsXiaomiAdv() {
		frameControl := p.parseFrameControl()

		dataMap["version"] = fmt.Sprint(p.parseVersion())
		dataMap["productId"] = fmt.Sprint(p.parseProductId())
		dataMap["frameCounter"] = fmt.Sprint(p.parseFrameCounter())

		if frameControl.HasMacAddress {
			dataMap["mac"] = p.parseMacAddress()
		}

		if frameControl.HasEvent {
			eventType := p.parseEventType()
			if eventType == EventTypeBattery {
				dataMap["battery"] = fmt.Sprint(p.parseBatteryEvent())
			} else if eventType == EventTypeTemperature {
				dataMap["temp"] = fmt.Sprint(p.parseTemperatureEvent())
			} else if eventType == EventTypeHumidity {
				dataMap["hum"] = fmt.Sprint(p.parseHumidityEvent())
			} else if eventType == EventTypeTemperatureAndHumidity {
				dataMap["temp"] = fmt.Sprint(p.parseTemperatureEvent())
				dataMap["hum"] = fmt.Sprint(p.parseHumidityEventWhenTemperatureAsWell())
			}
		}
	}
}
