package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"log"
)

// 弄臣技能，首次将要死亡的时候不会死亡,但是如果此时醉酒或中毒就会死亡。
// 如果其他角色的能力保护弄臣免于死亡，弄臣不会使用掉自己的能力。只有当弄臣将要真正的死去时，他的能力才会触发。
func WillNotDie(character model.BaseCharacter, curTime int) (string, bool) {
	msg := ""
	//确认弄臣是否被技能保护以及保护时间
	if util.FindElementInSlice(model.ProtectedStatus, character.CharacterStatus) {
		//弄臣被别人的技能保护
		if character.CharacterStatus[model.ProtectedStatus] > curTime {
			log.Println("弄臣由于被技能保护，免于死亡，且自身技能仍然有效")
			msg += "弄臣由于被技能保护，免于死亡，且自身技能仍然有效\n"
			return msg, false
		}
	}
	//醉酒或中毒会直接死亡
	if util.FindElementInSlice(model.PoisonedStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.PoisonedStatus] > curTime {
			log.Println("弄臣中毒状态，技能释放失败")
			msg += "弄臣中毒状态，技能释放失败\n"
			return msg, false
		}
	}
	if util.FindElementInSlice(model.DrunkStatus, character.CharacterStatus) {
		if character.CharacterStatus[model.DrunkStatus] > curTime {
			log.Println("弄臣醉酒状态，技能释放失败")
			msg += "弄臣醉酒状态，技能释放失败\n"
			return msg, false
		}
	}
	if model.OncePerGameSkillMap[character.CharacterName] {
		log.Println("弄臣自身技能触发。")
		model.OncePerGameSkillMap[character.CharacterName] = false
		msg += "弄臣自身技能触发，免于死亡\n"
		return msg, false
	}
	log.Println("弄臣死亡")
	msg += "弄臣死亡\n"
	//character.IsDead = true
	//character.DeadReason = deadReason
	return msg, true
}
