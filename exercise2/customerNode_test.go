package exercise2_test

import (
	"strconv"
	"testing"

	"github.com/jzipfler/htw-ava/exercise2"
)

func TestBuyDecisionForCustomer(t *testing.T) {
	customerNode := exercise2.NewCustomerNode()
	if exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD > 0 {
		for i := 1; i < exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD; i++ {
			customerNode.IncreaseHeardAdvertisementsByCompanyId(4711)
			if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(4711); shouldBuy {
				t.Error("The customer would buy the product, but the advertisement threshold is " + strconv.Itoa(exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD) + " and only " + strconv.Itoa(i) + " advertisements are heard.")
				break
			}
		}
		customerNode.IncreaseHeardAdvertisementsByCompanyId(4711)
		if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(4711); !shouldBuy {
			t.Error("The customer would not buy the product, but the advertisement threshold is " + strconv.Itoa(exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD) + " and " + strconv.Itoa(customerNode.HeardAdvertisementsByCompanyId(4711)) + " advertisements are heard.")
		}
	} else {
		if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(4711); !shouldBuy {
			t.Error("The threshold for the advertisement is less/equal 0 and the function said it should not buy the item.")
		}
	}

	if exercise2.THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED > 0 {
		for i := 1; i < exercise2.THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED; i++ {
			customerNode.IncreaseHeardFriendBuyedItemByCompanyId(1337)
			if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(1337); shouldBuy {
				t.Error("The customer would buy the product, but the friend info threshold is " + strconv.Itoa(exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD) + " and only " + strconv.Itoa(i) + " advertisements are heard.")
				break
			}
		}
		customerNode.IncreaseHeardFriendBuyedItemByCompanyId(1337)
		if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(1337); !shouldBuy {
			t.Error("The customer would not buy the product, but the friend info threshold is " + strconv.Itoa(exercise2.THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD) + " and " + strconv.Itoa(customerNode.HeardAdvertisementsByCompanyId(1337)) + " advertisements are heard.")
		}
	} else {
		if shouldBuy := customerNode.WouldTheCustomerBuyProductFromCompanyWithId(1337); !shouldBuy {
			t.Error("The threshold for the friend info is less/equal 0 and the function said it should not buy the item.")
		}
	}
}
