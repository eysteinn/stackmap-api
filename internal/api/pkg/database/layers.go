package database

import (
	"fmt"

	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/logger"
)

func GetUniqueProducts() (UniqueProducts, error) {
	var uniqueProducts UniqueProducts
	result := GetDB().Model(Raster_geoms{}).Select("product").Distinct("product").Find(&uniqueProducts.Products)
	logger.GetLogger().Debug("Lines: ", len(uniqueProducts.Products))
	return uniqueProducts, result.Error
}

func GetAvailableTimes(product string) (ProductTimes, error) { // (map[string]time.Time, error) {
	fmt.Println("Prod: ", product)
	prodTimes := ProductTimes{Product: product}
	result := GetDB().Model(Raster_geoms{}).Select("datetime").Find(&prodTimes.Times, "product = ?", product)

	fmt.Println(result.RowsAffected)
	return prodTimes, result.Error
	/*for _, layer := range layers {
		fmt.Println("Layer: ", layer)
		prodTimes := ProductTimes{Product: layer}
		GetDB().Model(Raster_geoms{}).Select("datetime").Find(&prodTimes.Times, "product = ?", layer)
		fmt.Println(len(prodTimes.Times))
	}*/
}
