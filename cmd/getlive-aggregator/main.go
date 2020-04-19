package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	gourl "net/url"
)

const apiKey = "AIzaSyAGCgz7P0NrM46y4MEBUAPFslWaiX3AUVk"
const ytDataURL = "https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet"

type YTInfo struct {
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"snippet"`
	} `json:"items"`
}

/* {
	"kind": "youtube#videoListResponse",
	"etag": "\"nxOHAKTVB7baOKsQgTtJIyGxcs8/QsNeS22fOTcerU5zRuSvXmMZKVY\"",
	"pageInfo": {
	 "totalResults": 1,
	 "resultsPerPage": 1
	},
	"items": [
	 {
	  "kind": "youtube#video",
	  "etag": "\"nxOHAKTVB7baOKsQgTtJIyGxcs8/dNkLpWFalDZ5m8Vb0awokCebC-g\"",
	  "id": "87-ZFjLfBAQ",
	  "snippet": {
	   "publishedAt": "2020-04-15T21:39:30.000Z",
	   "channelId": "UCg3_C7BwcV0kBlJbBFHTPJQ",
	   "title": "One World: Together At Home Special to Celebrate COVID-19 Workers",
	   "description": "Join Global Citizen, the World Health Organization, Lady Gaga, Taylor Swift, Billie Eilish, Lizzo, and many more artists and healthcare experts as we raise funds for global COVID-19 response efforts. Learn more at GlobalCitizen.org/TogetherAtHome.\n\nARTIST Lineup:\nHours 1 & 2 \nAdam Lambert\nAndra Day\nBlack Coffee\nCharlie Puth\nEason Chan\nHozier & Maren Morris\nHussain Al Jassmi\nJennifer Hudson\nJessie Reyez\nKesha\nLang Lang\nLiam Payne\nLisa Mishra\nLuis Fonsi\nMilky Chance\nNiall Horan\nPicture This\nRita Ora\nSofi Tukker\nThe Killers\nVishal Mishra \n\nHours 3 & 4 \nAdam Lambert\nAnnie Lennox\nBen Platt\nCassper Nyovest\nChristine And The Queens\nCommon\nDelta Goodrem\nEllie Goulding\nFinneas\nJack Johnson\nJacky Cheung\nJess Glynne\nJessie J\nJuanes\nKesha\nMichael Bublé\nRita Ora\nSebastián Yatra\nSheryl Crow\nSho Madjozi\nSofi Tukker\nThe Killers\nZucchero\n \nHours 5 & 6 \nAngèle\nAnnie Lennox\nBen Platt\nBilly Ray Cyrus\nCharlie Puth\nChristine And The Queens\nCommon\nEason Chan\nEllie Goulding\nHozier\nJennifer Hudson\nJessie J\nJohn Legend\nJuanes\nLady Antebellum\nLeslie Odom Jr.\nLuis Fonsi\nNiall Horan\nPicture This\nSebastián Yatra\nSheryl Crow\nSuperM\n\n_____________________________________________________________________\nGlobal Citizen is a social action platform for a global generation that aims to solve the world’s biggest challenges. On our platform, you can learn about issues, take action on what matters most, and join a community committed to social change. We believe we can end extreme poverty because of the collective actions of Global Citizens across the world.\n\nRegister to become a Global Citizen and start taking action today: https://www.globalcitizen.org/\n\nYou can also find us at: \nWebsite: https://www.globalcitizen.org/\nFacebook: https://www.facebook.com/GLBLCTZN\nTwitter: https://twitter.com/glblctzn\nInstagram: https://www.instagram.com/glblctzn/\nTumblr: http://glblctzn.tumblr.com/\nGoogle+: https://plus.google.com/+GLBLCTZN",
	   "thumbnails": {
		"default": {
		 "url": "https://i.ytimg.com/vi/87-ZFjLfBAQ/default_live.jpg",
		 "width": 120,
		 "height": 90
		},
		"medium": {
		 "url": "https://i.ytimg.com/vi/87-ZFjLfBAQ/mqdefault_live.jpg",
		 "width": 320,
		 "height": 180
		},
		"high": {
		 "url": "https://i.ytimg.com/vi/87-ZFjLfBAQ/hqdefault_live.jpg",
		 "width": 480,
		 "height": 360
		},
		"standard": {
		 "url": "https://i.ytimg.com/vi/87-ZFjLfBAQ/sddefault_live.jpg",
		 "width": 640,
		 "height": 480
		},
		"maxres": {
		 "url": "https://i.ytimg.com/vi/87-ZFjLfBAQ/maxresdefault_live.jpg",
		 "width": 1280,
		 "height": 720
		}
	   },
	   "channelTitle": "Global Citizen",
	   "tags": [
		"Global Citizen",
		"Global Citizenship",
		"Music Festival"
	   ],
	   "categoryId": "29",
	   "liveBroadcastContent": "live",
	   "localized": {
		"title": "One World: Together At Home Special to Celebrate COVID-19 Workers",
		"description": "Join Global Citizen, the World Health Organization, Lady Gaga, Taylor Swift, Billie Eilish, Lizzo, and many more artists and healthcare experts as we raise funds for global COVID-19 response efforts. Learn more at GlobalCitizen.org/TogetherAtHome.\n\nARTIST Lineup:\nHours 1 & 2 \nAdam Lambert\nAndra Day\nBlack Coffee\nCharlie Puth\nEason Chan\nHozier & Maren Morris\nHussain Al Jassmi\nJennifer Hudson\nJessie Reyez\nKesha\nLang Lang\nLiam Payne\nLisa Mishra\nLuis Fonsi\nMilky Chance\nNiall Horan\nPicture This\nRita Ora\nSofi Tukker\nThe Killers\nVishal Mishra \n\nHours 3 & 4 \nAdam Lambert\nAnnie Lennox\nBen Platt\nCassper Nyovest\nChristine And The Queens\nCommon\nDelta Goodrem\nEllie Goulding\nFinneas\nJack Johnson\nJacky Cheung\nJess Glynne\nJessie J\nJuanes\nKesha\nMichael Bublé\nRita Ora\nSebastián Yatra\nSheryl Crow\nSho Madjozi\nSofi Tukker\nThe Killers\nZucchero\n \nHours 5 & 6 \nAngèle\nAnnie Lennox\nBen Platt\nBilly Ray Cyrus\nCharlie Puth\nChristine And The Queens\nCommon\nEason Chan\nEllie Goulding\nHozier\nJennifer Hudson\nJessie J\nJohn Legend\nJuanes\nLady Antebellum\nLeslie Odom Jr.\nLuis Fonsi\nNiall Horan\nPicture This\nSebastián Yatra\nSheryl Crow\nSuperM\n\n_____________________________________________________________________\nGlobal Citizen is a social action platform for a global generation that aims to solve the world’s biggest challenges. On our platform, you can learn about issues, take action on what matters most, and join a community committed to social change. We believe we can end extreme poverty because of the collective actions of Global Citizens across the world.\n\nRegister to become a Global Citizen and start taking action today: https://www.globalcitizen.org/\n\nYou can also find us at: \nWebsite: https://www.globalcitizen.org/\nFacebook: https://www.facebook.com/GLBLCTZN\nTwitter: https://twitter.com/glblctzn\nInstagram: https://www.instagram.com/glblctzn/\nTumblr: http://glblctzn.tumblr.com/\nGoogle+: https://plus.google.com/+GLBLCTZN"
	   },
	   "defaultAudioLanguage": "en"
	  }
	 }
	]
   }
*/

func ExampleScrape(url string) {
	// Request the HTML page.
	u, err := gourl.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	v, ok := u.Query()["v"]
	if !ok {
		log.Fatalf("incorrect youtube URL")
	}

	reqURL := fmt.Sprintf(ytDataURL, v[0], apiKey)
	fmt.Println(reqURL)
	res, err := http.Get(reqURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var info YTInfo
	json.NewDecoder(res.Body).Decode(&info)
	fmt.Printf("%v", info)

	// Load the HTML document
	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// doc.Find("h1").Each(func(i int, sel *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	title := strings.Trim(sel.First().Text(), " ")
	// 	// title := sel.Find(".ytd-video-primary-info-renderer").Text()
	// 	// band := s.Find("a").Text()
	// 	// title := s.Find("i").Text()
	// 	fmt.Printf("Title %d: %s\n", i, title)
	// })

	// doc.Find("div #description").Each(func(i int, sel *goquery.Selection) {
	// 	fmt.Println("Found description")

	// 	// For each item found, get the band and title
	// 	desc := sel.Find("span .style-scope .yt-formatted-string").Text()
	// 	// band := s.Find("a").Text()
	// 	// title := s.Find("i").Text()
	// 	fmt.Printf("Desc %d: %s\n", i, desc)
	// })
}

func main() {
	ExampleScrape("https://www.youtube.com/watch?v=87-ZFjLfBAQ&feature=youtu.be")
}
