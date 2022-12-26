package model

import (
    "log"
    "time"
    "database/sql"
    "martian/config"
)

type Ration struct {
    Id          int    `db:"id,omitempty"`
    PacketType  string `db:"packet_type"`
    Content     string `db:"content"`
    Calories    int64  `db:"calories"`
    ExpiryDate  sql.NullString `db:"expiry_date"`
    Litre       float64  `db:"litre"`
}

type Food struct {
    Id          int    `db:"id,omitempty"`
    PacketType  string `db:"packet_type"`
    Content     string `db:"content"`
    Calories    int64  `db:"calories"`
    ExpiryDate  sql.NullString `db:"expiry_date"`
}

type Water struct {
    Id          int    `db:"id,omitempty"`
    PacketType  string `db:"packet_type"`
    Litre       float64  `db:"litre"`   
}

type ScheduleItems struct {
    ScheduleDate string
    TotalRow     int
    Foods        []Food
    Water        []Water
}

type ScheduleResponse struct {
    Items        []ScheduleItems
}

// get all foods for show
func GetRation() ([]Ration, error){
    db, err := config.ConnectDB()
    if err != nil {
        log.Printf("%s", err)
    }
    defer db.Close()

    var ration []Ration
    err = db.Select(&ration, "SELECT * FROM ration")
    if err != nil {
        log.Printf("%s", err)
    }

    return ration, nil
}

// add ration
func AddRation(packet_type string, content string, calories int64, expiry_date string, litre float64) (bool, error){
    db, err := config.ConnectDB()
    if err != nil {
        log.Printf("%s", err)
    }
    defer db.Close()
    
    if expiry_date != ""{
        _, err = db.Exec(`INSERT INTO ration (packet_type, content, calories, expiry_date, litre) VALUES (?,?,?,?,?)`, packet_type, content, calories, expiry_date, litre)
    } else {
        _, err = db.Exec(`INSERT INTO ration (packet_type, content, calories, expiry_date, litre) VALUES (?,?,?,?,?)`, packet_type, content, calories, nil, litre)
    }
    if err != nil {
        log.Printf("%s", err)
        return false, err
    }

    return true, nil
}

// delete ration
func DeleteRation(id string) (bool, error){
    db, err := config.ConnectDB()
    if err != nil {
        log.Printf("%s", err)
    }
    defer db.Close()

    _, err = db.Exec(`DELETE FROM ration WHERE id = ?`, id)
    if err != nil {
        log.Printf("%s", err)
        return false, err
    }

    return true, nil
}

// function for get foods and water & scheduling of ration
func Schedule() (ScheduleResponse, error) {
    db, err := config.ConnectDB()
    if err != nil {
        log.Printf("%s", err)
    }
    defer db.Close()

    var foods []Food
    err = db.Select(&foods, "SELECT id, packet_type, content, calories, expiry_date FROM `ration` WHERE packet_type = 'F' AND expiry_date >= CURDATE() ORDER BY expiry_date")
    if err != nil {
        log.Printf("%s", err)
    }

    var water []Water
    err = db.Select(&water, "SELECT id, packet_type, litre FROM `ration` WHERE packet_type = 'W'")
    if err != nil {
        log.Printf("%s", err)
    }

    res, err := MakeSchedule(foods, water)
    if err != nil {
        log.Printf("%s", err)
    }

    return res, nil
}

// make schedule of ration datewise
func MakeSchedule(foods []Food, water []Water) (ScheduleResponse, error) {
    var scheduleItems ScheduleItems
    var scheduleResponse ScheduleResponse
    var tmpFood []Food
    var usedFood []int
    var usedWater []int
    var tmpCal int64
    var scheduleDate string

    // start date is current date
    scheduleDate = time.Now().Format("2006-01-02")
    // loop over total food packets
    LOOP :
    for _, food := range foods {
        //check if food is not consumed
        _, found := find(usedFood, food.Id)
        if found {
            continue
        }

        if len(tmpFood) == 0 {
            if food.Calories < 2500 {
                tmpFood = append(tmpFood, food)
                tmpCal = food.Calories + tmpCal
                usedFood = append(usedFood, food.Id)
                // if single food packet fulfill the needed calories
                } else if food.Calories == 2500 {
                    tmpFood = append(tmpFood, food)
                    scheduleItems.Foods = tmpFood
                    day := len(scheduleResponse.Items)
                    scheduleItems.ScheduleDate =  time.Now().AddDate(0,0,day).Format("2006-01-02")
                    scheduleDate = time.Now().AddDate(0,0,day+1).Format("2006-01-02")
                    scheduleResponse.Items = append(scheduleResponse.Items, scheduleItems)
                    scheduleItems = ScheduleItems{}
                    tmpFood = nil
                    tmpCal = 0
                    usedFood = append(usedFood, food.Id)
                }
        } else {
            if (2500 - tmpCal) > 0 {
                //best match
                if (tmpCal + food.Calories == 2500) {
                    tmpFood = append(tmpFood, food)
                    scheduleItems.Foods = tmpFood
                    day := len(scheduleResponse.Items)
                    scheduleItems.ScheduleDate =  time.Now().AddDate(0,0,day).Format("2006-01-02")
                    scheduleDate = time.Now().AddDate(0,0,day+1).Format("2006-01-02")
                    scheduleResponse.Items = append(scheduleResponse.Items, scheduleItems)
                    scheduleItems = ScheduleItems{}
                    tmpFood = nil
                    tmpCal = 0
                    usedFood = append(usedFood, food.Id)
                } else if (tmpCal + food.Calories) < 2500 {
                    tmpFood = append(tmpFood, food)
                    tmpCal = food.Calories + tmpCal
                    usedFood = append(usedFood, food.Id)
                } else if (tmpCal + food.Calories) > 2500 {
                    // get optimal food packet if any
                    leastPacket, tc, fi := leastCalFood(foods, (tmpCal + food.Calories - 2500), tmpCal, usedFood, scheduleDate)
                    if (Food{}) != leastPacket {
                        if tc >= 2500 {
                            tmpFood = append(tmpFood, leastPacket)
                            scheduleItems.Foods = tmpFood
                            usedFood = append(usedFood, fi)
                            day := len(scheduleResponse.Items)
                            scheduleItems.ScheduleDate =  time.Now().AddDate(0,0,day).Format("2006-01-02")
                            scheduleDate = time.Now().AddDate(0,0,day+1).Format("2006-01-02")
                            scheduleResponse.Items = append(scheduleResponse.Items, scheduleItems)
                            tmpCal = 0
                            tmpFood = nil
                            scheduleItems = ScheduleItems{}
                        }
                    // if not found least calories packet then consume current packet    
                    } else {
                        tmpFood = append(tmpFood, food)
                        tmpCal = food.Calories + tmpCal
                        usedFood = append(usedFood, food.Id)
                        scheduleItems.Foods = tmpFood
                        day := len(scheduleResponse.Items)
                        scheduleItems.ScheduleDate =  time.Now().AddDate(0,0,day).Format("2006-01-02")
                        scheduleDate = time.Now().AddDate(0,0,day+1).Format("2006-01-02")
                        scheduleResponse.Items = append(scheduleResponse.Items, scheduleItems)
                        tmpCal = 0
                        tmpFood = nil
                        scheduleItems = ScheduleItems{}
                    }
                }
            }
        }
    }

    // iterate over remaining foods packet
    if len(foods) > len(usedFood) {
        var remainingCal int64
        for _, food := range foods {
            _, found := find(usedFood, food.Id)
            if found {
                continue
            }
            remainingCal = food.Calories + remainingCal
        }

        if remainingCal >= (2500 - tmpCal) {
            goto LOOP
        }
    }

    // add water to each day
    for index, res := range scheduleResponse.Items {

        wat, usedWat, litre := getWater(water, usedWater)
        usedWater = usedWat
        if (litre >= 2) {
            scheduleResponse.Items[index].Water = wat
            scheduleResponse.Items[index].TotalRow = len(res.Foods) + len(wat)
        } else {
            scheduleResponse.Items = removeFood(scheduleResponse.Items, index)
        }
    }

    return scheduleResponse, nil
}

// get water for survive
func getWater(water []Water, usedWater []int) ([]Water, []int, float64) {
    var tmpWater float64
    var resWater []Water
    var lastDiff float64
    lastDiff = 0
    hold := Water{}

    for _, wat := range water {
        _, found := find(usedWater, wat.Id)
        if found {
            continue
        }
        if (tmpWater < 2){
            if (tmpWater + wat.Litre == 2) {
                tmpWater = wat.Litre + tmpWater
                usedWater = append(usedWater, wat.Id)
                resWater = append(resWater, wat)
                hold = Water{}
                break
            } else if tmpWater + wat.Litre < 2 {
                tmpWater = wat.Litre + tmpWater
                usedWater = append(usedWater, wat.Id)
                resWater = append(resWater, wat)
                hold = Water{}
            } else if tmpWater + wat.Litre > 2 {
                if lastDiff == 0 {
                    lastDiff = tmpWater + wat.Litre - 2
                    hold = wat
                } else if lastDiff >  (tmpWater + wat.Litre - 2){
                    lastDiff = tmpWater + wat.Litre - 2
                    hold = wat
                }
                
            }
        } else {
            break
        }
    }
    if (hold != Water{}) {
        usedWater = append(usedWater, hold.Id)
        resWater = append(resWater, hold)
        tmpWater = tmpWater + hold.Litre
    }
    
    return resWater, usedWater, tmpWater
}

// find consumed food packet or water
func find(slice []int, val int) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

// find lest calories difference food packet
func leastCalFood(foods []Food, lastCalDiff int64, tmpCal int64, usedFood []int, scheduleDate string) (Food, int64, int) {
    leastFood := Food{}
    var foodId int

    for _, food := range foods {
        _, found := find(usedFood, food.Id)
        if found {
            continue
        }
        
        if food.ExpiryDate.String == scheduleDate {
            leastFood = food
            tmpCal = food.Calories + tmpCal
            foodId = food.Id
            break
        } else if (food.Calories + tmpCal -2500) < lastCalDiff {
            leastFood = food
            tmpCal = food.Calories + tmpCal
            foodId = food.Id
            break
        }
    }

    return leastFood, tmpCal, foodId
}

// remove food from response packet if insufficient
func removeFood(s []ScheduleItems, i int) []ScheduleItems {
    return append(s[:i], s[i+1:]...)
}