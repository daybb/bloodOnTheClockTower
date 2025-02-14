package util

import (
	"bloodOnTheClockTower/model"
	"fmt"
	"testing"
)

func Test_calculateNeighbor(t *testing.T) {
	all := []*model.BaseCharacter{}
	characterMap := make(map[int]*model.BaseCharacter)
	for i := 0; i < 8; i++ {
		if i == 7  {
			all = append(all, &model.BaseCharacter{
				Id:              i,
				CharacterName:   "",
				CharacterKind:   0,
				CharacterStatus: nil,
				IsDead:          true,
				DeadReason:      0,
				DeadTime:        0,
				ExactTime:       0,
				Neighbors:       nil,
			})
			characterMap[i] = &model.BaseCharacter{
				Id:              i,
				CharacterName:   "",
				CharacterKind:   0,
				CharacterStatus: nil,
				IsDead:          true,
				DeadReason:      0,
				DeadTime:        0,
				ExactTime:       0,
				Neighbors:       nil,
			}
		}else {
			all = append(all, &model.BaseCharacter{
				Id:              i,
				CharacterName:   "",
				CharacterKind:   0,
				CharacterStatus: nil,
				IsDead:          false,
				DeadReason:      0,
				DeadTime:        0,
				ExactTime:       0,
				Neighbors:       nil,
			})
			characterMap[i] = &model.BaseCharacter{
				Id:              i,
				CharacterName:   "",
				CharacterKind:   0,
				CharacterStatus: nil,
				IsDead:          false,
				DeadReason:      0,
				DeadTime:        0,
				ExactTime:       0,
				Neighbors:       nil,
			}
		}

	}
	for k,v := range characterMap{
		fmt.Println(k,*v)
	}
	Execute(characterMap[2],8,1,characterMap)
	fmt.Println("done")
	for k,v := range characterMap{
		fmt.Println(k,*v)
	}
	//resp := CalculateNeighbors(all)
	//fmt.Println(resp)
}
