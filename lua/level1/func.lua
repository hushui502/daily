-- case 1
function add( a )
    local sum = 0
    for i = 1, #a do 
        sum = sum + a[i]
    end
    return sum
end

-- case 2
globalCounter = 0
function incCount( n )
    -- for avoiding a "blank" value
    n = n or 1
    globalCounter = globalCounter + n
end

incCount(2)
incCount()
print(globalCounter)    --> 3

-- multiple retun values
s, e = string.find( "hello lua", "lua")
print(s, e)     --> 7, 9

function maximum( a )
    local mi = 1
    local m = a[mi]
    for i = 1, #a do 
        if a[i] > m then 
            mi = i
            m = a[mi]
        end
    end
    return m, mi
end

print(maximum({1, 2, 3, 4}))

-- multiple arguments
function add( ... )
    local s = 0
    for _, v in ipairs{...} do
        s = s + v
    end
    return s
end

-- vararg expression
local a, b = ...
function foo( ... )
    local a, b, c = ...
end

function fwrite( fmt, ... )
    return io.write(string.format( fmt, ... ))
end

-- table.pack 
function nonils( ... )
    local arg = table.pack(...)
    -- arg.n ==> len of the table
    for i = 1, arg.n do
        if arg[i] == nil then return false end
    end 
    return true
end
print(nonils(1, 2))
print(nonils(1, 2, nil))

-- table.unpack
-- genic programing
print(table.unpack({1, 2, 3}))
a, b = table.unpack({1, 2, 3}) --> 3 will be discarded

print(table.unpack({1, 2, 3}, 2, 3))

function unpack( t, i, n )
    i = i or 1
    n = n or #t

    if i <= n then
        retun t[i], unpack(t, i+1, n)
    end
end

-- selector
print(select(1, "a", "b", "c"))
print(select(2, "a", "b", "c"))
print(select("#", "a", "b", "c"))

function add2( ... )
    local s = 0
    for i = 1, select("#", ...) do
        -- print("num:", select(i, ...))
        s = s + select(i, ...)
    end
    return s
end

print(add2(2, 3))

-- tail-call elimination
function foor( n )
    if n > 0 then retun foor(n-1) end
end