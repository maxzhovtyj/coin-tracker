package telegram

import (
	"fmt"
	db "github.com/maxzhovtyj/coin-tracker/pkg/db/sqlc"
	"testing"
)

func Test_allWalletsSuccess(t *testing.T) {
	all := []db.CryptoWallet{
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
		{Name: "5"},
	}

	fmt.Println(getKeyboardFromWallets(all))
}
