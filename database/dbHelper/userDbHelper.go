package dbHelper

import (
	"06_RMS-chi-db/database"
	"06_RMS-chi-db/errorHandling"
	"06_RMS-chi-db/models"
	"fmt"
	"math"
)

func CreateUserNewAddress(a *models.Addresses, token string) (*models.Addresses, error) {
	id, err := CheckUserId(token)
	if err != nil {
		return nil, err
	}
	var add models.Addresses
	SQL := `INSERT INTO addresses (address_name, address_lat, address_lng,user_id) VALUES ($1,$2,$3,$4) RETURNING id, address_name, address_lat, address_lng`
	err = database.RMS.Get(&add, SQL, a.AddressName, a.AddressLat, a.AddressLng, id)
	if err != nil {
		return nil, errorHandling.UnableToAccessDB()
	}
	return a, nil
}

func GetAllRestaurants() ([]models.Restaurants, error) {
	rests := make([]models.Restaurants, 0)
	SQL := `SELECT restaurant_name, restaurant_address, user_id FROM restaurants WHERE archived_at IS NULL`
	err := database.RMS.Select(&rests, SQL)
	if err != nil {
		return nil, errorHandling.UnableToAccessDB()
	}
	return rests, nil
}

func GetRestaurantDishes(resId int) ([]models.Dishes, error) {
	dishes := make([]models.Dishes, 0)
	SQL := `SELECT dish_name, dish_cost, restaurant_id, user_id, preparing_time FROM dishes WHERE restaurant_id = $1 AND archived_at IS NULL`
	err := database.RMS.Select(&dishes, SQL, resId)
	if err != nil {
		return nil, errorHandling.UnableToAccessDB()
	}
	return dishes, nil
}

func GetAddDistance(resId int, userAddId int) (*models.AddressDistance, error) {
	var uAdd models.Addresses
	SQL1 := `SELECT address_name, address_lat, address_lng FROM addresses WHERE id=$1`
	err := database.RMS.Get(&uAdd, SQL1, userAddId)
	if err != nil {
		return nil, errorHandling.UnableToAccessDB()
	}
	var rAdd models.Addresses
	SQL2 := `SELECT a.address_name, a.address_lat, a.address_lng FROM addresses a INNER JOIN restaurants r ON a.id = r.restaurant_address WHERE r.id=$1`
	err = database.RMS.Get(&rAdd, SQL2, resId)
	if err != nil {
		return nil, errorHandling.UnableToAccessDB()
	}
	dlong := rAdd.AddressLng - uAdd.AddressLng
	dlat := rAdd.AddressLat - uAdd.AddressLat
	ans := (dlong + dlat) / 2.0
	ans = math.Abs(ans) * 6.371
	fmt.Println(uAdd)
	fmt.Println(rAdd)
	fmt.Println("distance", ans)
	d := models.AddressDistance{
		UserAddId: userAddId,
		RestId:    resId,
		Distance:  ans,
	}

	return &d, nil
}
