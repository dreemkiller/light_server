package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "sync"
        "github.com/gorilla/mux"
)

type Program struct {
    Number int `json:"number"`
}

var mutex *sync.Mutex

var currentProgram Program


func main() {
    currentProgram.Number = 1
    mutex = &sync.Mutex{}
    router := mux.NewRouter()
    router.HandleFunc("/CurrentProgram", GetCurrentProgram).Methods("GET")
    router.HandleFunc("/CurrentProgram", PutCurrentProgram).Methods("PUT")
    router.Handle("/", http.FileServer(http.Dir("./static/")))
    log.Fatal(http.ListenAndServe("192.168.0.110:8000", router))
}

func GetCurrentProgram(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("GetCurrentProgram\n")
    fmt.Printf("currentProgram:%v\n", currentProgram)

    var tempProgram Program
    mutex.Lock()
    tempProgram.Number = currentProgram.Number
    mutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(tempProgram); err != nil {
        fmt.Printf("Failed to encode:%v\n", err)
        return
    }
    fmt.Printf("Returning w:%v\n", w)
    return
}

func PutCurrentProgram(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("PutCurrentProgram\n")
    var newProgram Program

    if err := json.NewDecoder(r.Body).Decode(&newProgram); err != nil {
        fmt.Printf("Failed to decode:%v\n", err)
        return
    }
    fmt.Printf("newProgram.Number:%v\n", newProgram.Number)
    if (newProgram.Number != currentProgram.Number) {
        mutex.Lock()
        currentProgram.Number = newProgram.Number
        mutex.Unlock()
    }
}
