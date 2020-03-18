/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 19-07-2019
 * |
 * | File Name:     default.go
 * +===============================================
 */

package config

// Default represents default configuration in YAML format with 2-space
const Default = `
debug: true
database:
  name: lanserver
  url: mongodb://127.0.0.1:27017
app:
  broker:
    addr: tcp://127.0.0.1:1883
node:
  broker:
    addr: tcp://127.0.0.1:1883
`
