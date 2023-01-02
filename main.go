package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func SendAndWait(client MQTT.Client, topic string, contents string, retain bool) {
	token := client.Publish(topic, byte(1), retain, contents)

	_ = token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}
}

func main() {
	topic := os.Getenv("BLE_TOPIC")
	id := os.Getenv("BLE_ID")
	name := os.Getenv("BLE_NAME")

	enableLogging := os.Getenv("LOGGING") == "true"

	config := NewConfig()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID") + "-" + os.Getenv("BLE_ID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetCleanSession(true)

	choke := make(chan [2]string)

	opts.SetAutoReconnect(true)
	opts.SetBinaryWill(topic+id+"/LWT", []byte("Offline"), byte(1), true)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// set some standard stuff
	SendAndWait(client, topic+id+"/name", name, true)
	SendAndWait(client, topic+id+"/ip", "192.168.1.X", true)
	SendAndWait(client, topic+id+"/LWT", "Online", true)

	go func() {
		for {
			StartScan(func(results map[string]string) {

				if addr, ok := results["addr"]; ok {
					addr = strings.ToUpper(strings.ReplaceAll(addr, ":", ""))

					// check against whitelist and blacklists and such
					if !config.MatchesAgainstConfig(strings.ToUpper(addr), results["name"]) {
						// fmt.Println("Ignoring", results)
						return // no match
					}

					// example of what could be in the results
					// map[addr:e580a83705c2ce2dd163bda2ce589a8d connectable:yes
					// 		frameCounter:161 mac:582d34104766 name:ClearGrass Temp & RH productId:839
					//		rssi:-76 temp:20 version:3]
					allowedFields := []string{"rssi", "temp", "name", "hum", "battery", "uuid", "major", "minor", "txPower"}

					transitMessage := make(map[string]string)

					for key, value := range results {
						for _, v := range allowedFields {
							if v == key {
								transitMessage[key] = value
							}
						}
					}

					jsonStr, err := json.Marshal(transitMessage)
					if err != nil {
						fmt.Printf("Error: %s", err.Error())
					} else if enableLogging {
						fmt.Println(string(jsonStr))
					}

					SendAndWait(client, topic+id+"/bt/"+addr, fmt.Sprintf("%s", jsonStr), false)
				}
			})

			fmt.Println("Stopped scanning... trying again in a few seconds")
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		fmt.Println("Still processing")
		time.Sleep(60 * time.Second)
	}
}
