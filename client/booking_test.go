package client_test

import (
	"feastogether/client"
	"feastogether/config"
	"fmt"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	if cfg, err := config.GetConfig(".."); err != nil {
		t.Error(err)
	} else {
		fmt.Println(client.GetToken(cfg.UserConfig))
	}
}

// func TestGetSaveSaets(t *testing.T) {
// 	if cfg, err := config.GetConfig(".."); err != nil {
// 		log.Println(err)
// 	} else {
// 		fmt.Println(client.GetSaveSaets(
// 			cfg.UserConfig.Account,
// 			client.GetToken(cfg.UserConfig)))
// 	}
// }

func TestGetSaveSeats(t *testing.T) {
	if cfg, err := config.GetConfig(".."); err != nil {
		t.Error(err)
	} else {
		nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		cfg.BrandId = "BR00001"
		cfg.StoreID = "S2212290008"
		cfg.StoreCode = "TNXM"
		cfg.MealDate = nextWeek
		ret := client.GetSaveSeats(cfg.UserConfig, cfg.RestaurantConfig)
		fmt.Println(ret)
	}
}

func TestSaveBooking(t *testing.T) {
	if cfg, err := config.GetConfig(".."); err != nil {
		t.Error(err)
	} else {
		nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		cfg.BrandId = "BR00001"
		cfg.StoreID = "S2212290008"
		cfg.StoreCode = "TNXM"
		cfg.MealDate = nextWeek
		secs := client.GetSaveSeats(cfg.UserConfig, cfg.RestaurantConfig)
		fmt.Println(secs)

		ret := client.SaveBooking(cfg.UserConfig.Account, client.GetToken(cfg.UserConfig), cfg.RestaurantConfig)
		fmt.Println(ret)
	}
}
