package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/arjanvaneersel/getlive/internal/entry"
	"log"
	"net/http"
	gourl "net/url"
	"strings"
	"time"
)

// type alog struct {
// 	*log.Logger
// }

// func (log *alog) Critical(args ...interface{})                 { log.Error(args...) }
// func (log *alog) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
// func (log *alog) Notice(args ...interface{})                   { log.Print(args...) }
// func (log *alog) Noticef(format string, args ...interface{})   { log.Printf(format, args...) }

type Twitter struct {
	api    *anaconda.TwitterApi
	logger *log.Logger
	topics []string
	ytKey  string
}

func New(consumerKey, consumerSecret, accessToken, accessSecret, ytKey string, logger *log.Logger, topics ...string) (*Twitter, error) {
	if consumerKey == "" || consumerSecret == "" ||
		accessToken == "" || accessSecret == "" {
		return nil, fmt.Errorf("empty credentials")
	}

	if len(topics) == 0 {
		return nil, fmt.Errorf("no topics")
	}

	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessSecret, consumerKey, consumerSecret)
	// api.SetLogger(alog{logger})

	return &Twitter{
		api:    api,
		logger: logger,
		topics: topics,
		ytKey:  ytKey,
	}, nil
}

func (tw *Twitter) Aggregate(entryChan chan entry.NewEntry) error {
	if len(tw.topics) == 0 {
		return fmt.Errorf("at least one topic need to be provided")
	}

	stream := tw.api.PublicStreamFilter(gourl.Values{
		"track": tw.topics,
	})
	defer stream.Stop()

	tw.logger.Print("Started worker")
	for t := range stream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			// Ignore tweets without URLs
			if len(v.Entities.Urls) == 0 {
				continue
			}

			// resp, err := http.Get(v.Entities.Urls[0].Url)
			// if err != nil {
			// 	tw.logger.Printf("Couldn't get final URL for %q: %v", v.Entities.Urls[0].Url, err)
			// 	continue
			// }
			// defer resp.Body.Close()
			// if resp.StatusCode != 200 {
			// 	return e, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
			// }

			// Try to get a NewEntry, redirected URL from the URL.
			ne, rurl, err := getNewEntryFromMedia(v.Entities.Urls[0].Url, tw.ytKey)
			if err != nil {
				if err == ErrUnsupportedMedia {
					entryChan <- entry.NewEntry{
						Time:  time.Now(),
						Title: v.Text,
						URL:   rurl,
					}
				} else {
					switch err.(type) {
					case MediaNotFoundError:
						tw.logger.Printf("MediaNotFoundError: %v", err)
						continue
					default:
						tw.logger.Printf("getNewEntryFromMedia error: %v", err)
						continue
					}
				}
			}
			ne.Time = time.Now()
			entryChan <- ne

			// tw.logger.Printf("twitter : received tweet: %s... %s", v.Text[:100])
		case anaconda.EventTweet:
			tw.logger.Printf("received event tweet: %+v", v)
			// switch v.Event.Event {
			// case "favorite":
			// 	sn := v.Source.ScreenName
			// 	tw := v.TargetObject.Text
			// 	fmt.Printf("Favorited by %-15s: %s\n", sn, tw)
			// case "unfavorite":
			// 	sn := v.Source.ScreenName
			// 	tw := v.TargetObject.Text
			// 	fmt.Printf("UnFavorited by %-15s: %s\n", sn, tw)
			// }
		}
	}
	return nil
}

type MediaNotFoundError struct {
	URL string
}

func (err MediaNotFoundError) Error() string {
	return fmt.Sprintf("media not available at %q", err.URL)
}

var ErrUnsupportedMedia = errors.New("unsupported media platform")

func getNewEntryFromMedia(url, ytKey string) (e entry.NewEntry, u string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return e, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return e, "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	lcu := strings.ToLower(resp.Request.URL.String())
	switch {
	case strings.Contains(lcu, "youtube.com"):
		ne, err := getNewEntryFromYouTube(resp.Request.URL.String(), ytKey)
		if err != nil {
			return e, "", err
		}
		return ne, resp.Request.URL.String(), nil
	default:
		return e, resp.Request.URL.String(), ErrUnsupportedMedia
	}
}

const ytDataURL = "https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet"

//TODO: Implement all fields for twitter info
//TODO: Move to platform as a separate package
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

type YTInfo struct {
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"snippet"`
	} `json:"items"`
}

func getNewEntryFromYouTube(url, ytKey string) (e entry.NewEntry, err error) {
	u, err := gourl.Parse(url)
	if err != nil {
		return e, err
	}

	v, ok := u.Query()["v"]
	if !ok {
		return e, fmt.Errorf("incorrect youtube URL")
	}

	reqURL := fmt.Sprintf(ytDataURL, v[0], ytKey)
	// fmt.Println(reqURL)
	res, err := http.Get(reqURL)
	if err != nil {
		return e, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return e, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var info YTInfo
	json.NewDecoder(res.Body).Decode(&info)
	if len(info.Items) == 0 {
		return e, MediaNotFoundError{url}
	}
	e.Title = info.Items[0].Snippet.Title
	e.Description = info.Items[0].Snippet.Description
	e.Time = time.Now()
	e.URL = url
	return e, nil
}
