echo "Should be able to proxy many concurrent HTTP requests"
set -e

echo "* Starting tunneling server..."
timeout 20s rtun-server -f rtun-server.yml &
pid_server=$!

sleep 1

echo "* Starting tunneling agent..."
timeout 20s rtun -f rtun.yml &
pid_agent=$!

sleep 1

echo "* Starting an HTTP server..."
timeout 20s go run server.go 127.0.0.1:8080 &
pid_http_server=$!

sleep 1

echo "* Testing concurrent requests..."
timeout 20s go run client.go 127.0.0.1:18080 100

echo "* Terminating servers..."
kill -TERM ${pid_agent}
kill -TERM ${pid_server}
kill -TERM ${pid_http_server}
wait
