local userKey = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local windowStart = now - window

redis.call('ZREMRANGEBYSCORE', userKey, 0, windowStart)
local count = redis.call('ZCARD', userKey)

if count < limit then 
    redis.call('ZADD', userKey, now, now)
    redis.call('EXPIRE', userKey, window/1000000000)
    return 1
end

return 0