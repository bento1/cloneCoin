package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go/github.com/bento1/cloneCoin/utils"

	"go/github.com/bento1/cloneCoin/blockchain"
)

type URL string
type AddBlockBody struct {
	Message string
}

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
} // MarshalText(endodeing package에서 옴) Marshal을 사용(json.NewEncoder(rw).Encode(data)사용)== json으로 엔코딩 할떄 필드가 json string으로 서 어떻게 보일지 결정하는 메소드 (완전한 URL 표시를위해)

type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func (u URLDescription) String() string {
	return "Hello Im the URL"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "POST",
			Description: "Add A block",
			Payload:     "data:string",
		},
		{
			URL:         URL("/blocks/{id}"),
			Method:      "GET",
			Description: "See A block",
		},
	}

	rw.Header().Add("Content-Type", "application/json") // 브라우저에게 보낸 string이 json임을 알려줌
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)    3줄이 json.NewEncoder(rw).Encode(data) 로 대체할수 있음
	json.NewEncoder(rw).Encode(data)
	// 구조체를 Json으로 바꿩줘야해 Marshal 사용
	// json은 대부분 소문자로 이루어져있음  GO는 public이 대문자임 어떻게 할까?
	// field struct tag를 대신 사용함
	//URLDescription 뒤에 추가표시, omitempty 표시는  비어있으면 생략한다. 무시하고싶으면 `json:-`

}

const port string = ":4000"

//documatation을 만들고싶다
// API에서 할수 있는 일들의 목록을 보게해줌
// .../GET 들어가면 해당 documentation이 있음
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().ListBlocks())
	case "POST":
		// POST요청이온 메시지를 GOLANG struct로 decode해줘야함. NewDecoder는 reader를 받는데
		// r.Body로  부를수 있고 Decode는 포인터를 받음
		rw.Header().Add("Content-Type", "application/json")
		var addBlockBody AddBlockBody
		fmt.Println(addBlockBody)
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody)) //원본이 아닐수 있으니 원본을 보내야지
		fmt.Println(addBlockBody)
		blockchain.GetBlockChain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}
func main() {
	// explorer.Start()

	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
