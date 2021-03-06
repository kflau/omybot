package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	webhook string
	botName string
)

type Quote struct {
	Price   big.Float
	Webhook string
	Ric     string
}

type QuoteResponse struct {
	Chart Chart `json:"chart"`
}
type Chart struct {
	Result []Result `json:"result"`
	Error  string   `json:"error"`
}
type Result struct {
	Indicator Indicator `json:"indicators"`
}
type Indicator struct {
	Quote []Statistics `json:"quote"`
}
type Statistics struct {
	Close  []float64 `json:"close"`
	Volume []float64 `json:"volume"`
	High   []float64 `json:"high"`
	Open   []float64 `json:"open"`
	Low    []float64 `json:"low"`
}

// --------------------------------- Begin of Quote ---------------------------------

func NewQuote() *Quote {
	return &Quote{}
}

func (m *Quote) Type() string {
	return "Quote"
}

func (m *Quote) Parse(args []string) bool {

	if len(args) <= 0 || len(args) > 1 {
		return false
	}
	return true
}

func (m *Quote) MemberJoin(args *discordgo.MessageCreate) (string, error) {
	return "", nil
}

func (m *Quote) Forward(args []string) error {
	m.Ric = args[0]
	pm := url.Values{}
	pm.Set("region", "HK")
	pm.Add("lang", "zh-Hant-HK")
	pm.Add("range", "1d")
	pm.Add("includePrePost", "false")
	pm.Add("interval", "2m")
	pm.Add("corsDomain", "hk.finance.yahoo.com")
	pm.Add(".tsrc", "finance")
	url := "http://query1.finance.yahoo.com/v8/finance/chart/" + m.Ric + "?" + pm.Encode()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not fetch quote\n")
		return err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		fmt.Printf("HTTP StatusCode is invalid: %v\n", resp.StatusCode)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unknown response body\n")
		return err
	}
	qp := &QuoteResponse{}
	if err := json.Unmarshal(body, &qp); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	if r := qp.Chart.Result; len(r) > 0 {
		ind := r[0].Indicator
		if q := ind.Quote; len(q) > 0 {
			for _, v := range q[0].Close {
				if v > 0 {
					m.Price = *big.NewFloat(v)
					break
				}
			}
		}
	}
	return nil
}

func (m *Quote) Reply(args []string) error {
	m.Webhook = args[0]
	resp, err := http.PostForm(m.Webhook, url.Values{"content": {m.Price.String()}, "tts": {"true"}})
	if err != nil {
		fmt.Printf("Couldn't send message %v\n", err)
		return err
	} else if resp.StatusCode != 200 && resp.StatusCode != 204 {
		fmt.Printf("HTTP StatusCode %v, %v\n", resp.StatusCode, resp)
		return err
	}
	return nil
}

func (m *Quote) String() string {
	return fmt.Sprintf("%#v", m)
}

// --------------------------------- End of Quote ---------------------------------

func main() {
	token := ""
	flag.StringVar(&token, "token", "", "Discord Bot token")
	flag.StringVar(&webhook, "webhook", "", "Discord Webhook URL")
	flag.StringVar(&botName, "botName", "", "Bot Name")
	flag.Parse()
	if token == "" || webhook == "" || botName == "" {
		flag.PrintDefaults()
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Error creating Discord session: %v\n", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Printf("Error opening Discord session: %v\n", err)
	}

	fmt.Printf("QuoteBot is now running.  Press CTRL-C to exit.\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "QuoteBot")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Messages are boardcasted even to myself, ignore message sent by myself
	if botName == m.Author.Username && m.Author.Bot {
		return
	}

	qt := NewQuote()
	if m.Type == discordgo.MessageTypeGuildMemberJoin {
		if str, err := qt.MemberJoin(m); err != nil {
			fmt.Printf("Cannot handle MemberJoin [%#v] due to %#v\n", str, err)
		}
		return
	}
	if len(m.Content) <= 0 {
		return
	}
	args := strings.Fields(m.Content)[1:]
	if !qt.Parse(args) {
		fmt.Printf("Not a quote format\n")
		return
	}
	if err := qt.Forward(args); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	args = append([]string{webhook}, args...)
	if err := qt.Reply(args); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
