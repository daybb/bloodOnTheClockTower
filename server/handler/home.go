package handler

import (
	"bloodOnTheClockTower/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

func LoadHome(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("Client disconnected:", err)
				return
			}
			log.Println("Read error Home:", err)
			return
		}

		var reqBody model.HomeReqBody
		if err = json.Unmarshal(p, &reqBody); err != nil {
			log.Println("JSON unmarshal error:", err)
		}

		switch reqBody.Action {
		case "list_rooms":
			var reqBody model.ListRoomsReqBody
			if err = json.Unmarshal(p, &reqBody); err != nil {
				log.Println("JSON unmarshal error:", err)
			}
			listRooms(reqBody.Payload, conn)
		case "create_room":
			//var reqBody model.CreateRoomReqBody
			//if err = json.Unmarshal(p, &reqBody); err != nil {
			//	log.Println("JSON unmarshal error:", err)
			//}
			initDefaultRoom()
		case "join_room":
			var reqBody model.JoinRoomReqBody
			if err = json.Unmarshal(p, &reqBody); err != nil {
				log.Println("JSON unmarshal error:", err)
			}
			//加入房间并设为准备状态
			joinRoomAndReady(reqBody.Payload, conn)
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func listRooms(playerId string, conn *websocket.Conn) {
	cfg := model.GetConfig()
	cfg.HomeConnPool.Store(playerId, conn)
	// 发送房间列表给请求者
	marshalRooms, err := json.Marshal(cfg.Rooms)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	connVal, _ := cfg.HomeConnPool.LoadOrStore(playerId, conn)
	if err = connVal.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRooms); err != nil {
		log.Println("Write error:", err)
		return
	}
	cfg.HomeConnPool.Range(func(key, value any) bool {
		return true
	})
}

// 点击创建房间创建一个默认房间
func initDefaultRoom() {
	cfg := model.GetConfig()
	CfgMutex.Lock()
	defer CfgMutex.Unlock()
	room := model.Room{
		Id:           "2025",
		Name:         "摸摸鱼",
		Host:         "Server",
		CreatedAt:    time.Now().Format(time.RFC3339),
		Status:       model.Wait,
		Init:         false,
		Result:       "",
		Log:          "",
		Players:      []model.Player{},
		State:        model.GameState{},
		Executed:     nil,
		Nominated:    nil,
		CastPool:     nil,
		VotePool:     nil,
		VoteLogs:     nil,
		GameConnPool: &sync.Map{},
		Mux:          &sync.Mutex{},
		ResMux:       &sync.Mutex{},
	}
	for i := range cfg.Rooms {
		if cfg.Rooms[i].Id == room.Id {
			return
		}
	}
	room.Status = model.Wait
	room.CreatedAt = time.Now().Format(time.RFC3339)
	cfg.Rooms = append(cfg.Rooms, room)
	// 发送房间列表给所有人
	//marshalRooms, err := json.Marshal(cfg.Rooms)
	//if err != nil {
	//	log.Println("JSON marshal error:", err)
	//	return
	//}
	//cfg.HomeConnPool.Range(func(id, conn any) bool {
	//	//// 关闭创建房间者的连接
	//	//if id == room.Players[0].Id {
	//	//	conn.(*websocket.Conn).Close()
	//	//	cfg.HomeConnPool.Delete(id)
	//	//	return true
	//	//}
	//	if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRooms); err != nil {
	//		log.Println("Write error:", err)
	//		return false
	//	}
	//	return true
	//})
}

// 加入默认房间 用户输入昵称和位置id
func joinRoomAndReady(joinRoomPayload model.JoinRoomPayload, conn *websocket.Conn) {
	cfg := model.GetConfig()
	CfgMutex.Lock()
	defer CfgMutex.Unlock()

	room, roomIndex := findRoom("2025")
	msg, exclusive := nicknameAndPositionExclusive(joinRoomPayload.Player.Name, joinRoomPayload.Player.PositionId, *room)
	if !exclusive {
		//有重复的，需要发送对应消息到对应客户端
		cfg.HomeConnPool.Range(func(id, conn any) bool {
			if id == joinRoomPayload.Player.Id {
				if err := conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("Write error:", err)
					return false
				}
				return true
			}
			log.Println("still找不到匹配的")
			return false
		})
		return
	}
	cfg.Rooms[roomIndex].Players = append(room.Players, joinRoomPayload.Player)
	//加入game的长连接
	room.GameConnPool.Store(joinRoomPayload.Player.Id, conn)
	if room == nil {
		return
	}
	//}

	// 发送房间列表给home页所有人
	marshalRooms, err := json.Marshal(cfg.Rooms)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	cfg.HomeConnPool.Range(func(id, conn any) bool {
		if id == joinRoomPayload.Player.Id {
			conn.(*websocket.Conn).Close()
			cfg.HomeConnPool.Delete(id)
			return true
		}
		if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRooms); err != nil {
			log.Println("Write error:", err)
			return false
		}
		return true
	})

	// 发送房间给所有人
	marshalRoom, err := json.Marshal(*room)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	room.GameConnPool.Range(func(id, conn any) bool {
		if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRoom); err != nil {
			log.Println("Write error:", err)
			return false
		}
		return true
	})
	//ready := checkRoomStatus(room)
	//if ready {
	//	startGame(room)
	//}
	//return ready
}

func nicknameAndPositionExclusive(nickName string, pos int, config model.Room) (string, bool) {
	if config.Id != "2025" {
		log.Println("no config")
		return "配置错误", false
	}
	for _, v := range config.Players {
		if v.PositionId == pos {
			return "位置已存在", false
		}
		if v.Name == nickName {
			return "昵称已存在", false
		}
	}
	return "", true

}

// update game status and broadcast
func startGame(room *model.Room) {
	goodToStart := true
	for _, player := range room.Players {
		goodToStart = goodToStart && player.Waiting
	}
	if goodToStart {
		room.Status = model.Processing
		room.Init = false
		room.Result = ""
		room.Log = ""
		room.CastPool = map[string][]string{}
		room.VoteLogs = map[string]string{}
		room.VotePool = map[string]int{}
		room.Nominated = nil
		room.Executed = nil
		room.State = model.GameState{}
	}

	// 发送房间给所有人
	marshalRoom, err := json.Marshal(*room)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	room.GameConnPool.Range(func(id, conn any) bool {
		room.ResMux.Lock()
		defer room.ResMux.Unlock()
		if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRoom); err != nil {
			log.Println("Write error:", err)
			return false
		}
		return true
	})

	// 将room list 发给所有homeConn池里的人
	//marshalRooms, err := json.Marshal(cfg.Rooms)
	//if err != nil {
	//	log.Println("JSON marshal error:", err)
	//	return
	//}
	//cfg.HomeConnPool.Range(func(id, conn any) bool {
	//	if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshalRooms); err != nil {
	//		log.Println("Write error:", err)
	//		return false
	//	}
	//	return true
	//})
}

func checkRoomStatus(room *model.Room) bool {
	if len(room.Players) == 8 {
		length := 0
		room.GameConnPool.Range(func(key, value interface{}) bool {
			length++
			return true // 继续遍历
		})
		if length == 8 {
			return true
		}
	}
	return false
}
