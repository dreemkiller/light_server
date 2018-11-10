package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "sync"
        "github.com/gorilla/mux"
        "net/http/httputil"
)

type Program struct {
    Number int `json:"number"`
}

var mutex *sync.Mutex

var currentProgram Program

func logRequest(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
            log.Printf("Failed to dump Request:%v\n", err)
        }
        log.Printf("Request:%v\n", string(dump))
        h.ServeHTTP(w, r)
    })
}

func main() {
    currentProgram.Number = 1
    mutex = &sync.Mutex{}
    router := mux.NewRouter()
    router.HandleFunc("/CurrentProgram", GetCurrentProgram).Methods("GET")
    router.HandleFunc("/CurrentProgram", PutCurrentProgram).Methods("PUT")
    fs := http.FileServer(http.Dir("static"))
    router.PathPrefix("/").Handler(fs)
    log.Fatal(http.ListenAndServe("0.0.0.0:8000", logRequest(router)))
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
