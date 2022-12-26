package main

import (
    "fmt"
    "time"
    "testing"
    "database/sql"
    "martian/model"
)

type Response struct {
    Success string    `json:"success"`
    Mesaage string    `json:"message"`
}

// food and water both are for equal days
func TestEqualFoodAndWater(t *testing.T) {

    day := time.Now().AddDate(0,0,5).Format("2006-01-02")
    date := sql.NullString{day, true}

    var food = []model.Food{{1, "F" ,"MRE", 1500, date}, {2, "F", "Protein Bar", 1000, date}}
    var water = []model.Water{{1, "W", 1.5}, {2, "W", 0.5}}

    res, err := model.MakeSchedule(food, water)

    fmt.Println("food and water both for 1 day")
    
    if err != nil {
        t.Errorf("Fail, expected days is 1 and got %d", len(res.Items))
    } else {
        fmt.Printf("Expected days is 1 and got %d\n", len(res.Items))
    }
}

// food is for two days but water is insufficient
func TestInsufficientWater(t *testing.T) {

    day := time.Now().AddDate(0,0,5).Format("2006-01-02")
    date := sql.NullString{day, true}

    var food = []model.Food{{1, "F" ,"MRE", 1500, date}, {2, "F", "Protein Bars", 1000, date}, {13, "F", "Dry fruits", 2500, date}}
    var water = []model.Water{{1, "W", 1.5}, {2, "W", 0.5},  {3, "W", 1.5}}

    res, err := model.MakeSchedule(food, water)
    
    fmt.Println("food is for 2 days and water for 1 day")

    if err != nil || len(res.Items) > 1 {
        t.Errorf("Fail, expected days is 1 and got %d", len(res.Items))
    } else {
        fmt.Printf("Expected days is 1 and got %d\n", len(res.Items))
    }
}

// water if for two days but food is insufficient
func TestInsufficientFood(t *testing.T) {

    day := time.Now().AddDate(0,0,5).Format("2006-01-02")
    date := sql.NullString{day, true}

    var food = []model.Food{{1, "F" ,"MRE", 1500, date}, {2, "F", "Protein Bars", 1000, date}, {13, "F", "Dry fruits", 2400, date}}
    var water = []model.Water{{1, "W", 1.5}, {2, "W", 0.5},  {3, "W", 2}}

    res, err := model.MakeSchedule(food, water)
 
    fmt.Println("food is for 1 day and water for 2 day")

    if err != nil || len(res.Items) > 1 {
        t.Errorf("Fail, expected days is 1 and got %d", len(res.Items))
    } else {
        fmt.Printf("Expected days is 1 and got %d\n", len(res.Items))
    }
}
