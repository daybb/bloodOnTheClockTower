package handler

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"log"
	"math"
	"sync"
)

// 执行有顺序性，不可修改执行顺序
func checkoutNight(mux *sync.Mutex, game *model.Room) {
	mux.Lock()
	defer mux.Unlock()

	//var msgPlayer = "您"
	//var msgAll = ""

	// 承载技能释放者对象的池
	castPoolObj := map[*model.Player][]int{}
	for fromPlayerId, toPlayerIdSlice := range game.CastPool {
		for i, player := range game.Players {
			if player.Id == fromPlayerId {
				castPoolObj[&game.Players[i]] = []int{}
				for _, toPlayerId := range toPlayerIdSlice {
					for _, player := range game.Players {
						if player.Id == toPlayerId {
							castPoolObj[&game.Players[i]] = append(castPoolObj[&game.Players[i]], player.Index)
							break
						}
					}
				}
				break
			}
		}
	}
	// 结算第一夜信息
	if game.State.Stage == 0 {
		//FirstNight(game)
		msg := ""
		msg += "昨夜是 平安夜\n"
		for i := range game.Players {
			game.Players[i].Log += msg
		}
		game.Log += msg
		game.State.Stage++
		game.State.Night = false
		game.CurTime = model.Morning
		broadcast(game)
		return
	}
	// 结算除第一夜信息
	if game.Result == "" {
		game.State.Night = false
		game.State.Stage += 1
		msg := fmt.Sprintf("第%d天，入夜~\n", game.State.Day+1)
		OtherNight(game)
		game.Log += msg
		broadcast(game)
	}
}

// 白天切换到夜晚的时候重置所有技能施放状态
func resetSkillCast(game *model.Room) {
	for i := range game.Players {
		if !game.Players[i].BaseCharacter.IsDead && (game.Players[i].BaseCharacter.CharacterName == "Gambler" || game.Players[i].BaseCharacter.CharacterName == "Shabaloth") {
			game.Players[i].State.Casted = false
		} else {
			game.Players[i].State.Casted = true
		}
	}
}

func checkoutDay(mux *sync.Mutex, game *model.Room) {
	mux.Lock()
	defer mux.Unlock()
	// 结算处决
	execute(game)
	// 重置技能施放状态
	resetSkillCast(game)
	// 结算本局
	//checkout(game, game.Executed)
}

// checkout 结算本局
func checkout(game *model.Room, executed *model.Player) {
	msg := ""
	var realDemonCount int     // 恶魔数量，被占卜认定的不算
	var hasSlayerBullet bool   // 有杀手且杀手有子弹
	var aliveCount int         // 活人数量
	var canVote int            // 可投票数量
	var evilAliveCount int     // 邪恶玩家存活数量
	var mayorAlive bool        // 市长是否存活
	var scarletWomanAlive bool // 魅魔是否存活
	var poisonerAlive bool     // 下毒者是否存活
	var demonCount int         // 恶魔数量（不论死活）
	for _, player := range game.Players {
		// fixme 对应平民胜利条件1
		if player.BaseCharacter.CharacterKind == model.Devil && !player.State.Dead {
			realDemonCount++
		}
		// 可投票数
		if player.Ready.Vote {
			canVote += 1
		}
	}
	// 平民胜利条件1（恶魔受不了了自杀情况），这里有三种铲除恶魔的可能：1、杀手，2、处决，3、自刀。
	// 处决在结算投票时判定，枪杀在枪手施法后判定，自刀在判定刀人时判定，所以realDemonCount不可能为0
	// 魅魔再判定为双保险可删没测
	if game.Result == "" && realDemonCount == 0 && (!scarletWomanAlive || scarletWomanAlive && aliveCount < 5) {
		msg += "达成平民胜利条件一：恶魔被铲除干净\n"
		msg += "本局结束，平民胜利\n"
		game.Result = "平民阵营胜利"
	}
	// 平民胜利条件2
	if game.Result == "" && aliveCount == 3 && mayorAlive && executed == nil && !game.State.Night && !poisonerAlive {
		msg += "达成平民胜利条件二：白天剩三人，其中一个是市长，且当日无人被处决，且三人不是下毒者市长小恶魔的组合\n"
		msg += "本局结束，平民胜利\n"
		game.Result = "平民阵营胜利"
	}
	// 邪恶胜利条件2
	if game.Result == "" && evilAliveCount == aliveCount {
		msg += "达成邪恶胜利条件二：平民阵营被屠城\n"
		msg += "本局结束，邪恶胜利\n"
		game.Result = "邪恶阵营胜利"
	}
	// 邪恶胜利条件3
	halfAlive := int(math.Ceil(float64(aliveCount / 2)))
	if game.Result == "" && aliveCount <= 4 && demonCount == 1 && canVote-evilAliveCount <= halfAlive && evilAliveCount >= halfAlive && !hasSlayerBullet && !mayorAlive {
		msg += "达成邪恶胜利条件三：活人数小于等于4，未发生爪牙转化为恶魔，平民可投的票数不大于活人的半数（向上取整），且存活的邪恶玩家数量不小于活人的半数（向上取整），且没有杀手或有杀手没有子弹，且没有市长或市长已死或酒鬼市长\n"
		msg += "本局结束，邪恶胜利\n"
		game.Result = "邪恶阵营胜利"
	}

	if game.Result != "" {
		// 拼接日志
		for i := range game.Players {
			game.Players[i].Log += msg
		}
		game.Log += msg
		game.Status = model.After
		// 发送game 以便前端跳转review
		broadcast(game)
	}
}

// 首夜给爪牙展示恶魔，给恶魔展示爪牙
func ShowDevilAndMinion(game *model.Room) (minionId []string, devilId []string) {
	for _, v := range game.Players {
		if v.BaseCharacter.CharacterKind == model.Devil {
			devilId = append(devilId, v.Id)
		} else if v.BaseCharacter.CharacterKind == model.Minion {
			minionId = append(minionId, v.Id)
		}
	}
	for k, v := range game.Players {
		if v.BaseCharacter.CharacterKind == model.Devil {
			game.Players[k].Log += fmt.Sprintf("爪牙是座位号为%v，请互相确认身份。\n", minionId)
			emit(game, v.Id)
		} else if v.BaseCharacter.CharacterKind == model.Minion {
			game.Players[k].Log += fmt.Sprintf("恶魔是座位号为%v，请互相确认身份。\n", devilId)
			emit(game, v.Id)
		}
	}
	log.Printf("首夜确认恶魔和爪牙身份。恶魔是%v,爪牙是%v", devilId, minionId)
	return minionId, devilId
}

// 首个夜晚。恶魔和爪牙互相展示身份
// 水手发动技能（暂无）
// 侍臣发动技能（暂无）
// 对教父展示所有外来者标记（暂无）
// 魔鬼代言人选择存活玩家（暂无）
// 普卡选择玩家（暂无）
// 给祖母展示玩家（暂无）
// 侍女技能发动（暂无）
// 睁眼
func FirstNight(game *model.Room) {
	// 恶魔和爪牙互相展示身份
	ShowDevilAndMinion(game)
	game.State.Night = true
	game.CurTime = model.Night
	// todo 水手技能
	// todo 侍臣技能
	// todo 教父技能
	// todo 魔鬼代言人技能
	// todo 普卡技能
	// todo 祖母技能
	// todo 侍女技能
	//游戏进度+1
	broadcast(game)
	return
}

// 直接入夜
func DirectToNight(game *model.Room) {
	game.State.Night = true
	game.CurTime = model.Night
	game.Executed = nil
	game.Nominated = nil
	game.VotePool = map[string]int{}
	game.VoteLogs = map[string]string{}
	game.CurLogs = ""
	game.Log += "入夜"
	broadcast(game)
	return
}

// 夜晚 纯粹的状态变更，不含技能释放
func OtherNight(game *model.Room) {
	//设定current time
	game.CurTime = 23
	//距离游戏开始的时间。游戏是从晚上开始的，第二个晚上就是24，第三个晚上是48
	fromStartTime := game.State.Stage * 24
	//消除可以消除的负面状态或者守护效果
	for i := range game.Players {
		if len(game.Players[i].BaseCharacter.CharacterStatus) > 0 {
			for k, v := range game.Players[i].BaseCharacter.CharacterStatus {
				if v <= fromStartTime {
					delete(game.Players[i].BaseCharacter.CharacterStatus, k)
				}
			}
		}
	}

}
