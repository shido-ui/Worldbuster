package simulation

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

type ProductionRecipe struct {
	ID             string
	Name           string
	BusinessType   BusinessType
	MinLevel       int
	InputItem      string
	InputQuantity  int
	OutputItem     string
	OutputQuantity int
	LaborCost      int64
}

type ProductionResult struct {
	BusinessID     string
	RecipeID       string
	Produced       bool
	InputItem      string
	InputQuantity  int
	OutputItem     string
	OutputQuantity int
	LaborCost      int64
	Reason         string
}

func CanProduce(b BusinessState, r ProductionRecipe, inputQuantity int64) bool {
	if !b.Active || b.Type != r.BusinessType || b.Level < r.MinLevel || inputQuantity < int64(r.InputQuantity) {
		return false
	}
	return b.Cash >= r.LaborCost
}

var ErrBusinessIDEntropyFailure = errors.New("business ID entropy failure")

func NewBusinessID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", ErrBusinessIDEntropyFailure
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmtUUID(b), nil
}

func fmtUUID(b [16]byte) string {
	const h = "0123456789abcdef"
	out := make([]byte, 36)
	j := 0
	for i := 0; i < 16; i++ {
		if i == 4 || i == 6 || i == 8 || i == 10 {
			out[j] = '-'
			j++
		}
		out[j] = h[b[i]>>4]
		out[j+1] = h[b[i]&15]
		j += 2
	}
	return string(out)
}

func DeterministicBusinessID(ownerID string, businessType BusinessType) string {
	sum := sha256.Sum256([]byte("worldbuster:business:" + ownerID + ":" + string(businessType)))
	var b [16]byte
	copy(b[:], sum[:16])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmtUUID(b)
}
