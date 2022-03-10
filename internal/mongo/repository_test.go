package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func fakeClient(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
}
