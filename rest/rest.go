package rest

import (
	"encoding/json"

	"fmt"

	"log"

	"net/http"

	"github.com/github.com/bento1/cloneCoin/blockchain"
	"github.com/github.com/bento1/cloneCoin/p2p"
	"github.com/github.com/bento1/cloneCoin/utils"
	"github.com/github.com/bento1/cloneCoin/wallet"

	mux "github.com/gorilla/mux"
)

type url string

type addPeerPayload struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}
type errorResponse struct {
	ErrorMessage string `json:"errormessage"`
}
type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}
type AddTxPayLoad struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}
type myWalletResponse struct {
	Address string `json:"address"`
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the BlockChain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get balance about address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}

	// rw.Header().Add("Content-Type", "application/json") // 브라우저에게 보낸 string이 json임을 알려줌
	json.NewEncoder(rw).Encode(data)
}

// const port string = ":4000"
var port string

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
	case "POST":
		newblock := blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
		p2p.BroadcastNewBlock(newblock)
	}
}
func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	// id, err_conversion := strconv.Atoi(vars["height"])hash 페이지로 변경하였음
	id := vars["hash"]

	// block, err_getblock := blockchain.GetBlockChain().GetBlock(id)
	block, err_getblock := blockchain.FindBlock(id)
	encoder := json.NewEncoder(rw)
	if err_getblock == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err_getblock)})
	} else {
		encoder.Encode(block)
	}

}
func status(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// // rw.Header().Add("Content-Type", "application/json")
		// return
		blockchain.Status(blockchain.BlockChain(), rw)

	}
}
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch r.Method {
	case "GET":
		switch total {
		case "true":
			amount := blockchain.BalanceByAddress(blockchain.BlockChain(), address)
			utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{Address: address, Balance: amount}))
		default:
			utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(blockchain.BlockChain(), address)))
		}

	}
}
func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	utils.HandleErr(json.NewEncoder(rw).Encode(myWalletResponse{Address: address}))
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload AddTxPayLoad
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)                  //헤더에도 돈이 부족하다고
		json.NewEncoder(rw).Encode(errorResponse{err.Error()}) //다양한 상황을
		return

	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)

}
func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		fmt.Println(payload.Address, payload.Port)
		p2p.AddPeer(payload.Address, payload.Port, port, true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.ALLPeers(&p2p.Peers)) // 데이터 레이스 만듬
	}
}

// middleware 설계 adapter 패턴
// 모든 함수에 rw.Header().Add("Content-Type", "application/json") 가 들어간다. 이함수는 json 타입을 rw에 써주는것을 알려주는 역할임
//documantation, blocks, block 수행전에 불려짐
//NewRouter 의 객체 handler.Use(Middlewarefunction) 으로 사용
//HandlerFunc 는 타입인데,  http.Handler 는 인터페이스 이다.
//타입괄호안에 것으로 변경해준다.
//MarshalText 사용하기 위해 URL type을 만들어 준것처럼
//HandlerFunc 은 어댑터 패턴을 수행하기 위한 중간 타입 어댑터에 적절한 args를 주면 어댑터는 http.Handler 에서 필요한 것 을 구현해준다.
//여기서는 rw.Header().Add("Content-Type", "application/json") 	next.ServeHTTP(rw, r)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json") //여기에 공통 적용되는 코드를 넣음
		next.ServeHTTP(rw, r)
	})
}
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}
func Start(intport int) {
	port = fmt.Sprintf(":%d", intport)
	// handler_rest := http.NewServeMux()
	// handler_rest := mux.NewRouter()
	handler_rest := mux.NewRouter()
	handler_rest.Use(jsonContentTypeMiddleware, loggerMiddleware)
	handler_rest.HandleFunc("/", documentation).Methods("GET")
	handler_rest.HandleFunc("/status", status).Methods("GET")
	handler_rest.HandleFunc("/balance/{address}", balance).Methods("GET") //거래목록을 봄
	// handler_rest.HandleFunc("/balance/{address}?total=ture", balance).Methods("GET") //total balance를 본다 ?total-true는 따로 만드는게아니라 옵션이있음
	handler_rest.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler_rest.HandleFunc("/mempool", mempool).Methods("GET")
	handler_rest.HandleFunc("/transactions", transactions).Methods("POST")
	handler_rest.HandleFunc("/wallet", myWallet).Methods("GET")
	handler_rest.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	handler_rest.HandleFunc("/peers", peers).Methods("GET", "POST")
	handler_rest.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET") //[0-9]숫자 hexadecimal은 [a-f]까지 가지는 형식
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler_rest))

}
