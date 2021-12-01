package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go/github.com/bento1/cloneCoin/blockchain"
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
	data := homeData{"Home", blockchain.GetBlockChain().ListBlocks()}
	templates.ExecuteTemplate(rw, "home", data) // 위에 전역변수로 설정하였음
}
func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()                             // Post에 검색해서 찾으면 ParseForm, Form 순으로 부르면 됨
		data := r.Form.Get("blockData")           //Form은 Value라고 나와있음 add page에 input에 name설정한 부분
		blockchain.GetBlockChain().AddBlock(data) // 찾은 값으로 BLock을 추가해줌
		//redirection을 하고싶다. 공식 문서 검색->  찾아봄 ->
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
	// 동시에 실행할 수없다. 먼저 한개만 하고있음 포트가 달라도 url이 같음
	// go explorer.Start(3000)
	// rest.Start(4000) http가  multipleresigistration이라고 표시되어있음 HandleFunc이 같은 url안에 작동되어있음
	// ListenAndServe() 보면 multipDefaultServeMux multiplexer는 리퀘스트 보내면 url을 보고 있다가 핸들러를 호출
	// 같은 멀티플레서를 rest와 explorer에서 사용하니깐
	// 새로운 멀티플렉서 설계해줌 서로다른 url 핸들러를 사용하게한다.
	// http.NewServeMux()
}
