-- 矩阵
function creatematrix1( N, M )
    local mt = {}
    for i = 1, N do
        local row = {}
        mt[i] = row
        for j = 1, M do
            row[j] = 0
        end
    end
end

function creatematrix2( N, M )
    local mt = {}
    for i = 1, N do
        local aux = (i - 1) * M
        for j = 1, M do
            mt[aux + j] = 0
        end
    end
end

-- 稀疏矩阵相乘
function mult( a, b )
    local c = {}
    for i = 1, #a do
        local resultline = {}
        for k, va in pairs(a[i]) do
            for j, vb in pairs(b[k]) do
                local res = (resultline[j] or 0) + va * vb
                resultline[j] = (res ~= 0) and res or nil
            end
        end
        c[i] = resultline
    end
    return c
end


-- 链表
list = nil
list = {next = list, value = v}
local l = list
while l do
    if l.value then l = l.next else break end
end

-- 双端队列
function listNew()
    return {first = 0, last = -1}
end

function pushFirst(list, value)
    local first = list.first - 1
    list.first = first
    list[first] = value
end

function pushLast(list, value)
    local last = list.last + 1
    list.last = last
    list[last] = value
end

function popFirst(list)
    local first = list.first
    if first > list.last then error("list is empty") end
    local value = list[first]
    list[first] = nil
    list.first = first + 1
    return value 
end

function popLast(list)
    local last = list.last
    if list.first > last then error("list is empty") end 
    local value = list[last]
    list[last] = nil
    line.last = last - 1
    return value
end

-- 反向表
days = {"Sun", "Mon", "Tue", "Web", "Thu", "Fri", "Sat"}
revDays = {}
for k, v in pairs(days) do
    revDays[v] = k
end
x = "Tue"
print(revDays[x])   --> 3

-- 集合和包
reserved = {
    ["while"] = true,
    ["if"] = true,
    ["else"] = true,
    ["do"] = true,
}

s = "[if]"
for w in string.gmatch(s, "[%a_][%w_]*") do
    if not reserved[w] then 
        -- do something with 'w'
    end
end

function Set( list )
    local set = {}
    for _, l in ipairs(list) do 
        set[l] = true
    end
    return set
end
reserved2 = Set{"whike", "if", "else", "do"}
print(reserved2["if"])

local ids = {}
for w in string.gmatch(s, "[%a_][%w_]*") do
    if not reserved[w] then 
        ids[w] = true
    end
end

function insert( bag, element )
    bag[element] = (bag[element] or 0) + 1
end

function remove( bag, element )
    local count = bag[element]
    bag[element] = (count and count > 1) and count - 1 or nil
end

-- huge performance problem 会导致不停的移动,尽量避免 .. "str"操作
local buff = ""
for line in io.lines() do 
    buff  = buff .. line .. "\n"
end

-- 用表当一个字符串缓冲区
local t = {}
for line in io.lines() do
    t[#t+1] = line
end
local s = table.concat(t, "\n")
t[#t + 1] = ""
s = table.concat(t, "\n")   -- 在尾部加一个"\n", 这里的操作完全是为了避免..字符串拼接

-- 图
local function name2node (graph, name)
    local node = graph[name]
    if not node then 
        node = {name = name, adj = {}}
        graph[name] = node
    end
    return node
end

function readgraph()
    local graph = {}
    for line in io.lines() do
        local namefrom, nameto = string.match(line, "(%S+)%s+(%S)")
        local from = name2node(graph, namefrom)
        local to = name2node(graph, nameto)
        from.adj[to] = true
    end
    return graph
end

function findpath( curr, to, path, visited )
    path = path or {}
    visited = visited or {}
    if visited[curr] then
        return nil
    end
    visited[curr] = true
    path[#path + 1] = curr
    if curr == to then 
        return path
    end
    for node in pairs(curr.adj) do
        local p = findpath(node, to, path, visited)
        if p then return p end
    end
    table.remove(path)
end

function printpath( path )
    for i = 1, #path do
        print(path[i].name)
    end
end