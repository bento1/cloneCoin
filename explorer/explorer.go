package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/github.com/bento1/cloneCoin/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var port string
var templates *template.Template // 모든 template는 templates 변수로 컨트롤 하겠음

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block // 템플릿에서 읽을수 있어야함 export해야함
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data) // 위에 전역변수로 설정하였음
}
func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.BlockChain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)

	}
}
func Start(intport int) {
	port = fmt.Sprintf(":%d", intport)
	handler_explorer := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))     //home, add page 로드
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // ultities  gohtml  파일들 로드
	handler_explorer.HandleFunc("/", home)
	handler_explorer.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler_explorer))

}
