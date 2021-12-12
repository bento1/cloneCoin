package blockchain

import (
	"errors"
	"time"

	"github.com/github.com/bento1/cloneCoin/utils"
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
	TxID  string `json:"txid"` // 이전의 transaction output을 가져올수 있는 방법이됨
	Index int    `json :"index"`
	Owner string `json:"owner"`
	// Amount int    `json:"amount"`
}
type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}
type UTxOut struct {
	TxID   string `json:"txid"` // 이전의 transaction output을 가져올수 있는 방법이됨
	Index  int    `json :"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
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
	tx.getId()
	return &tx
}

// func makeTx(from string, to string, amount int) (*Tx, error) {
// 	// 현재금액을 알아야함 from의 모든 output을 추적한다.
// 	// 송금이 충분한 금액을 가지고 있는지 확인한다.
// 	// output > input을 만족하는 block 수만으로 송금액을 만족하는 block을 찾는다.
// 	if BlockChain().BalanceByAddress(from) < amount {
// 		return nil, errors.New("Not enough mody")
// 	}
// 	//돈이 충분하다면 해당 주소를 가지고 있는 블럭들을 가져오고
// 	// amount를 만족시키는 block 만큼 모아서 input block에 넣는다.
// 	var total int = 0
// 	var txIns []*TxIn
// 	var txOuts []*TxOut
// 	oldTxOuts := BlockChain().TxOutsByAddress(from)
// 	for _, txOut := range oldTxOuts {
// 		if total >= amount {
// 			break
// 		}
// 		txIn := &TxIn{txOut.Owner, txOut.Amount} //from의 owner임.   .//추가만할 뿐 사용할 output을 검증 하지 않았고, output을 사용했다고 체크하지 않음
// 		txIns = append(txIns, txIn)
// 		total += txOut.Amount
// 	}
// 	//잔돈
// 	change := total - amount
// 	if change != 0 {
// 		changeTx := &TxOut{from, change}
// 		txOuts = append(txOuts, changeTx)

// 	}
// 	//송금
// 	txOut := &TxOut{to, amount}
// 	txOuts = append(txOuts, txOut)
// 	tx := &Tx{Id: "", Timestamp: int(time.Now().Unix()), TxIns: txIns, TxOuts: txOuts}
// 	tx.getId()

// 	return tx, nil
// }
func makeTx(from string, to string, amount int) (*Tx, error) {
	// 현재금액을 알아야함 from의 모든 output을 추적한다.
	// 송금이 충분한 금액을 가지고 있는지 확인한다.
	// output > input을 만족하는 block 수만으로 송금액을 만족하는 block을 찾는다.
	if BalanceByAddress(BlockChain(), from) < amount {
		return nil, errors.New("Not enough money")
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

	return tx, nil
}
func isOnMempool(UTxOut *UTxOut) bool {
	exist := false
Outer:
	for _, tx := range Mempool.Txs {
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

func (m *mempool) AddTx(to string, amount int) error {
	//누가 보냈는지는 알필요가없다.
	// 지갑 (to)에서 정보를 받아오면된다.
	tx, err := makeTx("dongun", to, amount) //나중에 wallet이 들어감
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil

}

//컨펌되지 않은 transactions 가져오기 from mempool

func (m *mempool) txToConfirm() []*Tx {
	coinbase := makeCoinBaseTx("dongun")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
