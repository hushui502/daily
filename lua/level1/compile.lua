-- -- load()
-- f = load("i = i + 1")
-- i = 0
-- f(); print(i)
-- f(); print(i)

-- i = 32
-- local i = 0
-- f = load("i = i + 1; print(i)")
-- g = function () i = i + 1; print(i) end
-- f()     --> 33  load 是在全局环境中编译代码段
-- g()     --> 1

-- -- error
-- n = io.read("n")
-- if not n then error("invalid input") end

-- -- assert 便捷处理error
-- assert(io.read("*n"), "invalid input")

-- -- handle       error throw an exception, pcall can capture catch
-- local ok, msg = pcall(function () 
--     -- some code
--     if unexpected_condition then error() end
--     -- some code
--     print(a[i]) -- a maybe not a table
-- end)

-- if ok then
--     -- 
-- else
--     --
-- end


local status, err = pcall(function () error({code=121}) end)
print(err.code)         --> 121

function foo( str )
    if type(str) ~= "string" then
        error("string expected")
    end
    -- regular code
end

foo({x=1})
