package omybot

import (
	"fmt"
	"strings"
	// "math/big"
	"flag"
	// "net/http"
	// "net/url"
	// "io/ioutil"
	// "encoding/json"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var (
	webhook string 		= ""
	hkexToken string 	= ""
)

func main() {
	tokenPtr := flag.String("token", "", "Discord Bot token")
	webhookPtr := flag.String("webhook", "", "Discord Webhook URL")
	hkexTokenPtr := flag.String("hkexToken", "", "HKEX Token")
	flag.Parse()

	if *tokenPtr == "" || *webhookPtr == "" || *hkexTokenPtr == "" {
		flag.PrintDefaults()
		return
	}

	token := *tokenPtr
	webhook = *webhookPtr
	hkexToken = *hkexTokenPtr
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

	fmt.Printf("OMyBot is now running.  Press CTRL-C to exit.\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "OMyBot")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Messages are boardcasted even to myself, ignore message sent by myself
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Printf("Received Message: %v\n", m.Content)
	if strings.HasPrefix(m.Content, "!quote ") || strings.HasPrefix(m.Content, "!q ") {
		args := strings.Fields(m.Content)[1:]
		quote := new(Quote)
		quote.webhook = webhook
		err := quote.getQuote(args)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		err = quote.sendQuote()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}
}

// func sendQuote(quote string) (err error) {
// 	resp, err := http.PostForm(webhook, url.Values{"content": {quote}, "tts": {"true"}})
// 	if err != nil {
// 		fmt.Printf("Couldn't send message %v\n", err)
// 		return err
// 	} else if resp.StatusCode != 200 {
// 		fmt.Printf("HTTP StatusCode %v, %v\n", resp.StatusCode, resp)
// 		return err
// 	}
// 	return nil
// }

// func getQuote(args []string) (string, error) {
// 	param := url.Values{}
// 	param.Set("hchart", "1")
// 	param.Add("span", "0")
// 	param.Add("int", "0")
// 	param.Add("qid", "1524020346220")
// 	param.Add("ric", args[0])
// 	param.Add("token", hkexToken)
// 	param.Add("callback", "a")
// 	resp, err := http.Get("http://www1.hkex.com.hk/hkexwidget/data/getchartdata2?" + param.Encode())
// 	if err != nil {
// 		fmt.Printf("Could not fetch quote\n")
// 		return "nil", err
// 	}
// 	if resp.StatusCode != 200 {
// 		fmt.Printf("HTTP StatusCode not OK\n")
// 		return "nil", err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("Unknown response body\n")
// 		return "nil", err
// 	}
// 	bodyStr := string(body)
//     runes := []rune(bodyStr)
//     jsonBlob := string(runes[2:strings.LastIndex(bodyStr, ")")])
// 	price := *big.NewFloat(0)
//     quoteResponse := *new(map[string]interface{})
// 	if err := json.Unmarshal([]byte(jsonBlob), &quoteResponse); err != nil {
//         fmt.Printf("%v\n", err)
//         return "nil", err
//     }
//     for key, value := range quoteResponse["data"].(map[string]interface{}) {
// 	    if key == "datalist" {
// 	    	prices := value.([]interface{})
// 	    	if len(prices) < 1 {
// 	    		fmt.Printf("No prices from HKEX\n")
// 	    		return "nil", err
// 	    	}
// 	    	prices = prices[1].([]interface{})
// 	    	if len(prices) < 1 {
// 	    		fmt.Printf("No prices from HKEX\n")
// 	    		return "nil", err
// 	    	}
// 	    	price = *big.NewFloat(prices[1].(float64))
// 		    break
// 	    }
// 	}
// 	fmt.Printf("Quoted price: %v\n", price.String())
// 	return price.String(), nil
// }
