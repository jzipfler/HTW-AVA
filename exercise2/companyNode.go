package exercise2

import (
	"errors"
	"fmt"
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
	return CompanyNode{server.New(), 0, "", 0, make(map[int]server.NetworkServer)}
}

func NewCompanyNodeWithServerObject(serverObject server.NetworkServer) CompanyNode {
	return CompanyNode{serverObject, 0, "", 0, make(map[int]server.NetworkServer)}
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

func (companyNode CompanyNode) CompanyId() int {
	return companyNode.companyId
}

func (companyNode *CompanyNode) SetCompanyId(companyId int) {
	companyNode.companyId = companyId
}

func (companyNode CompanyNode) Product() string {
	return companyNode.product
}

func (companyNode *CompanyNode) SetProduct(product string) {
	companyNode.product = product
}

func (companyNode CompanyNode) AdvertisingBudget() int {
	return companyNode.advertisingBudget
}

func (companyNode *CompanyNode) SetAdvertisingBudget(advertisingBudget int) {
	companyNode.advertisingBudget = advertisingBudget
}

func (companyNode CompanyNode) RegularCustomers() map[int]server.NetworkServer {
	return companyNode.regularCustomers
}

func (companyNode *CompanyNode) SetRegularCustomers(regularCustomers map[int]server.NetworkServer) {
	companyNode.regularCustomers = regularCustomers
}

func (companyNode CompanyNode) IsRegularCustomer(customerId int) bool {
	_, available := companyNode.regularCustomers[customerId]
	return available
}

func (companyNode *CompanyNode) AddRegularCustomer(customerId int, serverObjectInformation server.NetworkServer, override bool) error {
	if _, available := companyNode.regularCustomers[customerId]; available && !override {
		return errors.New(fmt.Sprintf("The customer with the ID %d is already available and override is set to false.", customerId))
	}
	companyNode.regularCustomers[customerId] = serverObjectInformation
	return nil
}

func (companyNode CompanyNode) String() string {
	return companyNode.NetworkServer.String()
}
