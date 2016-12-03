package config

import "github.com/jinzhu/configor"
import "fmt"

var Data struct {
  Treshold     int64
  NotifyPeriod int64
  Notifiers    struct {
    Slack struct {
      Webhook string
    }
  }
}

func init() {
  fmt.Println("Loading config........")
  configor.Load(&Data, "src/config/config.json")
}
