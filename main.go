package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Block in a Blockchain
type Block struct {
	Pos       int
	Data      ProductCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

// Product Info to be added to a block in a Blockchain
type Product struct {
	Id      string `json:"id"`
	Desc    string `json:"desc"`
	Comapny string `json:"company"`
	DoM     string `json:"date_manufacturing"`
}

// Actual payload in a Block
type ProductCheckout struct {
	ProductId    string `json:"product_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

// The complete blockchain which is a slice of Blocks
type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

// Create new product to be added in a Blockchain
func newProduct(w http.ResponseWriter, r *http.Request) {

	var prod Product
	if err := json.NewDecoder(r.Body).Decode((&prod)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("coould not create:%v", err)
		w.Write([]byte("could not create new Product"))
		return
	}

	h := md5.New()
	io.WriteString(h, prod.Desc+prod.DoM)
	prod.Id = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(prod, "", " ")
	if err != nil {
		w.WriteHeader((http.StatusInternalServerError))
		log.Printf("could not marshal payload:%v", err)
		w.Write([]byte("Could not save Product"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// Write the Block to Blockchain after decding json request data
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem ProductCheckout
	if err := json.NewDecoder(r.Body).Decode((&checkoutItem)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("coould not write block:%v", err)
		w.Write([]byte("could not write new Product"))
		return
	}

	BlockChain.AddBlock(checkoutItem)

}

// Add to Blockchain
func (bc *Blockchain) AddBlock(data ProductCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}

}

// Create a new Block in Blockchain
func CreateBlock(prevBlock *Block, data ProductCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.PrevHash = prevBlock.Hash
	block.TimeStamp = time.Now().String()
	block.generateHash()

	return block

}

// Generate SHA-256 Hash for a Block
func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))

}

// check Block validity
func validBlock(block *Block, prev *Block) bool {

	if prev.Hash != block.Hash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prev.Pos+1 != block.Pos {
		return false
	}

	return true

}

// Check hash validity
func (b *Block) validateHash(hash string) bool {

	b.generateHash()
	if b.Hash != hash {
		return false
	}

	return true

}

// Initialise Blockchain by create a Genesins block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}

}

// Create genesis Block
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, ProductCheckout{IsGenesis: true})
}

// Get entire Blockchain
func getBlockchain(w http.ResponseWriter, r *http.Request) {

	jbtyes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader((http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err)
		return
	}

	io.WriteString(w, string(jbtyes))

}

// Client which accepts Http requests for writing new product, writing Block and retreving entire Blockchain
func main() {

	BlockChain = NewBlockchain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newProduct).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("prev. has:%x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data:%v\n", string(bytes))
			fmt.Printf("Hash:%v\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))

}
