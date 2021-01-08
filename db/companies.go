package db

import (
	"github.com/namnd/stockvn-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func FindCompanies(searchParams *model.CompanySearchParams) ([]*model.Company, error) {
	filter := bson.M{}
	if exchange := searchParams.Exchange; exchange != nil && *exchange != "all" {
		filter = bson.M{"exchange": *exchange}
	}
	sectorQuery := []bson.M{}
	if sectorIds := searchParams.SectorIds; sectorIds != nil && len(sectorIds) > 0 {
		ids := []string{}
		for _, sectorId := range sectorIds {
			ids = append(ids, *sectorId)
		}
		filter["$and"] = append(sectorQuery, bson.M{"sector_id": bson.M{"$in": ids}})
	}
	db := Connect()
	cursor, err := db.Companies.Find(db.Ctx, filter)
	defer cursor.Close(db.Ctx)
	defer db.Disconnect()

	if err != nil {
		return nil, err
	}
	var companies []*model.Company
	if err = cursor.All(db.Ctx, &companies); err != nil {
		return nil, err
	}
	return companies, nil
}

func FindCompany(exchange string, code string) (company *model.Company, err error) {
	filter := bson.M{"exchange": exchange, "code": code}
	db := Connect()
	err = db.Companies.FindOne(db.Ctx, filter).Decode(&company)
	return
}
