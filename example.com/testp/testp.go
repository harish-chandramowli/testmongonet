package testp
import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
const (
	MongoURI="mongodb+srv://<user>:<pass>@cls1.ibqaj.mongodb-qa.net/myFirstDatabase?retryWrites=true&w=majority&apiVersion=1"

	)

func getConn(ctx context.Context, t *testing.T) *mongo.Client {
	opts := options.Client().ApplyURI(MongoURI)
	//opts.ServerAPIOptions = options.ServerAPI(options.ServerAPIVersion1).
	//	SetStrict(false).SetDeprecationErrors(false)
	client, err := mongo.Connect(ctx, opts )
	if err != nil { t.Fatal(err) }
	return client
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
