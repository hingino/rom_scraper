package main

import (
	"fmt"
	"romking/scrape"
)

func main() {
	v := scrape.GetGameList("Game Cube", "https://vimm.net/vault/?p=list&system=GameCube&section=number")
	fmt.Println(v.ConsoleName + "Scrape Done")
}
