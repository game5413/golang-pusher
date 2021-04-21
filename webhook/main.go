package main

import (
  "flag"
  "fmt"
  "github.com/pusher/pusher-http-go"
  "io/ioutil"
  "log"
  "net/http"
)

var addr = flag.String("addr", ":8000", "http service address")

var pusherClient = pusher.Client{
  AppID: "",
  Key: "",
  Secret: "",
  Cluster: "ap1",
  Secure: true,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "Not found", http.StatusNotFound)
    return
  }
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  http.ServeFile(w, r, "webhook/index.html")
}

func pusherAuth(res http.ResponseWriter, req *http.Request) {
  if req.URL.Path != "/auth" {
    http.Error(res, "Not found", http.StatusNotFound)
    return
  }
  if req.Method != "POST" {
    http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  params, _ := ioutil.ReadAll(req.Body)
  response, err := pusherClient.AuthenticatePrivateChannel(params)

  if err != nil {
    http.Error(res, "Forbidden", http.StatusForbidden)
    return
  }

  res.Header().Set("Content-Type", "application/json")

  fmt.Fprintf(res, string(response))
  return
}

func webhookEndpoint(res http.ResponseWriter, req *http.Request) {
  if req.URL.Path != "/webhook" {
    http.Error(res, "Not found", http.StatusNotFound)
    return
  }
  if req.Method != "POST" {
    http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  body, _ := ioutil.ReadAll(req.Body)

  hook, err := pusherClient.Webhook(req.Header, body)
  if err != nil {
    http.Error(res, "Internal server error", http.StatusInternalServerError)
    return
  }

  for _, event := range hook.Events {
    fmt.Println(event.Event)
    switch event.Event {
    case "client-testEvent":
      pusherClient.Trigger("private-client-channel", "private-client-event", "trigger from client")
      break
    default:
      break
    }
  }

  res.WriteHeader(http.StatusOK)
  return
}

func main() {
  flag.Parse()
  http.HandleFunc("/", serveHome)
  http.HandleFunc("/auth", pusherAuth)
  http.HandleFunc("/webhook", webhookEndpoint)
  http.HandleFunc("/ajax", func(w http.ResponseWriter, r *http.Request) {
    pusherClient.Trigger("private-client-channel", "private-client-event", "hello from private channel")
  })
  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
