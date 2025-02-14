package model

//type GoodCharacter struct {
//	Id int
//	BaseCharacter
//}
//
//type MinionCharacter struct {
//	Id int
//	BaseCharacter
//}
//
//type DevilCharacter struct {
//	Id int
//	BaseCharacter
//}
//
//type OutsiderCharacter struct {
//	Id int
//	BaseCharacter
//}

type BaseCharacter struct {
	Id              int
	CharacterName   string          //角色名称。 如魄
	CharacterKind   int             //角色阵营。如村民
	CharacterStatus map[string]int  //todo 角色自身状态，如醉酒，中毒{"drunk":"时间戳"}
	IsDead          bool            //是否死亡
	DeadReason      int             //死亡原因 未死亡为“”
	DeadTime        int             //死亡时间 早上，黄昏，晚上
	ExactTime       int             //死亡具体时间
	Neighbors       []*BaseCharacter //邻居
}

var CharacterKindMap = map[string]int{
	"Chambermaid": Good, "Grandmother": Good, "Courtier": Good, "Gossip": Good, "TeaLady": Good, "Gambler": Good, "Sailor": Good,
	"Fool": Good, "Professor": Good, "Exorcist": Good, "Pacifist": Good, "Minstrel": Good, "Innkeeper": Good,
	"Godfather": Minion, "Mastermind": Minion, "Assassin": Minion, "DevilsAdvocate": Minion,
	"Po": Devil, "Pukka": Devil, "Shabaloth": Devil, "Zombuul": Devil,
	"Goon": Outsider, "Lunatic": Outsider, "Moonchild": Outsider, "Tinker": Outsider,
}

var GoodCharacterMap = map[int]string{
	1: "Chambermaid", 2: "Courtier", 3: "TeaLady", 4: "Sailor", 5: "Professor", 6: "Pacifist", 7: "Minstrel", 8: "Innkeeper",
	9: "Grandmother", 10: "Gossip", 11: "Gambler", 12: "Fool", 13: "Exorcist",
}

var MinionCharacterMap = map[int]string{
	1: "Godfather", 2: "Mastermind", 3: "Assassin", 4: "DevilsAdvocate",
}

var DevilCharacterMap = map[int]string{
	1: "Po", 2: "Pukka", 3: "Shabaloth", 4: "Zombuul",
}

var OutsiderCharacterMap = map[int]string{
	1: "Goon", 2: "Lunatic", 3: "Moonchild", 4: "Tinker",
}

var OncePerGameSkillMap = map[string]bool{
	"Professor": true, "Pacifist": true, "Fool": true, "Assassin": true, "Shabaloth": true,
}

// 阵营
const (
	Good     = 1 //村民
	Minion   = 2 //爪牙
	Outsider = 3 //外来者
	Devil    = 4 //恶魔
)

// 状态
const (
	AllGood   = 1 //无状态
	Drunk     = 2 //醉酒
	Poisoned  = 3 //中毒
	Protected = 4 //被保护
	Executed  = 5 //被处决
)

// 自身状态
const (
	DrunkStatus     = "Drunk"
	PoisonedStatus  = "Poisoned"
	ProtectedStatus = "Protected"
)

// 死亡原因
const (
	DeadReasonExecuted       = 1 //死于处决
	DeadReasonKilledByMinion = 2 //死于爪牙
	DeadReasonKilledByDevil  = 3 //死于恶魔
	KilledByStoryTeller      = 4 //死于说书人
	DeadByShabaloth          = 5
)

// 死亡时间
const (
	Morning = 1
	Dawn    = 2
	Night   = 3
)
