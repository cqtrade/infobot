package ftx

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/cqtrade/infobot/src/ftx/structs"
)

type NewOrder structs.NewOrder
type NewOrderResponse structs.NewOrderResponse
type OpenOrders structs.OpenOrders
type OrderHistory structs.OrderHistory
type NewTriggerOrder structs.NewTriggerOrder
type ModifyTriggerOrder structs.ModifyTriggerOrder
type NewTriggerOrderResponse structs.NewTriggerOrderResponse
type OpenTriggerOrders structs.OpenTriggerOrders
type TriggerOrderHistory structs.TriggerOrderHistory
type Triggers structs.Triggers

func (client *FtxClient) GetOpenOrders(market string) (OpenOrders, error) {
	var openOrders OpenOrders
	resp, err := client._get("orders?market="+market, []byte(""))
	if err != nil {
		log.Printf("Error GetOpenOrders", err)
		return openOrders, err
	}
	err = _processResponse(resp, &openOrders)
	return openOrders, err
}

func (client *FtxClient) GetOrderHistory(market string, startTime float64, endTime float64, limit int64) (OrderHistory, error) {
	var orderHistory OrderHistory
	requestBody, err := json.Marshal(map[string]interface{}{
		"market":     market,
		"start_time": startTime,
		"end_time":   endTime,
		"limit":      limit,
	})
	if err != nil {
		log.Printf("Error GetOrderHistory", err)
		return orderHistory, err
	}
	resp, err := client._get("orders/history?market="+market, requestBody)
	if err != nil {
		log.Printf("Error GetOrderHistory", err)
		return orderHistory, err
	}
	err = _processResponse(resp, &orderHistory)
	return orderHistory, err
}

func (client *FtxClient) GetOpenTriggerOrdersOriginal(market string, _type string) (OpenTriggerOrders, error) {
	var openTriggerOrders OpenTriggerOrders
	requestBody, err := json.Marshal(map[string]string{"market": market, "type": _type})
	if err != nil {
		log.Printf("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	resp, err := client._get("conditional_orders?market="+market, requestBody)
	if err != nil {
		log.Printf("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	err = _processResponse(resp, &openTriggerOrders)
	return openTriggerOrders, err
}

func (client *FtxClient) GetOpenTriggerOrders(market string) (OpenTriggerOrders, error) {
	var openTriggerOrders OpenTriggerOrders
	requestBody, err := json.Marshal(map[string]string{"market": market})
	if err != nil {
		log.Printf("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	resp, err := client._get("conditional_orders?market="+market, requestBody)
	if err != nil {
		log.Printf("Error GetOpenTriggerOrders", err)
		return openTriggerOrders, err
	}
	err = _processResponse(resp, &openTriggerOrders)
	return openTriggerOrders, err
}

func (client *FtxClient) GetTriggers(orderId string) (Triggers, error) {
	var trigger Triggers
	resp, err := client._get("conditional_orders/"+orderId+"/triggers", []byte(""))
	if err != nil {
		log.Printf("Error GetTriggers", err)
		return trigger, err
	}
	err = _processResponse(resp, &trigger)
	return trigger, err
}

func (client *FtxClient) GetTriggerOrdersHistory(market string, startTime float64, endTime float64, limit int64) (TriggerOrderHistory, error) {
	var triggerOrderHistory TriggerOrderHistory
	requestBody, err := json.Marshal(map[string]interface{}{
		"market":     market,
		"start_time": startTime,
		"end_time":   endTime,
	})
	if err != nil {
		log.Printf("Error GetTriggerOrdersHistory", err)
		return triggerOrderHistory, err
	}
	resp, err := client._get("conditional_orders/history?market="+market, requestBody)
	if err != nil {
		log.Printf("Error GetTriggerOrdersHistory", err)
		return triggerOrderHistory, err
	}
	err = _processResponse(resp, &triggerOrderHistory)
	return triggerOrderHistory, err
}

func (client *FtxClient) PlaceOrder(market string, side string, price float64,
	_type string, size float64, reduceOnly bool, ioc bool, postOnly bool) (NewOrderResponse, error) {
	var newOrderResponse NewOrderResponse
	requestBody, err := json.Marshal(NewOrder{
		Market:     market,
		Side:       side,
		Price:      price,
		Type:       _type,
		Size:       size,
		ReduceOnly: reduceOnly,
		Ioc:        ioc,
		PostOnly:   postOnly})
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		fmt.Println("Error Place Order", err)
		return newOrderResponse, err
	}
	resp, err := client._post("orders", requestBody)
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		return newOrderResponse, err
	}
	err = _processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}

func (client *FtxClient) PlaceMarketOrder(market string, side string,
	_type string, size float64) (NewOrderResponse, error) {
	var newOrderResponse NewOrderResponse
	requestBody, err := json.Marshal(NewOrder{
		Market:     market,
		Side:       side,
		Type:       _type,
		Size:       size,
		ReduceOnly: false,
		Ioc:        false,
		PostOnly:   false})
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		fmt.Println("Error Place Market Order", err)
		return newOrderResponse, err
	}
	resp, err := client._post("orders", requestBody)
	if err != nil {
		log.Printf("Error PlaceOrder", err)
		fmt.Println("Error Place Market Order", err)
		return newOrderResponse, err
	}
	err = _processResponse(resp, &newOrderResponse)
	return newOrderResponse, err
}

func (client *FtxClient) PlaceTriggerOrder(market string, side string, size float64,
	_type string, reduceOnly bool, retryUntilFilled bool, triggerPrice float64,
	orderPrice float64, trailValue float64) (NewTriggerOrderResponse, error) {

	var newTriggerOrderResponse NewTriggerOrderResponse
	var newTriggerOrder NewTriggerOrder

	switch _type {
	case "stop":
		if orderPrice != 0 {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
				OrderPrice:   orderPrice,
			}
		} else {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
			}
		}
	case "trailingStop":
		newTriggerOrder = NewTriggerOrder{
			Market:     market,
			Side:       side,
			Type:       _type,
			Size:       size,
			ReduceOnly: reduceOnly,
			TrailValue: trailValue,
		}
	case "takeProfit":
		if orderPrice != 0 {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
				OrderPrice:   orderPrice,
			}
		} else {
			newTriggerOrder = NewTriggerOrder{
				Market:       market,
				Side:         side,
				TriggerPrice: triggerPrice,
				Type:         _type,
				Size:         size,
				ReduceOnly:   reduceOnly,
			}
		}

	default:
		log.Printf("Trigger type is not valid")
	}
	requestBody, err := json.Marshal(newTriggerOrder)
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	resp, err := client._post("conditional_orders", requestBody)
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	err = _processResponse(resp, &newTriggerOrderResponse)
	return newTriggerOrderResponse, err
}

func (client *FtxClient) CancelOrder(orderId int64) (Response, error) {
	var deleteResponse Response
	id := strconv.FormatInt(orderId, 10)
	resp, err := client._delete("orders/"+id, []byte(""))
	if err != nil {
		log.Printf("Error CancelOrder", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}

func (client *FtxClient) CancelTriggerOrder(orderId int64) (Response, error) {
	var deleteResponse Response
	id := strconv.FormatInt(orderId, 10)
	resp, err := client._delete("conditional_orders/"+id, []byte(""))
	if err != nil {
		log.Printf("Error CancelTriggerOrder", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}

func (client *FtxClient) CancelAllOrders() (Response, error) {
	var deleteResponse Response
	resp, err := client._delete("orders", []byte(""))
	if err != nil {
		log.Printf("Error CancelAllOrders", err)
		return deleteResponse, err
	}
	err = _processResponse(resp, &deleteResponse)
	return deleteResponse, err
}

func (client *FtxClient) ModifyTriggerOrder(orderId int64, size float64, triggerPrice float64) (NewTriggerOrderResponse, error) {

	var newTriggerOrderResponse NewTriggerOrderResponse
	var modifyTriggerOrder ModifyTriggerOrder
	modifyTriggerOrder = ModifyTriggerOrder{
		TriggerPrice: triggerPrice,
		Size:         size,
	}
	id := strconv.FormatInt(orderId, 10)
	requestBody, err := json.Marshal(modifyTriggerOrder)
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	resp, err := client._post("conditional_orders/"+id+"/modify", requestBody)
	fmt.Println("Modify order CODE", resp.Status)
	// resp, err := client._post("conditional_orders/"+id, []byte(""))
	if err != nil {
		fmt.Println("Error PlaceTriggerOrder", err)
		return newTriggerOrderResponse, err
	}
	err = _processResponse(resp, &newTriggerOrderResponse)
	return newTriggerOrderResponse, err
}
