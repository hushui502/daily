function quote( s )
    local n = -1
    for w in string.gmatch(s, "]=*") do
        n = math.max(n, #w - 1)
    end

    local eq = string.rep("=", n + 1)

    return string.format("[%s[\n%s]%s]", eq, s, eq)
end

-- 序列化
function serialize(o)
    local t = type(o)
    if t == "number" or t == "string" or t == "boolean" or t == "nil" then
        io.write(string.format( "%q", o))
    elseif t == "table" then
        io.write("{\n")
        for k, v in pairs(o) do
            io.write("  ", k, " = ")
            serialize(v)
            io.write(",\n")
        end
        io.write("}\n")
    else
        error("cannot serialize a" .. type(o))
    end
end
serialize{a=12, b='lua', key="ss"}


-- 保存带有循环的表
function basicSerialize( o )
    return string.format( "%q", o )
end

function save( name, value, saved )
    saved = saved or {}
    io.write(name, " = ")
    if type(value) == "number" or type(value) == "string" then 
        io.write(basicSerialize(value), "\n")
    elseif type(value) == "table" then
        if saved[value] then
            io.write(saved[value], "\n")
        else
            saved[value] = name
            io.write("{}\n")
            for k, v in pairs(value) do
                k = basicSerialize(k)
                local fname = string.format( "%s[%s]", name, k )
                save(fname, v, saved)
            end
        end
    else 
        error("cannot save a " .. type(value))
    end
end