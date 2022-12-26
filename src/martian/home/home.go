package home

import(
    "fmt"
    "strconv"
    "net/http"
    "html/template"
    "martian/model"
    "github.com/acoshift/flash"
    "github.com/gorilla/sessions"
    log "github.com/gophish/gophish/logger"
)

type Flash struct{
    FlashMessage string
}

var sessionStore sessions.Store

var f = flash.New()

// home page
func IndexAction(w http.ResponseWriter, r *http.Request) {
    tmpl,_ := template.ParseFiles("./layout/layout.html", "./home/view/index.html")
    var fm Flash
    if f.Has("message") == true {
        s := fmt.Sprintf("%s", f.Get("message"))
        fm.FlashMessage = s
        f.Clear()
    }
    
    tmpl.Execute(w, fm)
}

// show add ration form
func AddRation(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./layout/layout.html", "./home/view/add.html"))
    tmpl.Execute(w, nil)
}

// post request to save ration in DB
func SaveRation(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    packet_type := r.FormValue("packet_type")
    content := r.FormValue("content")
    var message string

    calories, _ := strconv.ParseInt(r.FormValue("calories"), 10, 64)
    
    expiry_date := r.FormValue("expiry_date")
    
    var litre float64
    if r.FormValue("litre") != ""{
        litre, _ = strconv.ParseFloat(r.FormValue("litre"), 64)
    } else {
        litre = 0
    }

    res, err := model.AddRation(packet_type, content, calories, expiry_date, litre)

    if res && err == nil {
        message = "Ration added successfully"
    } else {
        message = "Could not added ration!"
    }

    f.Set("message", message)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// show inventory
func ViewInventory(w http.ResponseWriter, r *http.Request) {
    ration, err := model.GetRation()
    if err != nil {
        log.Error(err)
    }

    tmpl := template.Must(template.ParseFiles("./layout/layout.html", "./home/view/view.html"))
    data := []interface{}{
        ration,
    }
    tmpl.Execute(w, data)
}

// delete ration
func DeleteRation(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    id := r.FormValue("id")
    var message string

    res, err := model.DeleteRation(id)

    if err == nil && res {
        message = "Ration deleted successfully"
    } else {
        message = "Could not delete ration!"
    }

    f.Set("message", message)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// view schedule of ration from current date
func Schedule(w http.ResponseWriter, r *http.Request) {
    tmpl,_ := template.ParseFiles("./layout/layout.html", "./home/view/schedule.html")
    
    res, err := model.Schedule()

    if err != nil {
        log.Error(err)
    }

    data := map[string]interface{}{
        "Results": res.Items,
        "Day":   len(res.Items),
    }

    tmpl.Execute(w, data)
}