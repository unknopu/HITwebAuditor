package entities

import (
	"auditor/core/mongodb"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/sethvargo/go-password/password"
)

// UserDevicePattern user device pattern
type UserDevicePattern struct {
	Pattern string `json:"pattern" bson:"pattern"`
}

// Device mqtt device
type Device struct {
	mongodb.Model `bson:",inline"`
	BrokerHost    string               `json:"broker_host" bson:"-"`
	BrokerHostWss string               `json:"broker_host_wss" bson:"-"`
	MountPoint    string               `json:"-" bson:"mountpoint"`
	ClientID      string               `json:"client_id" bson:"client_id"`
	Username      string               `json:"username" bson:"username"`
	PassHash      string               `json:"-" bson:"passhash"`
	Password      string               `json:"password" bson:"-"`
	RoomID        string               `json:"-" bson:"room_id"`
	PublishAcl    []*UserDevicePattern `json:"-" bson:"publish_acl"`
	SubscribeAcl  []*UserDevicePattern `json:"-" bson:"subscribe_acl"`
}

// NewUserDevice new user device
func NewUserDevice(roomID primitive.ObjectID) *Device {
	res, _ := password.Generate(10, 3, 0, false, false)
	ps, _ := bcrypt.GenerateFromPassword([]byte(res), 8)
	return &Device{
		BrokerHost:    MqttBroker,
		BrokerHostWss: MqttWSSBroker,
		MountPoint:    "",
		ClientID:      uuid.NewV4().String(),
		Username:      uuid.NewV4().String(),
		PassHash:      string(ps),
		Password:      res,
		RoomID:        roomID.Hex(),
		PublishAcl:    []*UserDevicePattern{},
		SubscribeAcl: []*UserDevicePattern{
			{
				Pattern: fmt.Sprintf("rooms/%s", roomID.Hex()),
			},
		},
	}
}

// NewAlertDevice new alert device
func NewAlertDevice() *Device {
	res, _ := password.Generate(10, 3, 0, false, false)
	ps, _ := bcrypt.GenerateFromPassword([]byte(res), 8)
	return &Device{
		BrokerHost:    MqttBroker,
		BrokerHostWss: MqttWSSBroker,
		MountPoint:    "",
		ClientID:      uuid.NewV4().String(),
		Username:      uuid.NewV4().String(),
		PassHash:      string(ps),
		Password:      res,
		PublishAcl:    []*UserDevicePattern{},
		SubscribeAcl: []*UserDevicePattern{
			{
				Pattern: "alerts",
			},
		},
	}
}

// NewPublicDevice new public device
func NewPublicDevice() *Device {
	res, _ := password.Generate(10, 3, 0, false, false)
	ps, _ := bcrypt.GenerateFromPassword([]byte(res), 8)
	return &Device{
		MountPoint: "",
		ClientID:   uuid.NewV4().String(),
		Username:   uuid.NewV4().String(),
		PassHash:   string(ps),
		Password:   res,
		PublishAcl: []*UserDevicePattern{
			{
				Pattern: "rooms/#",
			},
			{
				Pattern: "alerts",
			},
		},
		SubscribeAcl: []*UserDevicePattern{},
	}
}
