package exercise2

import (
	"fmt"

	"github.com/jzipfler/htw-ava/server"
)

const (
	THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD  = 5
	THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED = 8
)

// This node represents a customer for exercise 2. The customer gets informed
// from companies about new products and buys one, if he hears the advertisement
// from a company. He also buys the product if enough friends told him, that
// they buyed the product already.
type CustomerNode struct {
	server.NetworkServer
	customerId           int
	friends              map[int]server.NetworkServer //Also neighbors.
	heardAdvertisements  map[int]int                  //An entry for each company
	heardFriendBuyedItem map[int]int                  //An entry for each company
}

//Creates a new CustomerNode with initialized values.
func NewCustomerNode() CustomerNode {
	return CustomerNode{server.New(), 0, make(map[int]server.NetworkServer), make(map[int]int), make(map[int]int)}
}

//Creates a new CustomerNode with initialized values,
//where the server object is set instead of setting an empty one.
func NewCustomerNodeWithServerObject(serverObject server.NetworkServer) CustomerNode {
	return CustomerNode{serverObject, 0, make(map[int]server.NetworkServer), make(map[int]int), make(map[int]int)}
}

//Returns the customer id
func (customerNode CustomerNode) CustomerId() int {
	return customerNode.customerId
}

//Set the customer id
func (customerNode *CustomerNode) SetCustomerId(customerId int) {
	customerNode.customerId = customerId
}

//Returns the friend map
func (customerNode CustomerNode) Friends() map[int]server.NetworkServer {
	return customerNode.friends
}

//Set the friend map
func (customerNode *CustomerNode) SetFriends(friends map[int]server.NetworkServer) {
	customerNode.friends = friends
}

//Returns the number of heared advertisements of the specific company.
func (customerNode CustomerNode) HeardAdvertisementsByCompanyId(companyId int) int {
	return customerNode.heardAdvertisements[companyId]
}

//Increases the number of heared advertisements of the specific company by one.
func (customerNode *CustomerNode) IncreaseHeardAdvertisementsByCompanyId(companyId int) {
	if _, available := customerNode.heardAdvertisements[companyId]; !available {
		customerNode.heardAdvertisements[companyId] = 0
	}
	customerNode.heardAdvertisements[companyId]++
}

//Returns the number of friends that buyed a product of the specified company.
func (customerNode CustomerNode) HeardFriendBuyedItemByCompanyId(companyId int) int {
	return customerNode.heardFriendBuyedItem[companyId]
}

//Increases the number of friends that buyed a product of the specified company by one.
func (customerNode *CustomerNode) IncreaseHeardFriendBuyedItemByCompanyId(companyId int) {
	if _, available := customerNode.heardFriendBuyedItem[companyId]; !available {
		customerNode.heardFriendBuyedItem[companyId] = 0
	}
	customerNode.heardFriendBuyedItem[companyId]++
}

//Checks if the customer would buy a product of the specifie company.
//It will check the values of the advertisements and the recommendations of the friends.
func (customerNode CustomerNode) WouldTheCustomerBuyProductFromCompanyWithId(companyId int) bool {
	if THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD <= 0 || THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED <= 0 {
		return true
	}
	if numberOfAdvertisements := customerNode.HeardAdvertisementsByCompanyId(companyId); numberOfAdvertisements >= THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD {
		return true
	}
	if numberOfFriends := customerNode.HeardFriendBuyedItemByCompanyId(companyId); numberOfFriends >= THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED {
		return true
	}
	return false
}

//The string representation
func (customerNode CustomerNode) String() string {
	return fmt.Sprintf("CustomerID: %d, Server-Settings: %v", customerNode.CustomerId(), customerNode.NetworkServer.String())
}
