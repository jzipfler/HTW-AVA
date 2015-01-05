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

// Creates / initializes a new CompanyNode
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

// Return the company id.
func (companyNode CompanyNode) CompanyId() int {
	return companyNode.companyId
}

// Set the value of the company id
func (companyNode *CompanyNode) SetCompanyId(companyId int) {
	companyNode.companyId = companyId
}

// Return the product name
func (companyNode CompanyNode) Product() string {
	return companyNode.product
}

// Set the name of the product
func (companyNode *CompanyNode) SetProduct(product string) {
	companyNode.product = product
}

// Return the budget for the advertisements
func (companyNode CompanyNode) AdvertisingBudget() int {
	return companyNode.advertisingBudget
}

// Set the budget for the advertiesements
func (companyNode *CompanyNode) SetAdvertisingBudget(advertisingBudget int) {
	companyNode.advertisingBudget = advertisingBudget
}

// Return the map of regular customers
func (companyNode CompanyNode) RegularCustomers() map[int]server.NetworkServer {
	return companyNode.regularCustomers
}

// Set the map of regular customers
func (companyNode *CompanyNode) SetRegularCustomers(regularCustomers map[int]server.NetworkServer) {
	companyNode.regularCustomers = regularCustomers
}

// Check if the customer with the given id is a regular customer.
func (companyNode CompanyNode) IsRegularCustomer(customerId int) bool {
	_, available := companyNode.regularCustomers[customerId]
	return available
}

// Add a customer to the regular customers. Returns a error if the customer is already available.
func (companyNode *CompanyNode) AddRegularCustomer(customerId int, serverObjectInformation server.NetworkServer, override bool) error {
	if _, available := companyNode.regularCustomers[customerId]; available && !override {
		return errors.New(fmt.Sprintf("The customer with the ID %d is already available and override is set to false.", customerId))
	}
	companyNode.regularCustomers[customerId] = serverObjectInformation
	return nil
}

// The string representation of this type.
func (companyNode CompanyNode) String() string {
	return fmt.Sprintf("CompanyID: %d, Product: %s, Budget: %d, Server-Settings: %v", companyNode.CompanyId(), companyNode.Product(), companyNode.AdvertisingBudget(), companyNode.NetworkServer.String())
}
