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
