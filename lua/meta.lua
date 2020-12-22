-- 新建表时不带元表
t = {}
print(getmetatable(t))  --> nil

t1 = {}
setmetatable(t, t1)
print(getmetatable(t) == t1)    --> true

-- table: 0000000000e99c60
-- table: 0000000000e99c60
-- 字符串标准库为所有字符串都设置了同一个元表，其他类型在默认情况下没有元表
print(getmetatable("hello"))
print(getmetatable("rool"))

-- 集合
local Set = {}
local mt = {}
function Set.new( l )
    local set = {}
    setmetatable(set, mt)
    for _, v in ipairs(l) do set[v] = true end
    return set
end


function Set.union( a, b )
    local res = Set.new{}
    if getmetatable(a) ~= mt or getmetatable(b) ~= mt then
        error("attempt to 'add' a set with a non-set value", 2)
    end
    for k in pairs(a) do res[k] = true end
    for k in pairs(b) do res[k] = true end
    return res
end


function Set.intersection( a, b )
    local res = Set.new{}
    for k in pairs(a) do
        res[k] = b[k]
    end
    return res
end

function Set.tostring( set )
    local l = {}
    for e in pairs(set) do
        l[#l + 1] = tostring(e)
    end
    return "{" .. table.concat(l, ", ") .. "}"
end


s1 = Set.new{1, 2, 3, 4}
s2 = Set.new{2, 3, 5, 6}
print(getmetatable(s1))     --> table: 00000000001e9760
print(getmetatable(s2))     --> table: 00000000001e9760

mt.__add = Set.union
s3 = s1 + s2
print(Set.tostring(s3))     --> {1, 2, 3, 4, 5, 6}

mt.__mul = Set.intersection
print(Set.tostring((s1 + s2) * s3))

-- s4 = s2 + 3      -- error

mt.__le = function (a, b)           -- 子集
    for k in pairs(a) do
        if not b[k] then return false end
    end
    return true
end

mt.lt = function (a, b)             -- 真子集
    return a <= b and not (b <= a)
end

mt.eq = function (a, b)
    return a <= b and b <= a
end

s1 = Set.new{2, 4}
s2 = Set.new{2, 4, 10}
print(s1 <= s2)
-- ...

-- lib meta method
print({})       -- table: 0000000000cf2d80

-- index
local mt = {__index = function (t) return t.__ end}
function setDefault(t, d)
    t.__ = d
    setmetatable(t, mt)
end

tab = {x=10, y=10}
print(tab.x, tab.z)     --> 10 nil
setDefault(tab, 0)
print(tab.x, tab.z)     --> 10 0

-- 防止命名冲突
local key = {}  -- 唯一的键
local mt = {__index = function (t) return t[key] end}
function setDefault2( t, d )
    t[key] = d
    setmetatable(t, mt)
end

-- proxy
function track( t )
    local proxy = {}
    local mt = {
        __index = function (_, k) 
            print("*access to element " .. tostring(k))
            return t[k]
        end,

        __newindex= function (_, k, v)
            print("*update of element " .. tostring(k) .. " to " .. tostring(v))
            t[k] = v
        end,

        __pairs = function () 
            return function (_, k)
                local nextkey, nextvalue = next(t, k)
                if nextkey ~= nil then 
                    print("*traversing element " .. tostring(nextkey))
                end
                return nextkey, nextvalue
            end
        end,

        __len = function () return #t end
    }

    setmetatable(proxy, mt)
    
    return proxy
end

t = {}
t = track(t)
t[2] = "hello"
print(t[2])

t = track({12, 23, 44})
print(#t)
for k, v in pairs(t) do print(k, v) end

-- read only
function readOnly( t )
    local proxy = {}
    local mt = {
        __index = t,
        __newindex = function (t, k, v)
            error("attempt to update a read-only table", 2)
        end        
    }
    setmetatable(proxy, mt)
    return proxy
end

days = readOnly{"sun", "mon", "tue", "wed", "thu", "fri", "sta"}
print(days[1])
days[1] = "ss"