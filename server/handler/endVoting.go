package handler

import (
	"bloodOnTheClockTower/model"
	"sync"
)

func endVoting(mux *sync.Mutex, game *model.Room) {
	mux.Lock()
	defer mux.Unlock()

	if !game.State.VotingStep {
		return
	}

	var msg string

	// 打印所有投票成功的票型
	msg += game.VoteLogs[game.Nominated.Id]
	for i := range game.Players {
		game.Players[i].Log += msg
	}
	game.Log += msg
	// 退出投票处决环节
	game.State.VotingStep = false
	// 发送日志
	broadcast(game)
}
