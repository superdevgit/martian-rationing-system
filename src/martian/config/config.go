package config

import(
    "os"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "github.com/jmoiron/sqlx"
  _ "github.com/go-sql-driver/mysql"
)

// Env stores the environment
var Env string

// AppPort stores port configuration for the server
var AppPort string = SetConfigValue("GOPORT", "8000")

// ElogFile has the name of the file to write container logs
var ElogFile string = SetConfigValue("ELOGFILE", "elog.txt")

// Esyslog check if on debug mode or not
var Esyslog string = SetConfigValue("ESYSLOG", "debug")

// DomainName is the domain name will allow to connect with the Dal
var DomainName string = SetConfigValue("DOMAINNAME", "https://localhost")

// mysql database params
var User              string = SetConfigValue("db_user", "YOUR_USER")

var Password          string = SetConfigValue("db_password", "YOUR_PASSWORD")

var Host              string = SetConfigValue("db_host", "DB_HOST")

var Database          string = SetConfigValue("db", "DATABASE")

var Type              string = SetConfigValue("TYPE", "mysql")

var Connectiondetails string = "" + User + ":" + Password + "@tcp(" + Host + ")/" + Database + ""

// SetConfigValue set environment by value and string
func SetConfigValue(e string, d string) string {
    jsonFile, err := os.Open("./config.json")
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()
    byteValue, err := ioutil.ReadAll(jsonFile)
    if err != nil {
        fmt.Println(err)
    }
    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)

    r := result[e]
    if r != nil {
        if v := result[e].(string); v != ""{
            return string(v)
        }
    }

    return d
}

// ConnectDB creates a connection and ping to the specific database
func ConnectDB() (*sqlx.DB, error) {

    return sqlx.Connect(Type, Connectiondetails)
}
