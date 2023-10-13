local key = KEYS[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key .. ":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt == nil or cnt <= 0 then
    -- 说明用户一直输错
    return -1
elseif code == expectedCode then
    -- 输入正确
    -- 用完，不能再用了
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户手一抖输错了
    -- 可验证次数减一
    redis.call("decr", cntKey)
    return -2
end