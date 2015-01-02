package exercise2_test

import (
	"testing"

	"github.com/jzipfler/htw-ava/exercise2"
)

func TestInitBudgetForCompanies(t *testing.T) {
	threshold := 11
	companyNode := exercise2.NewCompanyNode()
	for i := 1; i < 10; i++ {
		companyNode.InitAdvertisingBudgetWithThreashold(threshold)
		if companyNode.AdvertisingBudget() >= threshold {
			t.Error("The budget was greater or equal the threasold.")
		}
	}
}
