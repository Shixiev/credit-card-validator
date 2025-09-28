package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// --- Структуры и функции ---

// Bank — хранит информацию о банке и его диапазоне BIN
type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

// loadBankData загружает данные банков из файла
func loadBankData(path string) ([]Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", path, err)
	}
	defer f.Close()

	var banks []Bank
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 3)
		if len(parts) != 3 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		from, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		to, _ := strconv.Atoi(strings.TrimSpace(parts[2]))

		banks = append(banks, Bank{
			Name:    name,
			BinFrom: from,
			BinTo:   to,
		})
	}

	return banks, scanner.Err()
}

// getUserInput считывает ввод от пользователя
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите номер карты (Enter для выхода): ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// validateInput проверяет базовый формат (13–19 цифр, только числа)
func validateInput(cardNumber string) bool {
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}
	for _, r := range cardNumber {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// validateLuhn проверяет номер карты по алгоритму Луна
func validateLuhn(cardNumber string) bool {
	sum := 0
	double := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}

// extractBIN возвращает первые 6 цифр как int
func extractBIN(cardNumber string) int {
	if len(cardNumber) < 6 {
		return 0
	}
	bin, _ := strconv.Atoi(cardNumber[:6])
	return bin
}

// identifyBank ищет банк по BIN
func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Неизвестный эмитент"
}

// --- Основная программа ---
func main() {
	fmt.Println("Добро пожаловать в программу валидации карт! 🚀")

	// Загружаем банки
	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки банков:", err)
		return
	}
	fmt.Println("✅ Данные банков успешно загружены.\n")

	// Основной цикл
	for {
		cardNumber := getUserInput()

		// Проверка на выход
		if cardNumber == "" {
			fmt.Println("Программа завершена. Спасибо за использование! 👋")
			break
		}

		// Проверка формата
		if !validateInput(cardNumber) {
			fmt.Println("❌ Ошибка: номер карты должен содержать только цифры (13–19 символов).")
			continue
		}

		// Проверка алгоритмом Луна
		if !validateLuhn(cardNumber) {
			fmt.Println("❌ Ошибка: номер карты не прошёл проверку по алгоритму Луна.")
			continue
		}

		// Определение банка
		bin := extractBIN(cardNumber)
		bankName := identifyBank(bin, banks)

		// Вывод результата
		fmt.Println("✅ Номер карты валиден.")
		fmt.Println("🏦 Банк:", bankName)
		fmt.Println()
	}
}