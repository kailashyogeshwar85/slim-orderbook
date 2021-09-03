package engine

import (
	"log"
)

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order Order) []Trade {
	if order.Side.String() == "buy" {
		return book.processLimitBuyOrder(order)
	}
	return book.processLimitSellOrder(order)
}

func (book *OrderBook) processLimitBuyOrder(order Order) []Trade {
	log.Println("Processing LIMIT BUY ORDER ")
	// create a trade object
	trades := make([]Trade, 0, 1)

	n := len(book.Asks)

	if n == 0 {
		book.addBuyOrder(order)
		return trades
	}
	// check if we have atleast one matching order
	// first compare with last sell price
	// loop only when the new order price is less than last sell price
	if n != 0 || book.Asks[n-1].Price.LessThanOrEqual(order.Price) {
		// traverse all orders that match
		for i := n - 1; i >= 0; i-- {
			sellOrder := book.Asks[i]

			// last sell price is 10 limit = 9
			// return as the highest sell order is gt limit price
			if sellOrder.Price.GreaterThan(order.Price) {
				break
			}
			// fill the entire order when sellOrder has higher quantity
			if sellOrder.Quantity.GreaterThanOrEqual(order.Quantity) {
				trades = append(
					trades,
					Trade{
						TakerOrderID: order.ID,     // TakerOrderID
						MakerOrderID: sellOrder.ID, // Maker OrderID
						Quantity:     order.Quantity.BigInt().Uint64(),
						Price:        sellOrder.Price.BigInt().Uint64(),
						Timestamp:    order.Timestamp,
					},
				)
			}
			// partially fill the order and continue
			if sellOrder.Quantity.LessThan(order.Quantity) {
				trades = append(
					trades,
					Trade{
						TakerOrderID: order.ID,
						MakerOrderID: sellOrder.ID,
						Quantity:     sellOrder.Quantity.BigInt().Uint64(),
						Price:        sellOrder.Price.BigInt().Uint64(),
						Timestamp:    order.Timestamp,
					},
				)
				order.Quantity = order.Quantity.Sub(sellOrder.Quantity)
				// remove the sell Order as all quantities are filled by bid
				book.removeSellOrder(i)
				continue
			}
		}
	}
	// finally add the order with remaining qty to book
	book.addBuyOrder(order)
	return trades
}

//
func (book *OrderBook) processLimitSellOrder(order Order) []Trade {
	log.Println("Processing LIMIT SELL ORDER ")

	trades := make([]Trade, 0, 1)
	n := len(book.Bids)

	if n == 0 {
		book.addSellOrder(order)
		return trades
	}

	// Proceed only if the sell Price is Greather than user highest buy Price
	if n != 0 || book.Bids[n-1].Price.GreaterThanOrEqual(order.Price) {
		// travers all bids that match
		for i := n - 1; i >= 0; i-- {
			buyOrder := book.Bids[i]

			if buyOrder.Price.LessThan(order.Price) {
				break // exit
			}

			// fill the entire order of buy order is gte
			if buyOrder.Price.GreaterThanOrEqual(order.Price) {
				trades = append(
					trades,
					Trade{
						TakerOrderID: order.ID,
						MakerOrderID: buyOrder.ID,
						Quantity:     order.Quantity.BigInt().Uint64(),
						Price:        buyOrder.Price.BigInt().Uint64(),
						Timestamp:    order.Timestamp,
					},
				)
				buyOrder.Quantity = buyOrder.Quantity.Sub(order.Quantity)

				// if buyOrder.Quantity = 0
				if buyOrder.Quantity.IsZero() {
					book.removeBuyOrder(i)
				}
				return trades
			}

			// fill a partial order and continue
			if buyOrder.Quantity.LessThan(order.Quantity) {
				trades = append(
					trades,
					Trade{
						TakerOrderID: order.ID,
						MakerOrderID: buyOrder.ID,
						Quantity:     buyOrder.Quantity.BigInt().Uint64(),
						Price:        buyOrder.Price.BigInt().Uint64(),
						Timestamp:    order.Timestamp,
					},
				)
				order.Quantity = order.Quantity.Sub(buyOrder.Quantity)
				book.removeBuyOrder(i)
				continue
			}
		}
	}
	book.addSellOrder(order)
	return trades
}
