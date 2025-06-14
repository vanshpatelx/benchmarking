// package main

// import (
// 	"context"
// 	"database/sql"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"runtime"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// 	"google.golang.org/grpc"

// 	pb "test/prot/signuppb"
// 	_ "github.com/lib/pq"
// )

// var (
// 	db *sql.DB

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

// 	dbOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "db_open_connections",
// 		Help: "Current open DB connections",
// 	})

// 	okResponse       = []byte("OK")
// 	helloWorldBuffer = []byte("Hello, World!")
// )

// func recordMetrics() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	defer ticker.Stop()

// 	go func() {
// 		for range ticker.C {
// 			if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
// 				cpuUsage.Set(cpuPercents[0])
// 			}

// 			if vMem, err := mem.VirtualMemory(); err == nil {
// 				memoryUsage.Set(float64(vMem.Used) / (1024 * 1024))
// 			}

// 			if db != nil {
// 				stats := db.Stats()
// 				dbOpenConnections.Set(float64(stats.OpenConnections))
// 			}
// 		}
// 	}()
// }

// type authServer struct {
// 	pb.UnimplementedAuthServiceServer
// 	db *sql.DB
// }

// func (s *authServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
// 	username := req.GetUsername()
// 	password := req.GetPassword()

// 	if username == "" || password == "" {
// 		return &pb.SignupResponse{
// 			Message: "Username and password required",
// 			Success: false,
// 		}, nil
// 	}

// 	_, err := s.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
// 	if err != nil {
// 		log.Printf("DB insert error: %v", err)
// 		return &pb.SignupResponse{
// 			Message: "Error inserting user",
// 			Success: false,
// 		}, nil
// 	}

// 	return &pb.SignupResponse{
// 		Message: "Signup successful",
// 		Success: true,
// 	}, nil
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

// 	// PostgreSQL connection
// 	dsn := "host=pg-benchmark port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
// 	var err error
// 	db, err = sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to open DB: %v", err)
// 	}
// 	db.SetMaxOpenConns(500)
// 	db.SetMaxIdleConns(100)
// 	db.SetConnMaxLifetime(time.Hour)

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Failed to connect to DB: %v", err)
// 	}
// 	log.Println("Connected to PostgreSQL!")

// 	defer db.Close()

// 	// Metrics registration
// 	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage, dbOpenConnections)
// 	recordMetrics()

// 	// Start HTTP server for Prometheus and health routes
// 	go func() {
// 		mux := http.NewServeMux()

// 		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 			httpRequests.Inc()
// 			w.WriteHeader(http.StatusOK)
// 			w.Write(helloWorldBuffer)
// 		})

// 		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 			httpRequests.Inc()
// 			w.WriteHeader(http.StatusOK)
// 			w.Write(okResponse)
// 		})

// 		mux.Handle("/metrics", promhttp.Handler())

// 		server := &http.Server{
// 			Addr:              ":8080",
// 			Handler:           mux,
// 			ReadTimeout:       5 * time.Second,
// 			WriteTimeout:      10 * time.Second,
// 			IdleTimeout:       120 * time.Second,
// 			ReadHeaderTimeout: 3 * time.Second,
// 			MaxHeaderBytes:    1 << 20,
// 		}

// 		log.Println("Starting HTTP server on :8080 for metrics and health")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %v", err)
// 		}
// 	}()

// 	// Start gRPC server
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen on :50051: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	pb.RegisterAuthServiceServer(grpcServer, &authServer{db: db})

// 	go func() {
// 		log.Println("Starting gRPC server on :50051")
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatalf("gRPC server error: %v", err)
// 		}
// 	}()

// 	// Graceful shutdown
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt)
// 	<-stop

// 	log.Println("Shutting down servers...")
// 	grpcServer.GracefulStop()
// 	log.Println("gRPC server stopped.")
// }

// package main

// import (
// 	"context"
// 	"database/sql"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"runtime"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/keepalive"

// 	pb "test/prot/signuppb"
// 	_ "github.com/lib/pq"
// )

// var (
// 	db *sql.DB

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

// 	dbOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "db_open_connections",
// 		Help: "Current open DB connections",
// 	})

// 	okResponse       = []byte("OK")
// 	helloWorldBuffer = []byte("Hello, World!")
// )

// func recordMetrics() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	defer ticker.Stop()

// 	go func() {
// 		for range ticker.C {
// 			if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
// 				cpuUsage.Set(cpuPercents[0])
// 			}

// 			if vMem, err := mem.VirtualMemory(); err == nil {
// 				memoryUsage.Set(float64(vMem.Used) / (1024 * 1024))
// 			}

// 			if db != nil {
// 				stats := db.Stats()
// 				dbOpenConnections.Set(float64(stats.OpenConnections))
// 				log.Printf("DB Stats => Open: %d | InUse: %d | Idle: %d",
// 					stats.OpenConnections, stats.InUse, stats.Idle)
// 			}
// 		}
// 	}()
// }

// type authServer struct {
// 	pb.UnimplementedAuthServiceServer
// 	db *sql.DB
// }

// func (s *authServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
// 	start := time.Now()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Printf("Signup panic: %v", r)
// 		}
// 		log.Printf("Handled signup in %v", time.Since(start))
// 	}()

// 	select {
// 	case <-ctx.Done():
// 		log.Println("Signup request canceled or timed out")
// 		return &pb.SignupResponse{
// 			Message: "Request timeout or canceled",
// 			Success: false,
// 		}, nil
// 	default:
// 	}

// 	username := req.GetUsername()
// 	password := req.GetPassword()

// 	if username == "" || password == "" {
// 		return &pb.SignupResponse{
// 			Message: "Username and password required",
// 			Success: false,
// 		}, nil
// 	}

// 	_, err := s.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
// 	if err != nil {
// 		log.Printf("DB insert error for user %s: %v", username, err)
// 		return &pb.SignupResponse{
// 			Message: "Error inserting user",
// 			Success: false,
// 		}, nil
// 	}

// 	return &pb.SignupResponse{
// 		Message: "Signup successful",
// 		Success: true,
// 	}, nil
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

// 	dsn := "host=pg-benchmark port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
// 	var err error
// 	db, err = sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to open DB: %v", err)
// 	}
// 	db.SetMaxOpenConns(500)
// 	db.SetMaxIdleConns(100)
// 	db.SetConnMaxLifetime(time.Hour)

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Failed to connect to DB: %v", err)
// 	}
// 	log.Println("Connected to PostgreSQL!")

// 	defer db.Close()

// 	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage, dbOpenConnections)
// 	recordMetrics()

// 	// HTTP server for Prometheus + health
// 	go func() {
// 		mux := http.NewServeMux()

// 		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 			httpRequests.Inc()
// 			w.WriteHeader(http.StatusOK)
// 			w.Write(helloWorldBuffer)
// 		})

// 		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 			httpRequests.Inc()
// 			w.WriteHeader(http.StatusOK)
// 			w.Write(okResponse)
// 		})

// 		mux.Handle("/metrics", promhttp.Handler())

// 		server := &http.Server{
// 			Addr:              ":8080",
// 			Handler:           mux,
// 			ReadTimeout:       5 * time.Second,
// 			WriteTimeout:      10 * time.Second,
// 			IdleTimeout:       120 * time.Second,
// 			ReadHeaderTimeout: 3 * time.Second,
// 			MaxHeaderBytes:    1 << 20,
// 		}

// 		log.Println("Starting HTTP server on :8080")
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %v", err)
// 		}
// 	}()

// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen on :50051: %v", err)
// 	}

// 	grpcServer := grpc.NewServer(
// 		grpc.KeepaliveParams(keepalive.ServerParameters{
// 			MaxConnectionIdle: 5 * time.Minute,
// 			MaxConnectionAge:  30 * time.Minute,
// 			Time:              15 * time.Second,
// 			Timeout:           10 * time.Second,
// 		}),
// 	)

// 	pb.RegisterAuthServiceServer(grpcServer, &authServer{db: db})

// 	go func() {
// 		log.Println("Starting gRPC server on :50051")
// 		if err := grpcServer.Serve(lis); err != nil {
// 			log.Fatalf("gRPC server error: %v", err)
// 		}
// 	}()

// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt)
// 	<-stop

// 	log.Println("Shutting down servers...")
// 	grpcServer.GracefulStop()
// 	log.Println("gRPC server stopped.")
// }

package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	_ "github.com/lib/pq"
	_ "net/http/pprof"
	pb "test/prot/signuppb"
)

var (
	db  *sql.DB
	sem = make(chan struct{}, 10000) // bounded concurrency

	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total", Help: "Total number of HTTP requests",
	})
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent", Help: "CPU usage percentage",
	})
	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_mb", Help: "Memory usage in MB",
	})
	dbOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_open_connections", Help: "Current open DB connections",
	})

	okResponse       = []byte("OK")
	helloWorldBuffer = []byte("Hello, World!")
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func (s *authServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	case <-ctx.Done():
		return &pb.SignupResponse{Message: "Server busy", Success: false}, nil
	}

	if req.GetUsername() == "" || req.GetPassword() == "" {
		return &pb.SignupResponse{Message: "Username and password required", Success: false}, nil
	}

	if runtime.NumGoroutine() > 10000 {
		log.Println("High goroutine count detected, forcing GC")
		runtime.GC()
	}

	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(dbCtx, "INSERT INTO users (username, password) VALUES ($1, $2)", req.GetUsername(), req.GetPassword())
	if err != nil {
		log.Printf("DB error: %v", err)
		return &pb.SignupResponse{Message: "Database error", Success: false}, nil
	}

	return &pb.SignupResponse{Message: "Signup successful", Success: true}, nil
}

func recordMetrics() {
	ticker := time.NewTicker(5 * time.Second)
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
				log.Printf("DB connections: Open=%d InUse=%d Idle=%d Goroutines=%d",
					stats.OpenConnections, stats.InUse, stats.Idle, runtime.NumGoroutine())
			}
		}
	}()
}

func startMetricsHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write(helloWorldBuffer)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write(okResponse)
	})
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		log.Println("HTTP metrics server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
}

func setupDatabase() *sql.DB {
	dsn := "host=pg-benchmark port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB open failed: %v", err)
	}

	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}

func startGRPCServer(db *sql.DB) *grpc.Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		MaxConnectionAge:  30 * time.Minute,
		Time:              15 * time.Second,
		Timeout:           10 * time.Second,
	}),
		grpc.MaxConcurrentStreams(10000),
		grpc.MaxRecvMsgSize(1024*1024), // if needed
	)

		pb.RegisterAuthServiceServer(server, &authServer{db: db})

	go func() {
		log.Println("gRPC server running on :50051")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	return server
}

func startCustomGC() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Log memory usage stats
			log.Printf("HeapAlloc: %.2f MB | Sys: %.2f MB | NumGC: %d | Goroutines: %d",
				float64(m.HeapAlloc)/1024/1024,
				float64(m.Sys)/1024/1024,
				m.NumGC,
				runtime.NumGoroutine(),
			)

			// Force garbage collection if heap usage exceeds 100MB
			if m.HeapAlloc > 100*1024*1024 {
				log.Println("HeapAlloc > 100MB, forcing GC")
				runtime.GC()
			}
		}
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	prometheus.MustRegister(httpRequests, cpuUsage, memoryUsage, dbOpenConnections)

	db = setupDatabase()
	defer db.Close()

	recordMetrics()
	startCustomGC()
	startMetricsHTTPServer()
	grpcServer := startGRPCServer(db)

	// Start pprof server on port 6060
	go func() {
		log.Println("pprof server running on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("pprof server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
