print(string.format( "%x", 0xff & 0xabcd))
print(string.format( "%x", 0xff | 0xabcd))
print(string.format( "%x", 0xaaaa ~ -1))
print(string.format( "%x", ~0))
print(string.format( "%x", 0xff << 12))

function udiv( n, d )
    if d < 0 then
        if math.ult(n, d) then return 0
        else return 1
        end
    end
    local q = ((n >> 1) // d) << 1
    local r = n - q * d
    if not math.ult(r, d) then q = q + 1 end
    return q
end

-- file
local inp = assert(io.open(arg[1], "rb"))
local out = assert(io.open(arg[2]), "wb")

local data = inp:read("a")
data = string.gsub(data, "\r\n", "\n")
out:write(data)

assert(out:close())

local f = assert(io.open(arg[1], "rb"))
local data = f:read("a")
local validchars = "[%g%s]"
local pattern = "(" .. string.rep(validchars, 6) .. "+)\0"
for w in string.gmatch(data, pattern) do 
    print(w)
end

-- dump file
--[[
    2D2D2075736520696F2E777269746528 -- use io.write(
    22222C2022222C2022222920696E7374 "", "", "") inst
    656164206F6620696F2E777269746528 ead of io.write(
    2222202E2E202222202E2E202222290D "" .. "" .. "").
    0A0D0A2D2D2020202020202020202020 ...--
    20202020202020202053494D504C4520          SIMPLE
    494F204D4F44454C0D0A2D2D20636173 IO MODEL..-- cas
    652031202020666F726D6174206F7574 e 1   format out
]]
local f = assert(io.open(arg[1], "rb"))
local blocksize = 16
for bytes in f:lines(blocksize) do 
    for i = 1, #bytes do 
        local b = string.unpack("B", bytes, i)
        io.write(string.format("%02X", b))
    end
    io.write(string.rep("   ", blocksize - #bytes))
    bytes = string.gsub(bytes, "%c", ".")
    io.write(" ", bytes, "\n")
end
