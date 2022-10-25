package user

import "atus/backend/websocket"

func (user *User) SendNotification(clientHub *websocket.Hub, theType, title, message string) {
	for _, client := range clientHub.GetClientsByUID(user.UID) {
		client.MarshalAndSend("NOTIFICATION", map[string]interface{}{
			"type":    theType,
			"title":   title,
			"message": message,
		})
	}
}

func BroadcastNotification(clientHub *websocket.Hub, theType, title, message string) {
	allUsers, _ := GetAll()
	for _, user := range allUsers {
		user.SendNotification(clientHub, theType, title, message)
	}
}
