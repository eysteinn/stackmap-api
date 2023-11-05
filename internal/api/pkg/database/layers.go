package database

import (
	"fmt"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/EysteinnSig/stackmap-api/internal/api/pkg/logger"
)

func GetUniqueProjects() (UniqueProjects, error) {
	var uniqueProjects UniqueProjects
	gorm.DefaultTableNameHandler = func(db *gorm.DB, tableName string) string {
		return "public." + tableName
	}
	logger.GetLogger().Debug("Executing projects code")
	result := GetDB().Raw("select regexp_replace(n1.schema_name, '^project_', '') as project from (select schema_name from information_schema.schemata where schema_name ~ '^project*') n1;").Scan(&uniqueProjects.Projects)
	logger.GetLogger().Debug(result.Statement.SQL.String())
	//select regexp_replace(n1.schema_name, '^project_', '') as project from (select schema_name from information_schema.schemata where schema_name ~ '^project*') n1;
	//result := GetDB().Model(Raster_geoms{}).Select("product").Distinct("product").Find(&uniqueProducts.Products)
	logger.GetLogger().Debug("Lines: ", len(uniqueProjects.Projects))
	return uniqueProjects, result.Error
}

func GetUniqueProducts(project string) (UniqueProducts, error) {
	var uniqueProducts UniqueProducts
	/*gorm.DefaultTableNameHandler = func(db *gorm.DB, tableName string) string {
		fmt.Println("Setting schema1")
		return "project_" + project + "." + tableName
	}*/
	//result := GetDB().Model(Raster_geoms{}).Select("product").Distinct("product").Find(&uniqueProducts.Products)
	result := GetDB().Table("project_" + project + ".raster_geoms").Select("product").Distinct("product").Find(&uniqueProducts.Products)
	logger.GetLogger().Debug("Lines: ", len(uniqueProducts.Products))
	return uniqueProducts, result.Error
}

func GetAvailableTimes(project string, product string) (ProductTimes, error) { // (map[string]time.Time, error) {
	/*gorm.DefaultTableNameHandler = func(db *gorm.DB, tableName string) string {
		fmt.Println("Setting schema1")
		return "project_" + project + "." + tableName
	}*/
	fmt.Println("Prod: ", product)
	prodTimes := ProductTimes{Product: product}
	//result := GetDB().Model(Raster_geoms{}).Select("datetime").Find(&prodTimes.Times, "product = ?", product)
	result := GetDB().Table("project_"+project+".raster_geoms").Select("datetime").Find(&prodTimes.Times, "product = ?", product)
	fmt.Println(result.RowsAffected)
	tmpkey := map[time.Time]bool{}
	tmplst := []time.Time{}
	for _, t := range prodTimes.Times {
		if _, value := tmpkey[t]; !value {
			tmpkey[t] = true
			tmplst = append(tmplst, t)
		}
	}
	prodTimes.Times = tmplst

	sort.Slice(prodTimes.Times, func(i, j int) bool {
		return prodTimes.Times[i].Before(prodTimes.Times[j])
	})

	return prodTimes, result.Error
	/*for _, layer := range layers {
		fmt.Println("Layer: ", layer)
		prodTimes := ProductTimes{Product: layer}
		GetDB().Model(Raster_geoms{}).Select("datetime").Find(&prodTimes.Times, "product = ?", layer)
		fmt.Println(len(prodTimes.Times))
	}*/
}
