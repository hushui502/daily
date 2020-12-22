Account = {
    balance = 0,
    withdraw = function (self, v)
        self.balance = self.balance - v
    end
}

function Account:withdraw(v)
    self.balance = self.balance - v
end

function Account:deposit (v)
    self.balance = self.balance + v
end

function Account:new (o)
    o = o or {}
    self.__index = self
    setmetatable(o, self)
    return o
end

b = Account:new()
print(b.balance)

-- 继承
Account = {balance = 0}
function Account:new( o )
    o = o or {}
    self.__index = self
    setmetatable(o, self)
    return o
end

function Account:deposit( v )
    self.balance = self.balance - v
end

function Account:withdraw( v )
    if v > self.balance then error"insufficient funds" end
    self.balance = self.balance - v
end

SpecialAccount = Account:new()
s = SpecialAccount:new{limit=100.0}
print(s.limit)
s:deposit(11)
print(s.balance)

function SpecialAccount:withdraw( v )
    if v - self.balance >= self:getLimit() then
        error"insufficient funds"
    end
    self.balance = self.balance - v
end

function SpecialAccount:getLimit(  )
    return sele.limit or 0
end

s = SpecialAccount:new()
function s:getLimit(  )
    return self.balance * 0.1
end


-- private
function newAccount( initBalance )
    local self = {
        balance = initBalance,
        LIM = 10000.0
    }
    local extra = function () 
        if self.balance > self.LIM then
            return self.balance * 0.10
        else
            return 0
        end
    end

    local withdraw = function (v)
        self.balance = self.balance - v
    end

    local deposit = function (v)
        self.balance = self.balance + v
    end

    local getBalance = function () return self.balance + extra() end

    return {
        withdraw = withdraw,
        deposit = deposit,
        getBalance = getBalance
    }
end

acc1 = newAccount(100.2)
acc1.withdraw(12.3)
print(acc1.getBalance())


-- single-method object
function newObject( value )
    return function (action, v)
        if action == "get" then return value
        elseif action == "set" then value = v
        else error"invalid action"
        end
    end
end

d = newObject(10)
print(d("get"))
d("set", 22)
print(d("get"))

-- dual representation
local balance = {}
Account = {}

function Account:withdraw( v )
    balance[self] = balance[self] - v
end

function Account:deposit( v )
    balance[sele] = balance[self] + v
end

function Account:balance(  )
    return balance[self]
end

function Account:new( o )
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    balance[o] = 0
    return o
end

a = Account:new{}
a:deposit(100.0)
print(a:balance())