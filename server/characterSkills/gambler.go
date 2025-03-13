package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"log"
)

// 赌徒每晚可以选择一位玩家并猜测其角色，如果猜错了就会死亡。如果赌徒醉酒或中毒，则不发动技能，即不会死。
func GamblerGuessSkill(target string, guess string, gambler *model.BaseCharacter, allCharacters []*model.Player) (string, bool) {
	msg := ""
	if gambler.IsDead {
		return msg, false
	}
	//赌徒技能施放状态修改 fixme 需要优化
	for i := range allCharacters {
		if allCharacters[i].BaseCharacter.CharacterName == model.GoodCharacterMap[11] {
			allCharacters[i].State.Casted = true
			break
		}
	}
	//赌徒猜错了，但是处于中毒或醉酒状态，无事发生
	for i := range allCharacters {
		if allCharacters[i].Id == target {
			if allCharacters[i].BaseCharacter.CharacterName != guess && (util.IsCharacterDrunk(allCharacters[i].BaseCharacter) || util.IsCharacterDrunk(allCharacters[i].BaseCharacter)) {
				log.Println("赌徒猜错了，但是由于处于醉酒或中毒状态，无事发生")
				msg += "赌徒猜错了，但是由于处于醉酒或中毒状态，无事发生"
				return msg, false
				//赌徒猜对了，无事发生
			} else if allCharacters[i].BaseCharacter.CharacterName == guess {
				log.Println("赌徒猜对了，无事发生")
				msg += "赌徒猜对了，无事发生"
				return msg, false
				//赌徒猜错了，死亡
			} else {
				log.Println("赌徒猜错，死亡")
				gambler.IsDead = true
				msg += "赌徒猜错，死亡"
				return msg, true
			}
		}
	}
	return "", false
}
