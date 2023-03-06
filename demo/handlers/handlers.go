package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"

	"demo/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `json:"name"`
	Pin  string `json:"Pin" binding:"required"`
}
type deposite struct {
	Name        string `json:"name"`
	Account_num int    `json:"account_num" unique:"true"`
	Pin         string `json:"Pin" binding:"required"`
	Balance     int    `json:"balance"`
}

type transfer struct {
	From   int    `json:"from"`
	To     int    `json:"to"`
	Pin    string `json:"pin"`
	Amount int    `json:"amount"`
}

type state struct {
	Name         string   `json:"name"`
	Account_num  int      `json:"account_num"`
	Pin          string   `json:"Pin"`
	Balance      int      `json:"balance"`
	Transactions []string `json:"transactions"`
}

type reset struct {
	Account_num int    `json:"account_num"`
	Pin         string `json:"Pin"`
	New_pin     string `json:"new_pin"`
}

func hashPassword(password string) string {
	// Hash the password using SHA256
	hash := sha256.Sum256([]byte(password))

	// Encode the hash as a hex string and return it
	return hex.EncodeToString(hash[:])
}

func Create_acc(c *gin.Context) {
	var new_person Person

	if err := c.BindJSON(&new_person); err != nil {
		return
	}
	new_person.Pin = hashPassword(new_person.Pin)

	acc := rand.Intn(1000)

	per := bson.M{"name": new_person.Name, "account_num": acc, "pin": new_person.Pin, "balance": 0}

	collection := db.Connectdb()

	// Insert a document
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "account_num", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(context.Background(), index); err != nil {
		log.Fatal(err)
	}

	res, err := collection.InsertOne(context.Background(), per)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("accout created with num : %d", acc)})

	fmt.Println(res)

}

func Deposite_money(c *gin.Context) {

	var User deposite
	var user1 deposite
	if err := c.BindJSON(&User); err != nil {
		return
	}
	acc := User.Account_num
	pin1 := hashPassword(User.Pin)
	amount := User.Balance
	fmt.Println("reson", User.Account_num)

	collection := db.Connectdb()

	res := collection.FindOne(context.Background(), bson.M{"account_num": acc}).Decode(&user1)
	fmt.Println("reson1", user1.Pin)
	fmt.Println("deposited", res)

	if user1.Pin == pin1 {
		update := bson.M{"balance": (amount + user1.Balance)}
		update1 := bson.M{"transactions": fmt.Sprintf("Deposited %d by you", amount)}
		res, _ := collection.UpdateOne(context.Background(), bson.M{"account_num": acc}, bson.M{"$set": update})
		res1, _ := collection.UpdateOne(context.Background(), bson.M{"account_num": acc}, bson.M{"$addToSet": update1})
		fmt.Println("deposited", res)
		fmt.Println("deposited", res1)
		c.JSON(200, gin.H{"message": fmt.Sprintf("deposited %d in account  \n  total balance %d", amount, amount+user1.Balance)})
	}

}

func Transfer_money(c *gin.Context) {
	var User transfer
	var user1 deposite
	var user2 deposite

	if err := c.BindJSON(&User); err != nil {
		return
	}

	from := User.From
	to := User.To
	amount := User.Amount
	pin1 := hashPassword(User.Pin)

	collection := db.Connectdb()

	collection.FindOne(context.Background(), bson.M{"account_num": from}).Decode(&user1)
	collection.FindOne(context.Background(), bson.M{"account_num": to}).Decode(&user2)

	if user1.Pin == pin1 {
		update := bson.M{"balance": (user1.Balance - amount)}
		update1 := bson.M{"transactions": fmt.Sprintf("credited %d to %s ", amount, user2.Name)}
		update3 := bson.M{"balance": (user2.Balance + amount)}
		update4 := bson.M{"transactions": fmt.Sprintf("Deposited %d by %s", amount, user1.Name)}

		collection.UpdateOne(context.Background(), bson.M{"account_num": user1.Account_num}, bson.M{"$set": update})
		collection.UpdateOne(context.Background(), bson.M{"account_num": user1.Account_num}, bson.M{"$push": update1})
		collection.UpdateOne(context.Background(), bson.M{"account_num": user2.Account_num}, bson.M{"$set": update3})
		collection.UpdateOne(context.Background(), bson.M{"account_num": user2.Account_num}, bson.M{"$push": update4})

		c.JSON(200, gin.H{"message": "transfom completed"})
	}

}

func Statement(c *gin.Context) {
	var temp state
	var user1 state
	if err := c.BindJSON(&temp); err != nil {
		return
	}
	acc := temp.Account_num
	pin1 := hashPassword(temp.Pin)

	collection := db.Connectdb()

	collection.FindOne(context.Background(), bson.M{"account_num": acc}).Decode(&user1)

	if user1.Pin == pin1 {
		fmt.Println("deposited", user1.Transactions)
		c.JSON(200, gin.H{"Transactions": user1.Transactions})
	}

}

func Withdrawal(c *gin.Context) {

	var User deposite
	var user1 deposite
	if err := c.BindJSON(&User); err != nil {
		return
	}
	acc := User.Account_num
	pin1 := hashPassword(User.Pin)
	amount := User.Balance
	fmt.Println("reson", User.Account_num)

	collection := db.Connectdb()

	res := collection.FindOne(context.Background(), bson.M{"account_num": acc}).Decode(&user1)
	fmt.Println("reson1", user1.Pin)
	fmt.Println("deposited", res)

	if user1.Balance < amount {
		c.JSON(200, gin.H{"error": "insufficient funds "})
		return
	}

	if user1.Pin == pin1 {
		update := bson.M{"balance": (user1.Balance - amount)}
		update1 := bson.M{"transactions": fmt.Sprintf("Withdrawal %d by you", amount)}
		res, _ := collection.UpdateOne(context.Background(), bson.M{"account_num": acc}, bson.M{"$set": update})
		res1, _ := collection.UpdateOne(context.Background(), bson.M{"account_num": acc}, bson.M{"$addToSet": update1})
		fmt.Println("deposited", res)
		fmt.Println("deposited", res1)
		c.JSON(200, gin.H{"message": fmt.Sprintf("withdrawed %d from account  \n   total balance %d", amount, user1.Balance-amount)})
	}
}

func Resetting(c *gin.Context) {
	var temp reset
	var user1 deposite

	if err := c.BindJSON(&temp); err != nil {
		return
	}

	acc := temp.Account_num
	pin1 := hashPassword(temp.Pin)
	pin2 := temp.New_pin

	collection := db.Connectdb()

	collection.FindOne(context.Background(), bson.M{"account_num": acc}).Decode(&user1)

	if user1.Pin == pin1 {
		update := bson.M{"pin": hashPassword(pin2)}
		collection.UpdateOne(context.Background(), bson.M{"account_num": acc}, bson.M{"$set": update})
		c.JSON(200, gin.H{"message": "pin updated"})
	}

}
