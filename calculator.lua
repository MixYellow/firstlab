-- Определяем приоритеты операций
local priorities = {
    ['+'] = 1,
    ['-'] = 1,
    ['*'] = 2,
    ['/'] = 2,
    ['^'] = 3,
}

-- Ассоциативность: 'L' - левая, 'R' - правая
local associativity = {
    ['+'] = 'L',
    ['-'] = 'L',
    ['*'] = 'L',
    ['/'] = 'L',
    ['^'] = 'R',
}

-- Функция, определяющая, нужно ли заключить подвыражение в скобки
-- parent_op - оператор родительского выражения
-- child_op - оператор дочернего выражения
-- is_right_child - является ли дочерний элемент правым операндом
local function needs_parentheses(parent_op, child_op, is_right_child)
    if child_op == nil then return false end
    local p_prio = priorities[parent_op] or 0
    local c_prio = priorities[child_op] or 0
    if c_prio < p_prio then
        return true
    elseif c_prio == p_prio then
        if associativity[parent_op] == 'L' and is_right_child then
            return true   
        elseif associativity[parent_op] == 'R' and not is_right_child then
            return true   
        end
    end
    return false
end

-- Основная функция преобразования RPN в инфиксную запись и вычисления
function evaluate_rpn(rpn_str)
    local stack = {}   
    local tokens = {}
    for token in rpn_str:gmatch("%S+") do
        table.insert(tokens, token)
    end

    for _, token in ipairs(tokens) do
        
        local num = tonumber(token)
        if num then
            table.insert(stack, { expr = token, value = num })
        else
          
            if #stack < 2 then
                error("Недостаточно операндов для оператора '" .. token .. "'")
            end
            local right = table.remove(stack)
            local left = table.remove(stack)
            local res_value
            local expr

            -- Вычисляем результат и формируем инфиксное выражение
            if token == '+' then
                res_value = left.value + right.value
                expr = left.expr .. " + " .. right.expr
            elseif token == '-' then
                res_value = left.value - right.value
                expr = left.expr .. " - " .. right.expr
            elseif token == '*' then
                res_value = left.value * right.value
                expr = left.expr .. " * " .. right.expr
            elseif token == '/' then
                if right.value == 0 then error("Деление на ноль") end
                res_value = left.value / right.value
                expr = left.expr .. " / " .. right.expr
            elseif token == '^' then
                res_value = left.value ^ right.value
                expr = left.expr .. " ^ " .. right.expr
            else
                error("Неизвестный оператор: " .. token)
            end

            -- Добавляем скобки, если необходимо (на основе приоритетов)
            -- Проверяем левый операнд
            local left_op = nil
            if left.expr:match("^[%d.]+$") == nil then
                
                left_op = left.expr:match("[%+%-%*/%^]")   
            end
            local right_op = nil
            if right.expr:match("^[%d.]+$") == nil then
                right_op = right.expr:match("[%+%-%*/%^]")
            end

            if needs_parentheses(token, left_op, false) then
                expr = "(" .. expr .. ")"
            end
            if needs_parentheses(token, right_op, true) then
                expr = "(" .. expr .. ")"
            end

            table.insert(stack, { expr = expr, value = res_value })
        end
    end

    if #stack ~= 1 then
        error("Некорректное выражение: осталось " .. #stack .. " элементов в стеке")
    end
    return stack[1].expr, stack[1].value
end

-- Основная программа: ввод строки, обработка, вывод
print("Введите выражение в обратной польской записи (RPN):")
local input = io.read()
if not input or input == "" then
    print("Пустая строка.")
    return
end

local success, infix, result = pcall(evaluate_rpn, input)
if success then
    print("Инфиксное выражение:", infix)
    print("Результат:", result)
else
    print("Ошибка:", result)
end
