package main

import (
	"fmt"
	"math/big"
)

func CalcFeeWithRatio(feeRatio uint8, fee *big.Int) (feePayer, sender *big.Int) {
	feePayer = new(big.Int).Div(new(big.Int).Mul(fee, big.NewInt(int64(feeRatio))), big.NewInt(100))
	sender = new(big.Int).Sub(fee, feePayer)
	return
}

func toKAIA(f *big.Int) string {
	denom := new(big.Float).SetFloat64(1e18)
	bf := new(big.Float).SetInt(f)
	v := new(big.Float).Quo(bf, denom)
	s, _ := v.MarshalText()
	return string(s)
}

func main() {
	balance := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	gas := uint64(21000)
	gasPrice := new(big.Int).SetInt64(10000000000000)
	totalFee := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))
	feePayerFee, senderFee := CalcFeeWithRatio(60, totalFee)
	value := new(big.Int).Sub(new(big.Int).Set(balance), senderFee)

	fmt.Printf("[poc_negative_balance] Balance: %s | Value: %s | TotalFee: %s\n",
		toKAIA(balance), toKAIA(value), toKAIA(totalFee))

	check1 := balance.Cmp(new(big.Int).Add(value, senderFee)) >= 0
	check2 := balance.Cmp(feePayerFee) >= 0
	checkCombined := balance.Cmp(new(big.Int).Add(value, totalFee)) >= 0
	fmt.Printf("Validation: sender_ok=%v, feePayer_ok=%v | combined_ok=%v -> PASS but WRONG\n",
		check1, check2, checkCombined)

	bal := new(big.Int).Set(balance)
	bal.Sub(bal, value)
	bal.Sub(bal, feePayerFee)
	bal.Sub(bal, senderFee)

	if bal.Sign() < 0 {
		fmt.Printf("❌ NEGATIVE BALANCE: %s KAIA (STATE CORRUPTED)\n", toKAIA(bal))
	} else {
		fmt.Printf("Final balance: %s KAIA\n", toKAIA(bal))
	}
}
