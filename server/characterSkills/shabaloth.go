package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"log"
)

// 沙巴洛斯每个夜晚能够吞食两名玩家，但是可能会在下一个夜晚将其中一名吐出来
// 攻击茶艺师及临近玩家时临近玩家不会死亡，茶艺师会死亡
// 因反刍而复活的玩家重新获得自己的能力，即使是已经使用过的“每局游戏限一次”的能力也会重新获得。如果有“在你的首个夜晚”的能力，玩家可以再次使用该能力。
// todo 研究何时复活玩家 最好只是让沙巴洛斯偶尔进行反刍。每局游戏一次，或者两次，通常已经足够了。
func KillTwo(target, allCharacters []*model.BaseCharacter, character *model.BaseCharacter, curTime int, timeNow int) error {
	//如果沙巴罗斯醉酒，技能失效
	if util.FindElementInSlice(model.PoisonedStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.PoisonedStatus] > curTime {
			//沙巴罗斯中毒状态，技能发动失败
			log.Println("沙巴罗斯中毒状态，杀人失败")
			RegurgitateOne(character, allCharacters, curTime)
			return nil
		}
	}
	if util.FindElementInSlice(model.DrunkStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.DrunkStatus] > curTime {
			//沙巴罗斯醉酒状态，技能发动失败
			log.Println("沙巴罗斯醉酒状态，杀人失败")
			RegurgitateOne(character, allCharacters, curTime)
			return nil
		}
	}
	for i := range target {
		if util.FindElementInSlice(model.ProtectedStatus, target[i].CharacterStatus) {
			//主动被别人的技能保护
			//todo 确认保护是次数限制还是时间限制
			if target[i].CharacterStatus[model.ProtectedStatus] > curTime {
				log.Printf("%v由于被技能保护，免于死亡", target[i].CharacterName)
				continue
			}
		}
		//被动被茶艺师保护，且茶艺师不处于异常状态，不会死亡
		neighbors := target[i].Neighbors
		for _, neighbor := range neighbors {
			if neighbor.CharacterName == "TeaLady" {
				//check茶艺师状态
				if util.FindElementInSlice(model.PoisonedStatus, neighbor.CharacterStatus) {
					if character.CharacterStatus[model.PoisonedStatus] > curTime {
						//茶艺师中毒状态，技能发动失败
						log.Println("茶艺师中毒状态，保护失败")
						target[i].IsDead = true
						target[i].DeadReason = model.DeadByShabaloth
						target[i].DeadTime = timeNow
						target[i].ExactTime = curTime
						continue
					}
				}
				if util.FindElementInSlice(model.DrunkStatus, neighbor.CharacterStatus) {
					if character.CharacterStatus[model.DrunkStatus] > curTime {
						//沙巴罗斯醉酒状态，技能发动失败
						log.Println("茶艺师醉酒状态，保护失败")
						target[i].IsDead = true
						target[i].DeadReason = model.DeadByShabaloth
						target[i].DeadTime = timeNow
						target[i].ExactTime = curTime
						continue
					}
				}
				log.Println("被茶艺师保护，免死")
				continue
			}
		}
		//被自身技能保护，如弄臣
		if target[i].CharacterName == "Fool" && model.OncePerGameSkillMap[character.CharacterName] {
			log.Println("弄臣自身技能触发。")
			model.OncePerGameSkillMap[character.CharacterName] = false
			continue
		}
		//杀人成功
		target[i].IsDead = true
		target[i].DeadReason = model.DeadByShabaloth
		target[i].DeadTime = timeNow
		target[i].ExactTime = curTime
		log.Printf("沙巴罗斯杀人成功，%v死亡", target[i].CharacterName)
	}
	RegurgitateOne(character, allCharacters, curTime)
	return nil
}

func RegurgitateOne(character *model.BaseCharacter, allCharacters []*model.BaseCharacter, curTime int) {
	if model.OncePerGameSkillMap[character.CharacterName] {
		for _, v := range allCharacters {
			if v.DeadReason == model.DeadByShabaloth && curTime-v.ExactTime == 24 {
				//todo 研究谁该被救回来合适
				v.CharacterStatus = nil
				v.DeadReason = 0
				v.DeadTime = 0
				v.ExactTime = 0
				v.IsDead = false
				v.Neighbors = nil //todo set neighbors
				model.OncePerGameSkillMap[v.CharacterName] = true
				log.Printf("%v由于沙巴罗斯的技能复活", v.CharacterName)
				return
			}
		}
	}
}
