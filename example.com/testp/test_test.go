package testp

import (
	"context"
	"fmt"
	"github.com/mongodb/mongonet"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"testing"
)

type Count struct {
	i int
	lock sync.RWMutex
}

const (
DbName          = "test"
CollName        = "foo"
)

func simpleInsert(goctx context.Context, t *testing.T, conn *mongo.Client, c *Count) {
	_, err := conn.Database(DbName).Collection(CollName).InsertOne(goctx, bson.D{{"a", c.i}})
	if err != nil {
		t.Fatal(err)
	}
	c.lock.Lock()
	c.i = c.i+1
	c.lock.Unlock()
}

func BSONFindAndGetAsInt(doc bson.D, field string) (res int, err error) {
	idx := mongonet.BSONIndexOf(doc, field)
	if idx < 0 {
		return res, nil
	}
	res, tipe, err := mongonet.GetAsInt(doc[idx])
	if err != nil {
		return 0, fmt.Errorf("Expected '%v' to be int (or equivalent), but got %v instead. Doc = %v", field, tipe, doc)
	}
	return res, nil
}

func simpleFind(goctx context.Context, t *testing.T, conn *mongo.Client, c *Count) {
	c.lock.RLock()
	count := c.i
	c.lock.RUnlock()
	i := 0
	cur, _ := conn.Database(DbName).Collection(CollName).Find(goctx, bson.D{})
	for cur.Next(goctx) {
		elem := bson.D{}
		err := cur.Decode(&elem)
		if err != nil {
			t.Fatal(err)
		}
		a, err := BSONFindAndGetAsInt(elem, "a")
		if err != nil {
			t.Fatal(err)
		}
		if a != i {
			t.Fatalf("expected %v but got %v", i, a)
		}
		i += 1
	}
	if i < count {
		t.Fatalf("expected i to be greater or equal to count. i = %v, count = %v",i, count)
	}
}

func TestInsertM(t *testing.T) {
	c := Count{0, sync.RWMutex{}}
	goctx := context.Background()
	defer goctx.Done()
	insC := getConn(goctx, t)
	readC := getConn(goctx, t)
	for {
		simpleInsert(goctx, t, insC, &c)
		simpleFind(goctx, t, readC, &c)
	}
}