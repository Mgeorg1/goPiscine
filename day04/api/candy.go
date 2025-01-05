package api

/*
#include "cow.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BuyCandyRequest struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type BuyCandyResponse struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

const (
	uknownCandyType = -1
	CE              = iota
	AA
	NT
	DE
	YR
)

var (
	candyTypeToString = map[int]string{
		CE: "CE",
		AA: "AA",
		NT: "NT",
		DE: "DE",
		YR: "YR",
	}
	stringToCandyType = map[string]int{
		"CE": CE,
		"AA": AA,
		"NT": NT,
		"DE": DE,
		"YR": YR,
	}

	candyPrices = map[int]int{
		CE: 10,
		AA: 15,
		NT: 17,
		DE: 21,
		YR: 23,
	}
)

type thanksInterface interface {
	thanks(req *BuyCandyRequest) string
}

type simpleThanks struct{}

func (s simpleThanks) thanks(req *BuyCandyRequest) string {
	return fmt.Sprintf("Thank you for buying %d %s candy!", req.CandyCount, req.CandyType)
}

type cow struct{}

func (c cow) thanks(req *BuyCandyRequest) string {
	thanks := fmt.Sprintf("Thank you for buying %d %s candy!", req.CandyCount, req.CandyType)
	return fmt.Sprintf("Thank you for buying %d %s candy! %s", req.CandyCount, req.CandyType, C.GoString(C.ask_cow(C.CString(thanks))))
}

func handleBuyCandyRequest(w http.ResponseWriter, r *http.Request, thanks thanksInterface) {
	var req BuyCandyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, err)
		return
	}

	candyType, ok := stringToCandyType[req.CandyType]
	if !ok {
		errorResponse(w, fmt.Errorf("unknown candy type: %s", req.CandyType))
		return
	}
	candyPrice, ok := candyPrices[candyType]
	if !ok {
		errorResponse(w, fmt.Errorf("unknown candy type: %s", req.CandyType))
		return
	}
	if req.Money < 0 {
		errorResponse(w, fmt.Errorf("money should be positive"))
		return
	}
	if req.CandyCount < 0 {
		errorResponse(w, fmt.Errorf("candy count should be positive"))
		return
	}
	candyPriceSum := candyPrice * req.CandyCount
	if candyPriceSum <= req.Money {
		change := req.Money - candyPriceSum
		resp := BuyCandyResponse{
			Change: change,
			Thanks: thanks.thanks(&req),
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if candyPriceSum > req.Money {
		errorResponse(w, fmt.Errorf("you need %d more money", candyPriceSum-req.Money))
		return
	}
}

func BuyCandyHandler(w http.ResponseWriter, r *http.Request) {
	handleBuyCandyRequest(w, r, simpleThanks{})
}

func BuyCandyHandlerCow(w http.ResponseWriter, r *http.Request) {
	handleBuyCandyRequest(w, r, cow{})
}
