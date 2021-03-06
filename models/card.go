package models

import (
	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

type Card struct {
	Topic      string
	PastDayAdd int
	Feeds      []*BriefFeed
}

func GetHotCardFlows(row_num int) []*Card {
	topics := GetHotTopics()
	cards := make([]*Card, len(topics))
	for i, topic := range topics {
		cards[i] = GetCardByTopic(topic, row_num)
	}
	return cards
}

func GetHotTopics() (topics []string) {
	c := client.dbSess.DB(client.dbName).C(keywordColName)
	glog.V(5).Info(client.dbName, " ", keywordColName)
	pastday := GetPastDayRange()
	glog.V(5).Info(pastday)
	m1 := bson.M{"$match": bson.M{"lastmodified": pastday}}
	m2 := bson.M{"$group": bson.M{"_id": "$keyword", "count": bson.M{"$sum": 1}}}
	m3 := bson.M{"$sort": bson.M{"count": -1}}
	pipe := c.Pipe(
		[]bson.M{
			m1,
			m2,
			m3,
		})
	var results []struct {
		Id    string "_id"
		Count int
	}
	err := pipe.All(&results)
	if err != nil {
		beego.Error(err)
		return
	}
	for _, res := range results {
		topics = append(topics, res.Id)
		glog.V(5).Info("topic: ", res.Id)
	}
	glog.V(3).Info("hot topics: ", len(topics))
	return
}

func GetCardByTopic(topic string, row_num int) *Card {
	kfps, err := GetKeywordFeedPairs(topic, row_num)
	if err != nil {
		beego.Error(err)
		return nil
	}
	c := Card{topic, GetYesterdayAddByKeyword(topic), make([]*BriefFeed, 0)}
	for i := 0; i < len(kfps); i++ {
		feeds, err := GetBriefFeedById(kfps[i].Feedid)
		if err != nil {
			beego.Error(err)
			return nil
		}
		if len(feeds) != 1 {
			beego.Error("feeds illegal!!!")
			return nil
		}
		c.Feeds = append(c.Feeds, &feeds[0])
	}
	return &c
}
