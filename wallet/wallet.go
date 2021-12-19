package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/github.com/bento1/cloneCoin/utils"
)

// const (
// 	signature     string = "f8e3c21ea69927de108c9d6746c3b6a54ab7d13539925bce031f2adbe42259cad1fded31dce6e2a663b97ba8b8c2b32ee824eda343257e7cbe182ed2af99d4d9"
// 	privateKey    string = "30770201010420ba6081f8aabed983a7427cd4de97a8722f5f87b7775b182dfe8d7230403abf75a00a06082a8648ce3d030107a144034200047633f71af619c2d145b6d9c6e0d3c1ecbebe42938b4395322e8300fdb1626156b175ba157c89df45760af96d8c6e70c319c472c67f92b9e9a9029f79b4f6a87c"
// 	HashedMessage string = "c33084feaa65adbbbebd0c9bf292a26ffc6dea97b170d88e501ab4865591aafd"
// )
const (
	walletFileName string = "dongcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string //hexdecimal
}

var w *wallet

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}
func Sign(w *wallet, payload string) string {
	payloadAsByte, err := hex.DecodeString(payload) //[]byte{} 해도되는데 err 체크  길이나 이런것들을 체크해야함
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsByte)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}
func resotreBigInt(payload string) (*big.Int, *big.Int, error) {
	Byte, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	utils.HandleErr(err)
	firstHalf := Byte[:len(Byte)/2]
	secondHalf := Byte[len(Byte)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalf)
	bigB.SetBytes(secondHalf)
	return &bigA, &bigB, nil
}
func Verify(sign string, payload string, address string) bool {
	R, S, err := resotreBigInt(sign)
	utils.HandleErr(err)
	X, Y, err := resotreBigInt(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: X, Y: Y}
	payloadByte, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publicKey, payloadByte, R, S)
	return ok
}

func aFromK(key *ecdsa.PrivateKey) string {
	// publickey := key.PublicKey //public key는 bigint의 x,y로 되어있음
	//publickey는 address가 됨
	z := append(key.X.Bytes(), key.Y.Bytes()...)
	return fmt.Sprintf("%x", z)
}
func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)
	//err가 다양한타입으로 날라옴, 확인해야함

	return !os.IsNotExist(err) //존재하지않는지 체크 true면 없음 false면 있음
}
func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}
func persistKey(key *ecdsa.PrivateKey) {

	bytes, err := x509.MarshalECPrivateKey(key) //key를 byte로바꾸고 hexadecimal(string)으로 바꿀필요없다. 그냥 byte를 파일로 쓰면된다.
	utils.HandleErr(err)
	err = os.WriteFile(walletFileName, (bytes), 0644) //0644는 읽기와 쓰기허용
	utils.HandleErr(err)

}
func resotreKey() (key *ecdsa.PrivateKey) { //named return임 이미 출력할것을 정의함
	keyAsByte, err := os.ReadFile(walletFileName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsByte)
	utils.HandleErr(err)

	return

}
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		//has a wallet already?
		if hasWalletFile() {
			w.privateKey = resotreKey()

		} else {
			privateK := createPrivateKey()
			persistKey(privateK)
			w.privateKey = privateK
		}
		//yes-> resotre
		w.Address = aFromK(w.privateKey)
		//no->create privateKey

	}
	return w
}

// func Start() {
// 	// privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) //비공개키는 지갑에 저장되어있어야한다.
// 	// utils.HandleErr(err)
// 	// keyAsByte, err := x509.MarshalECPrivateKey(privateKey) // private키를 byte로 변환
// 	// utils.HandleErr(err)
// 	// fmt.Printf("%x\n", keyAsByte)
// 	// bytehashedMessage, err := hex.DecodeString(HashedMessage)
// 	// utils.HandleErr(err)
// 	// r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytehashedMessage)
// 	// utils.HandleErr(err)
// 	// signature := append(r.Bytes(), s.Bytes()...)
// 	// fmt.Printf("%x\n", signature)
// 	// ok := ecdsa.Verify(&privateKey.PublicKey, bytehashedMessage, r, s)
// 	// fmt.Println(ok)

// 	// bytePrivateKey, err := hex.DecodeString(privateKey)
// 	// utils.HandleErr(err)
// 	// // x509.ParseECPrivateKey([]byte(privateKey)) //16짜리 byte인지 모름 인코딩방식이 맞는지확인해야함
// 	// restoredKey, err := x509.ParseECPrivateKey(bytePrivateKey)
// 	// utils.HandleErr(err)
// 	// fmt.Println(restoredKey) //아직 비공개키(퍼블릭키를 포함한)를 가지고있지만 sign을 검증할순없다.
// 	// //복구되길원했던 비공개키와 같은게 맞는지 모름
// 	// //서명구조 [[32][32]]
// 	// signatureByte, err := hex.DecodeString(signature)
// 	// utils.HandleErr(err)
// 	// rByte := signatureByte[:int(len(signatureByte)/2)]
// 	// sByte := signatureByte[int(len(signatureByte)/2):]
// 	// fmt.Println(rByte)
// 	// fmt.Println(sByte)
// 	// var bigR, bigS = big.Int{}, big.Int{}
// 	// bigR.SetBytes(rByte)
// 	// bigS.SetBytes(sByte)
// 	// fmt.Println(bigR)
// 	// fmt.Println(bigS)
// 	// messageByte, err := hex.DecodeString(HashedMessage)
// 	// utils.HandleErr(err)
// 	// ok := ecdsa.Verify(&restoredKey.PublicKey, messageByte, &bigR, &bigS)
// 	// fmt.Println(ok)
// }
