package util

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"log"
	"sort"
)

func FindElementInSlice(element string, target map[string]int) bool {
	if target[element] != 0 {
		return true
	}
	return false
}

// 每次涉及死人或者复活都需要调用此接口
func CalculateNeighbors(characterMap map[int]*model.BaseCharacter) {
	var allCharacterId []int
	//result := make(map[int][]int)
	for _, v := range characterMap {
		fmt.Println(v.IsDead)
		if !v.IsDead {
			allCharacterId = append(allCharacterId, v.Id)
		}
	}
	sort.Ints(allCharacterId)
	for i, v := range allCharacterId {
		// 计算左邻居和右邻居的索引（环形结构）
		leftIndex := (i - 1 + len(allCharacterId)) % len(allCharacterId)
		rightIndex := (i + 1) % len(allCharacterId)
		characterMap[v].Neighbors = []*model.BaseCharacter{characterMap[leftIndex], characterMap[rightIndex]}
		//result[v] = []int{allCharacterId[leftIndex], allCharacterId[rightIndex]}
	}
	return
}

// 处决
func Execute(character *model.BaseCharacter, curTime, timeNow int, characterMap map[int]*model.BaseCharacter) {
	character.IsDead = true
	character.DeadReason = model.DeadReasonExecuted
	character.DeadTime = timeNow
	character.ExactTime = curTime
	CalculateNeighbors(characterMap)
}

// 首夜给爪牙展示恶魔，给恶魔展示爪牙
func ShowDevilAndMinion(players []model.Player) (minionId []string, devilId []string) {
	for _, v := range players {
		if v.BaseCharacter.CharacterKind == model.Devil {
			devilId = append(devilId, v.Id)
		} else if v.BaseCharacter.CharacterKind == model.Minion {
			minionId = append(minionId, v.Id)
		}
	}
	for _, v := range players {
		if v.BaseCharacter.CharacterKind == model.Devil {
			v.Log += fmt.Sprintf("爪牙是座位号为%v，请互相确认身份", minionId)
		} else if v.BaseCharacter.CharacterKind == model.Minion {
			v.Log += fmt.Sprintf("恶魔是座位号为%v，请互相确认身份", devilId)
		}
	}
	log.Printf("首夜确认恶魔和爪牙身份。恶魔是%v,爪牙是%v", devilId, minionId)

	return minionId, devilId
}

// 首个夜晚。恶魔和爪牙互相展示身份
// 水手发动技能（暂无）
// 侍臣发动技能（暂无）
// 对教父展示所有外来者标记（暂无）
// 魔鬼代言人选择存活玩家（暂无）
// 普卡选择玩家（暂无）
// 给祖母展示玩家（暂无）
// 侍女技能发动（暂无）
// 睁眼
func FirstNight(players []model.Player) ([]string, []string) {
	// 恶魔和爪牙互相展示身份
	minionId, devilId := ShowDevilAndMinion(players)
	// todo 水手技能
	// todo 侍臣技能
	// todo 教父技能
	// todo 魔鬼代言人技能
	// todo 普卡技能
	// todo 祖母技能
	// todo 侍女技能
	return minionId, devilId
}
