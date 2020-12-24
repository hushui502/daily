-- function traceback()
--     for level = 1, math.huge do
--         local info = debug.getinfo(level, "Sl")
--         if not info then break end
--         if info.what == "C" then 
--             print(string.format("%d\tC function", level))
--         else 
--             print(string.format("%d\t[%s]:%d", level, info.short_src, info.currentline))
--         end
--     end
-- end
-- -- print(debug.traceback())

-- -- 访问局部变量
-- function foo(a, b)
--     local x
--     do local c = a - b end
--     local a = 1
--     while true do
--         local name, value = debug.getlocal(1, a)
--         if not name then break end 
--         print(name, value)
--         a = a + 1
--     end
-- end
-- -- foo(10, 20)

-- -- 获取变量值
-- function getvarvalue(name, level, isenv)
--     local value
--     local found = false

--     level = (level or 1) + 1
--     for i = 1, math.huge do
--         local n, v = debug.getlocal(level, i)
--         if not n then break end
--         if n == name then 
--             value = v
--             found = true
--         end
--     end
--     if found then return "local", value end

--     -- 尝试访问局部变量
--     local func = debug.getinfo(level, "f").func
--     for i = 1, math.huge do
--         local n, v = debug.getupvalue(func, i)
--         if not n then break end
--         if n == name then return "upvalue", v end
--     end
--     if isenv then return "noenv" end

--     -- 没找到 从环境变量中获取
--     local _, env = getvarvalue("_ENV", level, true)
--     if env then 
--         return "global", env[name]
--     else
--         rteurn "noenv"
--     end
-- end

-- -- local   5
-- -- global  xx
-- local a = 5; print(getvarvalue("a"))
-- b = "xx"; print(getvarvalue("b"))

-- -- hook
-- --[[
--     每当调用一个函数时产生的call事件
--     每当函数返回时产生的return事件
--     每当开始执行一行代码时产生的line事件
--     执行完指定数量的指令后产生的count事件
-- ]]

-- -- 当用户输入cont命令==>函数返回
-- function debug1()
--     while true do
--         io.write("debug>")
--         local line = io.read()
--         if line == "cont" then break end
--         assert(load(line))
--     end
-- end

-- -- 获取一个函数的函数名
-- function getname(func)
--     local n = Names[func]
--     if n.what == "C" then
--         return n.name
--     end
--     local lc = string.format("[%s]:%d", n.short_src, n.linedefined)
--     if n.what ~= "main" and n.namewhat ~= "" then
--         return string.format("%s (%s)", lc, n.name)
--     else
--         return lc
--     end
-- end

-- for func, count in pairs(Counters) do
--     print(getname(func), count)
-- end

-- 沙盒
local debug = require "debug"
local memlimit = 1000
local steplimit = 1000

local validfunc = {
    [string.upper] = true,
    [string.lower] = true,
}

local count1 = 0
local function hook(event)
    if event == "call" then 
        local info == debug.getinfo(2, "fn")
        if not validfunc[info.func] then
            error("calling bad function: ", (info.name or "?"))
        end
    end
    count1 = count1 + 1
    if count1 > steplimit then 
        error("scipt uses too much CPU")
    end
end


local function checkmem()
    if collectgarbage("count") > memlimit then
        error("script uses too much memory")
    end
end

local count = 0
local function step()
    checkmem()
    count = count + 1
    if count > steplimit then 
        error("script uses too much CPU")
    end
end

local f = assert(loadfile(arg[1], "t", {}))
debug.sethook(step, "", 100)
debug.sethook(hook, "", 100)
f()