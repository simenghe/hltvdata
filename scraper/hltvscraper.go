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

const agent = "Mozilla/4.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1913.47 Safari/537.36"

// CSGOteam holds the information for the csgo team
type CSGOteam struct {
	TeamName   string   `json:"teamName"`
	Points     int      `json:"points"` // Points need to be int
	Ranking    int      `json:"ranking"`
	Date       string   `json:"date"`
	PlayerList []string `json:"playerList"`
}
type urlStatus struct {
	url    string
	status bool
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
	resp, _, _ := request.Get(url).Set("User-Agent", agent).End()
	fmt.Println(resp.StatusCode)
	return scrapeHltvTeamsByURL(url)
}

// function to test if Url -
func testRequest(url string) bool {
	request := gorequest.New()
	resp, _, _ := request.Get(url).Set("User-Agent", agent).End()
	return resp.StatusCode == 200
}
func testRequestAsync(url string, c chan urlStatus) {
	request := gorequest.New()
	resp, _, _ := request.Get(url).Set("User-Agent", agent).End()
	c <- urlStatus{url: url, status: resp.StatusCode == 200}
}

// scraping function
func scrapeHltvTeamsByURL(url string) []CSGOteam {
	c := colly.NewCollector()
	colly.Async(true)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", agent)

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
				PlayerList: players,
				Date:       time.Now().Format(layoutISO),
			}
			csgoteams = append(csgoteams, team)
		})
	})
	c.Visit(url)
	return csgoteams
}

// ScrapeHltvTeamsByURLAsync scrapes the rankings through channel
func ScrapeHltvTeamsByURLAsync(url string, csChannel chan []CSGOteam) {
	c := colly.NewCollector()
	colly.Async(true)
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", agent)

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
				PlayerList: players,
				Date:       time.Now().Format(layoutISO),
			}
			csgoteams = append(csgoteams, team)
		})
	})
	c.Visit(url)
	csChannel <- csgoteams
}


// URLTraverseAsync traverses
func URLTraverseAsync() []string {
	bench := time.Now()
	c := make(chan urlStatus)
	start := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	var urlCount int = 0
	// THe working URLS in the end
	var workingUrls []string
	for start.Before(now) {
		day := strconv.Itoa(start.Day())
		month := strings.ToLower(start.Month().String())
		year := strconv.Itoa(start.Year())
		url := fmt.Sprintf("https://www.hltv.org/ranking/teams/%s/%s/%s", year, month, day)
		go testRequestAsync(url, c)
		time.Sleep(time.Millisecond * 80)
		// Increment Day
		start = start.AddDate(0, 0, 1)
		urlCount++
	}
	// Exhaust the channel
	for i := 0; i < urlCount; i++ {
		curSite := <-c
		if curSite.status {
			workingUrls = append(workingUrls, curSite.url)
			// fmt.Println(curSite.status)
		}
	}
	fmt.Println(time.Since(bench))
	return workingUrls
}
