package characterSkills

import (
	"bloodOnTheClockTower/model"
	"log"
)

// 修补匠技能，随时可以死亡，不需要理由
// 虽然你可以在白天冷不丁地杀死修补匠，但是修补匠还是死在夜晚会更有意思，因为玩家们会想知道修补匠到底是因为自己的能力死亡的，还是因为别的什么原因。
// 你可以一直不杀死修补匠。这会让他看起来非常可疑。
// 如果杀死修补匠会导致游戏结束，那我们推荐你永远不要在这种情况下杀死他。游戏的获胜与否应当靠玩家自己的努力，而不是看说书人的脸色。
func DieNoReason(character *model.BaseCharacter, curTime int, deadTime int) {
	if character.CharacterStatus[model.ProtectedStatus] > curTime {
		log.Println("修补匠由于被技能保护，免于死亡")
		return
	}
	log.Println("修补匠被说书人杀掉")
	character.IsDead = true
	character.DeadReason = model.KilledByStoryTeller
	character.DeadTime = deadTime //死亡时间，早上、黄昏、晚上
	return
}
