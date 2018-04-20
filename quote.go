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

func (m *Quote) getQuote() (error) {
    pm := url.Values{}
    pm.Set("hchart", "1")
    pm.Add("span", "0")
    pm.Add("int", "0")
    pm.Add("qid", "1524020346220")
    pm.Add("ric", m.ric)
    pm.Add("token", m.hkexToken)
    pm.Add("callback", "a")
    resp, err := http.Get("http://www1.hkex.com.hk/hkexwidget/data/getchartdata2?" + pm.Encode())
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
    bd := string(body)
    rn := []rune(bd)
    bb := string(rn[2:strings.LastIndex(bd, ")")])

    type QuoteResponse = map[string]interface{}
    type NodeList = []interface{}
    qp := &QuoteResponse{}
    if err := json.Unmarshal([]byte(bb), &qp); err != nil {
        fmt.Printf("%v\n", err)
        return err
    }
    for k, v := range (*qp)["data"].(QuoteResponse) {
        if k == "datalist" {
            pcs := v.(NodeList)
            if len(pcs) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            pcs = pcs[1].(NodeList)
            if len(pcs) < 1 {
                fmt.Printf("No prices from HKEX\n")
                return err
            }
            m.price = *big.NewFloat(pcs[1].(float64))
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
