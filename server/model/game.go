package model

type ActionReq struct {
	Action  string   `json:"action"`
	Targets []string `json:"targets"`
	Extra   string   `json:"extra"`
}
