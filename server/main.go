package main

import (
	"bloodOnTheClockTower/handler"
	"log"
	"net/http"
)

// run main start a websocket server
// create room create a room from request body, can do it by default config instead
func main() {
	http.HandleFunc("/home", handler.LoadHome)
	http.HandleFunc("/ws", handler.HandleWebSocket)
	// 等待开始页
	http.HandleFunc("/room/", handler.LoadRoom)
	// 游戏中页
	http.HandleFunc("/game/", handler.LoadGame)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
