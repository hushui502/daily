a = 3
b = 2

-- mod algo
if a % b == a - ((a // b) * b) then
    print("ok")
end

local tolerance = 10
function isturnback( angle )
    -- body
    math.randomseed(os.time())
    angle = angle % 360
    return (math.abs(angle - 180) < tolerance)
end

print(math.maxinteger+1 == math.mininteger)
print(math.mininteger-1 == math.maxinteger)
print(math.mininteger == -math.mininteger)


function cond2int( x )
    return math.tointeger( x ) or x
end
-- 函数调用必须在函数定义之后
print(cond2int(22.2))
print(cond2int(22.0))