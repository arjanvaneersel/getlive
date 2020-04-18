package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/arjanvaneersel/getlive/internal/entry"
	"log"
	"net/url"
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
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string, logger *log.Logger, topics ...string) (*Twitter, error) {
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
	}, nil
}

func (tw *Twitter) Aggregate(entryChan chan *entry.Entry) error {
	if len(tw.topics) == 0 {
		return fmt.Errorf("at least one topic need to be provided")
	}

	stream := tw.api.PublicStreamFilter(url.Values{
		"track": tw.topics,
	})
	defer stream.Stop()

	tw.logger.Print("twitter : Started worker")
	for t := range stream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			tw.logger.Printf("twitter : received tweet: %+v", v)
		case anaconda.EventTweet:
			tw.logger.Printf("twitter : received event tweet: %+v", t.(anaconda.EventTweet))
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

		entryChan <- &entry.Entry{}

		// if t.RetweetedStatus != nil {
		// 	continue
		// }

		// _, err := api.Retweet(t.Id, false)
		// if err != nil {
		// 	log.Errorf("could not retweet %d: %v", t.Id, err)
		// 	continue
		// }
		// log.Infof("retweeted %d", t.Id)
	}
	return nil
}
