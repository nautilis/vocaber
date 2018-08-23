package main

import(
	"time"
	"fmt"
	"./vocaber"
	"net/http"
	"github.com/gin-gonic/gin/json"
	"log"
	"strconv"
)
func checkerr(err error){
	if err != nil{
		panic(err)
	}
}

const TOKEN string = "Iveryverysecret"

func isJsonErr(err error, jsonByte []byte) bool{
	if err != nil {
		fmt.Println("json parse error")
		return true
	}else{
		fmt.Println(string(jsonByte))
		return false
	}
}

func isValidToken(r *http.Request) bool{
	token := r.FormValue("token")
	if token != TOKEN{
		return false
	}
	return true
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")
	token := r.FormValue("token")
	res := make(map[string]interface{})
	if TOKEN != token {
		res["result"] = 0
	}else{
		vocaberItem := vocaber.VocabItem{Value: item, Created: time.Now(), Knownit: false}
		err := vocaber.Save(&vocaberItem)
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		end := start.Add(time.Hour * 24)
		count, err := vocaber.Count(start, end)
		if err != nil {
			res["result"] = 0
		}else{
			res["result"] = count
		}
	}
	log.Println(res)
	resJson, err := json.Marshal(res)
	if isJsonErr(err, resJson){
		res["result"] = 0
	}
	fmt.Fprint(w, string(resJson))
}

func getItemsBySubDay(w http.ResponseWriter, r *http.Request){
	subdayStr := (r.URL.Query().Get("subday"))
	subday, err:= strconv.Atoi(subdayStr)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now()
	todayBegin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	day := todayBegin.Add(time.Hour * -24 * time.Duration(subday))
	log.Println(todayBegin)
	items, err := vocaber.GetByDate(day)
	res := make(map[string]interface{})
	res["items"] = items
	resJson, err := json.Marshal(res)
	if(isJsonErr(err, resJson)){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(resJson))
}

func known(w http.ResponseWriter, r *http.Request){
	res := make(map[string]string)
	if !isValidToken(r){
		res["result"] = "failed"
		return
	}
	id := r.FormValue("itemid")
	idint, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	success, _ := vocaber.Know(idint)
	if success{
		res["result"] = "success"
	}
	resJson, err := json.Marshal(res)
	if isJsonErr(err, resJson){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(resJson))
}

func getNotMaster(w http.ResponseWriter, r *http.Request){
	items, err := vocaber.GetNoMaster()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respItems := items[0:20]
	res := make(map[string]interface{})
	res["items"] = respItems
	//TODO shuffle
	resJson, err := json.Marshal(res)
	if isJsonErr(err, resJson){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(resJson))
}

func getTodayCount(w http.ResponseWriter, r *http.Request){
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := now.Add(time.Hour * 24)
	count, err := vocaber.Count(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := make(map[string]int)
	res["result"] = count
	resJson, err := json.Marshal(res)
	if isJsonErr(err, resJson){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(resJson))
}

func deleteItem(w http.ResponseWriter, r *http.Request){
	res :=make(map[string]string)
	if !isValidToken(r){
		res["result"] = "failed"
	}
	idstr := r.FormValue("itemid")
	id, err := strconv.Atoi(idstr)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	success, _ := vocaber.Delete(id)
	if success {
		res["result"] = "success"
	}else{
		res["result"] = "success"
	}
	resJson, err := json.Marshal(res)
	if isJsonErr(err, resJson){
		res["result"] = "failed"
	}
	fmt.Fprint(w, string(resJson))
}

func main() {
	http.HandleFunc("/item", createHandler)
	http.HandleFunc("/items_by_subday", getItemsBySubDay)
	http.HandleFunc("/known_it", known)
	http.HandleFunc("/get_not_master", getNotMaster)
	http.HandleFunc("/get_today_count", getTodayCount)
	http.HandleFunc("/delete_item", deleteItem)
	log.Fatal(http.ListenAndServe(":8088", nil))
}


