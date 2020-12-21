-- find, sub
s = "hello lua"
i, j = string.find(s, "hello")
print(i, j)
print(string.sub(s, i, j))


-- match
date = "Today is 14/3/1990"
d = string.match(date, "%d+/%d+/%d+")
print(d)

-- gsub
res = "hello lua"
res = string.gsub(res, string.sub(res, string.find(res, "hello")), "fuck")
print(res)

-- gmatch
s = "some thing"
words = {}
for w in string.gmatch(s, "%a+") do
    words[#words + 1] = w
end
print(#words)

-- pattern
--[[
    .   任意字符
    %a  字母
    %c  控制字符
    %d  数字
    %g  除空格外的可打印字符
    %l  小写字母
    %p  标点符号
    %s  空白符号
    %u  大写字母
    %w  字母或者数字
    %x  十六进制数字    
    %A  任意非字母
    +   重复一次或多次
    *   重复零次或多次
    -   重复零次或多次(最小匹配)
    ？  可选（出现零次或多次）
]]
print((string.gsub("hello, up-down!", "%A", ".")))                  --> hello..up.down.
print((string.gsub("one, and two; and three", "%a+", "word")))      --> word, word word; word word
print(string.match("the number 1234 is even", "%d+"))               --> 1234
test = "int x; /* x */ int y; /* y */"
print(string.gsub(test, "/%*.*%*/", ""))
print(string.gsub(test, "/%*.-%*/", ""))
if string.find(test, "^%d") then
    print("number")
end

--[[
    %b()
    %b{}
    %b<>
    %b[]
]]
s = "a (hello world) line"
print(string.gsub(s, "%b()", ""))   --> a line

-- capture      ()
pair = "name = Anna"
key, value = string.match(pair, "(%a+)%s*=%s*(%a+)")
print(key, value)       --> name Anna

date = "Today is 12/2/2002"
d, m, y = string.match(date, "(%d+)/(%d+)/(%d+)")
print(d, m, y)      --> 12      2       2002

s = [[then he said:  "it's all right"]]
q, quotedPart = string.match(s, "([\"'])(.-)%1")
print(quotedPart)
print(q)

function expand( s )
    return (string.gsub(s, "$(%w+)", _G))
end

name = "lua"; status = "great"
print(expand("$name is $status, isnt it?"))

-- URL
function unescape( s )
    s = string.gsub(s, "+", " ")
    s = string.gsub(s, "%%(%x%x)", function (h)
        return string.char(tonumber(h, 16))
    end)
    return s
end
print(unescape("a%2Bb+%3D+c"))  --> a+b = c

cgi = {}
function decode (s)
    for name, vlaue in string.gmatch(s, "([^&=]+)=([^&=]+)") do 
        name = unescape(name)
        value = unescape(value)
        cgi[name] = value
    end
end

function escape( s )
    s = string.gsub(s, "[&=+%%%c]", function (c)
        return string.format("%%%02X", string.byte(c))
    end)
    s = string.gsub(s, " ", "+")
    return s
end

function encode(t)
    local b = {}
    for k,v in pairs(t) do
        b[#b + 1] = (escape(k) .. "=" .. escape(v))
    end
    return table.concat(b, "&")
end

t = {name="al", query="a+b = c", q="yes or no"}
print(encode(t))

-- pratice 出现频率最高的单词
--[[
    task:
    1. 读取文本并计算出每一个单词的出现次数
    2. 按照出现次数降序对单词列表进行排序
    3. 输出有序列表中的前n个元素
]]

local counter = {}
for line in io.lines() do 
    for word in string.gmatch(line, "%w+") do
        counter[word] = (counter[word] or 0) + 1
    end
end

local words = {}
for w in pairs(counter) do 
    words[#word + 1] = w
end

table.sort(words, function (w1, w2)
    return counter[w1] > counter[w2] or counter[w1] == counter[w2] and w1 < w2
end)

local n = math.min(tonumber(arg[1]) or math.huge, #words)

for i = 1, n do
    io.write(words[i], "\t", counter[words[i]], "\n")
end