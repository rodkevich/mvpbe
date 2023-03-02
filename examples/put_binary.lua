wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.body = "{\"status\": \"PENDING\"}"

-- file = io.open("sample_item_request.json", "rb")
-- wrk.body = file:read("*a")

--  wrk "http://localhost:8080/api/v1/items/1" -s ./examples/put_binary.lua --latency -t 5 -c 10 -d 10

