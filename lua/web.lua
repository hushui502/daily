local socket = require "socket"

host = "www.lua.org"
file = "/manual/5.3/manual.html"

c = assert(socket.connect(host, 80))


c:send(request)

repeat
    local s, status, partial = c:receive(2^10)
    io.write(s or partial)
until status == "closed"

c:close()

-- NIO download
function download( host, file )
    local c = assert(socket.connect(host, 80))
    local count = 0
    local request = string.format("GET %s HTTP/1.0\r\nhost: %s\r\n\r\n", file, host)
    s:send(request)
    while true do
        local s, status = receive(c)
        count = count + #s
        if status == "closed" then break end
    end
    c:close()
    print(file, count)
end

function receive( connection )
    connection:settimeout(0)    -- 不阻塞
    local s, status, partial = connection:receive(2^10)
    if status == "timeout" then 
        coroutine.yield(connection)
    end
    return s or partial, status
end

-- 调度器
tasks = {}              -- 所有活跃的任务列表
function get( host, file )
    -- 为任务创建协程
    local co = coroutine.wrap(function ()
        download(host, file)
    end)
    -- 将其插入列表
    table.insert(tasks, co)
end

function dispatch()
    local i = 1
    local timeout = {}
    while true do
        if tasks[i] == nil then         -- 没有其他任务了？
            if tasks[1] == nil then     -- 列表为空？
                break                   -- 从循环中退出
            end
            i = 1                       -- 否则继续循环
        end
        local res = tasks[i]()          -- 运行一个任务
        if not res then                 -- 任务结束？
            table.remove(tasks, i)
        else
            i = i + 1                   -- 处理下一个任务
            timeout[#timeout + 1] = res
            if #timeout == #tasks then  -- 所有任务都阻塞？ 
                socket.select(timeout)  -- 等待
            end
        end
    end
end

-- get("www.lua.org", "/ftp/...1")
-- get("www.lua.org", "/ftp/...2")
-- get("www.lua.org", "/ftp/...3")

