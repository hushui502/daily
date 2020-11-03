function fact (n) 
    if n == 0 then
        return 1
    else
        return n * fact(n-1)
    end
end

local counter = {}

for line in io.lines() do 
    for word in string.gmatch(line, "%w+") do 
        counter[word] = (counter[word] or 0) + 1
    end
end

