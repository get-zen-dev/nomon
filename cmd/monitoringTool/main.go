package main

import (
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type ServerStatus struct {
	CpuStatus   float64
	MemStatus   float64
	DiskStatus  float64
	LastUpdated time.Time
	duration    time.Duration
	err         error
	wg          sync.WaitGroup
}

var ServerStat *ServerStatus

func NewCloudronServerConnection(duration uint64) *ServerStatus {
	return &ServerStatus{CpuStatus: 0, MemStatus: 0, DiskStatus: 0, LastUpdated: time.Now(), err: nil, duration: time.Duration(duration), wg: sync.WaitGroup{}}
}

func main() {
	ServerStat = NewCloudronServerConnection(5)
	StartMonitoring()
	log.Printf("Cloudron stats:\n CPU Usage: %f\n Disk Usage: %f\n Memory Usage: %f\n", ServerStat.CpuStatus, ServerStat.DiskStatus, ServerStat.MemStatus)
	// log.Println("Starting server on http://127.0.0.1:8080/")

	// r := newRouter()
	// log.Fatal(http.ListenAndServe(":8080", r))
}

// func newRouter() *mux.Router {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", indexHandler).Methods("GET")
// 	staticFileDirectory := http.Dir("./")
// 	staticFileHandler := http.StripPrefix("./", http.FileServer(staticFileDirectory))
// 	r.PathPrefix("./").Handler(staticFileHandler).Methods("GET")
// 	r.HandleFunc("/breakpoints", createBreakpointsHandler).Methods("POST")

// 	return r
// }

func StartMonitoring() {
	ServerStat.wg.Add(3)

	ServerStat.getCpu()
	ServerStat.getMem()
	ServerStat.getDisk()

	ServerStat.wg.Wait()
}

func (serverStatus *ServerStatus) getCpu() {
	defer serverStatus.wg.Done()
	totalPercent, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		log.Println("Error getting CPU: ", err)
	}

	serverStatus.CpuStatus = totalPercent[0]
}

func (serverStatus *ServerStatus) getMem() {
	defer serverStatus.wg.Done()
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting Memory: ", err)
	}
	serverStatus.MemStatus = memInfo.UsedPercent
}

func (serverStatus *ServerStatus) getDisk() {
	defer serverStatus.wg.Done()
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting Disk: ", err)
	}
	serverStatus.DiskStatus = diskInfo.UsedPercent
}

// var tpl = template.Must(template.ParseFiles("./index.html"))

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	tpl.Execute(w, nil)
// }

// type Breakpoints struct {
// 	X string `json:"x"`
// 	Y string `json:"y"`
// }

// func createBreakpointsHandler(w http.ResponseWriter, r *http.Request) {
// 	b := Breakpoints{}

// 	err := r.ParseForm()

// 	if err != nil {
// 		fmt.Println(fmt.Errorf("Error: %v", err))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	b.X = r.Form.Get("x")
// 	b.Y = r.Form.Get("y")

// 	log.Println(b)
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
