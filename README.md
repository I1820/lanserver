# Lanserver.sh
[![Travis branch](https://img.shields.io/travis/com/I1820/lanserver/master.svg?style=flat-square)](https://travis-ci.com/I1820/lanserver)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/lanserver?style=flat-square)](https://goreportcard.com/report/github.com/I1820/lanserver)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Codacy Badge](https://img.shields.io/codacy/grade/28e224e07bec4eca96eb8c30b4535603.svg?style=flat-square)](https://www.codacy.com/project/i1820/lanserver/dashboard)

## Introduction

Implementation of Ad-hoc standard for managing Ethernet based things similar to LoRa specification.
It uses JWT Tokens for device validation on uplink. devices send uplink and receive downlink via mqtt
protocol. Lanserver provides two interface for application: HTTP for management and MQTT for data.
