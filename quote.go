package main

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

func (m *Quote) getQuote(args []string) (error) {
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
    if resp.StatusCode != 200 && resp.StatusCode != 204 {
        fmt.Printf("HTTP StatusCode is invalid: %v\n", resp.StatusCode)
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

    type QuoteResponse = map[string]interface{}
    type NodeList = []interface{}
    quoteResponse := &QuoteResponse{}
    if err := json.Unmarshal([]byte(jsonBlob), &quoteResponse); err != nil {
        fmt.Printf("%v\n", err)
        return err
    }
    for key, value := range (*quoteResponse)["data"].(QuoteResponse) {
        if key == "datalist" {
            prices := value.(NodeList)
            if len(prices) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            prices = prices[1].(NodeList)
            if len(prices) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            m.price = *big.NewFloat(prices[1].(float64))
            fmt.Printf("Quoted price: %v\n", m.price.String())
            break
        }
    }
    return nil
}

func (m *Quote) sendQuote() (error) {
    resp, err := http.PostForm(m.webhook, url.Values{"content": {m.price.String()}, "tts": {"true"}})
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
        return fmt.Sprintf("[price:%v, ric:%v]", m.price, m.ric)
}
