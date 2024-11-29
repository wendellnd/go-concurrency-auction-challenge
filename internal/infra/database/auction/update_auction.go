package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) UpdateAuction(ctx context.Context, auction auction_entity.Auction) *internal_error.InternalError {
	auctionMongo := &AuctionEntityMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := ar.Collection.ReplaceOne(ctx, bson.M{"_id": auction.Id}, auctionMongo)
	if err != nil {
		logger.Error("Error trying to update auction", err)
		return internal_error.NewInternalServerError("Error trying to update auction")
	}

	return nil
}
