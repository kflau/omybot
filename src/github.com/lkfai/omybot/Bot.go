package main

import (
	"flag"
	"math/big"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// var buffer = make([][]byte, 0)

var (
	webhook string 		= ""
)

type QuoteApiResponse struct {
	data QuoteData
	qid string
}

func (b QuoteApiResponse) String() string {
	return fmt.Sprintf("%+v", b)
}

type QuoteData struct {
	responsecode string
	responsemsg string
	datalist [][]int
	start_h int
	start_m int
	end_h int
	end_m int
}

func (b QuoteData) String() string {
	return fmt.Sprintf("%+v", b.datalist)
}

type JokeApiResponse struct {
	Value JokeBody `json:"value"`
	Type  string   `json:"type"`
}

type JokeBody struct {
	Id         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

func main() {
	tokenPtr := flag.String("token", "", "Discord Bot token")
	webhookPtr := flag.String("webhook", "", "Discord Webhook URL")
	flag.Parse()

	token := *tokenPtr
	if token == "" {
		fmt.Println("No token provided. -token <bot token>")
		return
	}
	webhook = *webhookPtr
	if webhook == "" {
		fmt.Println("No webhook provided. -webhook <webhook URL>")
		return
	}
	
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Discord Chuck Norris bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "!joke")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Recieved Message: ")
	fmt.Println(m.Content)
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!quote") {
		err := sendQuote()
		if err != nil {
			fmt.Println(err)
		}
	} else if strings.HasPrefix(m.Content, "!joke") {
		err := sendJoke()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sendQuote() (err error) {
	quote, err := getQuote()
	if err != nil {
		return err
	}
	resp, err := http.PostForm(webhook, url.Values{"content": {quote}, "tts": {"true"}})
	if err != nil {
		fmt.Println("Couldn't send message")
		fmt.Println(err)
		return err
	} else {
		fmt.Println(resp)
		return err
	}
	return nil
}

func getQuote() (string, error) {
	resp, err := http.Get("http://www1.hkex.com.hk/hkexwidget/data/getchartdata2?hchart=1&span=0&int=0&ric=0005.HK&token=evLtsLsBNAUVTPxtGqVeG9B8MG%2bJ3ol24O59BU10wpi6NUJBeKz0Uz2Zu%2bSElSqf&qid=1524020346220&callback=a")
	if err != nil {
		fmt.Println("Could not fetch quote")
		return "nil", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unknown response body")
		return "nil", err
	}

	bodyStr := string(body)
    runes := []rune(bodyStr)
    jsonBlob := string(runes[2:strings.LastIndex(bodyStr, ")")])
	
	var quoteResponse map[string]interface{}
	var price big.Float

	if err := json.Unmarshal([]byte(jsonBlob), &quoteResponse); err != nil {
        fmt.Println(err)
        return "nil", err
    }
    for key, value := range quoteResponse["data"].(map[string]interface{}) {
	    if key == "datalist" {
	    	prices := value.([]interface{})
	    	if len(prices) < 1 {
	    		fmt.Println("No prices from HKEX")
	    		return "nil", err
	    	}
	    	prices = prices[1].([]interface{})
	    	if len(prices) < 1 {
	    		fmt.Println("No prices from HKEX")
	    		return "nil", err
	    	}
	    	price = *big.NewFloat(prices[1].(float64))
		    break
	    }
	}
	fmt.Println("Quoted price: ", price.String())
	return price.String(), nil
}

func sendJoke() (err error) {
	joke, err := getJoke()
	if err != nil {
		return err
	}
	resp, err := http.PostForm(webhook, url.Values{"content": {joke}, "tts": {"true"}})
	fmt.Println(resp)
	if err != nil {
		fmt.Println("Couldn't send message")
		fmt.Println(err)
		return err
	} else {
		fmt.Println(resp)
		return err
	}
	return nil
}

func getJoke() (string, error) {
	resp, err := http.Get("http://api.icndb.com/jokes/random")
	if err != nil {
		fmt.Println("Could not fetch joke")
		return "nil", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unknown response body")
		return "nil", err
	}

	var jokeResp JokeApiResponse
	json.Unmarshal(body, &jokeResp)
	fmt.Println(jokeResp)
	return jokeResp.Value.Joke, nil
}
