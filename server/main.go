package main

import (
	"bloodOnTheClockTower/handler"
	"bloodOnTheClockTower/model"
	"fmt"
	"log"
	"net/http"
)

var character = []model.BaseCharacter{}
var characterMap = map[int]*model.BaseCharacter{}

// run main start a websocket server
// create room create a room from request body, can do it by default config instead
func main() {
	http.HandleFunc("/home", handler.LoadHome)
	http.HandleFunc("/ws", handler.HandleWebSocket)
	// 等待开始页
	http.HandleFunc("/room/", handler.LoadRoom)
	// 游戏中页
	http.HandleFunc("/game/", handler.LoadGame)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// set up characters,八人局为例子，随机抽取村民5人、外来者1人，爪牙1人，恶魔1人
func assign() ([]model.BaseCharacter, map[int]*model.BaseCharacter) {
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
	characterMap := make(map[int]*model.BaseCharacter)
	for k, v := range fixedGroupMap[1] {
		character = append(character, model.BaseCharacter{
			Id:              k,
			CharacterName:   v,
			CharacterKind:   model.CharacterKindMap[v],
			CharacterStatus: nil,
			IsDead:          false,
		})
		characterMap[k] = &model.BaseCharacter{
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
