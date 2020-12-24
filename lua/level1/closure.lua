-- 函数是first class
a = {p = print}
a.p("hello closure")
print = math.sin    -- print 指向 math.sin函数
a.p(print(1))
math.sin = a.p
math.sin(10, 20)    -- math.sin 指向 print函数

foo = function (x) return 2*x end
print(foo(2))

network = {
    {name = "libai", ip = "11.23.44.11"},
    {name = "dufu", ip = "12.44.11.22"}
}
table.sort(network, function (a, b) return a.name > b.name end)
for k, v in ipairs(network) do
    print(string.format( "index=%d, name=%s, ip=%s", k, v.name, v.ip ))
end

function derivate (f, delta)
    delta = delta or 1e-4
    return function (x)
        return (f(x+delta)-f(x)/delta)
    end
end
print(derivate(math.sin)(4))

-- no global function
Lib = {}
Lib.foo = function (x, y) return x + y end 
Lib.goo = function (x, y) return x - y end
print(Lib.foo(2, 3), Lib.goo(2, 3))


local fact = function (n)
    if n == 0 then return 1
    else return n*fact(n-1)     -- bug,此时的fact是局部的，并未定义，所以这里使用的是全局fact
    end
end
-- fix bug
local fact
fact = function (n)
    if n == 0 then return 1
    else return n*fact(n-1)
    end
end


-- 词法定界
function sortbygrade( names, grades )
    table.sort(names, function (n1, n2)
        return grades[n1] > grades[n2]
    end)
end

-- local value escape
function newCounter()
    local count = 0
    return function () 
        count = count + 1
        return count
    end
end

c1 = newCounter()
print(c1())
print(c1())


-- function digitButton( digit )
--     return Button{label = tostring(digit),
--                 action = function ()
--                     add_to_display(digit)
--                     end
--                 }
-- end

local oldSin = math.sin
math.sin = function (x)
    return oldSin(x * (math.pi / 180)) 
end

-- redeclare a function, do-end create a sandbox
do 
    local oldSin = math.sin
    local k = math.pi / 180
    math.sin = function (x)
        return oldSin(x * k)
    end
end

-- open file    自定义一个沙盒
do 
    local oldOpen = io.open
    local access_OK = function (filename, mode)
        -- check access
        local f = assert(io.open(filename, mode))
    end 
    io.open = function (filename, mode)
        -- if access ok
        if access_OK(filename, mode) then 
            return oldOpen(filename, mode)
        else
            retun nil, "access denied!"
        end
    end
end


-- function programming
function disk( cx, cy, r )
    return function(x, y) 
        (x - cx)^2 + (y - cy)^2 <= r^2
    end
end

-- rect
function rect(left, right, bottom, up)
    return function (x, y)
        return left <= x  and x <= right and 
            bottom <= y and y <= up
    end
end

-- union diff

function union( r1, r2 )
    return function (x, y)
        return r1(x, y) and r2(x, y)
    end
end

function intersection( r1, r2 )
    return function (x, y)
        retun r1(x, y) and r2(x, y)
    end
end

function difference( r1, r2 )
    return function(x, y)
        return r1(x, y) and not r2(x, y)
    end
end

function translate( r, dx, dy )
    return function (x, y)
        return r(x-dx, y-dy)
    end
end

function plot( r, M, N )
    io.write("P1\n", M, " ", N, "\n")
    for i = 1, N do 
        local y = (N - i*2)/N
        for j = 1, M do
            local x = (j*2-M)/M
            io.write(r(x, y) and "1" or "0")
        end
        io.write("\n")
    end
end

c1 = disk(0, 0, 1)
plot(difference(c1, translate(c1, 0.3, 0)), 500, 500)