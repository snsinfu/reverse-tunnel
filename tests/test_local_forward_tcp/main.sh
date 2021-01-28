echo "Should forward local connection"

echo "* Starting tunneling server..."
rtun-server -f rtun-server.yml &
pid_server=$!

sleep 1

echo "* Starting tunneling agent..."
rtun -f rtun.yml &
pid_agent=$!

sleep 1

echo "* Starting local echo server..."
expect="OK 10e03ca70fcaae2c"
echo "${expect}" | nc -l localhost 8080 &

sleep 1

echo "* Testing tunneled connection..."
actual="$(echo "${expect}" | nc localhost 18080)"

echo "* Terminating servers..."
kill -TERM ${pid_agent}
kill -TERM ${pid_server}
wait

echo "* Examininng the result..."
echo "expect: ${expect}"
echo "actual: ${actual}"
test "${expect}" = "${actual}" || exit 1
