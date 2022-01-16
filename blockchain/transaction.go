package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/github.com/bento1/cloneCoin/utils"
	"github.com/github.com/bento1/cloneCoin/wallet"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	TxID      string `json:"txid"` // 이전의 transaction output을 가져올수 있는 방법이됨
	Index     int    `json :"index"`
	Signature string `json:"signature"`
	// Amount int    `json:"amount"`
}
type mempool struct {
	// Txs []*Tx //배열형태는 탐색이 오래걸림
	Txs map[string]*Tx //key TxID
	m   sync.Mutex
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}
type UTxOut struct {
	TxID   string `json:"txid"` // 이전의 transaction output을 가져올수 있는 방법이됨
	Index  int    `json :"index"`
	Amount int    `json:"amount"`
}

var m *mempool
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

var ErrorNoMoney = errors.New("Not Enough Money")
var ErrorNotValid = errors.New("Not Valid")

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}
func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(wallet.Wallet(), t.Id)
	}
}
func validate(tx *Tx) bool {
	//아웃풋으로 다음거래의 input으로 사용하는데, 아웃풋의 오너 (address) 를 검증해야한다. 블럭안에 들어있는 돈이 너의 돈인지 확인해야함
	//서명도 있고, public key도 있음 사람들이 너에게 보내는 address 너의 주소
	//tranasaction input의 sig를 가지고 있고 output으로 등록해줄때 verification을 할 수 있음.
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(BlockChain(), txIn.TxID) //직전거래를 찾음
		if prevTx == nil {
			valid = false
			break
		}
		//tx input이 참조한 txoutput을 찾음 그곳의 address를 가져옴
		//txinput의 index는 txoutput을 알수 있게 해주기 떄문에
		//unspent tx output의 index로 간다

		address := prevTx.TxOuts[txIn.Index].Address //prev Tx는 찾았다. txouput을 가서 txinput의 index로 간다.
		valid = wallet.Verify(txIn.Signature, tx.Id, address)
		if !valid {
			break
		}
	}

	return valid
}
func makeCoinBaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId() //이부분에 signatrue가 들어감
	return &tx

}

func makeTx(from string, to string, amount int) (*Tx, error) {
	// 현재금액을 알아야함 from의 모든 output을 추적한다.
	// 송금이 충분한 금액을 가지고 있는지 확인한다.
	// output > input을 만족하는 block 수만으로 송금액을 만족하는 block을 찾는다.
	if BalanceByAddress(BlockChain(), from) < amount {
		return nil, ErrorNoMoney
	}
	//돈이 충분하다면 해당 주소를 가지고 있는 블럭들을 가져오고
	// amount를 만족시키는 block 만큼 모아서 input block에 넣는다.
	var total int = 0
	var txIns []*TxIn
	var txOuts []*TxOut
	uTxOuts := UTxOutsByAddress(BlockChain(), from)
	for _, txOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{txOut.TxID, txOut.Index, from} //from의 owner임.   .//추가만할 뿐 사용할 output을 검증 하지 않았고, output을 사용했다고 체크하지 않음
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}
	//잔돈
	change := total - amount
	if change != 0 {
		changeTx := &TxOut{from, change}
		txOuts = append(txOuts, changeTx)

	}
	//송금
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut) //from잔돈 , to 송금액으로 구성되어있음
	tx := &Tx{Id: "", Timestamp: int(time.Now().Unix()), TxIns: txIns, TxOuts: txOuts}
	tx.getId()
	tx.sign()
	if validate(tx) {
		return tx, nil
	} else {
		return nil, ErrorNotValid
	}
}

func isOnMempool(UTxOut *UTxOut) bool {
	exist := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == UTxOut.TxID && input.Index == UTxOut.Index {
				exist = true
				// break// 이렇게 하면 안쪽만 종료된다.
				break Outer //label 기능으로 바깥도
			}
		}
	}
	return exist
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	//누가 보냈는지는 알필요가없다.
	// 지갑 (to)에서 정보를 받아오면된다.
	m.m.Lock()
	defer m.m.Unlock()
	tx, err := makeTx(wallet.Wallet().Address, to, amount) //나중에 wallet이 들어감
	if err != nil {
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil

}

//컨펌되지 않은 transactions 가져오기 from mempool

func (m *mempool) txToConfirm() []*Tx {
	m.m.Lock()
	defer m.m.Unlock()
	coinbase := makeCoinBaseTx(wallet.Wallet().Address)
	var txs []*Tx
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase) //rewards
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()

	m.Txs[tx.Id] = tx
}
