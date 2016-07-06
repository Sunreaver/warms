package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/sunreaver/gotools/mail"
	sys "github.com/sunreaver/gotools/system"
	"github.com/sunreaver/mahonia"
)

var (
	// StockDataRegexp sina接口返回数据提取
	StockDataRegexp = regexp.MustCompile(`="(.*)";`)

	// mongodb
	mongo *mgo.Session
)

// sina stock接口返回的数据位置
const (
	Name             = 0
	TodayOpening     = 1
	YesterdayClosing = 2
	NowValue         = 3
	Date             = 30
	Time             = 31
)

// Config 用户配置
type Config struct {
	Mail   []string         `json:"emails"`
	Stocks sort.StringSlice `json:"stocks"`
}

// Stock sina stock返回数据后整理
type Stock struct {
	Name             string  `bson:"name"`
	Code             string  `bson:"c,omitempty"`
	Margin           float64 `bson:"margin"`
	Now              string  `bson:"now"`
	Opening          string  `bson:"opening"`
	YesterdayClosing string  `bson:"yopening"`
	Downup           int     `bson:"downup,omitempty"`
	Time             string  `bson:"timeFormat"`
	TimeUnix         int64   `bson:"time"`
}

// SaveDB save
func (s *Stock) SaveDB() error {
	sess := mongo.Copy()
	defer sess.Close()

	_, e := sess.DB("").C(s.Code).Upsert(bson.M{"timeFormat": s.Time}, s)
	return e
}

func (s Stock) String() string {
	format := "%s:%2.2f%%\t>>>>>当前:%s\t昨收:%s\t幅度:%0.3f\t今开:%s<<<<<"
	yestoday, _ := strconv.ParseFloat(s.YesterdayClosing, 64)
	return fmt.Sprintf(format, s.Name, s.Margin/yestoday*100, s.Now, s.YesterdayClosing, s.Margin, s.Opening)
}

func init() {
	var err error
	mongo, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"127.0.0.1", "localhost"},
		Direct:   true,
		Timeout:  0,
		Database: "Stocks",
		Username: "stocks",
		Password: "1111",
	})
	if err != nil {
		panic("mongo Dial Error")
	}
	mongo.SetMode(mgo.Monotonic, true)
}

func main() {
	configs, e := readFile(sys.CurPath() + sys.SystemSep() + "stock.json")
	if e != nil {
		log.Println(e.Error())
		return
	}

	for _, cfg := range configs {

		if len(cfg.Mail) == 0 {
			continue
		}
		sort.Sort(cfg.Stocks)
		list := strings.Join(cfg.Stocks, ",")
		resp, err := http.Get("http://hq.sinajs.cn/list=" + list)
		if err != nil {
			log.Println("Get Err:", err)
			continue
		}
		body, e := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if e != nil {
			log.Println("ioutil Err:", err)
			continue
		}

		str := string(body)
		enc := mahonia.NewDecoder("gbk")
		matchs := StockDataRegexp.FindAllStringSubmatch(enc.ConvertString(str), -1)
		if len(matchs) == 0 {
			continue
		}
		stocks := []Stock{}
		for index, v := range matchs {
			if len(v) < 2 {
				continue
			}
			numerical := strings.Split(v[1], ",")
			if len(numerical) < 32 {
				log.Println("Format32 Err:", err)
				continue
			}

			upDown := "未知"
			yestoday, e1 := strconv.ParseFloat(numerical[YesterdayClosing], 64)
			today, e2 := strconv.ParseFloat(numerical[NowValue], 64)
			if e1 == nil && e2 == nil {
				upDown = strconv.FormatFloat(today-yestoday, 'f', -1, 64)
				if len(upDown) > 6 {
					upDown = upDown[:6]
				}
			}

			t := numerical[Date] + " " + numerical[Time]
			ti, e3 := time.Parse("2006-01-02 15:04:05", t)
			if e3 != nil {
				ti = time.Now()
			}
			s := Stock{
				Code:             cfg.Stocks[index],
				Margin:           float64(today - yestoday),
				Name:             numerical[Name],
				Now:              numerical[NowValue],
				Opening:          numerical[TodayOpening],
				YesterdayClosing: numerical[YesterdayClosing],
				Time:             t,
				TimeUnix:         ti.Unix(),
			}
			s.SaveDB()
			stocks = append(stocks, s)
		}

		outStr := ""

		for _, s := range stocks {
			outStr += fmt.Sprintln(s)
		}
		outStr += "\r\nHappy day!\r\n"
		fmt.Println(outStr)

		e = mail.SendMail(sys.CurPath()+sys.SystemSep()+"auth.json", outStr, "每日一报", "行情播报："+time.Now().Format("2006-01-02 15:04\r\n"), cfg.Mail)
		if e != nil {
			log.Println(e.Error())
		}
	}
	log.Println(time.Now().Format("2006-01-02 15:04:05"))
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
