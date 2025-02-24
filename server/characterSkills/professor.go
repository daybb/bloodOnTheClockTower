package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"errors"
	"fmt"
	"log"
)

// 教授每局游戏一次复活死亡玩家的机会,被复活玩家会重新获取能力，即使是一局一次且已使用过
// 教授醉酒或中毒不能发动技能成功，且不能再次使用技能
func ResurrectDead(character *model.BaseCharacter, toRescue string, allCharacters []*model.Player) (string, error) {
	msg := ""
	if character.IsDead {
		return msg, errors.New("已死亡，无法使用技能")
	}
	if !model.OncePerGameSkillMap[character.CharacterName] {
		return msg, errors.New("已使用过复活技能，不能重复使用")
	}
	//如果教授醉酒或中毒，技能被消耗且无事发生
	if util.IsCharacterDrunk(*character) || util.IsCharacterPoisoned(*character) {
		model.OncePerGameSkillMap[character.CharacterName] = false
		log.Println("教授中毒或醉酒，技能发动失败")
		return msg, nil
	}
	for i := range allCharacters {
		if allCharacters[i].Id == toRescue {
			//如果玩家非村民，技能失效
			if allCharacters[i].BaseCharacter.CharacterKind != model.Good {
				model.OncePerGameSkillMap[character.CharacterName] = false
				log.Println("选中玩家非村民，技能发动失败")
				msg += "选中玩家非村民，技能发动失败"
				return msg, nil
			}
		}
		//成功复活
		allCharacters[i].BaseCharacter.IsDead = false
		allCharacters[i].BaseCharacter.CharacterStatus = nil
		if exist := model.OncePerGameSkillMap[allCharacters[i].BaseCharacter.CharacterName]; exist {
			log.Println("教授技能发动成功，玩家重新获取技能")
			msg += fmt.Sprintf("教授技能发动成功，%v重新获取技能", allCharacters[i].BaseCharacter.CharacterName)
			model.OncePerGameSkillMap[allCharacters[i].BaseCharacter.CharacterName] = true
			//todo 每局一次的技能恢复
		}
		log.Printf("%v复活", allCharacters[i].BaseCharacter.CharacterName)
		msg += fmt.Sprintf("%v复活", allCharacters[i].BaseCharacter.CharacterName)
	}

	return msg, nil
}
