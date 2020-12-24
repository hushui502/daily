-- lua无法区分常量还是变量

-- 输出全局环境中的所有全局变量
for n in pairs(_G) do print(n) end

print(_G[arg])

function getfield( f )
    local v = _G
    for w in string.gmatch( f, "[%a_][%w_]*" ) do
        v = v[w]
    end
    return v
end

function setfield( f, v )
    local t = _G
    for w, d in string.gmatch(f, "([%a_][%w_]*)(%.?)") do
        if d == "." then 
            t[w] = t[w] or {}
            t = t[w]
        else
            t[w] = v
        end
    end
end

setfield("t.x.y", 10)
print(t.x.y)
print(getfield("t.x.y"))

-- 全局变量声明
setmetatable(_G, {
    __newindex = function (_, n)
        error("undeclared variable" .. n, 2)
    end,
    __index = function (_, n) 
        error("undeclared variable" .. n, 2)
    end,
})
-- 访问全局会出错
print(a)

function declare( name, initval )
    rawset(_G, name, initval or false)
end

var = "pairs"
if rawget(_G, var) == nil then 
    -- ...
end


-- 检测全局变量声明
local declareNames = {}
setmetatable(_G, {
    __newindex = function (t, n, v) 
        if not declareNames[n] then
            local w = debug.getinfo(2, "S").what
            if w ~= "main" and w ~= "C" then
                error("attempt to write to undeclared variable" .. n, 2)
            end
            declareNames[n] = true
        end
        rawset(t, n, v)
    end,

    __index = function (_, n) 
        if not declareNames[n] then
            error("attempt to write to undeclared variable" .. n, 2)
        else
            return nil
        end
    end,
})

_ENV 
local print, sin = print, math.sin
_ENV = nil
print(sin(13))
print(math.cos(33)) --> error,只要_ENV=nil之后的代码就无法直接访问全局变量

a = 12
local a = 22
print(a)        --> 22
print(_ENV.a)   --> 12
print(_G.a)     --> 12  通常_G and _ENV 都指向同一个表

_ENV 改变代码段环境
_ENV = {}
a = 1           --> 在_ENV 中创建一个字段
print(a)        --> error 访问全局

先把一些有用的值放入新环境，比如全局环境
a = 13
_ENV = {g = _G}
a = 1
g.print(_ENV.a, g.a)        --> 1   13

a = 1
local newgt = {}
setmetatable(newgt, {__index = _G})
_ENV = newgt
print(a)