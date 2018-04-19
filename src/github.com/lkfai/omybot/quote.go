package omybot

import (
    "fmt"
    "strings"
    "math/big"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
)

type Quote struct {
    price           big.Float
    webhook         string
    ric             string
    hkexToken       string
}


func (m *Quote) String() string {
        return fmt.Sprintf("[price:%v, ric:%v]", m.price, m.ric)
}

func (m *Quote) sendQuote() (error) {
    resp, err := http.PostForm(m.webhook, url.Values{"content": {m.price.String()}, "tts": {"true"}})
    if err != nil {
        fmt.Printf("Couldn't send message %v\n", err)
        return err
    } else if resp.StatusCode != 200 {
        fmt.Printf("HTTP StatusCode %v, %v\n", resp.StatusCode, resp)
        return err
    }
    return nil
}

func (m *Quote)  getQuote(args []string) (error) {
    param := url.Values{}
    param.Set("hchart", "1")
    param.Add("span", "0")
    param.Add("int", "0")
    param.Add("qid", "1524020346220")
    param.Add("ric", args[0])
    param.Add("token", m.hkexToken)
    param.Add("callback", "a")
    resp, err := http.Get("http://www1.hkex.com.hk/hkexwidget/data/getchartdata2?" + param.Encode())
    if err != nil {
        fmt.Printf("Could not fetch quote\n")
        return err
    }
    if resp.StatusCode != 200 {
        fmt.Printf("HTTP StatusCode not OK\n")
        return err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Unknown response body\n")
        return err
    }
    bodyStr := string(body)
    runes := []rune(bodyStr)
    jsonBlob := string(runes[2:strings.LastIndex(bodyStr, ")")])
    price := *big.NewFloat(0)
    quoteResponse := *new(map[string]interface{})
    if err := json.Unmarshal([]byte(jsonBlob), &quoteResponse); err != nil {
        fmt.Printf("%v\n", err)
        return err
    }
    for key, value := range quoteResponse["data"].(map[string]interface{}) {
        if key == "datalist" {
            prices := value.([]interface{})
            if len(prices) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            prices = prices[1].([]interface{})
            if len(prices) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            price = *big.NewFloat(prices[1].(float64))
            break
        }
    }
    fmt.Printf("Quoted price: %v\n", price.String())
    m.price = price
    return nil
}
