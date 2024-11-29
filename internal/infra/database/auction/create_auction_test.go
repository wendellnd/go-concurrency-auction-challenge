package auction

import (
	"context"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTriggerCreateRoutine(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should process auctions correctly", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())

		collection := mt.Coll
		repo := &AuctionRepository{
			Collection:        collection,
			audictionInterval: time.Second,
			timer:             time.NewTimer(time.Second),
			auctionChannel:    make(chan auction_entity.Auction, 1),
			batch:             make([]auction_entity.Auction, 0),
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go repo.triggerCreateRoutine(ctx)

		auction := auction_entity.Auction{
			Id:          uuid.New().String(),
			ProductName: "Test Product",
			Category:    "Test Category",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		repo.auctionChannel <- auction

		time.Sleep(2 * time.Second)

		if len(repo.batch) != 0 {
			t.Errorf("Expected batch to be empty, got %d", len(repo.batch))
		}
	})
}
