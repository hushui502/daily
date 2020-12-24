-- 局部变量， 尽量使用local局部变量，避免因为命名造成的变量混乱
-- 缩小变量作用域有助于提高代码的可读性
x = 10
local i = 1
while i <= x do
    local x = i * 2
    print(x)
    i = i + 1
end

if i > 20 then
    local x
    x = 20
    print(x + 2)
else 
    print(x)    -- 10
end
print(x)        -- 10

local a, b = 1, 10
if a < b then
    print(a)    -- 1
    local a     -- nil
    print(a)
end
print(a, b)     -- 1, 10


-- 控制结构,除了false和nil都认为是true
-- if,不支持switch
if a < b then return a else return b end

if op == "+" then 
    r = a + b
elseif op == "-" then 
    r = a - b
elseif op == "*" then 
    r = a * b
elseif op == "/" then 
    r = a / b
else 
    error("invalid operation")
end

-- while
local i = 1
while a[i] do
    print(a[i])
    i = i + 1
end

-- repeat,重复执行某一逻辑直到...
local line 
repeat
    line = io.read()
until line ~= ""
print(line)

-- 计算x的平方根
function sqr( x )
    local sqr = x / 2
    repeat
        sqr = (sqr + x/sqr) / 2
        local error = math.abs( sqr^2 - x )
    until error < x/10000
    print(sqr)
end
sqr(6)

-- for 
-- 不给循环设置上限
for i = 1, math.huge do 
    print(i)
    if i == 10 then
        break
    end
end

-- goto
::s1:: do
    local c = io.read(1)
    if c == '0' then goto s2
    elseif c == nil then print "ok" return 
    else goto s1
    end 
end
::s2:: do
    local c = io.read(1)
    -- ...
end