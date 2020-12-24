local m = require("math")
print(m.sin(3))

-- package search

function search( modename, path )
    modename = string.gsub(modename, "%.", "/")
    local msg = {}
    for c in string.gmatch(path, "[^;]+") do
        local fname = string.gsub(c, "?", modename)
        local f = io.open(fname)
        if f then 
            f:close()
            return fname
        else
            msg[#msg + 1] = string.format( "\n\tno file %s", fname)
        end 
    end
    return nil, table.concat(msg)
end