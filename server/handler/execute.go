package handler

import (
	"bloodOnTheClockTower/characterSkills"
	"bloodOnTheClockTower/model"
	"fmt"
	"math"
)

func execute(game *model.Room) {
	var msg string
	var executeeVoteCount int
	var allCharacter []*model.BaseCharacter
	var isDead bool
	game.Executed, executeeVoteCount = findExecutee(game)
	if game.Executed != nil {
		isDead = true
		//game.Executed.State.Dead = true
		//game.Executed.Ready.Nominate = false
		//game.Executed.Ready.Nominated = false
		msg += fmt.Sprintf("处决结果：[%s] 获得 %d 票，被处决，死亡\n", game.Executed.Name, executeeVoteCount)
		//弄臣发动技能
		if game.Executed.BaseCharacter.CharacterName == "Fool" {
			msgs, foolDead := characterSkills.WillNotDie(game.Executed.BaseCharacter, game.CurTime)
			msg += msgs
			isDead = foolDead
		}
		//和平主义者发动技能
		if game.Executed.BaseCharacter.CharacterKind == model.Good || game.Executed.BaseCharacter.CharacterKind == model.Outsider {
			msgs, isDeadAfter := characterSkills.ExecutedGoodStillLive(game.Executed.BaseCharacter)
			isDead = isDeadAfter
			msg += msgs
		}
		//吟游诗人发动技能
		if game.Executed.BaseCharacter.CharacterKind == model.Minion {
			for i := range game.Players {
				allCharacter = append(allCharacter, &game.Players[i].BaseCharacter)
			}
			//吟游诗人发动技能
			msgs, isMinstrelAct := characterSkills.DrunkenEveryone(&game.Executed.BaseCharacter, allCharacter, game.CurTime)
			if isMinstrelAct {
				msg += msgs
			}
		}
		if isDead {
			game.Executed.State.Dead = true
			game.Executed.Ready.Nominate = false
			game.Executed.Ready.Nominated = false
			game.Executed.BaseCharacter.IsDead = true
			game.Executed.BaseCharacter.ExactTime = game.CurTime
			game.Executed.BaseCharacter.DeadTime = model.Morning
			game.Executed.BaseCharacter.DeadReason = model.Executed
		}
	} else {
		msg += "处决结果：无人被处决\n"
	}
	for i := range game.Players {
		game.Players[i].Log += msg
	}
	game.Log += msg

	// 发送日志
	broadcast(game)
	// 立即结算
	checkout(game, game.Executed)
}

func findExecutee(game *model.Room) (*model.Player, int) {
	// 无人被提名
	if len(game.VotePool) == 0 {
		return nil, 0
	}
	var aliveCount int // 活人数量
	for _, player := range game.Players {
		if !player.BaseCharacter.IsDead {
			aliveCount++
		}
	}
	var halfAliveCount = int(math.Ceil(float64(aliveCount) / 2))
	var executeeId string
	var executeeVoteCount int
	var isHighestRepeated bool
	for nominatedId, voteCount := range game.VotePool {
		if voteCount >= halfAliveCount && voteCount == executeeVoteCount {
			isHighestRepeated = true
			continue
		}
		if voteCount >= halfAliveCount && voteCount > executeeVoteCount {
			executeeId = nominatedId
			executeeVoteCount = voteCount
		}
	}
	if executeeId != "" && !isHighestRepeated {
		for i, player := range game.Players {
			if player.Id == executeeId {
				return &game.Players[i], executeeVoteCount
			}
		}
	}
	return nil, executeeVoteCount
}
