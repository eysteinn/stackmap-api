package database

import "time"

type UniqueProducts struct {
	Products []string `json:"uniqueproducts"`
}

type ProductTimes struct {
	Product string      `json:"product"`
	Times   []time.Time `json:"times"`
}
type AvailableTimes struct {
}

type Raster_geoms struct {
	Gid      uint `gorm:"primaryKey"`
	Location string
	Src_srs  string
	Datetime time.Time
	Product  string
	//geom     string
}
