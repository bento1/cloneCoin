package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bento1/cloneCoin/utils"

	"github.com/bento1/cloneCoin/blockchain"

	"github.com/gorilla/mux"
)

type url string
type addBlockBody struct {
	Message string
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
} // MarshalText(endodeing package에서 옴) Marshal을 사용(json.NewEncoder(rw).Encode(data)사용)== json으로 엔코딩 할떄 필드가 json string으로 서 어떻게 보일지 결정하는 메소드 (완전한 URL 표시를위해)

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u urlDescription) String() string {
	return "Hello Im the URL"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{id}"),
			Method:      "GET",
			Description: "See A block",
		},
	}

	rw.Header().Add("Content-Type", "application/json") // 브라우저에게 보낸 string이 json임을 알려줌
	json.NewEncoder(rw).Encode(data)
}

// const port string = ":4000"
var port string

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().ListBlocks())
	case "POST":
		rw.Header().Add("Content-Type", "application/json")
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody)) //원본이 아닐수 있으니 원본을 보내야지
		blockchain.GetBlockChain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}
func Start(intport int) {
	port = fmt.Sprintf(":%d", intport)
	// handler_rest := http.NewServeMux()
	handler_rest := mux.NewRouter()
	handler_rest.HandleFunc("/", documentation)
	handler_rest.HandleFunc("/blocks", blocks)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler_rest))
	// 동시에 실행할 수없다. 먼저 한개만 하고있음 포트가 달라도 url이 같음
	// go explorer.Start(3000)
	// rest.Start(4000) http가  multipleresigistration이라고 표시되어있음 HandleFunc이 같은 url안에 작동되어있음
	// ListenAndServe() 보면 multipDefaultServeMux multiplexer는 리퀘스트 보내면 url을 보고 있다가 핸들러를 호출
	// 같은 멀티플레서를 rest와 explorer에서 사용하니깐
	// 새로운 멀티플렉서 설계해줌 서로다른 url 핸들러를 사용하게한다.
	// http.NewServeMux()
}
