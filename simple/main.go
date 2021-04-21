package main

import (
  "flag"
  "log"
  "net/http"
  "github.com/pusher/pusher-http-go"
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
  log.Println(r.URL)
  if r.URL.Path != "/" {
    http.Error(w, "Not found", http.StatusNotFound)
    return
  }
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }
  http.ServeFile(w, r, "simple/index.html")
}

func main() {
  flag.Parse()
  http.HandleFunc("/", serveHome)
  http.HandleFunc("/ajax", func(w http.ResponseWriter, r *http.Request) {
    data := "Hello from server"
    pusherClient.Trigger("simple-channel", "simple-event", data)
  })
  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
