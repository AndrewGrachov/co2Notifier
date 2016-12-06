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
  "packages/config"
  "project"
  "strconv"
  "strings"
  "time"
  "unicode/utf8"
)

// const Threshold = 1500
// const notifyPeriod = 10

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

func notifySlack(value string, notificationText string) {
  fmt.Println("Co2 concentration is above threshold value..")
  text := strings.Join([]string{notificationText, value}, "")
  notification := Notification{text}
  fmt.Println(notification)
  json, err := json.Marshal(notification)
  if err != nil {
    fmt.Println("marshal err: ", err)
  }
  fmt.Println(string(json))
  response, err := Post(config.Data.Notifiers.Slack.Webhook, json)
  if err != nil {
    fmt.Println(err)
  }
  toLog, _ := ioutil.ReadAll(response.Body)
  fmt.Println(string(toLog))
}

func main() {
  firstNotification := false
  notifying := false
  fmt.Println("PATH_SEPARATOR:", project.PATH_SEPARATOR)

  lastNotification = time.Now().Unix()
  fmt.Println("last notification:", lastNotification)
  fmt.Println("Co2 alerter start.. ", config.Data)
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

    path := strings.Join([]string{".", stringMonth, stringDay, fileName}, project.PATH_SEPARATOR)
    f, err := os.Open(path)
    if err != nil {
      panic(err)
    }

    // Create a new reader.
    r := csv.NewReader(bufio.NewReader(f))
    _, _ = r.Read()
    record, _ := r.Read()
    fmt.Println(record)
    fmt.Println(len(record))
    fmt.Println("value:", record[1])
    co2Value, _ := strconv.ParseInt(record[1], 10, 32)
    intCo2Value := int(co2Value)
    if co2Value > config.Data.Treshold {
      notifying = true
      currentTime := time.Now().Unix()
      if currentTime-config.Data.NotifyPeriod > lastNotification || firstNotification == false {
        firstNotification = true
        lastNotification = currentTime
        notifySlack(strconv.Itoa(intCo2Value), "Co2 concentration is above treshold value -> ")
      }
    } else if co2Value < config.Data.LowerTreshold && notifying {
      notifying = false
      notifySlack(strconv.Itoa(intCo2Value), "Co2 concentration is back to normal -> ")
    }
    f.Close()
    time.Sleep(3 * time.Second)
  }
}
