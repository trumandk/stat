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

/*
   fmt.Fprintf(w, "<td>%d</td>", info.ContainersRunning)
   fmt.Fprintf(w, "<td>%d</td>", info.ContainersPaused)
   fmt.Fprintf(w, "<td>%d</td>", info.ContainersStopped)
*/
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
		cli, err := client.NewClientWithOpts(client.FromEnv)
	//cli, err := client.NewClientWithOpts(client.WithHost("tcp://192.168.1.132:2375"), client.WithAPIVersionNegotiation())

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
		//math.Round(times[0] * 100 / 100),
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
