package exercise2

import (
	"math/rand"
	"time"

	"github.com/jzipfler/htw-ava/server"
)

const (
	ADVERTISING_BUGDET_MAXIMUM = 10
)

// This node represents a compyna for the exercise 2. The company tries to sell
// a product and starts publicity campaigns in form of adverdisements. These
// advertisements are sent to all neighbors which are customers.
// Each company got only a specifiy budget to do adverdisements, where earch
// campain reduces the budget.
// If a customer buys a product, he will be inserted to a list of regular
// customers and receives now all advertisements. That means, that the customer
// is now a neighbor of the company.
type CompanyNode struct {
	server.NetworkServer
	companyId         int
	product           string
	advertisingBudget int
	regularCustomers  map[int]server.NetworkServer
}

func NewCompanyNode() CompanyNode {
	return CompanyNode{server.New(), 0, "", 0, nil}
}

// Sets the budget to a random value with the following range:
// [0, ADVERTISING_BUDGET_MAXIMUM)
func (companyNode *CompanyNode) InitAdvertisingBudget() {
	companyNode.InitAdvertisingBudgetWithThreashold(ADVERTISING_BUGDET_MAXIMUM)
}

// Sets the budget to a random value with the following range:
// [0, threshold)
func (companyNode *CompanyNode) InitAdvertisingBudgetWithThreashold(threshold int) {
	rand.Seed(time.Now().UnixNano())
	companyNode.advertisingBudget = rand.Intn(threshold)
}
