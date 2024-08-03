package server

import (
	"encoding/json"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

type JoinRoomData struct {
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
}

type NewRoomDetail struct {
	User    *User  `json:"user"`
	RoomId  string `json:"roomId"`
	Members []User `json:"members"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DrawProps struct {
	Ctx          interface{} `json:"ctx"`
	CurrentPoint Point       `json:"currentPoint"`
	PrevPoint    *Point      `json:"prevPoint"`
}

type DrawOptions struct {
	Ctx          interface{} `json:"ctx"`
	CurrentPoint Point       `json:"currentPoint"`
	PrevPoint    *Point      `json:"prevPoint"`
	StrokeColor  string      `json:"strokeColor"`
	StrokeWidth  []float64   `json:"strokeWidth"`
	DashGap      []float64   `json:"dashGap"`
}

type ReplyMessage struct {
	Message string  `json:"message"`
	Title   *string `json:"title"`
}

type CanvasState struct {
	CanvasState string `json:"canvasState"`
	RoomId      string `json:"roomId"`
}

type DrawData struct {
	DrawOptions *DrawOptions `json:"drawOptions"`
	RoomId      string       `json:"roomId"`
}

type UndoPoint struct {
	UndoPoint string `json:"undoPoint"`
	RoomId    string `json:"roomId"`
}

var users []User

var undoPoints = make(map[string][]string)

func addUndoPoint(roomId string, undoPoint string) {
	undoPointInRoom, ok := undoPoints[roomId]
	if ok {
		undoPoints[roomId] = append(undoPointInRoom, undoPoint)
	} else {
		undoPoints[roomId] = []string{undoPoint}
	}
}

func getLastUndoPoint(roomId string) string {
	roomUndoPoints, ok := undoPoints[roomId]
	if !ok {
		return ""
	}
	if len(roomUndoPoints) == 0 {
		return ""
	}
	return roomUndoPoints[len(roomUndoPoints)-1]
}

func deleteLastUndoPoint(roomId string) {
	room, ok := undoPoints[roomId]
	if !ok {
		return
	}
	if len(room) > 0 {
		undoPoints[roomId] = room[:len(room)-1]
	}
}

func getUser(userId string) *User {
	for _, user := range users {
		if user.ID == userId {
			return &user
		}
	}
	return nil
}

func getRoomMembers(roomId string) []User {
	var roomMembers []User
	for _, user := range users {
		if user.RoomId == roomId {
			roomMembers = append(roomMembers, user)
		}
	}
	return roomMembers
}

func addUser(user *User) {
	users = append(users, *user)
}

func removeUser(userId string) {
	users = filterUsers(users, func(user User) bool {
		return user.ID != userId
	})
}

func filterUsers(users []User, predicate func(User) bool) []User {
	var filteredUsers []User
	for _, user := range users {
		if predicate(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func RegisterEventHandlers(server *socketio.Server) {
	server.OnEvent("/", "create-room", func(conn socketio.Conn, message string) {
		var joinRoomData JoinRoomData
		err := json.Unmarshal([]byte(message), &joinRoomData)

		if err != nil {
			panic(err)
		}
		joinRoom(server, conn, &joinRoomData)
	})
	server.OnEvent("/", "join-room", func(conn socketio.Conn, message string) {
		var joinRoomData JoinRoomData
		err := json.Unmarshal([]byte(message), &joinRoomData)

		if err != nil {
			panic(err)
		}
		rooms := server.Rooms("/")

		isRoomCreated := false

		for _, room := range rooms {
			if room == joinRoomData.RoomId {
				isRoomCreated = true
				break
			}
		}
		if !isRoomCreated {
			message := &ReplyMessage{
				Message: "Oops! The Room ID you entered doesn't exist or hasn't been created yet.",
			}

			jsonMessage, err := json.Marshal(message)

			if err != nil {
				panic(err)
			}

			conn.Emit("room-not-found", string(jsonMessage))
			return
		}

		joinRoom(server, conn, &joinRoomData)
	})
	server.OnEvent("/", "client-ready", func(conn socketio.Conn, roomId string) {
		membersInRoom := getRoomMembers(roomId)

		if len(membersInRoom) == 1 {
			conn.Emit("client-loaded")
			return
		}

		adminMember := membersInRoom[0]
		server.BroadcastToRoom("/", adminMember.RoomId, "get-canvas-state")
	})
	server.OnEvent("/", "send-canvas-state", func(conn socketio.Conn, message string) {
		var canvasStateOfRoom CanvasState
		err := json.Unmarshal([]byte(message), &canvasStateOfRoom)
		if err != nil {
			panic(err)
		}
		membersInRoom := getRoomMembers(canvasStateOfRoom.RoomId)
		lastMemberInRoom := membersInRoom[len(membersInRoom)-1]
		server.BroadcastToRoom("/", lastMemberInRoom.RoomId, "canvas-state-from-server", canvasStateOfRoom.CanvasState)
	})
	server.OnEvent("/", "draw", func(conn socketio.Conn, drawOptions interface{}, roomId string) {
		server.BroadcastToRoom("/", roomId, "update-canvas-state", drawOptions)
	})
	server.OnEvent("/", "clear-canvas", func(conn socketio.Conn, roomId string) {
		server.BroadcastToRoom("/", roomId, "clear-canvas")

	})
	server.OnEvent("/", "undo", func(conn socketio.Conn, message string) {
		var canvasStateOfRoom CanvasState
		err := json.Unmarshal([]byte(message), &canvasStateOfRoom)
		if err != nil {
			panic(err)
		}
		server.BroadcastToRoom("/", canvasStateOfRoom.RoomId, "undo-canvas", canvasStateOfRoom.CanvasState)
	})
	server.OnEvent("/", "get-last-undo-point", func(conn socketio.Conn, roomId string) {
		lastUndoPoint := getLastUndoPoint(roomId)
		conn.Emit("last-undo-point-from-server", lastUndoPoint)
	})
	server.OnEvent("/", "add-undo-point", func(conn socketio.Conn, message string) {
		var undoPointInRoom UndoPoint
		err := json.Unmarshal([]byte(message), &undoPointInRoom)
		if err != nil {
			panic(err)
		}
		addUndoPoint(undoPointInRoom.RoomId, undoPointInRoom.UndoPoint)
	})
	server.OnEvent("/", "delete-last-undo-point", func(conn socketio.Conn, roomId string) {
		deleteLastUndoPoint(roomId)
	})
	server.OnEvent("/", "leave-room", func(conn socketio.Conn) {
		leaveRoom(server, conn)
	})
}

func joinRoom(server *socketio.Server, conn socketio.Conn, joinRoomData *JoinRoomData) {
	conn.Join(joinRoomData.RoomId)

	user := &User{
		Username: joinRoomData.Username,
		RoomId:   joinRoomData.RoomId,
		ID:       conn.ID(),
	}

	addUser(user)

	membersInRoom := getRoomMembers(joinRoomData.RoomId)

	newRoomDetails := &NewRoomDetail{
		User:    user,
		RoomId:  joinRoomData.RoomId,
		Members: membersInRoom,
	}

	payload, err := json.Marshal(newRoomDetails)

	if err != nil {
		panic(err)
	}

	conn.Emit("room-joined", string(payload))

	membersInRoomJson, err := json.Marshal(membersInRoom)

	if err != nil {
		panic(err)
	}

	server.BroadcastToRoom("/", joinRoomData.RoomId, "update-members", string(membersInRoomJson))

	notificationTitle := "New member arrived!"

	notification := &ReplyMessage{
		Title:   &notificationTitle,
		Message: fmt.Sprintf("%s joined the party.", joinRoomData.Username),
	}

	jsonNotification, err := json.Marshal(notification)

	if err != nil {
		panic(err)
	}

	server.BroadcastToRoom("/", joinRoomData.RoomId, "send-notification", string(jsonNotification))
}

func leaveRoom(server *socketio.Server, conn socketio.Conn) {
	user := getUser(conn.ID())
	if user == nil {
		return
	}
	removeUser(conn.ID())
	members := getRoomMembers(user.RoomId)

	membersData, err := json.Marshal(members)

	if err != nil {
		panic(err)
	}

	server.BroadcastToRoom("/", user.RoomId, "update-members", string(membersData))

	notificationTitle := "Member departure!"
	notification := &ReplyMessage{
		Title:   &notificationTitle,
		Message: fmt.Sprintf("%s left the party", user.Username),
	}

	jsonNotification, err := json.Marshal(notification)

	if err != nil {
		panic(err)
	}

	server.BroadcastToRoom("/", user.RoomId, "send-notification", string(jsonNotification))
}
