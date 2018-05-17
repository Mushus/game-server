package main

type RoomParam struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	MaxUsers       int    `json:"max_users"`
	IsAutoMatching bool   `json:"is_auto_matching"`
}
