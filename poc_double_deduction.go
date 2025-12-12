package main

import (
	"fmt"
	"math/big"
)

func CalcFeeWithRatio(feeRatio uint8, totalFee *big.Int) (feePayerFee, senderFee *big.Int) {
	feePayerFee = new(big.Int).Div(new(big.Int).Mul(totalFee, big.NewInt(int64(feeRatio))), big.NewInt(100))
	senderFee = new(big.Int).Sub(totalFee, feePayerFee)
	return
}

func toKAIA(f *big.Int) string {
	return new(big.Int).Div(f, big.NewInt(1e18)).String()
}

func main() {
	balance := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	value := new(big.Int).Mul(big.NewInt(50), big.NewInt(1e18))
	totalFee := new(big.Int).Mul(big.NewInt(1e9), new(big.Int).SetUint64(21000))
	feePayerFee, senderFee := CalcFeeWithRatio(60, totalFee)

	fmt.Printf("[poc_double_deduction] Balance: %s | Value: %s | Fee: %s+%s\n",
		toKAIA(balance), toKAIA(value), toKAIA(feePayerFee), toKAIA(senderFee))

	check1 := balance.Cmp(new(big.Int).Add(value, senderFee)) >= 0
	check2 := balance.Cmp(feePayerFee) >= 0
	fmt.Printf("Validation: sender_ok=%v, feePayer_ok=%v -> PASS\n", check1, check2)

	bal := new(big.Int).Set(balance)
	bal.Sub(bal, value)
	bal.Sub(bal, feePayerFee)
	bal.Sub(bal, senderFee)

	fmt.Printf("Final balance: %s KAIA\n", toKAIA(bal))
}
