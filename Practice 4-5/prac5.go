package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(task16("3BC9", 16, 2))
	fmt.Println(task17(2, 3, 1.125))
	fmt.Println(task18([]int{-5, 1, -10, 3, 2, -18}))
	fmt.Println(task19([]int{1, 5, 8, 19}, []int{0, 7, 9, 18, 21}))
	fmt.Println(task20("sub", "This is sub test string with substring"))
	fmt.Println(task21("9.111+2.2"))
	fmt.Println(task22("In girum imus, nocte et. consumimur igni"))
	fmt.Println(task23(1, 5, 3, 7, 2, 6))
	fmt.Println(task24("In girum imus, nocte et. consumimur igni"))
	fmt.Println(task25(2024))
	fmt.Println(task26(89))
	fmt.Println(task27(51, 100))
	fmt.Println(task28(300, 380))
	fmt.Println(task29("Hello world"))
	fmt.Println(task30(645, 381))
}

/*
Перевод чисел из одной системы счисления в другую
Напишите программу, которая принимает на вход число в произвольной системе счисления (от 2 до 36)
и переводит его в другую систему счисления.
*/
func task16(value string, inRadix int, outRadix int) string {
	if inRadix < 2 || inRadix > 36 || outRadix < 2 || outRadix > 36 {
		panic("Wrong radix!")
	}
	inter, err := strconv.ParseInt(value, inRadix, 64)
	if err != nil {
		panic(err)
	}
	return strconv.FormatInt(inter, outRadix)
}

/*
Решение квадратного уравнения
Реализуйте функцию для нахождения корней квадратного уравнения вида \( ax^2 + bx + c = 0 \).
Учтите случай комплексных корней.
*/
func task17(a float64, b float64, c float64) string {
	var D = math.Pow(b, 2) - 4*a*c
	var x1 = ""
	var x2 = ""
	if D > 0 {
		fmt.Println(D)
		x1 = strconv.FormatFloat((-b+math.Sqrt(D))/(2*a), 'f', -1, 64)
		x2 = strconv.FormatFloat((-b-math.Sqrt(D))/(2*a), 'f', -1, 64)
	} else if D == 0 {
		x1 = strconv.FormatFloat(-b/(2*a), 'f', -1, 64)
		x2 = x1
	} else {
		x1 = strconv.FormatFloat(-b/(2*a), 'f', -1, 64) + "+" + strconv.FormatFloat(
			math.Sqrt(-D)/(2*a), 'f', -1, 64) + "i"
		x2 = strconv.FormatFloat(-b/(2*a), 'f', -1, 64) + "-" + strconv.FormatFloat(
			math.Sqrt(-D)/(2*a), 'f', -1, 64) + "i"
	}
	return "x1 = " + x1 + " x2 = " + x2
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

/*
Сортировка чисел по модулю
Дан массив чисел. Напишите функцию, которая отсортирует его элементы по возрастанию их абсолютных значений.
*/
func task18(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		var minIdx = i
		for j := i; j < len(arr); j++ {
			if abs(arr[j]) < abs(arr[minIdx]) {
				minIdx = j
			}
		}
		if abs(arr[minIdx]) < abs(arr[i]) {
			var tmp = arr[i]
			arr[i] = arr[minIdx]
			arr[minIdx] = tmp
		}
	}
	return arr
}

/*
Слияние двух отсортированных массивов
Напишите функцию, которая объединяет два отсортированных массива в один, сохраняя их упорядоченность.
*/
func task19(arr1 []int, arr2 []int) []int {
	var out = make([]int, len(arr1)+len(arr2))
	var idx1 = 0
	var idx2 = 0
	for i := 0; i < len(arr1)+len(arr2); i++ {
		if (len(arr1)-1 < idx1) || ((len(arr2)-1 > idx2) && arr2[idx2] < arr1[idx1]) {
			out[i] = arr2[idx2]
			idx2++
		} else {
			out[i] = arr1[idx1]
			idx1++
		}
	}
	return out
}

/*
Нахождение подстроки в строке без использования встроенных функций
Реализуйте функцию, которая находит первую позицию вхождения одной строки в другую.
Если подстрока не найдена, вернуть -1.
*/
func task20(sub string, orig string) int {
	var idx = -1
	for i := 0; i < len(orig); i++ {
		if sub[0] != orig[i] {
			continue
		}
		if i+len(sub) > len(orig) {
			return -1
		}
		for j := 0; j < len(sub); j++ {
			if orig[i+j] != sub[j] {
				break
			} else if orig[i+j] == sub[j] && j == len(sub)-1 {
				idx = i
			}
		}
		if idx != -1 {
			return idx
		}
	}
	return idx
}

/*
Калькулятор с расширенными операциями
Напишите программу, которая выполняет различные математические операции (+, -, *, /, ^, %), заданные пользователем.
Реализуйте обработку ошибок, связанных с делением на ноль или недопустимой операцией.
*/
func task21(inp string) float64 {
	var num = regexp.MustCompile(`\d+\.\d+`)
	var name = regexp.MustCompile(`[+\-*/^%]`)
	var operator = name.FindAllString(inp, -1)
	if len(operator) != 1 {
		fmt.Println("Unknown operation!")
		return math.NaN()
	}
	var numbers = num.FindAllString(inp, -1)
	if len(numbers) != 2 {
		fmt.Println("Wrong input! Check your numbers!")
		return math.NaN()
	}
	var parsedNums = make([]float64, 2)
	for i := 0; i < len(numbers); i++ {
		parsedNums[i], _ = strconv.ParseFloat(numbers[i], 64)
	}
	var out = 0.0
	switch operator[0] {
	case "+":
		out = parsedNums[0] + parsedNums[1]
	case "-":
		out = parsedNums[0] - parsedNums[1]
	case "*":
		out = parsedNums[0] * parsedNums[1]
	case "/":
		out = parsedNums[0] / parsedNums[1]
	case "^":
		out = math.Pow(parsedNums[0], parsedNums[1])
	case "%":
		out = math.Mod(parsedNums[0], parsedNums[1])
	}
	return out
}

/*
Проверка палиндрома
Реализуйте функцию, которая проверяет, является ли строка палиндромом (игнорируя пробелы, знаки препинания и регистр)
*/
func task22(inp string) bool {
	var chars = regexp.MustCompile(`[^A-Za-z]`)
	inp = strings.ToLower(inp)
	var str = chars.ReplaceAllString(inp, "")
	var reverse = ""
	for _, v := range str {
		reverse = string(v) + reverse
	}
	return str == reverse
}

/*
Нахождение пересечения трех отрезков
Даны три отрезка на числовой оси (их начальные и конечные точки).
Нужно определить, существует ли область пересечения всех трех отрезков.
*/
func task23(x1 float64, x2 float64, x3 float64, x4 float64, x5 float64, x6 float64) bool {
	var maxStart = x1
	var minEnd = x2
	if x3 > maxStart {
		maxStart = x2
	}
	if x5 > maxStart {
		maxStart = x5
	}
	if x4 < minEnd {
		minEnd = x4
	}
	if x6 < minEnd {
		minEnd = x6
	}
	return maxStart <= minEnd
}

/*
Выбор самого длинного слова в предложении
Напишите программу, которая находит самое длинное слово в предложении, игнорируя знаки препинания.
*/
func task24(inp string) string {
	var chars = regexp.MustCompile(`[A-Za-z]+`)
	inp = strings.ToLower(inp)
	var str = chars.FindAllString(inp, -1)
	var maxIdx = 0
	for i := 0; i < len(str); i++ {
		if len(str[i]) > len(str[maxIdx]) {
			maxIdx = i
		}
	}
	return str[maxIdx]
}

/*
Проверка высокосного года
Реализуйте функцию, которая проверяет, является ли введенный год високосным по правилам григорианского календаря.
*/
func task25(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}

/*
Числа Фибоначчи до определенного числа
Напишите программу, которая выводит все числа Фибоначчи, не превышающие заданного значения.
*/
func task26(n int) []int {
	if n < 0 {
		panic(errors.New("wrong input"))
	}
	var result = []int{0}
	if n == 0 {
		return result
	}
	result = append(result, 1)
	var idx = 2
	for {
		var next = result[idx-2] + result[idx-1]
		if next > n {
			break
		}
		result = append(result, next)
		idx++
	}
	return result
}

/*
Определение простых чисел в диапазоне
Реализуйте функцию, которая принимает два числа и выводит все простые числа в диапазоне между ними.
*/
func task27(start int, end int) []int {
	var numbers = make([]int, 0)
	var isPrime = false
	for i := start; i < end; i++ {
		isPrime = true
		for j := 2; j < int(math.Ceil(math.Sqrt(float64(i)))); j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			numbers = append(numbers, i)
		}
	}
	return numbers
}

/*
Числа Армстронга в заданном диапазоне
Напишите программу, которая выводит все числа Армстронга в заданном диапазоне.
Число Армстронга – это число, равное сумме своих цифр, возведенных в степень, равную количеству цифр числа.
*/
func task28(start int, end int) []int {
	var numbers = make([]int, 0)
	for i := start; i < end; i++ {
		x := i
		n := 0
		for x != 0 {
			x /= 10
			n++
		}
		sum := 0
		x = i
		for x != 0 {
			digit := x % 10
			sum += int(math.Pow(float64(digit), float64(n)))
			x /= 10
		}
		if sum == i {
			numbers = append(numbers, i)
		}
	}
	return numbers
}

/*
Реверс строки
Напишите программу, которая принимает строку и возвращает ее в обратном порядке,
не используя встроенные функции для реверса строк.
*/
func task29(in string) string {
	var out = ""
	for i := 0; i < len(in); i++ {
		out += string(in[len(in)-1-i])
	}
	return out
}

/*
Нахождение наибольшего общего делителя (НОД)
Реализуйте алгоритм Евклида для нахождения наибольшего общего делителя двух чисел с использованием цикла.
*/
func task30(x int, y int) int {
	if x < y {
		for {
			z := y % x
			if z == 0 {
				break
			}
			y = x
			x = z
		}
		return x
	} else {
		for {
			z := x % y
			if z == 0 {
				break
			}
			x = y
			y = z
		}
		return y
	}
}
