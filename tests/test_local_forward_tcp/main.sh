echo "Should forward local connection"
set -e

echo "* Building echo server and client..."
go build -o _echoserver echoserver.go
go build -o _echoclient echoclient.go

echo "* Starting tunneling server..."
timeout 20s rtun-server -f rtun-server.yml &
pid_server=$!

sleep 1

echo "* Starting tunneling agent..."
timeout 20s rtun -f rtun.yml &
pid_agent=$!

sleep 1

echo "* Starting local echo server..."
./_echoserver 127.0.0.1:8080 &

sleep 1

echo "* Testing tunneled connection..."
expect="OK 10e03ca70fcaae2c"
actual="$(echo "${expect}" | ./_echoclient 127.0.0.1:18080)"

sleep 1

echo "* Terminating servers..."
kill -TERM ${pid_agent}
kill -TERM ${pid_server}
wait

echo "* Examininng the result..."
echo "expect: ${expect}"
echo "actual: ${actual}"
test "${expect}" = "${actual}" || exit 1
