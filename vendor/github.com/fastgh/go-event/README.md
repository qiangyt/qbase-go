# Type-safe in-process pub/sub event library for Go language

# Features
  - Type-safe: every kind of event is a concrete go type instead of error-prone interface{} or map[string]interface{}. Unmatched event is rejected.
  - Pub-sub (topic): publishes event to multiple listeners, thread-safely
  - Asynchronous: listeners wait on their separated go channels and go routines. Listener execution panic won't affect others, error is logged only instead
  - Fine-tuned logging: agnostic logging abstraction, easy to adopt various of logging implementations. Go std log and github.com/phuslu/log are built-in ones. Or turn on logging completely.

# Install
  ```go
  go get github.com/fastgh/go-event
  ```

# Example
  ```go
  package main

  import (
    "fmt"

    "github.com/fastgh/go-event"
    "github.com/fastgh/go-event/loggers/std"
  )

  type MyEvent struct {
    Name string
  }

  func main() {
    myHub := event.NewHub("default", std.NewDefaultGlobalStdLogger())

    myTopic := event.CreateTopic(myHub, "myTopic", MyEvent{})

    myTopic.Sub("listener1", func(e MyEvent) {
      fmt.Println("listener1 - got event from", e)
    }, 0)

    myTopic.Sub("listener2", func(e MyEvent) {
      fmt.Println("listener2 - got event from", e)
    }, 0)

    myTopic.Pub(event.PubModeAuto, MyEvent{"fastgh"})

    myHub.Close(true)
  }
  ```

  To run the example:
  ```shell
  go run ./examples/hello/main.go
  ```

  It will output:
  ```log
  2022/09/18 18:57:29 <INFO> hub=default, topic=myTopic, listener= --> topic register begin
  2022/09/18 18:57:29 <INFO> hub=default, topic=myTopic, listener= --> topic register ok
  2022/09/18 18:57:29 <INFO> hub=default, topic=myTopic, listener=listener1 --> listener sub ok
  2022/09/18 18:57:29 <INFO> hub=default, topic=myTopic, listener=listener2 --> listener sub ok
  listener1 - got event from {fastgh}
  listener2 - got event from {fastgh}
  2022/09/18 18:57:29 <INFO> hub=default, topic=, listener= --> hub close begin
  2022/09/18 18:57:29 <INFO> event={"Id":2,"Hub":"default","Topic":"myTopic","Close":true,"Data":null}, listener=listener1 --> listener close begin
  2022/09/18 18:57:29 <INFO> event={"Id":2,"Hub":"default","Topic":"myTopic","Close":true,"Data":null}, listener=listener2 --> listener close begin
  2022/09/18 18:57:29 <INFO> hub=default, topic=, listener= --> hub close ok
  ```

# Logging

  ```go
  type Logger interface {
    LogDebug(enm LogEnum, hub string, topic string, lsner string)
    LogInfo(enm LogEnum, hub string, topic string, lsner string)
    LogError(enm LogEnum, hub string, topic string, lsner string, err any)

    LogEventDebug(enm LogEnum, lsner string, evnt Event)
    LogEventInfo(enm LogEnum, lsner string, evnt Event)
    LogEventError(enm LogEnum, lsner string, evnt Event, err any)
  }
  ```


# License
  MIT License
