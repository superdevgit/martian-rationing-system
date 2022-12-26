# Martian Rationing System
You are part of the Ares III mission to Mars exploring “Acidalia Planitia” plain on Mars in the year 2035. Due to unfortunate circumstances involving a dust storm, you got stranded on Mars alone with no communication to your team or Earth’s command center.
Your surface habitat ("Hab") on Mars contains a limited inventory of food supplies. Each of these supplies is in the form of a packet containing either Food or Water. The food items have details of the content, expiry date and calories they provide. The water packet contains only the details of the quantity in liters and doesn’t expire.

#### Basic feature of the application

Add Ration : Record the details of the supply packet to a storage mechanism (DB, File.. etc)

View Inventory : Retrieve the details of all the supply packets in the inventory

Delete Ration : Ability to delete a supply packet from the inventory that has been consumed or needs update.

View Schedule : Retrieve the available inventory in the storage mechanism and generate the schedule.

## Prerequisites
    Need to have 'go' and 'mysql' installed on your machine. 

## Run on your machine

First setup config.json file with all your credentials. Import mysql database 'martian_ration'.

Install dependencies packages by following commands

    go get github.com/gorilla/mux
    go get github.com/jmoiron/sqlx
    go get github.com/acoshift/flash
    go get github.com/gorilla/sessions
    go get github.com/go-sql-driver/mysql
    go get github.com/codegangsta/negroni
    go get github.com/gophish/gophish/logger

Set up variables GOROOT, GOPATH and PATH and build project by following command in terminal, this will create an executable build file. To see output run executable file

    go build

and then run executable file

then visit in browser

    http://localhost:8000/

## For test cases run this command
    go test -v
