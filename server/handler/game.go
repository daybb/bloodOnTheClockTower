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
		//玩家点击load game开始游戏
		case "load_game":
			//check game ready
			//checkRoomStatus(game)
			initGame(mux, game, playerId, conn)
		case "toggle_night":
			//if playerId == game.Host {
			//	toggleNight(mux, game)
			//}
		case "cast":
			//cast(mux, game, playerId, actionReq.Targets)
		case "nominate":
			//nominate(mux, game, playerId, actionReq.Targets)
		case "vote":
			//vote(mux, game, playerId)
		//入夜，包括第一夜
		case "checkout_night":
			if playerId == game.Host {
				checkoutNight(mux, game)
			}
		case "checkout_day":
			//if playerId == game.Host {
			//	checkoutDay(mux, game)
			//}
		case "end_voting":
			//if playerId == game.Host {
			//	endVoting(mux, game)
			//}
		case "quit_game":
			quitGame(mux, game, playerId)
		}

		// 有结果则跳出循环
		if game.Result != "" {
			break
		}

		// 检测是否房间内所有人都退出游戏
		detectIfAllQuited(mux, game)

		time.Sleep(time.Millisecond * 50)
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
		characterMap[k] = model.BaseCharacter{
			Id:              k,
			CharacterName:   v,
			CharacterKind:   model.CharacterKindMap[v],
			CharacterStatus: nil,
			IsDead:          false,
		}
	}
	fmt.Println(character)
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
		// 初始化玩家状态 防止非法返回房间引起bug
		for i, player := range game.Players {
			newPlayer := model.Player{}
			newPlayer.Id = player.Id
			newPlayer.Name = player.Name
			newPlayer.Index = player.Index
			newPlayer.PositionId = player.PositionId
			game.Players[i] = newPlayer
		}
		// 初始化玩家状态 依赖身份
		for i := range game.Players {
			game.Players[i].BaseCharacter = characterMap[game.Players[i].PositionId]
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
		emit(game, playerId)
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

func allocateCharacter(players []model.Player) ([]model.Player, string) {
	playerNum := len(players)
	var characterNumSlice []int
	// 索引0是村民，索引1是外来者，索引2是爪牙，索引3是恶魔
	switch playerNum {
	case 5:
		characterNumSlice = []int{3, 0, 1, 1}
	case 6:
		characterNumSlice = []int{3, 1, 1, 1}
	case 7:
		characterNumSlice = []int{5, 0, 1, 1}
	case 8:
		characterNumSlice = []int{5, 1, 1, 1}
	case 9:
		characterNumSlice = []int{5, 2, 1, 1}
	case 10:
		characterNumSlice = []int{7, 0, 2, 1}
	case 11:
		characterNumSlice = []int{7, 1, 2, 1}
	case 12:
		characterNumSlice = []int{7, 2, 2, 1}
	case 13:
		characterNumSlice = []int{9, 0, 3, 1}
	case 14:
		characterNumSlice = []int{9, 1, 3, 1}
	case 15:
		characterNumSlice = []int{9, 2, 3, 1}
	default:
		characterNumSlice = []int{5, 2, 1, 1} // 9人
	}
	// 里面存放的是针对对应身份类型池的index
	// 这些index的数量取决于characterNumSlice中的对应位置的数字，比如本局有5个村民，那townsfolkRandNums的长度就是5
	// 这些index都是根据对应的身份类型池的总数量的随机数，比如本局有3个村民，[0,7,4]，对应到TownsfolkPool的0/7/4索引位置
	townsfolkRandNums := []int{}
	outsidersRandNums := []int{}
	minionsRandNums := []int{}
	demonsRandNums := []int{}
	for i, num := range characterNumSlice {
		if num != 0 {
			switch i {
			case 0:
				townsfolkRandNums = genRandomPositionSlice(townsfolkRandNums, TownsfolkPool, num)
			case 1:
				outsidersRandNums = genRandomPositionSlice(outsidersRandNums, OutsidersPool, num)
			case 2:
				minionsRandNums = genRandomPositionSlice(minionsRandNums, MinionsPool, num)
			case 3:
				demonsRandNums = genRandomPositionSlice(demonsRandNums, DemonsPool, num)
			}
		}
	}

	var characterPoolForSelection []string
	var replaceDrunk string
	var repeatFlag bool
	for _, randIdx := range demonsRandNums {
		characterPoolForSelection = append(characterPoolForSelection, DemonsPool[randIdx])
	}
	for _, randIdx := range minionsRandNums {
		characterPoolForSelection = append(characterPoolForSelection, MinionsPool[randIdx])
		// 处理男爵
		if MinionsPool[randIdx] == Baron {
			outsidersNumsLength := len(outsidersRandNums)
			for {
				if len(outsidersRandNums) == outsidersNumsLength+2 {
					break
				}
				repeatFlag = false
				randIdxOutsiders := rand.Intn(len(OutsidersPool))
				for _, num := range outsidersRandNums {
					if randIdxOutsiders == num {
						repeatFlag = true
						break
					}
				}
				if !repeatFlag {
					outsidersRandNums = append(outsidersRandNums, randIdxOutsiders)
				}
			}
			townsfolkRandNums = townsfolkRandNums[:len(townsfolkRandNums)-2]
		}
	}
	if len(outsidersRandNums) != 0 {
		for _, randIdx := range outsidersRandNums {
			characterPoolForSelection = append(characterPoolForSelection, OutsidersPool[randIdx])
			// 处理酒鬼 酒鬼还在池中，后续阶段由replaceDrunk替换
			if OutsidersPool[randIdx] == Drunk {
				for {
					if replaceDrunk != "" {
						break
					}
					repeatFlag = false
					randIdxTownsfolk := rand.Intn(len(TownsfolkPool))
					for _, num := range townsfolkRandNums {
						if randIdxTownsfolk == num {
							repeatFlag = true
							break
						}
					}
					if !repeatFlag {
						replaceDrunk = TownsfolkPool[randIdxTownsfolk]
						break
					}
				}
			}
		}
	}
	for _, randIdx := range townsfolkRandNums {
		characterPoolForSelection = append(characterPoolForSelection, TownsfolkPool[randIdx])
	}

	// 打乱顺序
	characterPoolForSelection = Shuffle(characterPoolForSelection)

	var characterTypePoolForSelection []string
	for _, elem := range characterPoolForSelection {
		for key, pool := range CharacterPool {
			if Contains(pool, elem) {
				characterTypePoolForSelection = append(characterTypePoolForSelection, key)
			}
		}
	}

	var newPlayers []model.Player
	var newPlayer model.Player
	for i := range players {
		newPlayer.Id = players[i].Id
		newPlayer.Name = players[i].Name
		newPlayer.Index = i
		newPlayer.Waiting = false
		newPlayer.Character = characterPoolForSelection[i]
		newPlayer.CharacterType = characterTypePoolForSelection[i]
		newPlayers = append(newPlayers, newPlayer)
	}

	return newPlayers, replaceDrunk
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

//func initStatus(players []model.Player, replaceDrunk string) []model.Player {
//	for i, player := range players {
//		players[i].Ready.Nominate = true
//		players[i].Ready.Nominated = true
//		players[i].Ready.Vote = true
//		switch player.Character {
//		case Imp:
//			players[i].State.Evil = true
//			players[i].State.Demon = true
//		case Poisoner:
//			players[i].State.Evil = true
//		case ScarletWoman:
//			players[i].State.Evil = true
//		case Baron:
//			players[i].State.Evil = true
//		case Virgin:
//			players[i].State.Blessed = true
//		case Slayer:
//			players[i].State.Bullet = true
//		case Recluse:
//			players[i].State.Evil = true
//			evilCharacter := getRandEvilCharacter()
//			players[i].State.RegardedAs = evilCharacter
//			if evilCharacter == Spy {
//				players[i].State.Evil = false
//			}
//			players[i].State.RegardedAsSaved = evilCharacter
//			if players[i].State.RegardedAsSaved == Imp {
//				players[i].State.Demon = true
//			}
//		case Drunk:
//			players[i].CharacterType = Townsfolk
//			players[i].Character = replaceDrunk
//			players[i].State.Drunk = true
//			if replaceDrunk == Slayer {
//				players[i].State.Bullet = true
//			}
//		case FortuneTeller:
//			for {
//				randIdx := rand.Intn(len(players))
//				if players[randIdx].CharacterType == Townsfolk || players[randIdx].CharacterType == Outsiders {
//					players[randIdx].State.Demon = true
//					break
//				}
//			}
//		case Spy:
//			players[i].State.RegisteredAsType = Minions
//		}
//	}
//
//	return players
//}

func getRandEvilCharacter() string {
	evils := append(MinionsPool, DemonsPool...)
	randInt := rand.Intn(len(evils))
	return evils[randInt]
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
