package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/parnurzeal/gorequest"
)

// CSGOteam holds the information for the csgo team
type CSGOteam struct {
	ID         int    `json:"ID"`
	TeamName   string `json:"teamName"`
	Points     int    `json:"points"` // Points need to be int
	Ranking    int    `json:"ranking"`
	URL        string `json:"url"`
	Date       string `json:"date"`
	PlayerList string `json:"playerList"`
}

// ScrapeHltvTeams scrapes the top teams, needs automation currently.
func ScrapeHltvTeams() []CSGOteam {
	url := "https://www.hltv.org/ranking/teams"
	return scrapeHltvTeamsByURL(url)
}

// HltvTest function for invalid url
func HltvTest() []CSGOteam {
	url := "https://www.hltv.org/ranking/teams/2018/december/30"
	request := gorequest.New()
	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1913.47 Safari/537.36"
	resp, _, _ := request.Get(url).Set("User-Agent", agent).End()
	fmt.Println(resp.StatusCode)
	return scrapeHltvTeamsByURL(url)
}

// function to test if Url -
func testRequest(url string) bool {
	request := gorequest.New()
	agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1913.47 Safari/537.36"
	resp, _, _ := request.Get(url).Set("User-Agent", agent).End()
	fmt.Println(resp.StatusCode)
	return resp.StatusCode == 200
}

// main function
func scrapeHltvTeamsByURL(url string) []CSGOteam {
	c := colly.NewCollector()
	var header = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1913.47 Safari/537.36"
	colly.Async(true)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", header)

	})
	var csgoteams = []CSGOteam{}
	c.OnHTML("div.ranked-team", func(e *colly.HTMLElement) {
		e.ForEach("div.ranking-header", func(_ int, e *colly.HTMLElement) {
			playerString := e.ChildText("div.playersLine")
			players := strings.Split(playerString, "\n")
			for i := range players {
				players[i] = strings.TrimSpace(players[i])
			}
			const layoutISO = "2006-01-02" // Format the date for golang
			ranking, err := strconv.Atoi(strings.ReplaceAll(e.ChildText("span.position"), "#", ""))
			if err != nil {
				ranking = 0
			}
			// Format the points
			re := regexp.MustCompile("[0-9]+")
			points, err := strconv.Atoi((re.FindAllString(e.ChildText("span.points"), 1)[0]))
			if err != nil {
				points = 0
			}
			var team = CSGOteam{
				TeamName:   e.ChildText("span.name"),
				Ranking:    ranking,
				Points:     points,
				PlayerList: strings.Join(players, " "),
				Date:       time.Now().Format(layoutISO),
			}
			csgoteams = append(csgoteams, team)
		})
	})
	c.Visit(url)
	return csgoteams
}

// RankingTraverse traverses
func RankingTraverse() {
	start := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	var workingUrls []string
	for start.Before(now) {
		day := strconv.Itoa(start.Day())
		month := strings.ToLower(start.Month().String())
		year := strconv.Itoa(start.Year())
		url := fmt.Sprintf("https://www.hltv.org/ranking/teams/%s/%s/%s", year, month, day)
		success := testRequest(url)
		if success {
			workingUrls = append(workingUrls, url)
			fmt.Println(url)
		}
		// Increment Day
		start = start.AddDate(0, 0, 1)
	}
	fmt.Println(workingUrls)
}
