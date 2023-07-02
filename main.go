package main

import (
	"fmt"
	hw02unpackstring "github.com/zamatay/learn-otus/hw02_unpack_string"
)

type MsgUserBalanceChanged struct {
	userID  string
	balance string
}

type MsgEventChanged struct {
	eventID string
}

func processMessage(msg interface{}) {
	switch v := msg.(type) {
	case MsgEventChanged:
		fmt.Printf("%s\n", v.eventID)
	case MsgUserBalanceChanged:
		fmt.Printf("%s, %s\b", v.userID, v.balance)
	default:
		fmt.Print("message unknown")

	}
}

/*
user "user-1" balance was changed to "1000"
event "event-1" was changed
unknown message: unknown
*/
func main() {
	//processMessage(MsgUserBalanceChanged{"user-1", "1000"})
	//processMessage(MsgEventChanged{"event-1"})
	//processMessage("unknown")

	hw02unpackstring.Unpack()
}
