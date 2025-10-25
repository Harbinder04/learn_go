package main

import "fmt"

// 1️⃣ Create a custom type for better type safety
type DeliveryStatus int

// 2️⃣ Define constants using iota
const (
	Pending DeliveryStatus = iota
	Shipped
	Delivered
	Failed
)

// 🔁 WE can understand this by:
/* const (
Pending DeliveryStatus = iota // 0
Shipped DeliveryStatus = iota // 1
Delivered DeliveryStatus = iota // 2
Failed DeliveryStatus = iota // 3
)
*/

// 3️⃣ Add a String() method so it can print meaningful text
func (s DeliveryStatus) String() string {
	return [...]string{"Pending", "Shipped", "Delivered", "Failed"}[s]
}

func (s *DeliveryStatus) MarkDelivered() {
	*s = Delivered
}

func main() {
	var status DeliveryStatus = Shipped // this means shipped = 1, internally

	fmt.Println("Order Status:", status) // prints: Order Status: Shipped
	fmt.Println(Pending, Delivered, Failed)

	status2 := Pending
	fmt.Println(status2)

	status.MarkDelivered() // internally (&status).MarkDelivered()
	fmt.Println(status)

	fmt.Println(status2) // value = ?
}
