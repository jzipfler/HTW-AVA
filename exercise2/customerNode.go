package exercise2

import (
	"github.com/jzipfler/htw-ava/server"
)

const (
	THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD  = 5
	THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED = 8
)

type CustomerNode struct {
	server.NetworkServer
	customerId          int
	friends             map[int]CustomerNode //Also neighbors.
	heardAdvertiements  map[int]int          //An entry for each company
	heardFriedBuyedItem map[int]int          //An entry for each friend
}
