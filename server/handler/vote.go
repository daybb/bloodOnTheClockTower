package handler

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"sync"
)

// 当前提名的投票
// 每名存活的玩家每天可以投票给任意数量的玩家。
// 每名死亡的玩家在他死亡后，只剩最后一次投票机会。
func vote(mux *sync.Mutex, game *model.Room, playerId string) {
	mux.Lock()
	defer mux.Unlock()

	var msgAll = ""
	var msgPlayer = "您"

	for i, player := range game.Players {
		if player.Id == playerId && player.Ready.Vote && game.State.VotingStep {
			msgAll += fmt.Sprintf("[%s] ", player.Name)
			if game.Nominated != nil && game.Players[i].Ready.Vote {
				if player.BaseCharacter.IsDead {
					game.Players[i].Ready.Vote = false // 死人投了就不能再投票了
				}
				game.VotePool[game.Nominated.Id] += 1
				msgPlayer += fmt.Sprintf("决意投给 [%s] \n", game.Nominated.Name)
				game.Players[i].Log += msgPlayer
				msgAll += fmt.Sprintf("投票 [%s] 成功\n", game.Nominated.Name)
				// 总日志加入票池
				game.VoteLogs[game.Nominated.Id] += msgAll

				// 发送个人日志
				emit(game, playerId)
				break
			}
		}
	}
}
