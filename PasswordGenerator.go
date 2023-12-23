package main

import (
	"strings"
)

const (
	NUMBERBET = iota
	ALPHABET
	FULLBET
)

var BetLength = 0
var SelectedBet []byte
var NumberBet = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var AlphaBet = [...]byte{'s', 't', 'e', 'h', 'j', 'c', 'f', 'w', 'q', 'o', 'a', 'g', 'p', 'x', 'r', 'v', 'd', 'i', 'u', 'l', 'b', 'z', 'y', 'k', 'm', 'n'}
var FullBet = [...]byte{'c', 'e', 'z', 'j', 'y', 'a', 't', 'p', '9', 'o', 'w', 'x', 'k', '2', '8', '6', 'i', '1', '3', 'v', 'm', '0', 'h', 'u', 'r', 'l', 's', '5', 'b', '7', '4', 'f', 'q', 'n', 'd', 'g'}

func getSelectBetLength() int {
	return BetLength
}

// SelectBet
// @return Set success
func SelectBet(bet int) bool {
	switch bet {
	case NUMBERBET:
		SelectedBet = NumberBet[:]
	case ALPHABET:
		SelectedBet = AlphaBet[:]
	case FULLBET:
		SelectedBet = FullBet[:]
	default:
		return false
	}
	BetLength = len(SelectedBet)
	return true
}

// NextPasswd
// @parameter IPassword 用int8数组作为bet索引序列代替密码
// @return password,isEnd
func NextPasswd(IPassword []int8) (string, bool) {
	var strBuilder strings.Builder
	//先根据当前ipassword生成字符串
	for j := 0; j < len(IPassword); j++ {
		strBuilder.WriteByte(SelectedBet[IPassword[j]])
	}
	//然后刷新生成下次的ipasswd
	i := len(IPassword) - 1
	for {
		IPassword[i]++
		if int(IPassword[i]) == BetLength {
			IPassword[i] = 0
			i--
			if i < 0 {
				return strBuilder.String(), true
			}
			continue
		}
		return strBuilder.String(), false
	}
}

func setIPasswordByNum(IPassword []int8, count int) {
	for j := 0; j < count; j++ {
		i := len(IPassword) - 1
		for {
			IPassword[i]++
			if int(IPassword[i]) == BetLength {
				IPassword[i] = 0
				i--
				if i < 0 {
					break
				}
				continue
			}
			break
		}
	}
	//fmt.Println(count, IPassword)
}
