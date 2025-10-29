package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func ValidateCPF(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	cpf = strings.TrimSpace(cpf)

	if len(cpf) != 11 {
		return false
	}

	if !regexp.MustCompile(`^\d{11}$`).MatchString(cpf) {
		return false
	}

	if cpf == "00000000000" || cpf == "11111111111" ||
		cpf == "22222222222" || cpf == "33333333333" ||
		cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" ||
		cpf == "88888888888" || cpf == "99999999999" {
		return false
	}

	var sum int
	var remainder int

	for i := 1; i <= 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i-1]))
		sum += num * (11 - i)
	}
	remainder = (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}
	digit1, _ := strconv.Atoi(string(cpf[9]))
	if remainder != digit1 {
		return false
	}

	sum = 0
	for i := 1; i <= 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i-1]))
		sum += num * (12 - i)
	}
	remainder = (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}
	digit2, _ := strconv.Atoi(string(cpf[10]))
	if remainder != digit2 {
		return false
	}

	return true
}

func NormalizeCPF(cpf string) string {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	return strings.TrimSpace(cpf)
}

func NormalizeRG(rg string) string {
	rg = strings.ReplaceAll(rg, ".", "")
	rg = strings.ReplaceAll(rg, "-", "")
	return strings.TrimSpace(rg)
}
