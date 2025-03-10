package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"log"
)

// 弄臣技能，首次将要死亡的时候不会死亡,但是如果此时醉酒或中毒就会死亡。
// 如果其他角色的能力保护弄臣免于死亡，弄臣不会使用掉自己的能力。只有当弄臣将要真正的死去时，他的能力才会触发。
func WillNotDie(character *model.BaseCharacter, curTime int, deadReason int) {
	//确认弄臣是否被技能保护以及保护时间
	if util.FindElementInSlice(model.ProtectedStatus, character.CharacterStatus) {
		//弄臣被别人的技能保护
		if character.CharacterStatus[model.ProtectedStatus] > curTime {
			log.Println("弄臣由于被技能保护，免于死亡，且自身技能仍然有效")
			return
		}
	}
	if model.OncePerGameSkillMap[character.CharacterName] {
		log.Println("弄臣自身技能触发。")
		model.OncePerGameSkillMap[character.CharacterName] = false
		return
	}
	log.Println("弄臣死亡")
	character.IsDead = true
	character.DeadReason = deadReason
	return
}
