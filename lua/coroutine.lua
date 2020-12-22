-- co = coroutine.create( function () print("hi") end )
-- print(type(co))     --> thread

-- --[[
--     suspended   挂起
--     running     运行
--     normal      正常
--     dead        死亡
-- ]]
-- print(coroutine.status( co ))
-- coroutine.resume( co )   
-- print(coroutine.status( co ))

-- co = coroutine.create(function ()
--         for i = 1, 4 do
--             print("co", i)
--             coroutine.yield()
--         end
--     end)

-- coroutine.resume(co)
-- coroutine.resume(co)
-- coroutine.resume(co)
-- coroutine.resume(co)
-- coroutine.resume(co)
-- print(coroutine.resume(co))     --> false   cannot resume dead coroutine
-- print(coroutine.status( co ))


-- -- consumer producer
-- function producer()
--     while true do
--         local x = io.read()
--         send(x)
--     end
-- end

-- function consumer()
--     while true do
--         local x = receive()
--         io.write(x, "\n")
--     end
-- end

-- function receive()
--     local status, value = coroutine.resume(producer)
--     return true
-- end

-- function send(x)
--     coroutine.yield(x)
-- end

-- producer = coroutine.create(producer)

-- -- filter
-- function receive( prod )
--     local status, value = coroutine.resume(prod)
--     return value
-- end

-- function send( x )
--     coroutine.yield(x)
-- end

-- function producer()
--     return coroutine.create(function ()
--         while true do
--             local x = io.read()
--             send(x)
--         end
--     end)
-- end

-- function filter( prod )
--     return coroutine.create(function ()
--         for line = 1, math.huge do
--             local x = receive(prod)
--             x = string.format("%5d %s", line, x)
--             send(x)
--         end
--     end)
-- end

-- function consumer( prod )
--     while true do
--         local x = receive(prod)
--         io.write(x, "\n")
--     end
-- end

-- consumer(filter(producer))

-- sort
function permgen(a, n)
    n = n or #a
    if n <= 1 then
        coroutine.yield(a)
    else
        for i = 1, n do
            a[n], a[i] = a[i], a[n]
            permgen(a, n-1)
            a[n], a[i] = a[i], a[n]
        end
    end
end

function permutations(a)
    return coroutine.wrap( function () permgen(a) end )
end

function printResult( a )
    for i = 1, #a do io.write(a[i], " ") end
    io.write("\n")
end

for p in permutations{"a", "b", "c"} do
    printResult(p)
end

-- asyn io
local cmdQueue = {}
local lib = {}

function lib.readline( stream, callback )
    local nextCmd = function () 
        callback(stream:read())
    end
    table.insert(cmdQueue, nextCmd)
end

function lib.writeline( stream, line, callback )
    local nextCmd = function ()
        callback(stream:write(line))
    end
    table.insert(cmdQueue, nextCmd)
end

function lib.stop(  )
    table.insert(cmdQueue, "stop")
end

function lib.runloop()
    while true do
        local nextCmd = table.remove(cmdQueue, 1)
        if nextCmd == "stop" then
            break
        else
            nextCmd()
        end
    end
end

return lib

-- 同步IO
local t = {}
local inp = io.input()
local out = io.output()

for line in inp:lines() do 
    t[#t + 1] = line
end

for i = #t, 1, -1 do 
    out:write(t[i], "\n")
end

-- 异步重写
function run(code)
    local co = coroutine.wrap(function ()
        code()
        lib.stop()  -- 结束时候关闭事件循环
    end)
    co()            -- 启动协程
    lib.runloop()   -- 启动事件循环
end

local function putline(stream, line)
    local co = coroutine.running()      -- 启动协程
    local callback = (function () coroutine.resume(co) end)
    lib.writeline(stream, line, callback)
    coroutine.yield()
end    

local function getline(strem, line)
    local co = coroutine.running()
    local callback = (function (l) coroutine.resume(co, l), end)
    lib.readline(stream, callback)
    local line = coroutine.yield()
    return line
end

run(function ()
    local t = {}
    local inp = io.input()
    local out = io.output()

    while true do 
        local line = getline(inp)
        if not line then break end
        t[#t + 1] = line
    end

    for i = #t, 1, -1 do
        putline(out, t[i] .. "\n")
    end
end)