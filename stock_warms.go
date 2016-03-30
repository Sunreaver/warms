package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/sunreaver/goTools/mail"
	sys "github.com/sunreaver/goTools/system"
	"github.com/sunreaver/mahonia"
)

var (
	// StockDataRegexp sina接口返回数据提取
	StockDataRegexp = regexp.MustCompile(`="(.*)";`)

	// mongodb
	mongo *mgo.Session
)

const (
	Name             = 0
	TodayOpening     = 1
	YesterdayClosing = 2
	NowValue         = 3
	Date             = 30
	Time             = 31
)

type Config struct {
	Mail   []string `json:"emails"`
	Stocks []string `json:"stocks"`
}

type Stock struct {
	Name             string  `bson:"name"`
	Code             string  `bson:",omitempty"`
	Margin           float64 `bson:"margin"`
	Now              string  `bson:"now"`
	Opening          string  `bson:"opening"`
	YesterdayClosing string  `bson:"yopening"`
	Time             string  `bson:"timeFormat"`
	TimeUnix         int64   `bson:"time"`
}

func (s *Stock) SaveDB() error {
	sess := mongo.Copy()
	defer sess.Close()

	_, e := sess.DB("").C(s.Code).Upsert(bson.M{"timeFormat": s.Time}, s)
	return e
}

func init() {
	var err error
	mongo, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"127.0.0.1", "localhost"},
		Direct:   true,
		Timeout:  0,
		Database: "Stocks",
	})
	if err != nil {
		panic("mongo Dial Error")
	}
	mongo.SetMode(mgo.Monotonic, true)
}

func main() {
	configs, e := readFile(sys.CurPath() + sys.SystemSep() + "stock.json")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	format := "%s:\t>>>>>昨收:%s\t今收:%s\t今开:%s\t幅度:%s<<<<<\r\n"

	for _, cfg := range configs {

		if len(cfg.Mail) == 0 {
			continue
		}
		outStr := ""
		for _, v := range cfg.Stocks {
			resp, err := http.Get("http://hq.sinajs.cn/list=" + v)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			body, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				continue
			}

			str := string(body)
			enc := mahonia.NewDecoder("gbk")
			matchs := StockDataRegexp.FindAllStringSubmatch(enc.ConvertString(str), -1)
			if len(matchs) == 0 || len(matchs[0]) < 2 {
				continue
			}

			numerical := strings.Split(matchs[0][1], ",")
			if len(numerical) < 32 {
				continue
			}

			upDown := "未知"
			yestoday, e1 := strconv.ParseFloat(numerical[YesterdayClosing], 64)
			today, e2 := strconv.ParseFloat(numerical[NowValue], 64)
			if e1 == nil && e2 == nil {
				upDown = strconv.FormatFloat(today-yestoday, 'f', -1, 64)[0:6]
			}
			outStr = outStr + fmt.Sprintf(format, numerical[Name], numerical[YesterdayClosing], numerical[NowValue], numerical[TodayOpening], upDown)

			t := numerical[Date] + " " + numerical[Time]
			ti, e3 := time.Parse("2006-01-02 15:04:05", t)
			if e3 != nil {
				ti = time.Now()
			}
			s := Stock{
				Code:             v[2:],
				Margin:           float64(today - yestoday),
				Name:             numerical[Name],
				Now:              numerical[NowValue],
				Opening:          numerical[TodayOpening],
				YesterdayClosing: numerical[YesterdayClosing],
				Time:             t,
				TimeUnix:         ti.Unix(),
			}
			s.SaveDB()
		}

		outStr = outStr + "\r\nHappy day!\r\n"

		e = mail.SendMail(sys.CurPath()+sys.SystemSep()+"auth.json", outStr, cfg.Mail)
		if e != nil {
			fmt.Println(e.Error())
		}
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func readFile(fileName string) ([]Config, error) {
	b, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil, e
	}

	var cfg []Config
	e = json.Unmarshal(b, &cfg)
	if e != nil {
		return nil, e
	}
	return cfg, e
}
