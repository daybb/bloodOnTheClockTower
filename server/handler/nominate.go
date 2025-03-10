package handler

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"log"
	"sync"
)

// 提名，被提名玩家会加入投票池
// 同一时间只能有一名玩家被提名。如果一项提名已经被发起但处决投票尚未结束，则不能再提出另一项提名。
// 只有存活的玩家可以发起提名。虽然死亡的玩家也可以被提名，但这在绝大多数情况下都不是一个明智的举动。
// 每名玩家每天只能发起一次提名，并且每名玩家每天只能被提名一次。
func nominate(mux *sync.Mutex, game *model.Room, playerId string, targets []string) {
	mux.Lock()
	defer mux.Unlock()
	if game.CurTime != model.Morning {
		log.Println("非白天，不能提名")
		return
	}
	// 如果有处决者产生 不能提名
	if game.Executed != nil {
		msgPlayer := "本轮已处决过人，您的提名无效\n"
		for i, player := range game.Players {
			if player.Id == playerId {
				game.Players[i].Log += msgPlayer
				break
			}
		}
		// 发送日志
		emit(game, playerId)
		return
	}

	var msg = ""

	for i, player := range game.Players {
		if player.Id == playerId && player.Ready.Nominate && !player.BaseCharacter.IsDead && !game.State.VotingStep {
			for j, target := range game.Players {
				if targets[0] == target.Id && target.Ready.Nominated { // 死了也能被提名
					msg += fmt.Sprintf("[%s] ", player.Name)
					game.Players[i].Ready.Nominate = false  // 发动提名者不能再提名
					game.Players[j].Ready.Nominated = false // 被提名者不能再被提名
					game.Nominated = &target
					game.VotePool[game.Nominated.Id] = 0
					msg += fmt.Sprintf("提名 [%s] 进行处决公投\n", target.Name)
					break
				}
			}
			break
		}
	}
	for i := range game.Players {
		game.Players[i].Log += msg
	}
	game.Log += msg
	//进入投票处决环节
	game.State.VotingStep = true
	// 发送日志
	broadcast(game)
}
