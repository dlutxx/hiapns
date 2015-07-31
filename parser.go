package hiapns

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlutxx/apns"
	"time"
)

func extractString(d map[string]interface{}, key string) (string, error) {
	ival, ok := d[key]
	if !ok {
		return "", fmt.Errorf("key missing: %v", key)
	}
	val, ok := ival.(string)
	if !ok {
		return "", fmt.Errorf("bad value: %v, %v", key, ival)
	}
	return val, nil
}

// sample msg body:
// {"expire":1438241451, "app": "appname", "token": "token-str", "payload": {"aps":{"alert":"msg","sound":"hi.caf"},"extra":2}}
func ParseRequestFromJSON(bs []byte) (req Request, err error) {
	var data map[string]interface{}

	if err = json.Unmarshal(bs, &data); err != nil {
		return
	}
	// check expiration
	if iexpire, ok := data["expire"]; ok {
		if expire, ok := iexpire.(float64); ok {
			if now := time.Now().Unix(); now > int64(expire) {
				err = errors.New("msg expired")
				return
			}
		} else {
			err = errors.New("bad expire")
			return
		}
	} else {
		err = errors.New("key missing: expire")
		return
	}
	// extract token
	token, err := extractString(data, "token")
	if err != nil {
		return
	}
	// extract app
	req.App, err = extractString(data, "app")
	if err != nil {
		return
	}
	// extract payload
	ipayload, ok := data["payload"]
	if !ok {
		err = errors.New("key missing: payload")
		return
	}
	payload, ok := ipayload.(map[string]interface{})
	if !ok {
		err = errors.New("bad payload")
		return
	}

	req.Notif, err = NotificationFromMsg(token, payload)
	if err != nil {
		err = errors.New("invalid payload")
		return
	}

	return
}

func NotificationFromMsg(tok string, msg map[string]interface{}) (*apns.Notification, error) {
	p, err := PayloadFromMsg(msg)
	if err != nil {
		return nil, err
	}
	n := apns.NewNotification()
	n.Payload = p
	n.DeviceToken = tok
	n.Priority = apns.PriorityImmediate
	return &n, nil
}

func parseAPSFromMsg(aps *apns.APS, msg map[string]interface{}) error {
	if alert, ok := msg["alert"]; ok {
		if v, ok := alert.(string); ok {
			aps.Alert.Body = v
		} else {
			return ErrBadAlert
		}
	} else {
		return ErrAlertMissing
	}
	badgeNo := 0
	if badge, ok := msg["badge"]; ok {
		switch i := badge.(type) {
		case float64:
			badgeNo = int(i)
		default:
			return ErrBadBadge
		}
	}
	aps.Badge = &badgeNo

	if sound, ok := msg["sound"]; ok {
		if s, ok := sound.(string); ok {
			aps.Sound = s
		} else {
			return ErrBadSound
		}
	}
	if v, ok := msg["content-available"]; ok {
		if ca, ok := v.(float64); ok && ca > 0 {
			aps.ContentAvailable = 1
		}
	}
	return nil
}

func PayloadFromMsg(msg map[string]interface{}) (*apns.Payload, error) {
	p := apns.NewPayload()
	for k, v := range msg {
		switch k {
		case "aps":
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, ErrBadFormat
			}
			if err := parseAPSFromMsg(&p.APS, m); err != nil {
				return nil, err
			}
		default:
			p.SetCustomValue(k, v)
		}
	}
	return p, nil
}
