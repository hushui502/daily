print(os.time() // 365.242 + 1970)
print(os.date())
--[[
    %a 星期几的简写
    %A 星期几的全写
    %b 月份
    %B 
    %c 日期和时间
    %d 一个月中的第几天
    %H 24小时中的小时数
    %I 12小时
    %j 一年中的第几天
    %m 月份
    %M 分钟
    %p am pm
    %S 秒数
    %w 星期
    %W 一年中的第几周
    %x 日期
    %X 时间
    %y 两位数的年份
    %Y 完整年份 
    %z 时区
    %% 百分号
]]

t = os.date("*t")
print(os.date("%Y/%m/%d", os.time(t)))
t.day = t.day + 40
print(os.date("%Y/%m/%d", os.time(t)))
print(t.day, t.month)

-- benchmark ==> os.clock 精度更高
local x = os.clock()
local s = 0
for i = 1, 100000 do s = s + i end
print(string.format("elapsed time: %.2f\n", os.clock() - x))
