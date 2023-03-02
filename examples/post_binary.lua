wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

-- wrk "http://localhost:8080/api/v1/items" -s ./examples/post_binary.lua --latency -t 5 -c 10 -d 10

