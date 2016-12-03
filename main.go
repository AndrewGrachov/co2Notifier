package main

import (
  "bufio"
  "encoding/csv"
  "encoding/json"
  "fmt"
  //"io"
  "bytes"
  "io/ioutil"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
  "unicode/utf8"
)

const Threshold = 1500
const notifyPeriod = 10
const webhook = "https://hooks.slack.com/services/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"

var lastNotification int64

type Notification struct {
  Text string `json:"text"`
}

func Post(url string, json []byte) (resp *http.Response, err error) {
  response, err := http.Post(url, "application/json", bytes.NewBuffer(json))
  if err != nil {
    return nil, err
  }
  return response, nil
}

func notifySlack(value string) {
  fmt.Println("Co2 concentration is above threshold value..")
  text := strings.Join([]string{"Co2 concentration is above treshold value -> ", value}, "")
  notification := Notification{text}
  fmt.Println(notification)
  json, err := json.Marshal(notification)
  if err != nil {
    fmt.Println("marshal err: ", err)
  }
  fmt.Println(string(json))
  response, err := Post(webhook, json)
  if err != nil {
    fmt.Println(err)
  }
  toLog, _ := ioutil.ReadAll(response.Body)
  fmt.Println(string(toLog))
  //{"text": "This is a line of text in a channel.\nAnd this is another line of text."}
}

func main() {
  firstNotification := false

  lastNotification = time.Now().Unix()
  fmt.Println("last notification:", lastNotification)
  fmt.Println("Co2 alerter start.. Threshold: 1500")
  for {
    t := time.Now()
    m := t.Month()
    d := t.Day()

    stringDay := strconv.Itoa(d)
    stringMonth := strconv.Itoa(int(m))

    fmt.Println(utf8.RuneCountInString(stringDay))
    if utf8.RuneCountInString(stringDay) == 1 {
      stringDay = strings.Join([]string{"0", stringDay}, "")
    }

    fileName := strings.Join([]string{stringDay, "CSV"}, ".")

    path := strings.Join([]string{".", stringMonth, stringDay, fileName}, "/")
    fmt.Println("Hello world: ", path)
    f, _ := os.Open(path)

    // Create a new reader.
    r := csv.NewReader(bufio.NewReader(f))
    _, _ = r.Read()
    record, _ := r.Read()
    fmt.Println(record)
    fmt.Println(len(record))
    fmt.Println("value:", record[1])
    co2Value, _ := strconv.ParseInt(record[1], 10, 32)
    if co2Value > Threshold {
      currentTime := time.Now().Unix()
      if currentTime-notifyPeriod > lastNotification || firstNotification == false {
        firstNotification = true
        lastNotification = currentTime
        intCo2Value := int(co2Value)
        notifySlack(strconv.Itoa(intCo2Value))
      }
    }
    f.Close()
    time.Sleep(3 * time.Second)
  }
}
