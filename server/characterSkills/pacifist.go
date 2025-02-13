package characterSkills

import (
	"bloodOnTheClockTower/model"
	"log"
)

// 和平主义者，技能是被处决的善良玩家可能不会死亡
// 每局游戏触发一次和平主义者的能力，通常来说是合理的。
// 只要你觉得合适，也可以多触发几次。
// 在少见的情形下，为了让和平主义者看起来很可疑，你可以永不触发他的能力。
// todo 研究什么时候触发，触发几次合适，先设置一局一次
func ExecutedGoodStillLive(executedCharacter *model.BaseCharacter) {
	if executedCharacter.CharacterKind == model.Good || executedCharacter.CharacterKind == model.Outsider && model.OncePerGameSkillMap["Pacifist"] {
		log.Println("和平主义者发动技能，被处决的好人仍然存活")
		model.OncePerGameSkillMap["Pacifist"] = false
		return
	}
}
