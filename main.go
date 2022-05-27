package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./Assets.db")
	checkErr(err)

	// Define statements
	stmtInsert, err := db.Prepare("INSERT INTO Assets(AssetName, UnitPrice, Count) values(?,?,?)")
	checkErr(err)

	stmtUpdatePrice, err := db.Prepare("update Assets set UnitPrice=UnitPrice*? where AssetName=?")
	checkErr(err)

	stmtUpdateCount, err := db.Prepare("update Assets set count=count+? where AssetName=?")
	checkErr(err)

	stmtDelete, err := db.Prepare("delete from assets where AssetName=?")
	checkErr(err)

	// clear table
	_, err = db.Exec("delete from Assets")
	checkErr(err)
	fmt.Println("Delete all assets")

	// SEQ 1
	fmt.Println("SEQ1:")

	_, err = stmtInsert.Exec("Peugeot2008", "650000000", "1")
	checkErr(err)
	fmt.Println("Added Peugeot2008")

	_, err = stmtInsert.Exec("USD", "20000", "17500")
	checkErr(err)
	fmt.Println("Added USD")

	_, err = stmtInsert.Exec("Shopping", "800000000", "1")
	checkErr(err)
	fmt.Println("Added Shopping")

	fmt.Println("")

	// SEQ2
	fmt.Println("SEQ2:")

	_, err = stmtUpdatePrice.Exec(1.45, "USD")
	checkErr(err)
	fmt.Println("The dollar price rose 45%")

	_, err = stmtUpdatePrice.Exec(2, "Peugeot2008")
	checkErr(err)
	fmt.Println("Price of Peugeot2008 doubled")

	_, err = stmtUpdatePrice.Exec(.25, "Shopping")
	checkErr(err)
	fmt.Println("The price of the Shopping became a quarter")

	fmt.Println("")

	// SEQ3
	fmt.Println("SEQ3:")

	_, totalPrice := getAssetTotalPrice("Shopping")
	_, err = stmtDelete.Exec("Shopping")
	checkErr(err)
	fmt.Println("Shopping sold.")
	_, err = stmtInsert.Exec("CashIRT", 1, totalPrice)
	checkErr(err)
	fmt.Println("Added CashIRT: ", totalPrice)

	_, totalPrice = getAssetTotalPrice("USD")
	_, err = stmtDelete.Exec("USD")
	checkErr(err)
	fmt.Println("USD sold.")
	_, err = stmtUpdateCount.Exec(totalPrice, "CashIRT")
	checkErr(err)
	fmt.Println("Added CashIRT: ", totalPrice)

	fmt.Println("")

	// SEQ4
	fmt.Println("SEQ4:")

	count := 5000
	price := 40000
	_, err = stmtInsert.Exec("EOS", price, count)
	checkErr(err)
	_, err = stmtUpdateCount.Exec(-1*count*price, "CashIRT")
	checkErr(err)
	fmt.Println("Buy 5000 EOS, Spend IRT: ", count*price)

	count = 7
	price = 60000000
	_, err = stmtInsert.Exec("ETH", price, count)
	checkErr(err)
	_, err = stmtUpdateCount.Exec(-1*count*price, "CashIRT")
	checkErr(err)
	fmt.Println("Buy 7 ETH, Spend IRT: ", count*price)

	fmt.Println("")

	// Result
	fmt.Println("Finally, Assets are as follows:")
	rows, err := db.Query("SELECT * FROM assets")
	checkErr(err)
	var n string
	var p int
	var c int

	for rows.Next() {
		err = rows.Scan(&n, &p, &c)
		checkErr(err)
		fmt.Print("asset: ", n, "\t")
		fmt.Print("price: ", p, "\t")
		fmt.Print("count: ", c, "\t")
		fmt.Println()
	}

	rows.Close()
	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getAssetTotalPrice(name string) (bool, int) {
	db, err := sql.Open("sqlite3", "./Assets.db")
	checkErr(err)

	// query
	rows, err := db.Query("SELECT * FROM Assets WHERE AssetName=?", name)
	checkErr(err)
	defer rows.Close()
	var assetName string
	var unitPrice int
	var count int

	for rows.Next() {
		err = rows.Scan(&assetName, &unitPrice, &count)
		return true, unitPrice * count
	}

	return false, 0

}
