package characterSkills

import (
	"bloodOnTheClockTower/model"
)

// 吟游诗人技能，当爪牙死于处决，除了本人和旅行者外所有玩家醉酒到明天黄昏
// 白天早上8点-黄昏下午17点-夜晚晚上21点
// 发生时间一定是早上
func DrunkenEveryone(executed *model.BaseCharacter, characters []*model.BaseCharacter, currentTime int) {
	if executed.CharacterKind == model.Minion && executed.DeadReason == model.DeadReasonExecuted {
		for i := range characters {
			//即使死亡也要上状态，以免复活
			//如果角色已有醉酒状态，判断状态时长
			if characters[i].CharacterStatus[model.DrunkStatus] != 0 && characters[i].CharacterStatus[model.DrunkStatus] < currentTime+33 {
				characters[i].CharacterStatus[model.DrunkStatus] = currentTime + 33
			}
			//如果没有醉酒状态，加上醉酒状态
			if characters[i].CharacterStatus[model.DrunkStatus] == 0 {
				characters[i].CharacterStatus[model.DrunkStatus] = currentTime + 33
			}
		}
	}
}
