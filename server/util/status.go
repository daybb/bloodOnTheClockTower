package util

import "bloodOnTheClockTower/model"

// 判断角色是否醉酒
func IsCharacterDrunk(character *model.BaseCharacter) bool {
	if character.CharacterStatus == nil {
		return false
	}
	if FindElementInSlice(model.Drunk, character.CharacterStatus) >= 0 {
		return true
	}
	return false
}

// 判断角色是否中毒
func IsCharacterPoisoned(character *model.BaseCharacter) bool {
	if character.CharacterStatus == nil {
		return false
	}
	if FindElementInSlice(model.Poisoned, character.CharacterStatus) >= 0 {
		return true
	}
	return false
}
