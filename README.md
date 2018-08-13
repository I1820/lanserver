# Lanserver.sh
[![Travis branch](https://img.shields.io/travis/com/I1820/lanserver/master.svg?style=flat-square)](https://travis-ci.com/I1820/lanserver)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/lanserver?style=flat-square)](https://goreportcard.com/report/github.com/I1820/lanserver)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Maintainability](https://api.codeclimate.com/v1/badges/5db031209d82d7354ae0/maintainability)](https://codeclimate.com/github/I1820/lanserver/maintainability)

## Introduction

Implementation of Ad-hoc standard for managing Ethernet based things similar to LoRa specification.
It uses JWT Tokens for device validation on uplink. devices send uplink and receive downlink via mqtt
protocol. Lanserver provides two interface for application: HTTP for management and MQTT for data.
