package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/I1820/lanserver/models"
	jwt "github.com/dgrijalva/jwt-go"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Log provides data gathering endpoint for devices.
// we call these data that are came from devices log
// from kaa old days.
// log/{deveui}/send
func Log(client paho.Client, message paho.Message) {
	deveui := strings.Split(message.Topic(), "/")[1]

	logger.Infof("Device: %s /log/send", deveui)

	result := db.Collection("devices").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("deveui", deveui),
	))
	var d models.Device
	if err := result.Decode(&d); err != nil {
		logger.Error(err)
		return
	}

	var log models.LogMessage
	if err := json.Unmarshal(message.Payload(), &log); err != nil {
		logger.Error(err)
		return
	}

	if log.Token != d.Token {
		logger.Error(fmt.Errorf("Mismatched token"))
		return
	}

	var key []byte
	copy(key[:], envy.Get("NODE_SECRET", ""))
	if _, err := jwt.ParseWithClaims(d.Token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		c := token.Claims.(*jwt.StandardClaims)

		if !c.VerifyIssuer("lanserver.sh", true) {
			return nil, fmt.Errorf("Unexpected issuer %v", c.Issuer)
		}
		if c.Id != deveui {
			return nil, fmt.Errorf("Mismatched identifier %s != %s", c.Id, deveui)
		}
		return key, nil
	}); err != nil {
		logger.Error(err)
		return
	}

	b, err := json.Marshal(models.RxMessage{
		DevEUI: deveui,
		Data:   log.Data,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	mqttApplication.Publish(fmt.Sprintf("device/%s/rx", deveui), 0, true, b)
}

// Notification provides data sending endpoint for devices.
// we call these data that are sent to devices notification
// from kaa old days.
// device/{deveui}/tx
func Notification(client paho.Client, message paho.Message) {
	deveui := strings.Split(message.Topic(), "/")[1]

	var notification models.TxMessage
	if err := json.Unmarshal(message.Payload(), &notification); err != nil {
		logger.Error(err)
		return
	}

	b, err := json.Marshal(models.NotificationMessage{
		Data: notification.Data,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	mqttNode.Publish(fmt.Sprintf("notification/%s/request", deveui), 0, true, b)
}
