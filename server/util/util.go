package util

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"sort"
)

func FindElementInSlice(element string, target map[string]int) bool {
	if target[element] != 0 {
		return true
	}
	return false
}

// 每次涉及死人或者复活都需要调用此接口
func CalculateNeighbors(characterMap map[int]*model.BaseCharacter)  {
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
		characterMap[v].Neighbors = []*model.BaseCharacter{characterMap[leftIndex],characterMap[rightIndex]}
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
