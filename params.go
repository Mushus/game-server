package main

type RoomParam struct {
	Name           string `json:"name"`
	Password       string `json:"password"`
	MaxUsers       int    `json:"maxUsers"`
	IsAutoMatching bool   `json:"isAutoMatching"`
}
