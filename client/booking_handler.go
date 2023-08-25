package client

import (
	"encoding/json"
	"feastogether/captchasolver"
	"feastogether/config"
	"feastogether/fetch"
	"fmt"
	"log"
	"time"
)

// api
const (
	LOGIN_API      = "https://www.feastogether.com.tw/api/994f5388-d001-4ca4-a7b1-72750d4211cf/custSignIn"
	SAVE_SEATS_API = "https://www.feastogether.com.tw/api/booking/saveSeats"
	SAVE_SAETS_API = "https://www.feastogether.com.tw/api/booking/saveSaets"
	BOOKING_API    = "https://www.feastogether.com.tw/api/booking/booking"
	SVG_API        = "https://www.feastogether.com.tw/api/994f5388-d001-4ca4-a7b1-72750d4211cf/get2FASvgByBrand"
)

// 取得 Token
func GetToken(user config.UserConfig) string {

	payload := Login{
		Act:         user.Account,
		Pwd:         user.Password,
		ICode:       "+886",
		CountryCOde: "TW",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal struct to JSON:%v\n", err)
		return ""
	}

	resp, err := fetch.Post(LOGIN_API, payloadBytes, "", "")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	var data Response
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode response body: %v\n", err)
		return ""
	}
	if data.StatusCode != 1000 {
		log.Println(data.Result.Msg)
		return ""
	}
	return data.Result.CustomerLoginResp.Token
}

// 不清楚是什麼
// func GetSaveSaets(act string, token string) string {

// 	resp, err := fetch.Post(SAVE_SAETS_API, nil, act, token)
// 	if err != nil {
// 		log.Println(err)
// 		return ""
// 	}

// 	defer resp.Body.Close()

// 	var data SaveSaetsResponse
// 	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		log.Printf("Failed to decode response body: %v\n", err)
// 		return ""
// 	}

// 	if data.StatusCode != 1000 {
// 		log.Println(data)
// 		return ""
// 	}
// 	return data.Result
// }

func solver(brandId string) (string, string) {
	var code, ans string
	for {
		payload := GetSvg{
			BrandId: brandId,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal struct to JSON:%v\n", err)
			continue
		}

		resp, err := fetch.Post(SVG_API, payloadBytes, "", "")
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()

		var data Response
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("Failed to decode response body: %v\n", err)
			continue
		}

		if data.StatusCode != 1000 {
			log.Println(data.Result.Msg)
			continue
		}

		if data.Result.Svg == "" {
			return "", "" // No Captcha is needed
		}

		ans = captchasolver.SolveCaptcha(data.Result.Svg)
		if ans == "" {
			log.Println("Solver Failed.")
			continue
		}
		code = data.Result.Code
		break
	}
	fmt.Printf("Captcha Solved: %s\n", ans)

	return code, ans
}

// 立即訂位
func GetSaveSeats(userConfig config.UserConfig, payload config.RestaurantConfig) int {
	code, str := solver(payload.BrandId)
	token := GetToken(userConfig)

	saveSeats := SaveSeats{
		StoreID:     payload.StoreID,
		PeopleCount: payload.PeopleCount,
		MealPeriod:  payload.MealPeriod,
		MealDate:    payload.MealDate,
		MealTime:    payload.MealTime,
		MealSeq:     payload.MealSeq,
		Zkde:        "1j6ul4y94ejru6xk7vu4vu4",
		SvgCode:     code,
		SvgStr:      str,
	}

	var exp int
	for {
		payloadBytes, err := json.Marshal(saveSeats)
		if err != nil {
			log.Printf("Failed to marshal struct to JSON:%v\n", err)
			continue
		}

		resp, err := fetch.Post(SAVE_SEATS_API, payloadBytes, userConfig.Account, token)
		if err != nil {
			log.Println(err)
			continue
		}

		defer resp.Body.Close()

		var data Response
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("Failed to decode response body: %v\n", err)
			continue
		}

		if data.StatusCode != 1000 {
			log.Println(data.StatusCode, data.Result.Msg)
			// 101010: 客滿
			// 105006: 系統忙碌中
			if data.StatusCode != 105006 && data.StatusCode != 101010 {
				saveSeats.SvgCode, saveSeats.SvgStr = solver(payload.BrandId)
				token = GetToken(userConfig)
			}
			time.Sleep(time.Second)
			continue
		}
		exp = data.Result.ExpirationTime
		break
	}

	return exp
}

// 送出訂位
func SaveBooking(act string, token string, payload config.RestaurantConfig) bool {

	booking := Booking{
		StoreID:    payload.StoreID,
		MealPeriod: payload.MealPeriod,
		MealDate:   payload.MealDate,
		MealTime:   payload.MealTime,
		MealSeq:    4,
		Special:    0,
		ChildSeat:  0,
		Adult:      payload.PeopleCount,
		Child:      0,
		ChargeList: []struct {
			Seq   int "json:\"seq\""
			Count int "json:\"count\""
		}{
			// 大人
			{
				Seq:   201,
				Count: payload.PeopleCount,
			},
			// 小孩
			{
				Seq:   202,
				Count: 0,
			},
		},
		StoreCode:    payload.StoreCode,
		RedirectType: "iEat_card",
		Domain:       "https://www.feastogether.com.tw",
		PathFir:      "booking",
		PathSec:      "result",
		Yuuu:         "892389djdj883831445",
	}

	payloadBytes, err := json.Marshal(booking)
	if err != nil {
		log.Printf("Failed to marshal struct to JSON:%v\n", err)
	}

	resp, err := fetch.Post(BOOKING_API, payloadBytes, act, token)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return false
	}
	return true
}
