-- 发送到的 key， 也就是 code:业务:手机号码, 例如 code:login:13800138000
local key = KEYS[1]
-- 使用次数，还可以验证几数，例如 code:login:13800138000:cnt
local cntKey = key .. ":cnt"
-- 验证码，例如 123456
local val = ARGV[1]
-- 验证码有效期，例如 60
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- 如果 ttl 为 -1，说明 key 存在，但没有过期时间
    return -2
    -- 如果 ttl 为 -2，说明 key 不存在
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    -- 返回 0，表示发送成功
    return 0
else
    -- 发送太频繁
    return -1
end