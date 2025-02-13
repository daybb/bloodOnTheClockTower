package characterSkills

import (
	"bloodOnTheClockTower/model"
	"bloodOnTheClockTower/util"
	"errors"
	"log"
)

// 教授每局游戏一次复活死亡玩家的机会,被复活玩家会重新获取能力，即使是一局一次且已使用过
// 教授醉酒或中毒不能发动技能成功，且不能再次使用技能
func ResurrectDead(character *model.BaseCharacter, target *model.BaseCharacter) error {
	if character.IsDead {
		return errors.New("已死亡，无法使用技能")
	}
	if !model.OncePerGameSkillMap[character.CharacterName] {
		return errors.New("已使用过复活技能，不能重复使用")
	}
	if !target.IsDead {
		return errors.New("角色未死亡，重新选择")
	}
	//如果教授醉酒或中毒，技能被消耗且无事发生
	if util.IsCharacterDrunk(character) || util.IsCharacterPoisoned(character) {
		model.OncePerGameSkillMap[character.CharacterName] = false
		log.Println("教授中毒或醉酒，技能发动失败")
		return nil
	}
	//如果玩家非村民，技能失效
	if target.CharacterKind != model.Good {
		model.OncePerGameSkillMap[character.CharacterName] = false
		log.Println("选中玩家非村民，技能发动失败")
		return nil
	}
	//成功复活
	target.IsDead = false
	target.CharacterStatus = nil
	if exist := model.OncePerGameSkillMap[target.CharacterName]; exist {
		log.Println("教授技能发动成功，玩家重新获取技能")
		model.OncePerGameSkillMap[target.CharacterName] = true
		//todo 每局一次的技能恢复
	}
	log.Printf("%v复活", target.CharacterName)
	return nil
}
