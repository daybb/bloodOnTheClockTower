package handler

import (
	"bloodOnTheClockTower/characterSkills"
	"bloodOnTheClockTower/model"
	"fmt"
	"sync"
)

// 前端控制技能发动顺序
func cast(mux *sync.Mutex, game *model.Room, playerId string, targets []string, extra string) {
	mux.Lock()
	defer mux.Unlock()

	var msgPlayer = "您"
	var msgAll = ""
	var characterToCast *model.BaseCharacter
	var allCharacters []*model.Player
	var err error
	for i := range game.Players {
		allCharacters = append(allCharacters, &game.Players[i])
		if game.Players[i].Id == playerId {
			characterToCast = &game.Players[i].BaseCharacter
		}
	}
	switch game.AllCharacters[playerId] {
	//水手发动（暂无）
	case "Sailor":
		return
		//旅店老板发动（暂无）
	case "Innkeeper":
		return
		// 侍臣发动（暂无）
	case "Courtier":
		return
		// 赌徒发动
	case "Gambler":
		//是否成功应该是白天来确定的
		var dead bool
		var pos int
		var curLog string
		msgAll, dead = characterSkills.GamblerGuessSkill(targets[0], extra, characterToCast, allCharacters)
		for i := range game.Players {
			if game.Players[i].Id == targets[0] {
				pos = game.Players[i].PositionId
			}
		}
		msgPlayer += fmt.Sprintf("发动技能，赌%v的身份是%v", pos, extra)
		if dead {
			curLog += fmt.Sprintf("发动技能，赌%v的身份是%v，结果死亡", pos, extra)
		} else {
			curLog += fmt.Sprintf("发动技能，赌%v的身份是%v，结果无事发生", pos, extra)
		}
		// 魔鬼代言人发动（暂无）
	case "DevilsAdvocate":
		return
		// 疯子发动（暂无）
	case "Lunatic":
		return
		// 驱魔人发动（暂无）
	case "Exorcist":
		return
		// 僵怖发动（暂无）
	case "Zombuul":
		return
		// 普卡发动（暂无）
	case "Pukka":
		return
		// 沙巴罗斯发动
	case "Shabaloth":
		err, msgAll = characterSkills.KillTwo(characterToCast, targets, allCharacters, game.CurTime, model.Night)
		// 魄发动（暂无）
	case "Po":
		return
		// 刺客发动
	case "Assassin":
		msgAll, err = characterSkills.KillAnyway(characterToCast, targets[0], allCharacters, game.CurTime, model.Night)
		// 教父发动（暂无）
	case "Godfather":
		return
		// 教授发动
	case "Professor":
		msgAll, err = characterSkills.ResurrectDead(characterToCast, targets[0], allCharacters)
		if err != nil {
			return
		}
		// 造谣者发动（暂无）
	case "Gossip":
		return
		// 修补匠发动
	case "Tinker":
		msgAll = characterSkills.DieNoReason(characterToCast, game.CurTime, model.Night)
		// 月之子发动（暂无）
	case "Moonchild":
		return
		// 孙子祖母发动（暂无）
	case "Grandmother":
		return
		// 侍女发动（暂无）
	case "Chambermaid":
		return
	default:
		return
	}

	for i, player := range game.Players {
		if player.Id == playerId {
			game.Players[i].Log += msgPlayer
			break
		}
	}
	game.Log += msgAll
	// 发送日志
	emit(game, playerId)
	// 发送game 更新casted状态
	broadcast(game)
}
