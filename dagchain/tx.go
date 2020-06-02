package dagchain

import (
	"Jight/config"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"strconv"

	"math/big"
	"time"
)

var txNum int = 0

var TxNumEth int64 = 0

type GeneralTx interface {
	FetchLatestValidateNum() [2]int
	FetchValidateNum() []int
	CheckVerification() bool
	SetVerification(v bool)
	Verify() bool
	FetchNumber() int
	FetchHash() [32]byte
	FetchSenderNum() int
	FetchCitedCount() int
	AddCitedCount()
	DecCitedCount()
	Serialize()  []byte
}

type GeneralTxEth interface {
	FetchLatestValidateNum() [2]int64
	FetchValidateNum() []int64
	CheckVerification() bool
	SetVerification(v bool)
	Verify() bool
	FetchNumber() int64
	FetchHash() [64]byte
	FetchSenderNum() int64
	FetchCitedCount() int
	AddCitedCount()
	DecCitedCount()
	Serialize()  []byte
}

// Transaction struct
type Transaction struct {
	Number int 			// just to make a transaction more readable
	ValidateNum [2]int		// make validate ref more readable
	Hash [32]byte			// Hash of the tx
	Parent [32]byte		// parent hash ref
	Validate [2][32]byte	// validate hash ref
	Income [32]byte		// income hash ref
	Sender [34]byte		// sender public key not address
	SenderNum int		// make the sender more readable
	Value int		// the money of transformation
	Receiver [34]byte		// receiver address
	Nonce int		// nonce of mining
	Timestamp int64		// Timestamp of tx
	Signature [64]byte	// Signature of tx
	Verification bool // if the tx verify the sample transaction
	CitedCount int	// how many times it is cited by other transactions or TUs
}

// Transaction struct
type TransactionEth struct {
	Number int64 			// just to make a transaction more readable
	ValidateNum [2]int64		// make validate ref more readable
	Hash [64]byte			// Hash of the tx
	Parent [64]byte		// parent hash ref
	Validate [2][64]byte	// validate hash ref
	Income [64]byte		// income hash ref
	Sender [40]byte		// sender public key not address
	SenderNum int64		// make the sender more readable
	Value int		// the money of transformation
	Receiver [40]byte		// receiver address
	Nonce int		// nonce of mining
	Timestamp int64		// Timestamp of tx
	Signature [64]byte	// Signature of tx
	Verification bool // if the tx verify the sample transaction
	CitedCount int	// how many times it is cited by other transactions or TUs
}

// struct PureTx is the slim version of Transaction with no extra fields (e.g., more readable)
// struct PureTx is used to be serialized to calculate the storage size
type PureTx struct {
	Parent [32]byte
	Validate [2][32]byte
	Income [32]byte
	Sender [34]byte
	Value int
	Receiver [34]byte
	Nonce int
	Timestamp int64
	Signature [64]byte
}

type PureTxEth struct {
	Parent [64]byte
	Validate [2][64]byte
	Income [64]byte
	Sender [40]byte
	Value int
	Receiver [40]byte
	Nonce int
	Timestamp int64
	Signature [64]byte
}

type TxContent struct {
	Receiver [34]byte
	Value int
	Timestamp int64
	Income	[32]byte
	Nonce int
}

type TxContentEth struct {
	Receiver [40]byte
	Value int
	Timestamp int64
	Income	[64]byte
	Nonce int
}

type TxContentList []TxContent

type TxContentEthList []TxContentEth

// TU (Transaction Union) struct
type TU struct {
	Number int
	ParNum [2]int
	ValidateNum []int
	SenderNum int // the number of the sender
	TCList []TxContent
	Validate [][32]byte
	Signature [64]byte
	CitedCount int		// how many times it is cited by other transactions or TUs
}

// TU (Transaction Union) struct
type TUEth struct {
	Number int64
	ParNum [2]int64
	ValidateNum []int64
	SenderNum int64 // the number of the sender
	TCList []TxContentEth
	Validate [][64]byte
	Signature [64]byte
	CitedCount int		// how many times it is cited by other transactions or TUs
}

// struct PureTU is the slim version of Transaction Union with no extra fields (e.g., more readable)
// struct PureTU is used to be serialized to calculate the storage size
type PureTU struct {
	Parent [2][32]byte
	Validate [][32]byte
	Sender [34]byte
	//TCList []TxContent
	Signature [64]byte
}

type PureTUEth struct {
	Parent [2][64]byte
	Validate [][64]byte
	Sender [40]byte
	//TCList []TxContent
	Signature [64]byte
}

var GXs = make(map[int]GeneralTx)

var GXsEth = make(map[int64]GeneralTxEth)

func PrintGXs() string {
	var returnString string
	returnString = returnString + "\n"
	for _, g := range GXs {
		returnString = returnString + "GX number: " +  strconv.Itoa(g.FetchNumber())+ " "
		returnString = returnString + "send from: " + strconv.Itoa(g.FetchSenderNum()) + " "
		returnString = returnString + "cited by: " + strconv.Itoa(g.FetchCitedCount()) + " "
		returnString = returnString + "ValidateNum: "
		for _, vn := range g.FetchValidateNum() {
			returnString = returnString + strconv.Itoa(vn) + " "
		}
		switch g.(type) {
		case *Transaction:
			returnString += "type: transaction"
		case *TU:
			returnString += "type: TU"
		}
		returnString = returnString + "\n"
	}
	/*for i:=config.GENESIS_ADDR_COUNT; i< len(GXs); i++ {
		returnString = returnString + "GX number: " +  strconv.Itoa(GXs[i].FetchNumber())+ " "
		returnString = returnString + "send from: " + strconv.Itoa(GXs[i].FetchSenderNum()) + " "
		returnString = returnString + "cited by: " + strconv.Itoa(GXs[i].FetchCitedCount()) + " "
		returnString = returnString + "ValidateNum: "
		for _, vn := range GXs[i].FetchValidateNum() {
			returnString = returnString + strconv.Itoa(vn) + " "
		}
		switch GXs[i].(type) {
		case *Transaction:
			returnString += "type: transaction"
		case *TU:
			returnString += "type: TU"
		}
		returnString = returnString + "\n"
	}*/
	return returnString
}

// TODO mining
func (tx Transaction) Pow() int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-config.POW_TARGET_BITS))

	/*var hashInt big.Int
	var hash [32]byte
	var nonce uint32 = 0

	for nonce < config.MAX_NONCE {
		data := getRawData(tx, nonce)

		hash = sha256.Sum256(data)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce*/
	return 0
}

func (tx TransactionEth) Pow() int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-config.POW_TARGET_BITS))
	return 0
}

// return the raw data of tx need in pow
func getRawData(tx Transaction, nonce uint32) []byte {
	var validate [][]byte
	for _, v := range tx.Validate {
		validate = append(validate, v[:])
	}
	return bytes.Join(
		[][]byte{
			IntToBytes(int64(tx.Number)),
			tx.Parent[:],
			bytes.Join(validate, []byte{}),
			tx.Income[:],
			tx.Sender[:],
			IntToBytes(int64(tx.Value)),
			tx.Receiver[:],
			IntToBytes(int64(nonce)),
			IntToBytes(tx.Timestamp),
			tx.Signature[:],
			BoolToBytes(tx.Verification),
		},
		[]byte{})
}

func BoolToBytes(b bool) []byte {
	if b {
		return []byte("1")
	} else {
		return []byte("0")
	}
}

func IntToBytes(data int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// calculate hash of tx
func (tx Transaction) HashTx() [32]byte {

	txCopy := tx
	txCopy.Hash = [32]byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	return hash
}

func (tx TransactionEth) HashTx() [64]byte {

	txCopy := tx
	txCopy.Hash = [64]byte{}

	hash := sha256.Sum256(txCopy.Serialize())

	var result [64]byte
	copy(result[:32], hash[:])
	copy(result[32:], hash[:])

	return result
}

// Sign a tx with private key
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsGenesisTx() {
		return
	}

	// TODO not complete copy
	txCopy := *tx

	txCopy.Hash = [32]byte{}
	txCopy.Signature = [64]byte{}

	//dataToSign := fmt.Sprintf("%x\n", txCopy)

	dataToSign := ""

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)

	var sign64 [64]byte
	copy(sign64[:], signature)

	tx.Signature = sign64
}

func (tx *TransactionEth) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsGenesisTx() {
		return
	}

	// TODO not complete copy
	txCopy := *tx

	txCopy.Hash = [64]byte{}
	txCopy.Signature = [64]byte{}

	//dataToSign := fmt.Sprintf("%x\n", txCopy)

	dataToSign := ""

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)

	var sign64 [64]byte
	copy(sign64[:], signature)

	tx.Signature = sign64
}

// verify a signature of tx
func (tx *Transaction) Verify() bool {
	/*
		// It just return true in the test
		if tx.IsGenesisTx() {
			return true
		}

		// TODO not complete copy
		txCopy := *tx

		r := big.Int{}
		s := big.Int{}
		sigLen := len(tx.Signature)
		r.SetBytes(tx.Signature[:(sigLen/2)])
		s.SetBytes(tx.Signature[(sigLen/2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(tx.Sender)
		x.SetBytes(tx.Sender[:(keyLen/2)])
		y.SetBytes(tx.Sender[(keyLen/2):])

		txCopy.Hash = nil
		txCopy.Signature = nil

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		curve := elliptic.P256()

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}*/

	return true
}

func (tx *TransactionEth) Verify() bool {
	return true
}

// serialize tx
// in fact, the pure version of tx will be serialized
func (tx Transaction) Serialize() []byte {
	var encode bytes.Buffer

	// construct the pure version of tx
	pureTx := PureTx{
		Parent: tx.Parent,
		Validate: tx.Validate,
		Income: tx.Income,
		Sender: tx.Sender,
		Value: tx.Value,
		Receiver: tx.Receiver,
		Nonce: tx.Nonce,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(pureTx)

	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	return encode.Bytes()
}

func (tx TransactionEth) Serialize() []byte {
	var encode bytes.Buffer

	// construct the pure version of tx
	pureTx := PureTxEth{
		Parent: tx.Parent,
		Validate: tx.Validate,
		Income: tx.Income,
		Sender: tx.Sender,
		Value: tx.Value,
		Receiver: tx.Receiver,
		Nonce: tx.Nonce,
		Timestamp: tx.Timestamp,
		Signature: tx.Signature,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(pureTx)

	if err != nil {
		log.Panic("txEth encode fail:", err)
	}

	return encode.Bytes()
}

// serialize TU
func (tu TU) Serialize() []byte {
	var encode bytes.Buffer

	// construct the pure version of tx
	// since the serialize() function is mainly used to evaluate the storage size
	// as a result, the fields in the pureTU is filled in the empty bytes
	pureTU := PureTU{
		Parent: [2][32]byte{config.MOCK_TX, config.MOCK_TX},
		Validate: tu.Validate,
		Sender: config.MOCK_ACCOUNT,
		//TCList: tu.TCList,
		Signature: tu.Signature,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(pureTU)

	if err != nil {
		log.Panic("TU encode fail:", err)
	}

	return encode.Bytes()
}

func (tu TUEth) Serialize() []byte {
	var encode bytes.Buffer

	// construct the pure version of tx
	// since the serialize() function is mainly used to evaluate the storage size
	// as a result, the fields in the pureTU is filled in the empty bytes
	pureTU := PureTUEth{
		Parent: [2][64]byte{config.MOCK_TX_Eth, config.MOCK_TX_Eth},
		Validate: tu.Validate,
		Sender: config.MOCK_ACCOUNT_ETH,
		//TCList: tu.TCList,
		Signature: tu.Signature,
	}

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(pureTU)

	if err != nil {
		log.Panic("TUEth encode fail:", err)
	}

	return encode.Bytes()
}

// serialize TxContentList
func (tcl TxContentList) Serialize() []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tcl)

	if err != nil {
		log.Panic("TxContent encode fail:", err)
	}

	return encode.Bytes()
}

func (tcl TxContentEthList) Serialize() []byte {
	var encode bytes.Buffer
	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tcl)

	if err != nil {
		log.Panic("TxContent encode fail:", err)
	}

	return encode.Bytes()
}

// Deserialize tx
func DeserializeTx(data []byte) Transaction {
	var tx Transaction

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tx)
	if err != nil {
		log.Panic("tx decode fail:", err)
	}

	return tx
}

func DeserializeTxEth(data []byte) TransactionEth {
	var tx TransactionEth

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tx)
	if err != nil {
		log.Panic("txEth decode fail:", err)
	}

	return tx
}

// serialize TxContentList
func DeserializeTCL(data []byte) TxContentList {
	var tcl TxContentList

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tcl)
	if err != nil {
		log.Panic("TxContentList decode fail:", err)
	}
	return tcl
}

func DeserializeTCLEth(data []byte) TxContentEthList {
	var tcl TxContentEthList

	decode := gob.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&tcl)
	if err != nil {
		log.Panic("TxContentEthList decode fail:", err)
	}
	return tcl
}

func (tx *Transaction) CheckVerification() bool {
	return tx.Verification
}

func (tx *TransactionEth) CheckVerification() bool {
	return tx.Verification
}

func (tx *Transaction) FetchLatestValidateNum() [2]int {
	return tx.ValidateNum
}

func (tx *TransactionEth) FetchLatestValidateNum() [2]int64 {
	return tx.ValidateNum
}

func (tx *Transaction) FetchValidateNum() []int {
	return tx.ValidateNum[:]
}

func (tx *TransactionEth) FetchValidateNum() []int64 {
	return tx.ValidateNum[:]
}

func (tx *Transaction) SetVerification(v bool) {
	tx.Verification = v
}

func (tx *TransactionEth) SetVerification(v bool) {
	tx.Verification = v
}

func (tx *Transaction) FetchNumber() int {
	return  tx.Number
}

func (tx *TransactionEth) FetchNumber() int64 {
	return  tx.Number
}

func (tx *Transaction) FetchHash() [32]byte {
	return  tx.Hash
}

func (tx *TransactionEth) FetchHash() [64]byte {
	return  tx.Hash
}

func (tx *Transaction) FetchSenderNum() int {
	return tx.SenderNum
}

func (tx *TransactionEth) FetchSenderNum() int64 {
	return tx.SenderNum
}

func (tx *Transaction) FetchCitedCount() int {
	return tx.CitedCount
}

func (tx *TransactionEth) FetchCitedCount() int {
	return tx.CitedCount
}

func (tx *Transaction) AddCitedCount() {
	tx.CitedCount++
}

func (tx *TransactionEth) AddCitedCount() {
	tx.CitedCount++
}

func (tx *Transaction) DecCitedCount() {
	tx.CitedCount--
}

func (tx *TransactionEth) DecCitedCount() {
	tx.CitedCount--
}

// create a new tx including mining
func NewTx(validateNum [2]int, par [32]byte, validate [2][32]byte, income [32]byte, sender string, senderNum int, value int, receiver string, verification bool) *Transaction {
	senderBytes := []byte(sender)
	var senderBytes34 [34]byte
	copy(senderBytes34[:], senderBytes)

	receiverBytes := []byte(receiver)
	var receiverBytes34 [34]byte
	copy(receiverBytes34[:], receiverBytes)


	tx := &Transaction{txNum, validateNum, [32]byte{}, par, validate,
		income,senderBytes34, senderNum, value, receiverBytes34,
		0, time.Now().Unix(), [64]byte{}, verification, 0}
	txNum ++
	return tx
}

// create a new tx including mining
func NewTxEth(validateNum [2]int64, par [64]byte, validate [2][64]byte, income [64]byte, sender string, senderNum int64,
	value int, receiver string, verification bool) *TransactionEth {
	senderBytes := []byte(sender)[2:]
	var senderBytes40 [40]byte
	copy(senderBytes40[:], senderBytes)

	receiverBytes := []byte(receiver)[2:]
	var receiverBytes40 [40]byte
	copy(receiverBytes40[:], receiverBytes)

	tx := &TransactionEth{TxNumEth, validateNum, [64]byte{}, par, validate,
		income,senderBytes40, senderNum, value, receiverBytes40,
		0, time.Now().Unix(), [64]byte{}, verification, 0}
	TxNumEth ++
	return tx
}

// create a genesis tx
func NewGenesisTx(value int, receiver []byte) *Transaction {

	return NewTx([2]int{-1, -1}, [32]byte{}, [2][32]byte{}, [32]byte{}, "", -1, value, "", false)

}

func NewGenesisTxEth(value int, num int64) *TransactionEth {

	return &TransactionEth{num, [2]int64{-1, -1}, [64]byte{}, [64]byte{}, [2][64]byte{},
		[64]byte{},[40]byte{}, -1, value, [40]byte{},
		0, time.Now().Unix(), [64]byte{}, false, 0}

}

// determine if a tx is a genesis tx
func (tx Transaction) IsGenesisTx() bool {
	return 0 == len(tx.Parent) && 0 == len(tx.Validate) && 0 == len(tx.Sender)
}

func (tx TransactionEth) IsGenesisTx() bool {
	return 0 == len(tx.Parent) && 0 == len(tx.Validate) && 0 == len(tx.Sender)
}

// create a new Transaction Union
func NewTU(validateNum [2]int, acc *Account, income [32]byte, value int, receiver string,) (*TU, *TxContent) {
	log.Println("The latest validateNum of a TU: ", validateNum)

	receiverBytes := []byte(receiver)
	var receiverBytes34 [34]byte
	copy(receiverBytes34[:], receiverBytes)

	newTU := &TU{
		Number:txNum,
		ParNum: [2]int{acc.LastIdNo, 0},
		SenderNum: acc.AccountNo,
		Signature: [64]byte{},
		CitedCount: 0,
	}

	// construct the new TCList
	TCList := make([]TxContent, config.MERGE_PERIOD)
	i:=0
	for i=0; i< config.MERGE_PERIOD-1; i++ {
		tmpTx, ok := GXs[acc.WithoutMergeIds[i]].(*Transaction)
		if !ok {
			panic("TmpTx is not a transaction")
		}
		TCList[i] = TxContent{
			Receiver: tmpTx.Receiver,
			Value: tmpTx.Value,
			Income: tmpTx.Income,
			Timestamp: tmpTx.Timestamp,
			Nonce: tmpTx.Nonce,
		}
	}

	lastTC := TxContent{
		Receiver: receiverBytes34,
		Value: value,
		Income: income,
		Timestamp: time.Now().Unix(),
		Nonce: 0,
	}
	TCList[i] = lastTC

	/*var TCListTotal []TxContent
	if acc.LatestTU != nil {
		TCListTotal = acc.LatestTU.TCList
	}

	for i=0; i < config.MERGE_PERIOD; i++ {
		TCListTotal = append(TCListTotal, TCList[config.MERGE_PERIOD-1-i])
	}*/

	var TCListTotal []TxContent

	newTU.TCList = TCListTotal


	/////////////////////////////////////////////////////////////////////////////////
	// construct the redundant validate references
	/////////////////////////////////////////////////////////////////////////////////
	// 1. merge the 10 previous unmerged transactions
	var periodValidateNumOriginal map[int]bool = make(map[int]bool)
	// 1.1 deal with the previous 9 transactions
	for _, id := range acc.WithoutMergeIds {
		//log.Println("id: ", id)
		vn := GXs[id].FetchLatestValidateNum()
		//log.Println("vn:", vn)
		validateAccNum0 := GXs[vn[0]].FetchSenderNum()
		validateAccNum1 := GXs[vn[1]].FetchSenderNum()
		for n ,_ := range periodValidateNumOriginal {
			nAccNum := GXs[n].FetchSenderNum()
			//log.Println("n: ", n)
			//log.Println("nAccNum: ", nAccNum)
			//log.Println("Before dec CitedCount: ", GXs[n].FetchCitedCount())
			if periodValidateNumOriginal[n] == true {
				if nAccNum == validateAccNum0 || nAccNum == validateAccNum1 {
					periodValidateNumOriginal[n] = false
					GXs[n].DecCitedCount()
				}
			}
			//log.Println("After dec CitedCount: ", GXs[n].FetchCitedCount())
		}
		periodValidateNumOriginal[vn[0]]=true
		periodValidateNumOriginal[vn[1]]=true
		//log.Println("periodValidateNumOriginal: ", periodValidateNumOriginal)
	}
	// 1.2 deal with the 10th transaction
	GXs[validateNum[0]].AddCitedCount()
	GXs[validateNum[1]].AddCitedCount()
	validateAccNum0 := GXs[validateNum[0]].FetchSenderNum()
	validateAccNum1 := GXs[validateNum[1]].FetchSenderNum()

	for n ,_ := range periodValidateNumOriginal {
		nAccNum := GXs[n].FetchSenderNum()
		//log.Println("n: ", n)
		//log.Println("nAccNum: ", nAccNum)
		//log.Println("Before dec CitedCount: ", GXs[n].FetchCitedCount())
		if periodValidateNumOriginal[n] == true {
			if nAccNum == validateAccNum0 || nAccNum == validateAccNum1 {
				periodValidateNumOriginal[n] = false
				GXs[n].DecCitedCount()
			}
		}
		//log.Println("After dec CitedCount: ", GXs[n].FetchCitedCount())
	}
	periodValidateNumOriginal[validateNum[0]]=true
	periodValidateNumOriginal[validateNum[1]]=true
	//log.Println("periodValidateNumOriginal: ", periodValidateNumOriginal)
	var periodValidateNum []int
	for n, b := range periodValidateNumOriginal {
		if b {
			periodValidateNum = append(periodValidateNum, n)
		}
	}

	// 2. merge the periodValidateNum with the past TU
	var validateNumOld []int
	var validateNumNew []int
	var toDeleteNum map[int]bool = make(map[int]bool)
	if acc.LatestTU == nil {
		validateNumNew = periodValidateNum
	} else {
		var vanMap map[int]bool = make(map[int]bool)
		for _, pvn := range periodValidateNum {
			vanMap[GXs[pvn].FetchSenderNum()]=true
		}
		validateNumOld = acc.LatestTU.ValidateNum
		for _, vn := range validateNumOld {
			vnAccNum := GXs[vn].FetchSenderNum()
			//log.Println("vn: ", vn)
			//log.Println("vnAccNum: ", vnAccNum)
			//log.Println("Before dec CitedCount: ", GXs[vn].FetchCitedCount())
			if !toDeleteNum[vn] {
				if vanMap[vnAccNum] {
					toDeleteNum[vn] = true
					GXs[vn].DecCitedCount()
				}
			}
			//log.Println("After dec CitedCount: ", GXs[vn].FetchCitedCount())
		}

		for _, vn := range validateNumOld {
			if !toDeleteNum[vn] {
				validateNumNew = append(validateNumNew, vn)
			}
		}
		for _, pvn := range periodValidateNum {
			validateNumNew = append(validateNumNew, pvn)
		}

	}

	newTU.ValidateNum = validateNumNew

	acc.LastIdNo = txNum
	acc.WithoutMergeIds = [config.MERGE_PERIOD-1]int{}

	txNum ++
	return newTU, &lastTC
}

// create a new Transaction Union
func NewTUEth(validateNum [2]int64, acc *AccountEth, income [64]byte, value int, receiver string) (*TUEth, *TxContentEth) {
	log.Println("The latest validateNum of a TU: ", validateNum)

	receiverBytes := []byte(receiver)[2:]
	var receiverBytes40 [40]byte
	copy(receiverBytes40[:], receiverBytes)

	newTU := &TUEth{
		Number: TxNumEth,
		ParNum: [2]int64{acc.LastIdNo, 0},
		SenderNum: acc.AccountNo,
		Signature: [64]byte{},
		CitedCount: 0,
	}

	// construct the new TCList
	TCList := make([]TxContentEth, config.MERGE_PERIOD)
	i:=0
	for i=0; i< config.MERGE_PERIOD-1; i++ {
		tmpTx, ok := GXsEth[acc.WithoutMergeIds[i]].(*TransactionEth)
		if !ok {
			panic("TmpTx is not a transaction")
		}
		TCList[i] = TxContentEth{
			Receiver: tmpTx.Receiver,
			Value: tmpTx.Value,
			Income: tmpTx.Income,
			Timestamp: tmpTx.Timestamp,
			Nonce: tmpTx.Nonce,
		}
	}

	lastTC := TxContentEth{
		Receiver: receiverBytes40,
		Value: value,
		Income: income,
		Timestamp: time.Now().Unix(),
		Nonce: 0,
	}
	TCList[i] = lastTC

	/*var TCListTotal []TxContent
	if acc.LatestTU != nil {
		TCListTotal = acc.LatestTU.TCList
	}

	for i=0; i < config.MERGE_PERIOD; i++ {
		TCListTotal = append(TCListTotal, TCList[config.MERGE_PERIOD-1-i])
	}*/

	var TCListTotal []TxContentEth

	newTU.TCList = TCListTotal


	/////////////////////////////////////////////////////////////////////////////////
	// construct the redundant validate references
	/////////////////////////////////////////////////////////////////////////////////
	// 1. merge the 10 previous unmerged transactions
	var periodValidateNumOriginal = make(map[int64]bool)
	// 1.1 deal with the previous 9 transactions
	for _, id := range acc.WithoutMergeIds {
		log.Println("id: ", id)
		vn := GXsEth[id].FetchLatestValidateNum()
		log.Println("vn:", vn)
		validateAccNum0 := GXsEth[vn[0]].FetchSenderNum()
		validateAccNum1 := GXsEth[vn[1]].FetchSenderNum()
		for n ,_ := range periodValidateNumOriginal {
			nAccNum := GXsEth[n].FetchSenderNum()
			//log.Println("n: ", n)
			//log.Println("nAccNum: ", nAccNum)
			//log.Println("Before dec CitedCount: ", GXsEth[n].FetchCitedCount())
			if periodValidateNumOriginal[n] == true {
				if nAccNum == validateAccNum0 || nAccNum == validateAccNum1 {
					periodValidateNumOriginal[n] = false
					GXsEth[n].DecCitedCount()
				}
			}
			//log.Println("After dec CitedCount: ", GXsEth[n].FetchCitedCount())
		}
		periodValidateNumOriginal[vn[0]]=true
		periodValidateNumOriginal[vn[1]]=true
		//log.Println("periodValidateNumOriginal: ", periodValidateNumOriginal)
	}
	// 1.2 deal with the 10th transaction
	GXsEth[validateNum[0]].AddCitedCount()
	GXsEth[validateNum[1]].AddCitedCount()
	validateAccNum0 := GXsEth[validateNum[0]].FetchSenderNum()
	validateAccNum1 := GXsEth[validateNum[1]].FetchSenderNum()

	for n ,_ := range periodValidateNumOriginal {
		nAccNum := GXsEth[n].FetchSenderNum()
		//log.Println("n: ", n)
		//log.Println("nAccNum: ", nAccNum)
		//log.Println("Before dec CitedCount: ", GXsEth[n].FetchCitedCount())
		if periodValidateNumOriginal[n] == true {
			if nAccNum == validateAccNum0 || nAccNum == validateAccNum1 {
				periodValidateNumOriginal[n] = false
				GXsEth[n].DecCitedCount()
			}
		}
		//log.Println("After dec CitedCount: ", GXsEth[n].FetchCitedCount())
	}
	periodValidateNumOriginal[validateNum[0]]=true
	periodValidateNumOriginal[validateNum[1]]=true
	//log.Println("periodValidateNumOriginal: ", periodValidateNumOriginal)

	// Deal with a special case: when a tx validate the txs issued by the same person, e.g., (A->B) --validate--> (A->C)
	for n ,b := range periodValidateNumOriginal {
		if (GXsEth[n].FetchSenderNum() == acc.AccountNo) && b {
			periodValidateNumOriginal[n] = false
			GXsEth[n].DecCitedCount()
		}
	}

	var periodValidateNum []int64
	for n, b := range periodValidateNumOriginal {
		if b {
			periodValidateNum = append(periodValidateNum, n)
		}
	}

	// 2. merge the periodValidateNum with the past TU
	var validateNumOld []int64
	var validateNumNew []int64
	var toDeleteNum = make(map[int64]bool)
	if acc.LatestTU == nil {
		validateNumNew = periodValidateNum
	} else {
		var vanMap = make(map[int64]bool)
		for _, pvn := range periodValidateNum {
			vanMap[GXsEth[pvn].FetchSenderNum()]=true
		}
		validateNumOld = acc.LatestTU.ValidateNum
		for _, vn := range validateNumOld {
			vnAccNum := GXsEth[vn].FetchSenderNum()
			//log.Println("vn: ", vn)
			//log.Println("vnAccNum: ", vnAccNum)
			//log.Println("Before dec CitedCount: ", GXsEth[vn].FetchCitedCount())
			if !toDeleteNum[vn] {
				if vanMap[vnAccNum] {
					toDeleteNum[vn] = true
					GXsEth[vn].DecCitedCount()
				}
			}
			//log.Println("After dec CitedCount: ", GXsEth[vn].FetchCitedCount())
		}

		for _, vn := range validateNumOld {
			if !toDeleteNum[vn] {
				validateNumNew = append(validateNumNew, vn)
			}
		}
		for _, pvn := range periodValidateNum {
			validateNumNew = append(validateNumNew, pvn)
		}

	}

	newTU.ValidateNum = validateNumNew

	acc.LastIdNo = TxNumEth
	acc.WithoutMergeIds = [config.MERGE_PERIOD-1]int64{}

	TxNumEth ++
	return newTU, &lastTC
}

func (tu * TU) FetchLatestValidateNum() [2]int {
	var latestValidateNum [2]int
	validateLen := len(tu.ValidateNum)
	latestValidateNum[0] = tu.ValidateNum[validateLen-1]
	latestValidateNum[1] = tu.ValidateNum[validateLen-2]
	return latestValidateNum
}

func (tu * TUEth) FetchLatestValidateNum() [2]int64 {
	var latestValidateNum [2]int64
	validateLen := len(tu.ValidateNum)
	latestValidateNum[0] = tu.ValidateNum[validateLen-1]
	latestValidateNum[1] = tu.ValidateNum[validateLen-2]
	return latestValidateNum
}

func (tu *TU) FetchValidateNum() []int {
	return tu.ValidateNum
}

func (tu *TUEth) FetchValidateNum() []int64 {
	return tu.ValidateNum
}

func (tu * TU) CheckVerification() bool {
	// to do
	return true
}

func (tu * TUEth) CheckVerification() bool {
	// to do
	return true
}

func (tu * TU) SetVerification(v bool) {
	// to do
}

func (tu * TUEth) SetVerification(v bool) {
	// to do
}

func (tu * TU) Verify() bool {
	// to do
	return true
}

func (tu * TUEth) Verify() bool {
	// to do
	return true
}

func (tu * TU) FetchNumber() int {
	return tu.Number
}

func (tu * TUEth) FetchNumber() int64 {
	return tu.Number
}

func (tu * TU) FetchHash() [32]byte {
	// to do
	return [32]byte{}
}

func (tu * TUEth) FetchHash() [64]byte {
	// to do
	return [64]byte{}
}

func (tu * TU) FetchSenderNum() int {
	return tu.SenderNum
}

func (tu * TUEth) FetchSenderNum() int64 {
	return tu.SenderNum
}

func (tu * TU) FetchCitedCount() int {
	return tu.CitedCount
}

func (tu * TUEth) FetchCitedCount() int {
	return tu.CitedCount
}

func (tu * TU) AddCitedCount() {
	tu.CitedCount ++
}

func (tu * TUEth) AddCitedCount() {
	tu.CitedCount ++
}

func (tu *TU) DecCitedCount() {
	tu.CitedCount--
}

func (tu *TUEth) DecCitedCount() {
	tu.CitedCount--
}

// prune the old tips, which is called when a new Transaction Union is created
func PruneOldTxs(dc *DagChain, acc *Account) error {
	log.Printf("Prune old transactions of account: %d\n", acc.AccountNo)
	txsByAccount := acc.WithoutPruneIds
	var toPruneNum = make(map[int]bool)
	for _, tn := range txsByAccount {
		//log.Printf("tx in txsByAccount: %d", tn)
		if GXs[tn].FetchCitedCount() == 0 {
			if dc.Tips[tn] == nil {
				toPruneNum[tn] = true
			}
		}
		/*if _, ok := GXs[tn].(*TU); ok {
			if tn != acc.LatestTU.Number {
				toPruneNum[tn] = true
			}
		}*/
	}
	// prune the tx from the WithoutPruneIds in the account
	var newTxsByAccount []int
	for _, txNum := range txsByAccount {
		if !toPruneNum[txNum] {
			newTxsByAccount = append(newTxsByAccount, txNum)
		}
	}
	acc.WithoutPruneIds = newTxsByAccount

	// prune the tx from the global slice []GXs
	var newGXs =make(map[int]GeneralTx)
	for _, gx := range GXs {
		if !toPruneNum[gx.FetchNumber()] {
			newGXs[gx.FetchNumber()] = gx
		}
	}
	GXs = newGXs

	// prune the txs from the DbMerging
	for n, b:= range toPruneNum {
		if b{
			if err := RemoveTx2(dc.DBMerging, n); err!=nil {
				return err
			}
		}
	}
	CompactDB(dc.DBMerging)
	CompactDB(dc.DB)
	return nil
}

func PruneOldTxsEth(dc *DagChainEth, acc *AccountEth) error {
	log.Printf("Prune old transactions of account: %d\n", acc.AccountNo)
	txsByAccount := acc.WithoutPruneIds
	var toPruneNum = make(map[int64]bool)
	for _, tn := range txsByAccount {
		//log.Printf("tx in txsByAccount: %d", tn)
		if GXsEth[tn].FetchCitedCount() == 0 {
			if dc.TipsEth[tn] == nil {
				toPruneNum[tn] = true
			}
		}
		/*if _, ok := GXs[tn].(*TU); ok {
			if tn != acc.LatestTU.Number {
				toPruneNum[tn] = true
			}
		}*/
	}
	// prune the tx from the WithoutPruneIds in the account
	var newTxsByAccount []int64
	for _, txNum := range txsByAccount {
		if !toPruneNum[txNum] {
			newTxsByAccount = append(newTxsByAccount, txNum)
		}
	}
	acc.WithoutPruneIds = newTxsByAccount

	// prune the tx from the global slice []GXs
	var newGXs =make(map[int64]GeneralTxEth)
	for _, gx := range GXsEth {
		if !toPruneNum[gx.FetchNumber()] {
			newGXs[gx.FetchNumber()] = gx
		}
	}
	GXsEth = newGXs

	// prune the txs from the DbMerging
	for n, b:= range toPruneNum {
		if b{
			if err := RemoveTxEth(dc.DBMerging, n); err!=nil {
				return err
			}
		}
	}
	CompactDB(dc.DBMerging)
	CompactDB(dc.DB)
	return nil
}