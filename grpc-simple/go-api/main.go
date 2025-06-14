package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	pb "test/prot/signuppb"

	"github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	port        = ":50051"
	postgresDSN = "host=pg-benchmark port=5432 user=youruser password=yourpassword dbname=yourdb sslmode=disable"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func (s *authServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	if username == "" || password == "" {
		return &pb.SignupResponse{
			Success: false,
			Message: "Username and password are required",
		}, nil
	}

	_, err := s.db.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		// Correct handling of PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return &pb.SignupResponse{
				Success: false,
				Message: "Username already exists",
			}, nil
		}

		log.Printf("Signup DB error: %v", err)
		return &pb.SignupResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	return &pb.SignupResponse{
		Success: true,
		Message: "Signup successful",
	}, nil
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("DB not reachable: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{db: db})

	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}

}
