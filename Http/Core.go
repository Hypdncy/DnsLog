package Http

import (
	"DnsLog/Core"
	"DnsLog/Dns"
	"crypto/subtle"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

//go:embed home
var home embed.FS

type resDnsJson struct {
	Code int64         `json:"code"`
	Data []Dns.LogInfo `json:"data"`
}

type resDomainJson struct {
	Code   int64  `json:"code"`
	Domain string `json:"domain"`
}

type application struct {
}

func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			expectedPassword, userOk := Core.Config.USER[username]
			passwordOk := subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1
			if passwordOk && userOk {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func ListingHttpManagementServer() {
	app := new(application)
	mux := http.NewServeMux()
	if !Core.Config.HTTP.ConsoleDisable {
		mux.Handle("/home/", http.FileServer(http.FS(home)))
		//mux.HandleFunc("/home", home)
	}
	mux.HandleFunc("/", app.basicAuth(index))
	//mux.HandleFunc("/index.html", app.basicAuth(getDnsData))
	//mux.HandleFunc("/api/verifyToken", app.basicAuth(verifyTokenApi))
	mux.HandleFunc("/api/get", app.basicAuth(getDnsData))
	mux.HandleFunc("/api/getDomain", app.basicAuth(getDomain))
	mux.HandleFunc("/api/clean", app.basicAuth(cleanDnsData))
	log.Println("Http Listing Start...")
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", Core.Config.HTTP.Ip, strconv.Itoa(Core.Config.HTTP.Port)),
		Handler: mux,
	}
	log.Println("Http address: http://" + Core.Config.HTTP.Ip + ":" + strconv.Itoa(Core.Config.HTTP.Port))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

//func home(w http.ResponseWriter, r *http.Request) {
//	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
//}

func index(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(301)
	//w.Header().Set("Location", "/home")
	////w.WriteHeader(200)
	http.Redirect(w, r, "/home/", http.StatusMovedPermanently)
}

func getDomain(w http.ResponseWriter, r *http.Request) {
	var username, _, _ = r.BasicAuth()
	err := json.NewEncoder(w).Encode(resDomainJson{Code: 200, Domain: username + "." + Core.Config.DNS.Domain})
	if err != nil {
		return
	}
}

func getDnsData(w http.ResponseWriter, r *http.Request) {
	var username, _, _ = r.BasicAuth()
	var infos, ok = Dns.L.Get(username)
	w.Header().Set("content-type", "text/json")
	w.WriteHeader(200)

	if ok {
		err := json.NewEncoder(w).Encode(makeResJson(200, infos))
		if err != nil {
			return
		}
	} else {
		err := json.NewEncoder(w).Encode(makeResJson(200, []Dns.LogInfo{}))
		if err != nil {
			return
		}
	}
}
func cleanDnsData(w http.ResponseWriter, r *http.Request) {

	var username, _, _ = r.BasicAuth()
	Dns.L.Clear(username)
	w.Header().Set("content-type", "text/json")
	w.WriteHeader(200)
	_, err := w.Write([]byte("{}"))
	if err != nil {
		return
	}
}

func makeResJson(code int64, infos []Dns.LogInfo) resDnsJson {
	var res = resDnsJson{Code: 200, Data: []Dns.LogInfo{}}
	res.Code = code
	for _, info := range infos {
		res.Data = append(res.Data, info)
	}

	return res
}
