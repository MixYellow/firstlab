use std::io::{self, Write};

// 1. Пузырьковая сортировка 
fn bubble_sort(arr: &mut Vec<i32>) {
    let n = arr.len();
    for i in 0..n {
        let mut swapped = false;
        for j in 0..n - i - 1 {
            if arr[j] > arr[j + 1] {
                arr.swap(j, j + 1);
                swapped = true;
            }
        }
        if !swapped {
            break;
        }
    }
}

// 2. Сортировка вставками
fn insertion_sort(arr: &mut Vec<i32>) {
    for i in 1..arr.len() {
        let key = arr[i];
        let mut j = i;
        while j > 0 && arr[j - 1] > key {
            arr[j] = arr[j - 1];
            j -= 1;
        }
        arr[j] = key;
    }
}

// 3. Быстрая сортировка — рекурсивная
fn quick_sort(arr: &mut Vec<i32>, low: isize, high: isize) {
    if low < high {
        let pi = partition(arr, low, high);
        quick_sort(arr, low, pi - 1);
        quick_sort(arr, pi + 1, high);
    }
}

fn partition(arr: &mut Vec<i32>, low: isize, high: isize) -> isize {
    let pivot = arr[high as usize];
    let mut i = low - 1;
    for j in low..high {
        if arr[j as usize] <= pivot {
            i += 1;
            arr.swap(i as usize, j as usize);
        }
    }
    arr.swap((i + 1) as usize, high as usize);
    i + 1
}

// Обёртка для удобного вызова quicksort
fn quick_sort_wrapper(arr: &mut Vec<i32>) {
    if !arr.is_empty() {
        let n = (arr.len() - 1) as isize;
        quick_sort(arr, 0, n);
    }
}

// Читает строку из stdin
fn read_line() -> String {
    let mut input = String::new();
    io::stdin().read_line(&mut input).expect("Ошибка чтения строки");
    input.trim().to_string()
}

// Парсит числа из строки, разделённые пробелами
fn parse_numbers(input: &str) -> Vec<i32> {
    input
        .split_whitespace()
        .filter_map(|s| s.parse::<i32>().ok())
        .collect()
}

// Выводит массив в одну строку
fn print_array(arr: &[i32]) {
    for (i, num) in arr.iter().enumerate() {
        if i > 0 {
            print!(" ");
        }
        print!("{}", num);
    }
    println!();
}


fn main() {
    println!("=== Сортировка массива ===");

    // Ввод массива
    println!("Введите целые числа через пробел (например: 5 2 8 1 9):");
    let input = read_line();
    let mut numbers = parse_numbers(&input);
    if numbers.is_empty() {
        println!("Не введено ни одного числа.");
        return;
    }

    println!("Исходный массив:");
    print_array(&numbers);

    // Выбор алгоритма
    println!("\nВыберите алгоритм сортировки:");
    println!("1. Пузырьковая сортировка (Bubble Sort)");
    println!("2. Сортировка вставками (Insertion Sort)");
    println!("3. Быстрая сортировка (Quick Sort)");
    print!("Ваш выбор (1-3): ");
    io::stdout().flush().unwrap();

    let choice = read_line();

    // Копия массива для каждого алгоритма, чтобы не перемешивать (можно и по месту)
    let mut sorted = numbers.clone();

    match choice.as_str() {
        "1" => {
            bubble_sort(&mut sorted);
            println!("\nРезультат (пузырьковая сортировка):");
        }
        "2" => {
            insertion_sort(&mut sorted);
            println!("\nРезультат (сортировка вставками):");
        }
        "3" => {
            quick_sort_wrapper(&mut sorted);
            println!("\nРезультат (быстрая сортировка):");
        }
        _ => {
            println!("Неверный выбор. Завершение.");
            return;
        }
    }
    print_array(&sorted);
}



