package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func main() {
	fmt.Println(task1(55))
	fmt.Println(task2(25))
	fmt.Println(task3([]int{1, 2, 3, 4}))
	fmt.Println(task4([]string{"Hello", "world"}))
	fmt.Println(task5(1, 1, 4, 5))
	fmt.Println(task6(4))
	fmt.Println(task7(2020))
	fmt.Println(task8(4, 9, 7))
	task9(25)
	fmt.Println(task10(15))
	fmt.Println(task11(5))
	fmt.Println(task12(7))
	fmt.Println(task13([]int{1, 2, 3, 4, 5}))
	fmt.Println(task14(20))
	fmt.Println(task15([]int{1, 2, 3, 4, 5}))
}

/*
Сумма цифр числа:
Напишите программу, которая принимает целое число и вычисляет сумму его цифр.
*/
func task1(x int) int {
	if x == 0 {
		return 0
	}
	return x%10 + task1(x/10)
}

/*
Преобразование температуры:
Напишите программу, которая преобразует температуру из градусов Цельсия в Фаренгейты и обратно.
*/
func task2(x float64) float64 {
	return x*1.8 + 32
}

/*
Удвоение каждого элемента массива:
Напишите программу, которая принимает массив чисел и возвращает новый массив, где каждое число удвоено.
*/
func task3(x []int) []int {
	x[0] = x[0] * 2
	if len(x) == 1 {
		return x
	}
	x = append(x[:1], task3(x[1:])...)
	return x
}

/*
Объединение строк:
Напишите программу, которая принимает несколько строк и объединяет их в одну строку через пробел.
*/
func task4(s []string) string {
	return strings.Join(s, " ")
}

/*
Расчет расстояния между двумя точками:
Напишите программу, которая вычисляет расстояние между двумя точками в 2D пространстве.
*/
func task5(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

/*
Проверка на четность/нечетность:
Напишите программу, которая проверяет, является ли введенное число четным или нечетным.
*/
func task6(x int) string {
	if x%2 == 0 {
		return "Четное"
	} else {
		return "Нечетное"
	}
}

/*
Проверка высокосного года:
Напишите программу, которая проверяет, является ли введенный год високосным
*/
func task7(year int) string {
	if year%400 == 0 {
		return "Високосный"
	}
	if year%100 == 0 {
		return "Невисокосный"
	}
	if year%4 == 0 {
		return "Високосный"
	}
	return "Невисокосный"
}

/*
Определение наибольшего из трех чисел:
Напишите программу, которая принимает три числа и выводит наибольшее из них.
*/
func task8(x, y, z int) int {
	var maxNum = x
	if y > maxNum {
		maxNum = y
	}
	if z > maxNum {
		maxNum = z
	}
	return maxNum
}

/*
Категория возраста:
Напишите программу, которая принимает возраст человека и выводит, к какой возрастной группе он относится
(ребенок (0 - 12), подросток (13 - 21), взрослый (22 - 61), пожилой ( > 61)).
*/
func task9(age int) {
	if age <= 12 {
		fmt.Println("Ребенок")
	} else if age <= 21 {
		fmt.Println("Подросток")
	} else if age <= 61 {
		fmt.Println("Взрослый")
	} else {
		fmt.Println("Пожилой")
	}
}

/*
Проверка делимости на 3 и 5:
Напишите программу, которая проверяет, делится ли число одновременно на 3 и 5.
*/
func task10(inp int) string {
	if inp%3 == 0 && inp%5 == 0 {
		return "Делится"
	} else {
		return "Не делится"
	}
}

/*
Факториал числа:
Напишите программу, которая вычисляет факториал числа.
*/
func task11(inp int) int {
	var acc = 1
	for inp > 1 {
		acc *= inp
		inp--
	}
	return acc
}

/*
Числа Фибоначчи:
Напишите программу, которая выводит первые `n` чисел Фибоначчи.
*/
func task12(inp int) []int {
	if inp <= 0 {
		panic(errors.New("wrong input"))
	}
	var result = make([]int, inp)
	result[0] = 0
	if inp >= 2 {
		result[1] = 1
	}
	for i := 2; i < len(result); i++ {
		result[i] = result[i-2] + result[i-1]
	}
	return result
}

/*
Реверс массива:
Напишите программу, которая переворачивает массив чисел.
*/
func task13(inp []int) []int {
	var result = make([]int, len(inp))
	for i := 0; i < len(inp); i++ {
		result[i] = inp[len(inp)-1-i]
	}
	return result
}

/*
Поиск простых чисел:
Напишите программу, которая выводит все простые числа до заданного числа.
*/
func task14(inp int) []int {
	var initial_numbers = make([]int, inp-2)
	var cnt = 0
	for i := 0; i < len(initial_numbers); i++ {
		initial_numbers[i] = i + 2
	}
	for i := 0; i < len(initial_numbers); i++ {
		if initial_numbers[i] == 0 {
			continue
		}
		var divider = initial_numbers[i]
		cnt++
		for j := i + 1; j < len(initial_numbers); j++ {
			if initial_numbers[j]%divider == 0 {
				initial_numbers[j] = 0
			}
		}
	}
	var out = make([]int, cnt)
	var idx = 0
	for i := 0; i < len(initial_numbers); i++ {
		if initial_numbers[i] != 0 {
			out[idx] = initial_numbers[i]
			idx++
		}
	}
	return out
}

/*
Сумма чисел в массиве:
Напишите программу, которая вычисляет сумму всех чисел в массиве.
*/
func task15(inp []int) int {
	var sum = 0
	for i := 0; i < len(inp); i++ {
		sum += inp[i]
	}
	return sum
}
