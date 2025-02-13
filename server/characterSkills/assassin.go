package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"errors"
	"log"
)

// 刺客技能，你可以选择一名玩家：他死亡，即使因为任何原因让他不会死亡。
func KillAnyway(character, target *model.BaseCharacter, curTime int, timeNow int) error {
	//刺客技能受中毒和醉酒影响
	if timeNow != model.Night {
		return errors.New("非夜晚，不能发动该技能")
	}
	if !model.OncePerGameSkillMap[character.CharacterName] {
		return errors.New("技能已使用，无法再次使用")
	}
	if util.FindElementInSlice(model.PoisonedStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.PoisonedStatus] > curTime {
			//刺客中毒状态，技能发动失败，但技能已使用，所以失效
			log.Println("刺客中毒状态，行刺失败且技能被消耗")
			model.OncePerGameSkillMap[character.CharacterName] = false
			return nil
		}
	}
	if util.FindElementInSlice(model.DrunkStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.DrunkStatus] > curTime {
			//刺客醉酒状态，技能发动失败，但技能已使用，所以失效
			log.Println("刺客醉酒状态，行刺失败且技能被消耗")
			model.OncePerGameSkillMap[character.CharacterName] = false
			return nil
		}
	}
	//刺客刺杀成功
	target.IsDead = true
	target.DeadTime = timeNow
	target.DeadReason = model.Executed
	log.Printf("刺客刺杀%v成功，%v死亡", target.CharacterName, target.CharacterName)
	return nil
}
