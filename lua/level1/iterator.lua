function values( t )
    local i = 0
    return function () i = i + 1; return t[i] end
end

t = {1, 2, 3}
iter = values(t)
while true do
    local element = iter()
    if element == nil then break end 
    print(element)
end

t = {3, 4, 5}
for element in values(t) do
    print(element)
end

-- 遍历标准输入所有单词
function allwords()
    local line = io.read()      -- 当前行
    local pos = 1               -- 当前行的当前位置
    return function ()          -- 迭代函数
        while line do           -- 当在前行时循环
            local w, e = string.match(line, "(%w+)()", pos)
            if w then           -- 发现一个单词
                pos = e         -- 下一个单词将位于该单词之后
                return w        -- 返回该单词
            else
                line = io.read()-- 没找到单词继续下一行
                pos = 1         -- 从下一行的第一个位置重新开始
            end
        end 
        return nil              -- 没有行了，迭代结束
    end
end

for word in allwords() do
    print(word)
end

-- ipairs
local function iter (t, i)
    i = i + 1
    local v = t[i]
    if v then 
        return i, v
    end
end

function ipairs (t) 
    return iter, t, 0
end

-- next
local function getnext( list, node )
    if not node then 
        return list
    else
        return node.next
    end
end

function traverse(list)
    return getnext, list, nil
end

-- traverse table
function pairsByKeys( t, f )
    local a = {}
    for n in pairs(t) do
        a[#a+1] = n
    end
    table.sort( a, f )
    local i = 0
    return function ()
        i = i + 1
        return a[i], t[a[i]]
    end
end

for name, line in pairsByKeys(lines) do
    print(name, line)
end

-- iterator
function allwords( f )
    for line in io.lines() do
        for word in string.gmatch(line, "%w+") do
            f(word)
        end
    end
end
allwords(print)

local count = 0
allwords(function (w)
    if w == "hello" then count = count + 1 end
end)
print(count)

-- 马尔科夫链
local statetab = {}

function prefix( w1, w2 )
    return w1 .. " " .. w2
end

function insert( prefix, value )
    local list = statetab[prefix]
    if list == nil then 
        statetab[prefix] = {value}
    else
        list[#list + 1] = value
    end
end

function allwords1()
    local line = io.read()
    local pos = 1
    return function () 
        while line do
            local w, e = string.match(line, "(%w+[,;.:]?)()", pos)
            if w then
                pos = e
                return w
            else
                line = io.read()
                pos = 1
            end
        end
        return nil
    end
end

local MAXGEN = 200
local NOWORD = "\n"
local w1, w2 = NOWORD, NOWORD
for nextword in allwords() do
    insert(prefix(w1, w2), nextword)
    w1 = w2; w2 = nextword;
end
insert(prefix(w1, w2), NOWORD)
                                    -- 开始生成文本
w1 = NOWORD; W2 = NOWORD            -- 重新初始化
for i = 1, MAXGEN do
    local list = statetab[prefix(w1, w2)]
    local r = math.random(#list)    -- 随机选出一个元素
    local nextword = list[r]
    if nextword == NOWORD then return end
    io.write(nextword, " ")
    w1 = w2; w2 = nextword
end