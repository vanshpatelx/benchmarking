<!-- 
docker network create monitoring-net

docker run -d \
  --name go-api \
  --cpus="1.0" \
  --memory="512m" \
  --network monitoring-net \
  -p 8080:8080 \
  benchmarktest-go-api

docker run -d \
  --name prometheus \
  --cpus="1.0" \
  --memory="512m" \
  --network monitoring-net \
  -p 9090:9090 \
  -v "$(pwd)/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml" \
  prom/prometheus \
  --config.file=/etc/prometheus/prometheus.yml

docker run -d \
  --name grafana \
  --cpus="1.0" \
  --memory="512m" \
  --network monitoring-net \
  -p 3000:3000 \
  -e GF_SECURITY_ADMIN_PASSWORD=admin \
  grafana/grafana

docker stats -->


# Create the Docker network
docker network create monitoring-net

# Run your Go API container
docker run -d \
  --name go-api \
  --cpus="1.0" \
  --memory="512m" \
  --network monitoring-net \
  benchmarktest-go-api


docker build -t go-api-i ./go-api

docker run -d \
  --name go-api \
  --cpus="1.0" \
  --memory="512m" \
  --network monitoring-net \
  go-api-i   



# Run Prometheus container
docker run -d \
  --name prometheus \
  --cpus="0.25" \
  --memory="512m" \
  --network monitoring-net \
  -p 9090:9090 \
  -v "$(pwd)/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml" \
  prom/prometheus \
  --config.file=/etc/prometheus/prometheus.yml

# Run Grafana container
docker run -d \
  --name grafana \
  --cpus="0.25" \
  --memory="512m" \
  --network monitoring-net \
  -p 3000:3000 \
  -e GF_SECURITY_ADMIN_PASSWORD=admin \
  grafana/grafana

# Run NGINX container
docker run -d \
  --name nginx \
  --cpus="2.0" \
  --memory="512m" \
  --network monitoring-net \
  -p 8080:80 \
  -v "$(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro" \
  nginx:alpine

# Show real-time container stats
docker stats


docker run -d \
  --name pg-benchmark \
  --network monitoring-net \
  -p 5432:5432 \
  -e POSTGRES_USER=youruser \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=yourdb \
  -v pg_data:/var/lib/postgresql/data \
  -v "$(pwd)/init.sql":/docker-entrypoint-initdb.d/init.sql \
  postgres:15 \
  -c max_connections=200 \
  -c shared_buffers=256MB \
  -c work_mem=16MB \
  -c maintenance_work_mem=64MB \
  -c effective_cache_size=512MB \
  -c logging_collector=on \
  -c log_destination=stderr
