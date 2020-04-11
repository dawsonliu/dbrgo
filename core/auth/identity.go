package core

type Identity struct {
	Id         int               `json:"id"`
	Username   string            `json:"username"`
	Alias      string            `json:"alias"`
	Role       string            `json:"role"`
	Properties map[string]string `json:"properties"`
}
