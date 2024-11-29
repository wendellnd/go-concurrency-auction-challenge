package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection

	batch             []auction_entity.Auction
	audictionInterval time.Duration
	timer             *time.Timer
	auctionChannel    chan auction_entity.Auction
}

func remove(s []auction_entity.Auction, i int) []auction_entity.Auction {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	batchUpdateInterval := getMaxBatchSizeInterval()

	auctionRepository := &AuctionRepository{
		Collection:        database.Collection("auctions"),
		audictionInterval: getAuctionInterval(),
		timer:             time.NewTimer(batchUpdateInterval),
		auctionChannel:    make(chan auction_entity.Auction, 100),
		batch:             make([]auction_entity.Auction, 0),
	}

	auctionRepository.triggerCreateRoutine(context.Background())

	return auctionRepository
}

func (ar *AuctionRepository) triggerCreateRoutine(ctx context.Context) {
	fmt.Println("triggerCreateRoutine")
	go func() {
		defer close(ar.auctionChannel)

		for {
			select {
			case auction, ok := <-ar.auctionChannel:
				if ok {
					logger.Info("Auction received " + auction.Id)
					ar.batch = append(ar.batch, auction)
				} else {
					logger.Info("Channel closed")
				}
				ar.timer.Reset(ar.audictionInterval)
			case <-ar.timer.C:
				if len(ar.batch) > 0 {
					for index, auction := range ar.batch {
						logger.Info("Checking auction " + auction.Id)
						if auction.Status == auction_entity.Active && time.Now().After(auction.Timestamp.Add(ar.audictionInterval)) {
							logger.Info("Auction completed " + auction.Id)
							auction.Status = auction_entity.Completed
							err := ar.UpdateAuction(ctx, auction)
							if err != nil {
								logger.Error("Error trying to update auction", err)
								return
							}

							ar.batch = remove(ar.batch, index)
						} else {
							logger.Info("Auction invalid " + auction.Id)
						}
					}

					ar.timer.Reset(ar.audictionInterval)
				}
			}
		}
	}()
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	fmt.Println(auctionEntity.Id)

	ar.auctionChannel <- *auctionEntity

	ar.triggerCreateRoutine(ctx)

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_UPDATE_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
