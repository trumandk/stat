package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	//"math"
	"net/http"
)

type Message struct {
	Total       string
	Available   string
	UsedPercent int
	CPU         int
	Containers  int
	Running     int
	Paused      int
	Stopped     int
}

func status(w http.ResponseWriter, req *http.Request) {
	v, err := mem.VirtualMemory()

	if err != nil {
		fmt.Println("err:", err)
	}
	times, err1 := cpu.Percent(0, false)
	if err1 != nil {
		log.Fatal(err1)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	info, err2 := cli.Info(ctx)
	if err2 != nil {
		panic(err2)
	}

	m := Message{
		fmt.Sprintf("%.2f GB", float64(v.Total)/1000000000),
		fmt.Sprintf("%.2f GB", float64(v.Available)/1000000000),
		int(v.UsedPercent),
		int(times[0]),
		info.Containers,
		info.ContainersRunning,
		info.ContainersPaused,
		info.ContainersStopped,
	}
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", status)
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
