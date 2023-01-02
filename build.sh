#!/bin/bash
docker build -t golang-ble-gateway .

docker tag golang-ble-gateway ghcr.io/gijsvandulmen/iot-ble-go-gateway:latest

# pushes
docker push ghcr.io/gijsvandulmen/iot-ble-go-gateway:latest
