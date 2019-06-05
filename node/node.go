package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/I1820/lanserver/models"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Node handles communication with nodes over the MQTT
type Node struct {
	db      *mongo.Database
	nodeCLI paho.Client
	appCLI  paho.Client
}

// Log provides data gathering endpoint for devices.
// we call these data that are came from devices log
// from kaa old days.
// thing -> log/{deveui}/send -> /device/{deveui}/rx -> app
func (n *Node) Log(client paho.Client, message paho.Message) {
	deveui := strings.Split(message.Topic(), "/")[1]

	logrus.Infof("Device: %s /log/send", deveui)

	result := n.db.Collection("devices").FindOne(context.Background(), bson.M{
		"deveui": deveui,
	})
	var d models.Device
	if err := result.Decode(&d); err != nil {
		logrus.Error(err)
		return
	}

	var log models.LogMessage
	if err := json.Unmarshal(message.Payload(), &log); err != nil {
		logrus.Error(err)
		return
	}

	if log.Token != d.Token {
		logrus.Error(fmt.Errorf("Mismatched token"))
		return
	}

	b, err := json.Marshal(models.RxMessage{
		DevEUI: deveui,
		Data:   log.Data,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	n.appCLI.Publish(fmt.Sprintf("device/%s/rx", deveui), 0, true, b)
}

// Notification provides data sending endpoint for devices.
// we call these data that are sent to devices notification
// from kaa old days.
// app -> device/{deveui}/tx -> /notification/{deveui}/request -> thing
func (n *Node) Notification(client paho.Client, message paho.Message) {
	deveui := strings.Split(message.Topic(), "/")[1]

	var notification models.TxMessage
	if err := json.Unmarshal(message.Payload(), &notification); err != nil {
		logrus.Error(err)
		return
	}

	b, err := json.Marshal(models.NotificationMessage{
		Data: notification.Data,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	var qos byte
	if notification.Confirmed {
		qos = 1
	}
	n.nodeCLI.Publish(fmt.Sprintf("notification/%s/request", deveui), qos, true, b)
}

// New creates mqtt client and connect into broker
func New(appBroker string, nodeBroker string, db *mongo.Database) (*Node, error) {
	n := &Node{
		db: db,
	}

	// application side mqtt
	{
		opts := paho.NewClientOptions()
		opts.AddBroker(appBroker)
		opts.SetOrderMatters(false)
		opts.SetOnConnectHandler(func(client paho.Client) {
			if t := client.Subscribe("device/+/tx", 0, n.Notification); t.Wait() && t.Error() != nil {
				logrus.Fatalf("MQTT subscription error: %s", t.Error())
			}
		})

		n.appCLI = paho.NewClient(opts)
		if t := n.appCLI.Connect(); t.Wait() && t.Error() != nil {
			return nil, t.Error()
		}
	}

	// node side mqtt
	{
		opts := paho.NewClientOptions()
		opts.AddBroker(nodeBroker)
		opts.SetOrderMatters(false)
		opts.SetOnConnectHandler(func(client paho.Client) {
			if t := client.Subscribe("log/+/send", 0, n.Log); t.Wait() && t.Error() != nil {
				logrus.Fatalf("MQTT subscription error: %s", t.Error())
			}
		})

		n.nodeCLI = paho.NewClient(opts)
		if t := n.nodeCLI.Connect(); t.Wait() && t.Error() != nil {
			return nil, t.Error()
		}
	}

	return n, nil
}
