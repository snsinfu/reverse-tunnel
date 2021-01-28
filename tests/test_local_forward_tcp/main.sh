echo "Should forward local connection"
set -e

echo "* Starting tunneling server..."
rtun-server -f rtun-server.yml &
pid_server=$!

sleep 1

echo "* Starting tunneling agent..."
rtun -f rtun.yml &
pid_agent=$!

sleep 1

echo "* Starting local echo server..."
go run echoserver.go 127.0.0.1:8080 &

sleep 3

echo "* Testing tunneled connection..."
expect="OK 10e03ca70fcaae2c"
actual="$({ echo "${expect}"; sleep 1; } | nc 127.0.0.1 18080)"
# XXX: nc cannot reliably receive response unless we sleep after echoing. This
# occurs when the connection is tunneled. This would be a timing bug of rtun,
# rtun-server or the tunneling protocol.

sleep 1

echo "* Terminating servers..."
kill -TERM ${pid_agent}
kill -TERM ${pid_server}
wait

echo "* Examininng the result..."
echo "expect: ${expect}"
echo "actual: ${actual}"
test "${expect}" = "${actual}" || exit 1
