package engine

// making field lowercase makes it private property
type OrderBook struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}

// APIs
// addBuyOrder(order)
// addSellOrder(order)
// removeBuyOrder(orderId)

// Add the new Order to end of orderbook in bids
func (book *OrderBook) addBuyOrder(order Order) {
	n := len(book.Bids)

	if n == 0 {
		book.Bids = append(book.Bids, order)
	} else {
		var i int

		for i := n - 1; i >= 0; i-- {
			buyOrder := book.Bids[i]

			// check the price of existing order
			// convert decimal to Signed int
			if buyOrder.Price.LessThan(order.Price) {
				break
			}
		}

		// if new order price is less than the last order price
		if i == n-1 {
			// append the new order at end
			book.Bids = append(book.Bids, order)
		} else {
			// add order to the index before the order which
			copy(book.Bids[i+1:], book.Bids[i:])
			book.Bids[i] = order
		}
	}
}

func (book *OrderBook) addSellOrder(order Order) {
	n := len(book.Asks)

	if n == 0 {
		book.Asks = append(book.Asks, order)
	} else {
		var i int
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.Asks[i]

			if sellOrder.Price.LessThan(order.Price) {
				break
			}
		}
		if i == n-1 {
			// append the new order at end
			book.Asks = append(book.Asks, order)
		} else {
			// add order to the index before the order which
			copy(book.Asks[i+1:], book.Asks[i:])
			book.Asks[i] = order
		}
	}
}

func (book *OrderBook) removeBuyOrder(index int) {
	book.Bids = append(book.Bids[:index], book.Bids[index+1:]...)
}

func (book *OrderBook) removeSellOrder(index int) {
	book.Asks = append(book.Asks[:index], book.Asks[index+1:]...)
}
