package main

import (
	"errors"
	"fmt"

	levelup "github.com/multiverse-os/levelup"
)

type TestObject struct {
	Name        string
	Description string
}

func main() {
	fmt.Println("LevelDB [LEVELUP] Wrapper Example")
	fmt.Println("=============================================================")
	fmt.Println("A level db wrapper providing both standard k/v store, and")
	fmt.Println("document storage. Focused around the concept of having only")
	fmt.Println("writing iteraction with the leveldb intstance, and all the")
	fmt.Println("reads occur against a cache maintained used a buffered channel")
	fmt.Println("that loops and does a mutex locked single write throttled write")
	fmt.Println("at a time.\n")
	fmt.Println("It supports COLLECTIONS, special types, a very flexible codec")
	fmt.Println("so that it can be used for cache, basic embedded database, ")
	fmt.Println("a message queue, and several other things; but while trying very")
	fmt.Println("hard to be atomic, it doesn't support transactions yet, and so")
	fmt.Println("that should be taken into consideration. This is mostly a first")
	fmt.Println("attempt at building a database; using leveldb as a base, and")
	fmt.Println("a lot was learned and we will be using this; it is very likely")
	fmt.Println("we will be taking what we learned and start again.\n")
	fmt.Println("=============================================================")

	leveldb, err := levelup.Open("test/kvstore")
	if err != nil {
		panic(errors.New("[error] failed to open leveldb datastore:" + err.Error()))
	}

	fmt.Println("=============================================================")
	fmt.Println("Creating new collection 'test-collection',...")
	collectionName := "test-collection"

	collection := leveldb.NewCollection(collectionName)

	fmt.Println("collection:", collection)

	//testKey := "test-key"
	//testValue := TestObjectFactory()

	// TODO: Should we do a PutString and GetString? And then one for each major
	// type for the general K/V storage that is not document based?
	fmt.Println("Directly adding test key(string)/value([]byte), to the")
	fmt.Println("to the collection.")
	err = collection.Put("test", []byte("value"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("no error in put!")
	}
	fmt.Println("Directly getting the test value we just inserted to test the")
	fmt.Println("lowest level key/value functionality.")
	get, err := collection.Get("test")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("no error in get!")
	}
	fmt.Println("get []byte:", get)
	fmt.Println("get string([]byte):", string(get))

	fmt.Println("=============================================================")
	fmt.Println("==---------------------------------------------------------==")
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("==---------------------------------------------------------==")
	fmt.Println("=============================================================")

	//fmt.Println("inserting a key/value into the created collect [ ", collectionName, "]")
	//fmt.Println("inserting test key: [", testKey, "] and test value: [", testValue, "]")
	//err = leveldb.Collection(collectionName).PutObject(testKey, testValue)
	//if err != nil {
	//	fmt.Println("[kv][ERROR:PutObject()] error inserting the item to into the kv collection:", err)
	//}

	//fmt.Println("Now getting the value out of the database using the test key...")

	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println("Standard GET / PUTs test")
	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println(" put 'test', []byte('test')")
	//kvStore.Put("test", []byte("test"))
	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println(" get 'test'")
	//fmt.Println(" WANT []byte('test')")
	//get, err := kvStore.Get("test")
	//if err != nil {
	//	fmt.Println("ERROR DURING BASIC GET:", err)
	//}
	//fmt.Println("get:", get)
	//fmt.Println("-------------------------------------------------------------")

	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println("Size of the collection is:", collection.Size())
	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println("Size of the collection is:", kvStore.Size())
	//fmt.Println("-------------------------------------------------------------")
	//fmt.Println("Each record:")
	//for _, record := range kvStore.All() {
	//	fmt.Println("record:", record)
	//}
	//fmt.Println("-------------------------------------------------------------")

	//returnValue := []byte{}

	//fmt.Println("returnValue:", returnValue)
	//err = kvStore.GetObject("test", &returnValue)
	//if err != nil {
	//	fmt.Println("[kv][ERROR:GetObject()] failed to use get on collection with test key:", err)
	//	fmt.Println("testKey:", testKey)
	//	fmt.Println("returnValue:", returnValue)
	//}
	//fmt.Println("The returned value is:", returnValue)
	//fmt.Println("The returned value as a string is:", string(returnValue))

}

func TestObjectFactory() TestObject {
	return TestObject{
		Name:        "test",
		Description: "description",
	}
}
