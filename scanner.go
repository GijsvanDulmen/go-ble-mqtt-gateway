package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"
	bleparsing "vandulmen.net/iot/ble-gateway/ble"
)

type ScanResultHandler func(map[string]string)

func StartScan(handler ScanResultHandler) {
	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	duration := (60 * 60 * 24 * 150) * time.Second
	// duration := 10 * time.Second

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), duration))

	fmt.Println("Start scanning")
	ScanResultCheckError(ble.Scan(ctx, true /* duplicates */, func(a ble.Advertisement) {
		results := ScanAdvertismentHandler(a)
		handler(results)
	}, nil))
}

func ScanAdvertismentHandler(a ble.Advertisement) map[string]string {
	dataMap := make(map[string]string)

	dataMap["addr"] = a.Addr().String()
	dataMap["rssi"] = fmt.Sprint(a.RSSI())

	if a.Connectable() {
		dataMap["connectable"] = "yes"
	} else {
		dataMap["connectable"] = "no"
	}

	if len(a.LocalName()) > 0 {
		dataMap["name"] = a.LocalName()
	}

	if len(a.ServiceData()) > 0 {
		for _, adv := range a.ServiceData() {
			parser := bleparsing.NewXiaomiParser(adv.Data, adv.UUID)
			if parser.IsXiaomiAdv() {
				parser.AddToMap(dataMap)
			}

			cgParser := bleparsing.NewCleargrassParser(adv.Data, adv.UUID)
			if cgParser.IsCleargrassAdv() {
				cgParser.AddToMap(dataMap)
			}
		}
	}
	if len(a.ManufacturerData()) > 0 {
		ibeacon := bleparsing.NewIBeaconParser(a.ManufacturerData())
		if ibeacon.IsIbeaconAdv() {
			ibeacon.AddToMap(dataMap)
		}
	}

	return dataMap
}

func ScanResultCheckError(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		log.Fatalf(("canceled"))
		// fmt.Printf("canceled\n")
	default:
		fmt.Printf("Error")
		fmt.Printf(err.Error())
		fmt.Printf("\n")
	}
}
