package handler

import (
	"bloodOnTheClockTower/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var character = []model.BaseCharacter{}
var characterMap = map[int]model.BaseCharacter{}

func LoadGame(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 获取URL参数
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	roomId := parts[2]
	playerId := parts[3]
	game, _ := findRoom(roomId)

	for {
		if game == nil {
			break
		}

		mux := game.Mux

		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("Client disconnected:", err)
				return
			}
			log.Println("Read error Gaming Process:", err)
			return
		}

		var actionReq model.ActionReq
		if err = json.Unmarshal(p, &actionReq); err != nil {
			log.Println("JSON unmarshal error:", err)
		}

		switch actionReq.Action {
		//host点击load game开始游戏
		case "load_game":
			fmt.Println("load game now")
			initGame(mux, game, playerId, conn)
		//玩家发动技能
		case "cast":
			cast(mux, game, playerId, actionReq.Targets, actionReq.Extra)
		//玩家发动提名
		//提名-投票-结束投票-提名-投票-结束投票.....
		case "nominate":
			nominate(mux, game, playerId, actionReq.Targets)
		//提名后发起投票
		case "vote":
			vote(mux, game, playerId)
		//入夜
		case "checkout_night":
			//if playerId == "1" {
			checkoutNight(mux, game)
		//}
		case "first_night":
			FirstNight(game)
		//直接入夜
		case "direct_night":
			DirectToNight(game)
			//公投结束之后进行处决
		case "execute":
			//if playerId == "1" {
			checkoutDay(mux, game)
			//}
		// 提名后结束投票,退出投票环节
		case "end_voting":
			//if playerId == "1" {
			endVoting(mux, game)
			//}
			//case "quit_game":
			//	quitGame(mux, game, playerId)
			//}

			// 有结果则跳出循环
			if game.Result != "" {
				break
			}

			time.Sleep(time.Millisecond * 50)
		}
	}
}

// set up characters,八人局为例子，随机抽取村民5人、外来者1人，爪牙1人，恶魔1人
func assign() ([]model.BaseCharacter, map[int]model.BaseCharacter) {
	//分配阵营
	var character []model.BaseCharacter

	//首版固定5-1-1-1和角色名称
	fixedGroupMap := map[int][]string{
		1: {model.GoodCharacterMap[11], model.GoodCharacterMap[5], model.GoodCharacterMap[7],
			model.GoodCharacterMap[6], model.GoodCharacterMap[12],
			model.OutsiderCharacterMap[4],
			model.MinionCharacterMap[3],
			model.DevilCharacterMap[3]},
		2: {model.GoodCharacterMap[11], model.GoodCharacterMap[2], model.GoodCharacterMap[1],
			model.GoodCharacterMap[9],
			model.OutsiderCharacterMap[4], model.OutsiderCharacterMap[3],
			model.MinionCharacterMap[1],
			model.DevilCharacterMap[4]},
		3: {model.GoodCharacterMap[1], model.GoodCharacterMap[8], model.GoodCharacterMap[7],
			model.GoodCharacterMap[13], model.GoodCharacterMap[12],
			model.OutsiderCharacterMap[2],
			model.MinionCharacterMap[3],
			model.DevilCharacterMap[2]},
		4: {model.GoodCharacterMap[8], model.GoodCharacterMap[3], model.GoodCharacterMap[6],
			model.GoodCharacterMap[5], model.GoodCharacterMap[12],
			model.OutsiderCharacterMap[4],
			model.MinionCharacterMap[4],
			model.DevilCharacterMap[3]},
		5: {model.GoodCharacterMap[4], model.GoodCharacterMap[1], model.GoodCharacterMap[7],
			model.GoodCharacterMap[2], model.GoodCharacterMap[11],
			model.OutsiderCharacterMap[1],
			model.MinionCharacterMap[2],
			model.DevilCharacterMap[1]},
	}
	//todo 修改fixedGroupMap的index来获取不同组合
	characterMap := make(map[int]model.BaseCharacter)
	for k, v := range fixedGroupMap[1] {
		character = append(character, model.BaseCharacter{
			Id:              k,
			CharacterName:   v,
			CharacterKind:   model.CharacterKindMap[v],
			CharacterStatus: nil,
			IsDead:          false,
		})
		characterMap[k+1] = model.BaseCharacter{
			Id:              k,
			CharacterName:   v,
			CharacterKind:   model.CharacterKindMap[v],
			CharacterStatus: nil,
			IsDead:          false,
		}
	}
	return character, characterMap
}

func init() {
	character, characterMap = assign()
}

func initGame(mux *sync.Mutex, game *model.Room, playerId string, conn *websocket.Conn) {
	mux.Lock()
	defer mux.Unlock()
	game.GameConnPool.Store(playerId, conn)
	if game == nil {
		return
	}

	if !game.Init {
		// 初始化
		game.Init = true
		game.CreatedAt = time.Now().Format(time.RFC3339)
		game.Result = ""
		game.Log = ""
		game.Executed = nil
		game.CastPool = map[string][]string{}
		game.VoteLogs = map[string]string{}
		game.VotePool = map[string]int{}
		game.Nominated = nil
		game.Executed = nil
		game.State = model.GameState{}
		game.AllCharacters = make(map[string]string)
		// 初始化玩家状态 防止非法返回房间引起bug
		for i, player := range game.Players {
			newPlayer := model.Player{}
			newPlayer.Id = player.Id
			newPlayer.Name = player.Name
			newPlayer.Index = player.Index
			newPlayer.PositionId = player.PositionId
			newPlayer.Ready.Nominate = true
			newPlayer.Ready.Nominated = true
			newPlayer.Ready.Vote = true
			//首轮都不需要发动技能
			newPlayer.State.Casted = true
			game.Players[i] = newPlayer
		}
		// 初始化玩家状态 依赖身份
		for i := range game.Players {
			game.Players[i].BaseCharacter = characterMap[game.Players[i].PositionId]
			game.AllCharacters[game.Players[i].Id] = characterMap[game.Players[i].PositionId].CharacterName
		}
		// 保存玩家身份到总日志
		//var hasRecluse bool
		game.Log = "本局配置：\n"
		for _, player := range game.Players {
			game.Log += fmt.Sprintf("玩家 [%s] 的身份是 {%s} \n", player.Name, player.BaseCharacter.CharacterName)
			if player.State.Drunk {
				game.Log += fmt.Sprintf("玩家 [%s] 的身份其实是 {%s} ~\n", player.Name, "Drunk")
			}
			//if player.Character == Recluse {
			//	game.Log += fmt.Sprintf("玩家 [%s] 的被当作的身份是 {%s} ~\n", player.Name, player.State.RegardedAs)
			//	hasRecluse = true
			//}
			//if player.Character == FortuneTeller {
			//	hasFortuneTeller = true
			//}
		}
		game.Log += "------------本--局--开--始------------\n"
	}

	// 群发game
	if game.Result == "" {
		broadcast(game)
	}
}

// emit 发送game到指定终端
func emit(game *model.Room, destinationId string) {
	game.ResMux.Lock()
	defer game.ResMux.Unlock()
	marshaledGame, err := json.Marshal(*game)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	game.GameConnPool.Range(func(id, conn any) bool {
		if id == destinationId {
			if err := conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshaledGame); err != nil {
				log.Println("Write error:", err)
				return false
			}
			return false
		}
		return true
	})
}

// broadcast 广播game
func broadcast(game *model.Room) {
	game.ResMux.Lock()
	defer game.ResMux.Unlock()
	marshaledGame, err := json.Marshal(*game)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	game.GameConnPool.Range(func(id, conn any) bool {
		if err = conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, marshaledGame); err != nil {
			log.Println("Write error:", err)
			return false
		}
		return true
	})
}

func genRandomPositionSlice(indexSliceForCharacterTypePool []int, characterByTypePool []string, num int) []int {
	randomInt := rand.Intn(len(characterByTypePool))
	indexSliceForCharacterTypePool = append(indexSliceForCharacterTypePool, randomInt)
	for {
		if len(indexSliceForCharacterTypePool) == num {
			break
		}
		randomInt = rand.Intn(len(characterByTypePool))
		repeatFlag := false
		for j := 0; j < len(indexSliceForCharacterTypePool); j++ {
			if indexSliceForCharacterTypePool[j] == randomInt {
				repeatFlag = true
				break
			}
		}
		if !repeatFlag {
			indexSliceForCharacterTypePool = append(indexSliceForCharacterTypePool, randomInt)
		}
	}
	return indexSliceForCharacterTypePool
}

func quitGame(mux *sync.Mutex, game *model.Room, playerId string) {
	mux.Lock()
	defer mux.Unlock()
	if game == nil {
		return
	}

	for i, player := range game.Players {
		if player.Id == playerId {
			game.Players[i].Quited = true
			break
		}
	}

	// 关闭退出者的game连接
	game.GameConnPool.Range(func(id, conn any) bool {
		// 关闭创建房间者的连接
		if id == playerId {
			conn.(*websocket.Conn).Close()
			game.GameConnPool.Delete(id)
			return true
		}
		return true
	})
}

func detectIfAllQuited(mux *sync.Mutex, game *model.Room) {
	cfg := model.GetConfig()
	mux.Lock()
	defer mux.Unlock()

	if game == nil {
		return
	}

	var allQuited = true
	for _, player := range game.Players {
		allQuited = allQuited && player.Quited
	}
	if allQuited {
		CfgMutex.Lock()
		var newRooms []model.Room
		for _, roomm := range cfg.Rooms {
			if game.Id != roomm.Id {
				newRooms = append(newRooms, roomm)
			}
		}
		cfg.Rooms = newRooms
		CfgMutex.Unlock()
	}
}
