# Lanserver
[![Travis branch](https://img.shields.io/travis/com/I1820/lanserver/master.svg?style=flat-square)](https://travis-ci.com/I1820/lanserver)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/lanserver?style=flat-square)](https://goreportcard.com/report/github.com/I1820/lanserver)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Codacy Badge](https://img.shields.io/codacy/grade/28e224e07bec4eca96eb8c30b4535603.svg?style=flat-square)](https://www.codacy.com/project/i1820/lanserver/dashboard)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/lanserver)

## Introduction

Implementation of Ad-hoc standard for managing Ethernet-based things similar to LoRa specification.
It uses JWT Tokens for device authorization on the uplink.
Devices send uplink and receive downlink via MQTT protocol.
Lanserver provides two interfaces for application: HTTP for management and MQTT for data.

Lanserver implementation did not cover any IoT protocol specification.
It assumes your device is connected directly into Ethernet without any gateway.
If you want to have specific IoT protocol and control on each aspect of that protocol
this project is not suitable for you, please search more. :see_no_evil:

## Protocol
First of all, you must register your device in Lanserver and after that Lanserver gives you the token.
Use the provided token on your LAN device to transmit data into Lanserver.

Following MQTT topics are used in Lanserver for communicating via devices:

- `/log/{deveui}/send`:
this topic is used to transmit data into Lanserver with the following structure:

[models/messages.go/LogMessage]
```
{
    data: {base64 encoded data}
    token: {token}
}
```
- `/notification/{deveui}/request`:
this topic is used to recieve data from Lanserver with the following structure:

[models/messages.go/NotificationMessage]
```
{
    data: {base64 encoded data}
}
```

In the application layer following MQTT topics are used in Lanserver for communicating via applications:

- `/device/{deveui}/rx`:
this topic used to transmit data into application with the following structure:

[models/messages.go/RxMessage]
```
{
    data: {base64 encoded data}
    deveui: {device identification}
}
```
