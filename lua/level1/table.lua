-- Table
a = {}
k = "x"
a[k] = 10
a[20] = "great"
print(a["x"])
k = 20
print(a[k])
a["x"] = a["x"] + 1
print(a["x"])

b = a
print(b[k])

b = nil
-- print(b[k])
print(a[k])

-- index of the table
c = {}
for i = 1, 100 do 
    c[i] = i*2
end
print(c[3])

c.x = 10    --> c["x"] = 10
print(c["x"])
print(c.y)  --> nil


-- 整形和字符串不代表相同key
i = 10; j = "10"; k = "+10"
a[i] = "number key"
a[j] = "string key"
a[k] = "another string key"
print(a[i], a[j], a[k])     -->number key      string key      another string key
s = a[tonumber(j)]  --> number key
s = a[tonumber(k)]  --> number key

-- 整形和浮点数不存在上述问题，所以注意bug
a[2] = 2
a[2.0] = 3
print(a[2]) --> 3

-- constructor
days = {"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
print(days[3])

a = {x = 1, y = 3}
print(a.x, a.y)

playline = {color="blue",
            thickness=2,
            {x=0, y=0},
            {x=1, y=10}
        }
print(playline[1].x)    --> 0

-- 方括号显示的指定每一个索引
opnames = {["+"]="add", ["-"]="sub"}
print(opnames["-"])

-- array table list
a = {}
for i = 1, 2 do
    a[i] = io.read()
end

for i = 1, #a do
    print(a[i])
end

-- pairs遍历表，但是顺序不一定
for k, v in pairs(a) do
    print(k, v)
end

-- ipairs遍历列表，顺序一定
l = {10, s, 12, "hello"}
for k, v in ipairs(l) do
    print(k, v)
end

-- for
for k = 1, #l do 
    print(k, l[k])
end

-- 安全访问
-- E = {}
-- zip = (((company or E).director or E).address or E).zipcode


-- standard lib
t = {}
for line in io.lines() do
    table.insert( t, line )
end
print(#t)

-- remove elem form table t, when you removed elem, then table will fill gap space by table.move,
-- because table's underlying layer is impled by C, so the performance overhead is not significant,
-- but if your table's cap is very large, you should not use remove or move casually
table.remove( t )
table.remove( t, 1 )
-- you can impl a stack or queue by table.insert and table.remove

-- insert a new elem to t
table.move(t, 1, #)
t[1] = "newNode"

-- in the computer world, a move is actually a copy of a value from one place to another. So as in
-- the example, we must explicity remove the last element after the move
table.move(t, 2, #a, 1)
a[#a] = nil