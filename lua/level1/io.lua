-- use io.write("", "", "") instead of io.write("" .. "" .. "")

--                    SIMPLE IO MODEL
-- case 1   format output
io.write("sin(3) = ", math.sin( 3 ), "\n")
io.write(string.format( "sin(3) = %.4f\n", math.sin( 3 )))


-- case 2   replace input
t = io.read("a")
t = string.gsub( t, "bad", "good")
io.write(t)

-- case 3   repleace input by func
t = io.read("all")
t = string.gsub( t, "[\128-\255=]", function(c) 
        return string.format( "=%02X", string.byte(c))
    end)
io.write(t)

-- case 4   read line
for count = 1, math.huge do 
    local line = io.read("L")
    if line == nil then break end
    io.write(string.format( "%6d    ", count), line)
end

-- case5    read line 
local count = 0
for line in io.lines() do 
    count = count + 1
    io.write(string.format( "%6d  ", count), line, "\n")
end

-- case 6   sort reading
local lines = {}
for line in io.lines() do
    lines[#line + 1] = line
end

table.sort(lines)

for _, l in ipairs(lines) do
    io.write(l, "\n")
end

-- case 7   block read
while true do
    local block = io.read(2^13)
    if not block then break end
    io.write(block)
end

-- case 8 max elem in the every line
while true do
    local n1, n2, n3 = io.read("n", "n", "n")
    if not n1 then break end
    print(math.max( n1, n2, n3 ))
end

             COMPLETE IO MODEL
print(io.open("func.lua", "r"))

local f = assert(io.open("func.lua", "r"))
local t = f:read("a")
f:close()

io.stderr:write("err")

local temp = io.input()
io.input("func.lua")
io.input():close()
io.input(temp)

for block in io.input():lines(2^13) do
    io.write(block)
end

function fsize( file )
    local current = file:seek() -- 保存当前的位置
    local size  = file:size()   -- 保存当前文件的大小
    file:seek("set", current)   -- 恢复当前位置
    return size
end

-- rename / remove file
os.remove()
os.rename()
os.exit(-1)
print(os.getenv("HOME"))

-- create dir by os.execute == shell command
function createDir( dirname )
    os.execute("mkdir " .. dirname)
end


-- popen
local f = io.popen("ls /LUA", "r")
local dir = {}
for entry in file:lines() do 
    dir[#dir + 1] = entry
end

-- send mail
local subject = "some news"
local address = "hufan@gmail.com"

local cmd = string.format( "mail -s '%s' '%s'", subject, address)
local f = io.popen(cmd, "w")
f:write([[
    Nothing important to say.
]])
f:close()