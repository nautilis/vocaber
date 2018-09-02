package main

import(
	"time"
	"fmt"
	"./vocaber"
	"net/http"
	"log"
	"strconv"
	"strings"
	"os/exec"
	"io/ioutil"
	"encoding/json"
	"math/rand"
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

func shuffle(list []vocaber.VocabItem, size int)([]vocaber.VocabItem){
	length := len(list)
    newList := make([]vocaber.VocabItem, size)
    fmt.Println(list)
	s := rand.NewSource(time.Now().Unix())
    r := rand.New(s)
	for i:=0; i< size; i++ {
	    index := r.Intn(length)
        fmt.Println(index)
        newList[i] = list[index]
	}
	return newList
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getUrlRes(url string, second int) string{
	_timeout := time.Duration( time.Duration(second) * time.Second)
	client := &http.Client{
		Timeout: _timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		resp, err = client.Do(req)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString
}


func createHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	item := r.FormValue("item")
	token := r.FormValue("token")
	res := make(map[string]interface{})
	if TOKEN != token {
		res["result"] = 0
	}else{
		vocaberItem := vocaber.VocabItem{Value: item, Created: time.Now(), Knownit: 0}
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
	enableCors(&w)
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
	enableCors(&w)
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
	enableCors(&w)
	items, err := vocaber.GetNoMaster()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	length := len(items)
	var size int
	if length < 20{
		size = length
	}else{
		size = 20
	}
	respItems := shuffle(items, size)
	res := make(map[string]interface{})
	res["items"] = respItems
	//TODO shuffle
	resJson, err := json.Marshal(res)
	//TODO shuffle
	if isJsonErr(err, resJson){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(resJson))
}

func getTodayCount(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
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
	enableCors(&w)
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

type Sentence struct {
	Trans string `json:"trans"`
	Orig string `json:"orig"`
	Backend int `json:"backend"`
}

type  GoogleTran struct{
	Sentences []Sentence `json:"sentences"`
}

func translate(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	strList := strings.Split(r.URL.String(), "/")
	phrase := strList[len(strList)-1]
	log.Printf("phrase is %s", phrase)
	out, err := exec.Command("bash", "-c", "./get_tk.py " + phrase).Output()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tk := strings.Replace(string(out), "\n", "", -1)
	log.Printf( "got tk %s", string(out))
	url := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=auto&tl=zh-CN&hl=zh-CN&dt=t&dt=bd&dj=1&source=input&tk=" + tk + "&q=" + phrase
	respBody := getUrlRes(url, 5)
	//respBody := `{"sentences":[{"trans":"你好","orig":"hello","backend":1}],"dict":[{"pos":"感叹词","terms":["你好!","喂!"],"entry":[{"word":"你好!","reverse_translation":["Hello!","Hi!","Hallo!"],"score":0.13323711},{"word":"喂!","reverse_translation":["Hey!","Hello!"],"score":0.020115795}],"base_form":"Hello!","pos_enum":9}],"src":"en","confidence":1,"ld_result":{"srclangs":["en"],"srclangs_confidences":[1],"extended_srclangs":["en"]}}`
	googleTran := GoogleTran{}
	if err := json.Unmarshal([]byte(respBody), &googleTran); err !=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trans := googleTran.Sentences[0].Trans
	log.Printf("google json %s", trans)
	result := make(map[string]string)
	result["result"] = trans
	respJson, err := json.Marshal(result)
	if isJsonErr(err, respJson){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(respJson))
}

func main() {
	http.HandleFunc("/item", createHandler)
	http.HandleFunc("/items_by_subday", getItemsBySubDay)
	http.HandleFunc("/known_it", known)
	http.HandleFunc("/get_not_master", getNotMaster)
	http.HandleFunc("/get_today_count", getTodayCount)
	http.HandleFunc("/delete_item", deleteItem)
	http.HandleFunc("/get_translate/", translate)
	log.Fatal(http.ListenAndServe(":8088", nil))
}


