package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"fmt"
	"log"
)

// 沙巴洛斯每个夜晚能够吞食两名玩家，但是可能会在下一个夜晚将其中一名吐出来
// 攻击茶艺师及临近玩家时临近玩家不会死亡，茶艺师会死亡
// 因反刍而复活的玩家重新获得自己的能力，即使是已经使用过的“每局游戏限一次”的能力也会重新获得。如果有“在你的首个夜晚”的能力，玩家可以再次使用该能力。
// todo 研究何时复活玩家 最好只是让沙巴洛斯偶尔进行反刍。每局游戏一次，或者两次，通常已经足够了。
func KillTwo(shabaloth *model.BaseCharacter, target []string, allCharacters []*model.Player, curTime int, timeNow int) (error, string) {
	msg := ""
	//沙巴洛斯技能施放状态修改 fixme 需要优化
	for i := range allCharacters {
		if allCharacters[i].BaseCharacter.CharacterName == model.DevilCharacterMap[3] {
			allCharacters[i].State.Casted = true
			break
		}
	}
	//如果沙巴罗斯醉酒，技能失效
	if util.FindElementInSlice(model.PoisonedStatus, shabaloth.CharacterStatus) {
		if shabaloth.CharacterStatus[model.PoisonedStatus] > curTime {
			//沙巴罗斯中毒状态，技能发动失败
			log.Println("沙巴罗斯中毒状态，杀人失败")
			RegurgitateOne(shabaloth, allCharacters, curTime)
			msg = "沙巴罗斯中毒状态，杀人失败"
			return nil, msg
		}
	}
	if util.FindElementInSlice(model.DrunkStatus, shabaloth.CharacterStatus) {
		if shabaloth.CharacterStatus[model.DrunkStatus] > curTime {
			//沙巴罗斯醉酒状态，技能发动失败
			log.Println("沙巴罗斯醉酒状态，杀人失败")
			RegurgitateOne(shabaloth, allCharacters, curTime)
			msg = "沙巴罗斯中毒状态，杀人失败"
			return nil, "沙巴罗斯中毒状态，杀人失败"
		}
	}
outerloop:
	for i := range target {
		for j := range allCharacters {
			if allCharacters[j].Id == target[i] {
				//主动被别人的技能保护成功
				if util.FindElementInSlice(model.ProtectedStatus, allCharacters[j].BaseCharacter.CharacterStatus) {
					//todo 确认保护是次数限制还是时间限制
					if allCharacters[j].BaseCharacter.CharacterStatus[model.ProtectedStatus] > curTime {
						log.Printf("%v由于被技能保护，免于死亡", allCharacters[j].BaseCharacter.CharacterName)
						break
					}
				}
				//被动被茶艺师保护，且茶艺师不处于异常状态，不会死亡
				neighbors := allCharacters[j].BaseCharacter.Neighbors
				for _, neighbor := range neighbors {
					if neighbor.CharacterName == "TeaLady" {
						//茶艺师中毒，保护失败
						if util.FindElementInSlice(model.PoisonedStatus, neighbor.CharacterStatus) {
							if shabaloth.CharacterStatus[model.PoisonedStatus] > curTime {
								log.Println("茶艺师中毒状态，保护失败")
								allCharacters[j].BaseCharacter.IsDead = true
								allCharacters[j].BaseCharacter.DeadReason = model.DeadByShabaloth
								allCharacters[j].BaseCharacter.DeadTime = timeNow
								allCharacters[j].BaseCharacter.ExactTime = curTime
								continue outerloop
							}
						}
						//茶艺师醉酒，保护失败
						if util.FindElementInSlice(model.DrunkStatus, neighbor.CharacterStatus) {
							if shabaloth.CharacterStatus[model.DrunkStatus] > curTime {
								log.Println("茶艺师醉酒状态，保护失败")
								allCharacters[j].BaseCharacter.IsDead = true
								allCharacters[j].BaseCharacter.DeadReason = model.DeadByShabaloth
								allCharacters[j].BaseCharacter.DeadTime = timeNow
								allCharacters[j].BaseCharacter.ExactTime = curTime
								continue outerloop
							}
						}
						log.Println("被茶艺师保护，免死")
						continue outerloop
					}

				}
				//被自身技能保护，如弄臣
				if allCharacters[j].BaseCharacter.CharacterName == "Fool" && model.OncePerGameSkillMap[shabaloth.CharacterName] {
					log.Println("弄臣自身技能触发。")
					model.OncePerGameSkillMap[shabaloth.CharacterName] = false
					continue outerloop
				}
				//杀人成功
				allCharacters[j].BaseCharacter.IsDead = true
				allCharacters[j].BaseCharacter.DeadReason = model.DeadByShabaloth
				allCharacters[j].BaseCharacter.DeadTime = timeNow
				allCharacters[j].BaseCharacter.ExactTime = curTime
				log.Printf("沙巴罗斯杀人成功，%v死亡", allCharacters[j].BaseCharacter.CharacterName)
				msg += fmt.Sprintf("沙巴罗斯杀人成功，%v死亡", allCharacters[j].BaseCharacter.CharacterName)
			}
		}

	}
	regurgitateMsg := RegurgitateOne(shabaloth, allCharacters, curTime)
	if regurgitateMsg != "" {
		msg += fmt.Sprintf("沙巴罗斯技能复活成功，%v复活", regurgitateMsg)
	}
	return nil, msg
}

// 复活一个人
func RegurgitateOne(shabaloth *model.BaseCharacter, allCharacters []*model.Player, curTime int) string {
	if model.OncePerGameSkillMap[shabaloth.CharacterName] {
		for _, v := range allCharacters {
			if v.BaseCharacter.DeadReason == model.DeadByShabaloth && curTime-v.BaseCharacter.ExactTime == 24 {
				//todo 研究谁该被救回来合适
				v.BaseCharacter.CharacterStatus = nil
				v.BaseCharacter.DeadReason = 0
				v.BaseCharacter.DeadTime = 0
				v.BaseCharacter.ExactTime = 0
				v.BaseCharacter.IsDead = false
				v.BaseCharacter.Neighbors = nil //todo set neighbors
				model.OncePerGameSkillMap[v.BaseCharacter.CharacterName] = true
				log.Printf("%v由于沙巴罗斯的技能复活", v.BaseCharacter.CharacterName)
				return v.BaseCharacter.CharacterName
			}
		}
	}
	return ""
}
