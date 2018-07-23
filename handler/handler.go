package handler

import (
	"io"
	"net/http"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"go-blockchain/model"
	"go-blockchain/helper"
)

var Blockchain []model.Block

//Handler for get block
func HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

//Handler for writing a block
func HandleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m model.Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		helper.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := helper.GenerateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		helper.RespondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if helper.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		helper.ReplaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	helper.RespondWithJSON(w, r, http.StatusCreated, newBlock)

}