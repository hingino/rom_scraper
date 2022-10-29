package scrape

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type game struct {
	title             string
	game_Info_Link    string
	download_MD5_Hash string
	download_Link     string
}

type GameList struct {
	BaseLink    string
	ConsoleName string
	Games       []game
}

func vaultLinks(start_url string) []string {
	links := []string{}

	c := colly.NewCollector(
		colly.AllowedDomains("vimm.net"),
	)

	c.OnHTML("#vaultMenu > a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		links = append(links, "https://vimm.net"+url)

		// debug print to terminal
		fmt.Println(url)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Getting Vault Links")
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Fatal("Unable to navigate to console base vault. Exiting")
		}
	})

	c.Visit(start_url)

	fmt.Println("Done")

	return links
}

func gameDetails(vaultLink string) []game {
	games := []game{}

	c := colly.NewCollector(
		colly.AllowedDomains("vimm.net"),
	)

	c.OnHTML(".even > td > a[onmouseover]", func(e *colly.HTMLElement) {
		game := game{}
		game.title = e.Text
		game.game_Info_Link = "https://vimm.net" + e.Attr("href")
		games = append(games, game)

		// debug print to terminal
		fmt.Println(game.title, "- indexed.")
	})

	c.OnHTML(".odd > td > a[onmouseover]", func(e *colly.HTMLElement) {
		game := game{}
		game.title = e.Text
		game.game_Info_Link = "https://vimm.net" + e.Attr("href")
		games = append(games, game)

		// debug print to terminal
		fmt.Println(game.title, "- indexed.")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Getting Game Details")
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Fatal("Unable to navigate to vault page. Exiting")
		}
	})

	c.Visit(vaultLink)

	fmt.Println("Done")

	return games
}

func downloadHash(detailsLink string) string {
	hash := ""

	c := colly.NewCollector(
		colly.AllowedDomains("vimm.net"),
	)

	c.OnHTML(".goodHash > td > #data-md5", func(e *colly.HTMLElement) {
		hash = e.Text

		//debug print hash
		fmt.Println(hash)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Getting Game File Hash")
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Fatal("Unable to navigate to game details page. Exiting")
		}
	})

	c.Visit(detailsLink)

	fmt.Println("Done")

	return hash
}

func downloadLink(detailsLink string) string {
	url := "https://download8.vimm.net/download/?mediaId="

	c := colly.NewCollector(
		colly.AllowedDomains("vimm.net"),
	)

	c.OnHTML("#download_form > input[name=mediaId]", func(e *colly.HTMLElement) {
		mediaId := e.Attr("value")
		url += mediaId
		//debug print hash
		fmt.Println(url)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Getting Game File Hash")
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Fatal("Unable to navigate to game details page. Exiting")
		}
	})

	c.Visit(detailsLink)

	fmt.Println("Done")

	return url
}

func setGames(baseLink string) []game {
	// get links to all vaults #-Z
	vault_links := vaultLinks(baseLink)

	// scrape entry page (games starting with numbers)
	games_links := gameDetails(vault_links[0])
	for j := 0; j < len(games_links); j++ {
		games_links[j].download_MD5_Hash = downloadHash(games_links[j].game_Info_Link)
		games_links[j].download_Link = downloadLink(games_links[j].game_Info_Link)
	}

	// crawl and scrape a-z
	for i := 1; i < len(games_links); i++ {
		new_links := gameDetails(vault_links[i])

		for j := 0; j < len(new_links); j++ {
			new_links[j].download_MD5_Hash = downloadHash(new_links[j].game_Info_Link)
			new_links[j].download_Link = downloadLink(new_links[j].game_Info_Link)
			games_links = append(games_links, new_links[j])
		}
	}

	return games_links
}

// public function to get list of games
func GetGameList(console string, base_link string) GameList {
	var gl GameList
	gl.ConsoleName = console
	gl.BaseLink = base_link
	gl.Games = setGames(gl.BaseLink)

	return gl
}
