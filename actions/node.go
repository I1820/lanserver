package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/I1820/lanserver/models"
	jwt "github.com/dgrijalva/jwt-go"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Log provides data gathering endpoint for devices.
// we call these data that are came from devices log
// from kaa old days.
func Log(client paho.Client, message paho.Message) {
	var deveui string
	fmt.Sscanf(message.Topic(), "/log/%s/send", &deveui)

	result := db.Collection("devices").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("deveui", deveui),
	))
	var d models.Device
	if err := result.Decode(&d); err != nil {
		return
		// return c.Error(http.StatusInternalServerError, err)
	}

	var log models.LogMessage
	if err := json.Unmarshal(message.Payload(), &log); err != nil {
		return
		// return c.Error(http.StatusBadRequest, err)
	}

	if log.Token != d.Token {
		return
		// return c.Error(http.StatusForbidden, fmt.Errorf("Mismatched token"))
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
		return
		// return c.Error(http.StatusForbidden, err)
	}

	b, err := json.Marshal(models.RxMessage{
		DevEUI: deveui,
		Data:   log.Data,
	})
	if err != nil {
		return
		// return c.Error(http.StatusInternalServerError, err)
	}
	mqttApplication.Publish(fmt.Sprintf("device/%s/rx", deveui), 0, true, b)
}

// Notification provides data sending endpoint for devices.
// we call these data that are sent to devices notification
// from kaa old days.
func Notification(client paho.Client, message paho.Message) {
	var deveui string
	fmt.Sscanf(message.Topic(), "/device/%s/tx", &deveui)

	var notification models.TxMessage
	if err := json.Unmarshal(message.Payload(), &notification); err != nil {
		return
		// return c.Error(http.StatusBadRequest, err)
	}

	b, err := json.Marshal(models.NotificationMessage{
		Data: notification.Data,
	})
	if err != nil {
		return
		// return c.Error(http.StatusInternalServerError, err)
	}
	mqttNode.Publish(fmt.Sprintf("notification/%s/request", deveui), 0, true, b)
}
