package main

import (
	"feastogether/client"
	"feastogether/config"
	"log"
)

func main() {
	// 讀取 config
	cfg, err := config.GetConfig("./")
	if err != nil {
		log.Println(err)
		return
	}

	// 取得 token
	// token = client.GetToken(cfg.UserConfig)
	// if Token == "" {
	// 	return
	// }

	// 立即訂位 , 取得訂位開始 - 過期時間
	if secs := client.GetSaveSeats(cfg.UserConfig, cfg.RestaurantConfig); secs == 0 {
		log.Println("訂位失敗 Expiration Time = ", secs)
		return
	}
	log.Println("訂位成功")

	// 判斷是否取得訂位開始時間
	// 確認訂位
	if ret := client.SaveBooking(cfg.UserConfig.Account, client.GetToken(cfg.UserConfig), cfg.RestaurantConfig); ret == false {
		log.Println("訂位送出失敗。")
		return
	}
	log.Println("訂位已送出，請盡速付訂金。")
}
