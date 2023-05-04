package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	port       = ":3001"
	binanceURL = "https://api.binance.com/api/v3/ticker/price"
)

type PairRequest struct {
	Pairs []string `json:"pairs" binding:"required"`
}

type priceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type binanceError struct {
	Message string `json:"msg"`
}

type binanceResponse []priceData

func (r binanceResponse) toMap() (map[string]float64, error) {
	m := make(map[string]float64)
	for _, v := range r {
		//parse to float to remove trailing zeros and check if valid value
		f, err := strconv.ParseFloat(v.Price, 64)
		if err != nil {
			return nil, err
		}

		m[v.Symbol] = f
	}
	return m, nil
}

func Start() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	v1.GET("/rates", pairInfoByGet)
	v1.POST("/rates", pairInfoByPost)

	r.Run(port)
}

func loadFromBinance(pairs []string) (binanceResponse, error) {
	for i, pair := range pairs {
		pairs[i] = strings.Replace(pair, "-", "", 1)
	}

	url := fmt.Sprintf("%s?symbols=[\"%s\"]", binanceURL, strings.Join(pairs, "\",\""))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var response binanceError
		if err = json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		return nil, errors.New(response.Message)
	}

	var response binanceResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func logAndAbort(ctx *gin.Context, code int, err error) {
	log.Println(err)
	ctx.AbortWithError(code, err)
}

func pairInfoByGet(ctx *gin.Context) {
	value, ok := ctx.GetQuery("pairs")
	if !ok {
		logAndAbort(ctx, http.StatusBadRequest, errors.New("parameter pairs is required"))
		return
	}

	pairs := strings.Split(value, ",")
	response, err := loadFromBinance(pairs)
	if err != nil {
		logAndAbort(ctx, http.StatusInternalServerError, err)
		return
	}

	r, err := response.toMap()
	if err != nil {
		logAndAbort(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(200, r)
}

func pairInfoByPost(ctx *gin.Context) {
	var request PairRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logAndAbort(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := loadFromBinance(request.Pairs)
	if err != nil {
		logAndAbort(ctx, http.StatusInternalServerError, err)
		return
	}

	r, err := response.toMap()
	if err != nil {
		logAndAbort(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(200, r)
}
