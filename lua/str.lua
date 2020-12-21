a = "one string"
b = string.gsub( a, "one", "another" )
print(a)
print(b)

-- len of the string
print(#a)

-- string connect
print(a .. " " .. b)

print("\u{3b1} \u{3b2}")

-- [[]]
page = [[
    <html>
    <head>
        <title>hello</title>
    </head>
    </html>
]]
io.write(page)

-- type conversion
print(type(11 .. 22))   --> string

print("10" + 1)   --> number

print(tonumber("100101", 2))
print(tostring(10) == "10")

-- standard lib
-- print(string.rep( "a", 2^20 ))  --> cash down(LOL
print(string.byte( "abc", 1, 2 ))
print(string.format( "x = %x", 200)) 

-- utf8
print(utf8.len("孔子"))

print(utf8.char(114, 233, 115, 117))
s = "我是一个abc"
print(utf8.char(utf8.codepoint(s, utf8.offset(s, 5))))

for i, c in utf8.codes(s) do
    print(i, c)
end
