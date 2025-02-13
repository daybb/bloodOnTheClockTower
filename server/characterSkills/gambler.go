package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"errors"

	"log"
)

// 赌徒每晚可以选择一位玩家并猜测其角色，如果猜错了就会死亡。如果赌徒醉酒或中毒，则不发动技能，即不会死。
func GamblerGuessSkill(character *model.BaseCharacter, guess string, gambler *model.BaseCharacter) error {
	if gambler.IsDead {
		return errors.New("你已死亡，不能猜角色")
	}
	//赌徒猜错了，但是处于中毒或醉酒状态，无事发生
	if character.CharacterName != guess && (util.IsCharacterDrunk(character) || util.IsCharacterDrunk(character)) {
		log.Println("赌徒猜错了，但是由于处于醉酒或中毒状态，无事发生")
		return nil
		//赌徒猜对了，无事发生
	} else if character.CharacterName == guess {
		log.Println("赌徒猜对了，无事发生")
		return nil
		//赌徒猜错了，死亡
	} else {
		log.Println("赌徒猜错，死亡")
		character.IsDead = true
		return nil
	}
}
