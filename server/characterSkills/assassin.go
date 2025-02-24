package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"errors"
	"fmt"
	"log"
)

// 刺客技能，你可以选择一名玩家：他死亡，即使因为任何原因让他不会死亡。
func KillAnyway(assassin *model.BaseCharacter, target string, allCharacters []*model.Player, curTime int, timeNow int) (string, error) {
	msg := ""
	//刺客技能受中毒和醉酒影响
	if timeNow != model.Night {
		return msg, errors.New("非夜晚，不能发动该技能")
	}
	if !model.OncePerGameSkillMap[assassin.CharacterName] {
		return msg, errors.New("技能已使用，无法再次使用")
	}
	//刺客中毒状态，技能发动失败，但技能已使用，所以失效
	if util.FindElementInSlice(model.PoisonedStatus, assassin.CharacterStatus) {
		if assassin.CharacterStatus[model.PoisonedStatus] > curTime {
			log.Println("刺客中毒状态，行刺失败且技能被消耗")
			model.OncePerGameSkillMap[assassin.CharacterName] = false
			msg += "刺客中毒状态，行刺失败且技能被消耗"
			return msg, nil
		}
	}
	//刺客醉酒状态，技能发动失败，但技能已使用，所以失效
	if util.FindElementInSlice(model.DrunkStatus, assassin.CharacterStatus) {
		if assassin.CharacterStatus[model.DrunkStatus] > curTime {
			log.Println("刺客醉酒状态，行刺失败且技能被消耗")
			model.OncePerGameSkillMap[assassin.CharacterName] = false
			msg += "刺客醉酒状态，行刺失败且技能被消耗"
			return "", nil
		}
	}
	//刺客刺杀成功
	for i := range allCharacters {
		if allCharacters[i].Id == target {
			allCharacters[i].BaseCharacter.IsDead = true
			allCharacters[i].BaseCharacter.DeadTime = timeNow
			allCharacters[i].BaseCharacter.DeadReason = model.Executed
			log.Printf("刺客刺杀%v成功，%v死亡", allCharacters[i].BaseCharacter.CharacterName, allCharacters[i].BaseCharacter.CharacterName)
			msg += fmt.Sprintf("刺客刺杀%v成功，%v死亡", allCharacters[i].BaseCharacter.CharacterName, allCharacters[i].BaseCharacter.CharacterName)
			return "", nil
		}
	}

	return "", nil
}
