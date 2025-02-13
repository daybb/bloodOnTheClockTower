package util

import "bloodOnTheClockTower/model"

// 判断角色是否醉酒
func IsCharacterDrunk(character *model.BaseCharacter) bool {
	if character.CharacterStatus == nil {
		return false
	}
	return FindElementInSlice(model.DrunkStatus, character.CharacterStatus)
}

// 判断角色是否中毒
func IsCharacterPoisoned(character *model.BaseCharacter) bool {
	if character.CharacterStatus == nil {
		return false
	}
	return FindElementInSlice(model.PoisonedStatus, character.CharacterStatus)
}
