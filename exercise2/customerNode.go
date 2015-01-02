package exercise2

import (
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

func NewCustomerNode() CustomerNode {
	return CustomerNode{server.New(), 0, make(map[int]server.NetworkServer), make(map[int]int), make(map[int]int)}
}

func NewCustomerNodeWithServerObject(serverObject server.NetworkServer) CustomerNode {
	return CustomerNode{serverObject, 0, make(map[int]server.NetworkServer), make(map[int]int), make(map[int]int)}
}

func (customerNode CustomerNode) CustomerId() int {
	return customerNode.customerId
}

func (customerNode *CustomerNode) SetCustomerId(customerId int) {
	customerNode.customerId = customerId
}

func (customerNode CustomerNode) Friends() map[int]server.NetworkServer {
	return customerNode.friends
}

func (customerNode CustomerNode) HeardAdvertisementsByCompanyId(companyId int) int {
	return customerNode.heardAdvertisements[companyId]
}

func (customerNode *CustomerNode) IncreaseHeardAdvertisementsByCompanyId(companyId int) {
	if _, available := customerNode.heardAdvertisements[companyId]; !available {
		customerNode.heardAdvertisements[companyId] = 0
	}
	customerNode.heardAdvertisements[companyId]++
}

func (customerNode *CustomerNode) IncreaseHeardFriendBuyedItemByCompanyId(companyId int) {
	if _, available := customerNode.heardFriendBuyedItem[companyId]; !available {
		customerNode.heardFriendBuyedItem[companyId] = 0
	}
	customerNode.heardFriendBuyedItem[companyId]++
}

func (customerNode CustomerNode) WouldTheCustomerBuyProductFromCompanyWithId(companyId int) bool {
	if THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD <= 0 || THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED <= 0 {
		return true
	}
	if numberOfAdvertisements := customerNode.heardAdvertisements[companyId]; numberOfAdvertisements >= THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD {
		return true
	}
	if numberOfFriends := customerNode.heardFriendBuyedItem[companyId]; numberOfFriends >= THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED {
		return true
	}
	return false
}

func (customerNode CustomerNode) String() string {
	return customerNode.NetworkServer.String()
}
