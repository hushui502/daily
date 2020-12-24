-- 只有对象可以从弱引用列表中删除，数字和布尔不可以
-- 字符串也是值不是对象，对于一个字符串类型的键来说，除非它对应的值被回收，否则是不会从弱引用列表中被移除的
a = {}
mt = {__mode = "k"}
setmetatable(a, mt)
key = {}
a[key] = 1
key = {}
a[key] = 2
collectgarbage()
-- 因为第二个键已经替代了第一个键，所以grabage会销毁第一个键,只剩下第二个键
for k, v in pairs(a) do print(v) end    --> 2

-- memorize function
local results = {}
setmetatable(results, {__mode = "v"})   --> 第一避免空间过大，垃圾回收。第二因为这里索引是字符串，所以没有必要弱引用键
function mem_loadstring( s )
    local res = results[s]
    if res == nil then 
        res = assert(load(s))
        results[s] = res
    end
    return res
end

local results = {}
setmetatable(results, {__mode = "v"})
function createRGB(r, g, b)
    local key = string.format( "%d-%d-%d", r, g, b)
    local color = results[key]
    if color == nil then
        color = {red = r, green = g, blue = b}
        results[key] = color
    end
    return color
end

-- object attribute
do 
    local mem = {}
    setmetatable(mem, {__mode = "k"})
    function factory (o)
        local res = mem[o]
        if not res then 
            res = (function () return o end)
            mem[o] = res
        end
        return res
    end
end

-- 析构函数
o = {x = "hi"}
setmetatable(o, {__gc = function (o) print(o.x) end})
o = nil
collectgarbage()    --> hi


o = {x = "hi"}
mt = {__gc = true}  --> 必须标记为需要析构处理，否则不会触发析构函数操作
setmetatable(o, mt)
mt.__gc = function (o) print(o.x) end
o = nil
collectgarbage()

-- 同一周期中析构多个对象时， 析构会逆序的调用这个对象的析构器
mt = {__gc = function (o) print(o[1]) end}
list = nil
for i = 1, 3 do
    list = setmetatable({i, link = list}, mt)
end
list = nil
collectgarbage()    --> 3 2 1

-- 复苏 resurrection 当一个析构器被调用时，它的参数正是需要被析构的对象，会临时变为活跃
A = {x = "this is A"}
B = {f = A}
setmetatable(B, {__gc = function (o) print(o.f.x) end})
A, B = nil
collectgarbage()

-- 每次gc后运行一个function
do 
    local mt = {__gc = function(o)
        print("new cycle")
        setmetatable({}, getmetatable(o))
    end}

    setmetatable({}, mt)
end

collectgarbage()
collectgarbage()
collectgarbage()
