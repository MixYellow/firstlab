local priorities = { ['+']=1, ['-']=1, ['*']=2, ['/']=2, ['^']=3 }
local associativity = { ['+']='L', ['-']='L', ['*']='L', ['/']='L', ['^']='R' }

local function need_paren(parent_op, child, is_right)
    local child_op = child.op
    if child_op == nil then return false end
    local pp = priorities[parent_op]
    local cp = priorities[child_op]
    if cp < pp then
        return true
    elseif cp == pp then
        if associativity[parent_op] == 'L' and is_right then
            return true
        elseif associativity[parent_op] == 'R' and not is_right then
            return true
        end
    end
    return false
end

function evaluate_rpn(rpn_str)
    local stack = {}   
    for token in rpn_str:gmatch("%S+") do
        local num = tonumber(token)
        if num then
            table.insert(stack, { expr=token, value=num, op=nil })
        else
            if #stack < 2 then error("Недостаточно операндов для '"..token.."'") end
            local right = table.remove(stack)
            local left = table.remove(stack)

            local value
            if token == '+' then
                value = left.value + right.value
            elseif token == '-' then
                value = left.value - right.value
            elseif token == '*' then
                value = left.value * right.value
            elseif token == '/' then
                if right.value == 0 then error("Деление на ноль") end
                value = left.value / right.value
            elseif token == '^' then
                value = left.value ^ right.value
            else
                error("Неизвестный оператор "..token)
            end

            
            local left_str = left.expr
            local right_str = right.expr

            if need_paren(token, left, false) then
                left_str = "(" .. left_str .. ")"
            end
            if need_paren(token, right, true) then
                right_str = "(" .. right_str .. ")"
            end

            local expr = left_str .. " " .. token .. " " .. right_str
            table.insert(stack, { expr=expr, value=value, op=token })
        end
    end
    if #stack ~= 1 then error("Некорректное выражение") end
    return stack[1].expr, stack[1].value
end

print("Введите ОПЗ выражение:")
local input = io.read()
if not input or input == "" then return end

local ok, infix, result = pcall(evaluate_rpn, input)
if ok then
    print("Инфиксная запись:", infix)
    print("Результат:", result)
else
    print("Ошибка:", result)
end
