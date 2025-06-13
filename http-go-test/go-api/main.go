// package main

// import (
// 	"net/http"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// )

// var (
// 	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Total number of HTTP requests",
// 	})

// 	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "cpu_usage_percent",
// 		Help: "Dummy CPU usage in percent",
// 	})

// 	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "memory_usage_mb",
// 		Help: "Dummy Memory usage in MB",
// 	})
// )

// func recordMetrics() {
// 	go func() {
// 		for {
// 			cpuPercents, _ := cpu.Percent(0, false)
// 			if len(cpuPercents) > 0 {
// 				cpuUsage.Set(cpuPercents[0])
// 			}

// 			vMem, _ := mem.VirtualMemory()
// 			usedMB := float64(vMem.Used) / (1024 * 1024)
// 			memoryUsage.Set(usedMB)

// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// }

// func main() {
// 	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage)

// 	recordMetrics()

// 	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	})

// 	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Hello, World!"))
// 	})

// 	http.Handle("/metrics", promhttp.Handler())

// 	http.ListenAndServe(":8080", nil)
// }



// package main

// import (
// 	"log"
// 	"net/http"
// 	"runtime"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// )

// var (
// 	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Total number of HTTP requests",
// 	})

// 	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "cpu_usage_percent",
// 		Help: "CPU usage percentage",
// 	})

// 	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "memory_usage_mb",
// 		Help: "Memory usage in MB",
// 	})
// )

// func recordMetrics() {
// 	go func() {
// 		for {
// 			if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
// 				cpuUsage.Set(cpuPercents[0])
// 			}

// 			if vMem, err := mem.VirtualMemory(); err == nil {
// 				memoryUsage.Set(float64(vMem.Used) / (1024 * 1024))
// 			}

// 			time.Sleep(1 * time.Second)
// 		}
// 	}()
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU() * 2) // Use more CPU cores

// 	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage)
// 	recordMetrics()

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Hello, World!"))
// 	})

// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	})

// 	mux.Handle("/metrics", promhttp.Handler())

// 	server := &http.Server{
// 		Addr:              ":8080",
// 		Handler:           mux,
// 		ReadTimeout:       5 * time.Second,
// 		WriteTimeout:      10 * time.Second,
// 		IdleTimeout:       120 * time.Second,
// 		ReadHeaderTimeout: 3 * time.Second,
// 		MaxHeaderBytes:    1 << 20, // 1 MB
// 	}

// 	log.Println("Starting server on :8080")
// 	log.Fatal(server.ListenAndServe())
// }


// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"runtime"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// )

// var (
// 	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Total number of HTTP requests",
// 	})

// 	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "cpu_usage_percent",
// 		Help: "CPU usage percentage",
// 	})

// 	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "memory_usage_mb",
// 		Help: "Memory usage in MB",
// 	})

// 	okResponse       = []byte("OK")
// 	helloWorldBuffer = []byte("Hello, World!")
// )

// func recordMetrics() {
// 	ticker := time.NewTicker(5 * time.Second) // Reduce frequency to lessen overhead
// 	defer ticker.Stop()

// 	go func() {
// 		for range ticker.C {
// 			if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
// 				cpuUsage.Set(cpuPercents[0])
// 			}

// 			if vMem, err := mem.VirtualMemory(); err == nil {
// 				memoryUsage.Set(float64(vMem.Used) / (1024 * 1024))
// 			}
// 		}
// 	}()
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

// 	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage)
// 	recordMetrics()

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(helloWorldBuffer)
// 	})

// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		httpRequests.Inc()
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(okResponse)
// 	})

// 	mux.Handle("/metrics", promhttp.Handler())

// 	server := &http.Server{
// 		Addr:              ":8080",
// 		Handler:           mux,
// 		ReadTimeout:       5 * time.Second,
// 		WriteTimeout:      10 * time.Second,
// 		IdleTimeout:       120 * time.Second,
// 		ReadHeaderTimeout: 3 * time.Second,
// 		MaxHeaderBytes:    1 << 20,
// 	}

// 	// Graceful shutdown
// 	go func() {
// 		log.Println("Starting server on :8080")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Server failed: %v", err)
// 		}
// 	}()

// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt)
// 	<-stop

// 	log.Println("Shutting down server...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err := server.Shutdown(ctx); err != nil {
// 		log.Fatalf("Graceful shutdown failed: %v", err)
// 	}
// 	log.Println("Server stopped.")
// }


package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	db *sql.DB

	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	})

	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "CPU usage percentage",
	})

	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_mb",
		Help: "Memory usage in MB",
	})

	dbOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Current open DB connections",
	})

	okResponse       = []byte("OK")
	helloWorldBuffer = []byte("Hello, World!")
)

func recordMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
				cpuUsage.Set(cpuPercents[0])
			}

			if vMem, err := mem.VirtualMemory(); err == nil {
				memoryUsage.Set(float64(vMem.Used) / (1024 * 1024))
			}

			if db != nil {
				stats := db.Stats()
				dbOpenConnections.Set(float64(stats.OpenConnections))
			}
		}
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Connect to PostgreSQL
	dsn := "host=pg-benchmark port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// Pool settings
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("Connected to PostgreSQL!")
	defer db.Close()

	// Register Prometheus metrics
	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage, dbOpenConnections)
	recordMetrics()

	mux := http.NewServeMux()

	// Hello route
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write(helloWorldBuffer)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write(okResponse)
	})

	// Prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

	// Signup route
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
		if err != nil {
			http.Error(w, "Error inserting user", http.StatusInternalServerError)
			log.Println("DB error:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Signup successful"))
	})

	// Server config
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// Graceful shutdown
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server stopped.")
}
