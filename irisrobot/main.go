package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/bank"
	irisInit "github.com/irisnet/irishub/init"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	cdc      = app.MakeCodec()
	nameFrom string
	nameTo   string
	addrFrom string
	addrTo   string
	seqFrom  int64
	seqTo    int64


	nameFrom2 string
	nameTo2   string
	addrFrom2 string
	addrTo2   string
	seqFrom2  int64
	seqTo2    int64

	nameFrom3 string
	nameTo3   string
	addrFrom3 string
	addrTo3  string
	seqFrom3  int64
	seqTo3    int64

	nameDel string
	addrDel string
	seqDel  int64
	valFrom string
	valTo   string

	nameVoter string
	addrVoter string
	seqVoter  int64
)

func VoteOnProposal(name string, voter string, option string) (resultTx ctypes.ResultBroadcastTxCommit) {
	//get Account
	acc := GetAccount(voter)
	accnum := acc.AccountNumber
	sequence := acc.Sequence
	chainID := viper.Get("chain")

	id := rand.Int()
	fmt.Println(acc.AccountNumber, acc.Sequence, chainID)

	jsonStr := []byte(fmt.Sprintf(`{
		"base_tx":{
			"name":"%s",
			"password":"12345678",
			"account_number":"%d",
			"sequence":"%d",
			"gas": "200000",
			"fee": "0.04iris",
			"chain_id":"%s",
            "gas_adjustment": "1.2"
       },
		"voter": "%s",
        "option": "%s"
	}`, name, accnum, sequence, chainID, voter, option))
	res, body, _ := Request("1317", "POST", fmt.Sprintf("/gov/proposals/%d/votes?async=true", id), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))
	if res.StatusCode == http.StatusOK {

		fmt.Println("success", resultTx.Hash)

	} else {

		fmt.Println("error", err)
	}

	return resultTx

}

func GetSequence(account string) int64 {

	seq := int64(-1)

	res, body, err := Request("1317", "GET", fmt.Sprintf("/auth/accounts/%s", account), nil)

	if res.StatusCode == http.StatusOK {

		var accInfo bank.BaseAccount

		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil {
			fmt.Println("error: ", err)
		} else {

			seq = accInfo.Sequence

		}

	}

	return seq
}

//get all the inforamtion
func GetAccount(account string) bank.BaseAccount {

	var accInfo bank.BaseAccount

	res, body, err := Request("1317", "GET", fmt.Sprintf("/auth/accounts/%s", account), nil)

	if res.StatusCode == http.StatusOK {

		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil {
			fmt.Println("error: ", err)
		}

	}

	return accInfo
}

func GetAccountByName(name string) keys.KeyOutput {
	var accInfo keys.KeyOutput

	res, body, err := Request("1317", "GET", fmt.Sprintf("/keys/%s", name), nil)

	if res.StatusCode == http.StatusOK {

		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil {
			fmt.Println("error: ", err)
		}

	}

	return accInfo
}

func Request(port, method, path string, payload []byte) (*http.Response, string, error) {
	var (
		res *http.Response
	)
	ip := viper.Get("IP")
	url := fmt.Sprintf("http://%v:%v%v", ip, port, path)

	fmt.Println(url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	res, err = http.DefaultClient.Do(req)

	output, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return res, string(output), err
}

func SendTransactionWithSequence(fromName string, toAddr string, seq int64) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTx) {

	// send
	coinbz := sdk.NewInt64Coin("iris", 1).String()
	port := "1317"
	chainID := viper.Get("chain")

	account := GetAccountByName(fromName)
	acc := GetAccount(account.Address)
	accnum := acc.AccountNumber

	//fmt.Println(coinbz," accum: ", acc.AccountNumber, acc.Sequence,chainID)

	jsonStr := []byte(fmt.Sprintf(`{
    "base_tx":{
            "name": "%s",
			"password": "12345678",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id": "%s",
            "gas_adjustment": "1.2"
       }, 
		"amount": "%s",
        "sender": "%s"
}`, fromName, accnum, seq, chainID, coinbz, account.Address))

	res, body, _ := Request(port, "POST", fmt.Sprintf("/bank/accounts/%s/transfers?async=true", toAddr), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))

	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	} else {
		fmt.Println("error", err)
	}

	return receiveAddr, resultTx
}

func DelegateTransactionWithSequence(fromName string, valAddr string, seq int64) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTx) {

	// send
	coinbz := sdk.NewInt64Coin("iris", 1).String()

	port := "1317"
	chainID := viper.Get("chain")

	account := GetAccountByName(fromName)
	acc := GetAccount(account.Address)
	accnum := acc.AccountNumber

	jsonStr := []byte(fmt.Sprintf(`{
    "base_tx":{
            "name": "%s",
			"password": "12345678",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id": "%s",
            "gas_adjustment": "1.2"
       }, 
		 "delegate": {
    "validator_addr": "%s",
    "delegation": "%s"
  }
}`, fromName, accnum, seq, chainID, valAddr, coinbz))

	res, body, _ := Request(port, "POST", fmt.Sprintf("/stake/delegators/%s/delegate?async=true", account.Address), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))

	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	} else {
		fmt.Println("error", err)
	}

	return receiveAddr, resultTx
}

func RedelegateTransactionWithSequence(fromName string, valFromAddr string, valToAddr string, seq int64) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTx) {

	port := "1317"
	chainID := viper.Get("chain")

	account := GetAccountByName(fromName)
	acc := GetAccount(account.Address)
	accnum := acc.AccountNumber

	jsonStr := []byte(fmt.Sprintf(`{
    "base_tx":{
            "name": "%s",
			"password": "12345678",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id": "%s",
            "gas_adjustment": "1.2"
       }, 
		  "redelegate": {
             "validator_src_addr": "%s",
             "validator_dst_addr": "%s",
             "shares": "1"
           }
}`, fromName, accnum, seq, chainID, valFromAddr, valToAddr))

	res, body, _ := Request(port, "POST", fmt.Sprintf("/stake/delegators/%s/redelegate?async=true", account.Address), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))

	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	} else {
		fmt.Println("error", err)
	}

	return receiveAddr, resultTx
}

func WidthdrawTransactionWithSequence(fromName string, seq int64) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTx) {

	port := "1317"
	chainID := viper.Get("chain")

	account := GetAccountByName(fromName)
	acc := GetAccount(account.Address)
	accnum := acc.AccountNumber

	jsonStr := []byte(fmt.Sprintf(`{
    "base_tx":{
            "name": "%s",
			"password": "12345678",
			"account_number": "%d",
			"sequence": "%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id": "%s",
            "gas_adjustment": "1.2"
       }, 
		  "withdraw_address": "%s"
}`, fromName, accnum, seq, chainID, account.Address))

	res, body, _ := Request(port, "POST", fmt.Sprintf("/distribution/%s/withdrawReward?async=true", account.Address), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))

	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	} else {
		fmt.Println("error", err)
	}

	return receiveAddr, resultTx
}

func SendTransactionBackforth(fromName string, toName string, fromAddr string, toAddr string, fromSeq int64, toSeq int64) {

	// send

	_, result1 := SendTransactionWithSequence(fromName, toAddr, fromSeq)

	_, result2 := SendTransactionWithSequence(toName, fromAddr, toSeq)

	fmt.Println(result1.Hash)
	fmt.Println(result2.Hash)

}

func DelegateTransaction(fromName string, valFrom string, valTo string) {

	// delegate

	_, result1 := DelegateTransactionWithSequence(fromName, valFrom, seqDel)

	seqDel++

	_, result1 = WidthdrawTransactionWithSequence(fromName, seqDel)

	seqDel++

	fmt.Println(result1.Hash)
}

func init() {
	viper.SetConfigName("config")            //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("./config") // 第一个搜索路径
	//viper.AddConfigPath("$HOME/.appname")  // 可以多次调用添加路径

	err := viper.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//init irishub bech32
	initBech32()

	initAccounts()

}

func initAccounts() {
	//transfer
	nameFrom = viper.GetString("NameFrom")
	addrFrom = viper.GetString("AddrFrom")

	nameTo = viper.GetString("NameTo")
	addrTo = viper.GetString("AddrTo")

	//transfer
	nameFrom2 = viper.GetString("NameFrom2")
	addrFrom2 = viper.GetString("AddrFrom2")

	nameTo2 = viper.GetString("NameTo2")
	addrTo2 = viper.GetString("AddrTo2")

	//transfer
	nameFrom3 = viper.GetString("NameFrom3")
	addrFrom3 = viper.GetString("AddrFrom3")

	nameTo3 = viper.GetString("NameTo3")
	addrTo3 = viper.GetString("AddrTo3")

	//get sequence
	seqFrom = GetSequence(addrFrom)
	seqTo = GetSequence(addrTo)

	//get sequence
	seqFrom2 = GetSequence(addrFrom2)
	seqTo2 = GetSequence(addrTo2)

	//get sequence
	seqFrom3 = GetSequence(addrFrom3)
	seqTo3 = GetSequence(addrTo3)

	//delegation
	nameDel = viper.GetString("NameDel")
	addrDel =  viper.GetString("AddrDel")
	seqDel = GetSequence(addrDel)
	valFrom = viper.GetString("ValFrom")
	valTo = viper.GetString("ValTo")

	nameVoter = viper.GetString("NameVoter")
	addrVoter = viper.GetString("AddrVoter")
	seqVoter = GetSequence(addrVoter)

}

//initial bech32 HRL:faa
func initBech32() {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(irisInit.Bech32PrefixAccAddr, irisInit.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(irisInit.Bech32PrefixValAddr, irisInit.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(irisInit.Bech32PrefixConsAddr, irisInit.Bech32PrefixConsPub)
	config.Seal()
}

func main() {

	fmt.Println("Starting the application...")

	//
	c := cron.New()
	//freq :=viper.Get("frequency")

	//fmt.Println(freq)
	feqTransfer := "0-59/5 * * * * ?"      //@every second
	feqDelegate := "0-59/10 * * * * ?" //每分钟执行一次，30s的时候
	feqGov := "0-59/10 * * * * ?"             //每分钟时执行一次

	//
	////@every second send 2 transfer txs
	c.AddFunc(feqTransfer, func() {

		SendTransactionBackforth(nameFrom, nameTo, addrFrom, addrTo, seqFrom, seqTo)

		seqFrom++
		seqTo++

		log.Println("transfer cron running:")

	})
	c.AddFunc(feqTransfer, func() {

		SendTransactionBackforth(nameFrom2, nameTo2, addrFrom2, addrTo2, seqFrom2, seqTo2)

		seqFrom2++
		seqTo2++

		log.Println("transfer cron2 running:")

	})

	c.AddFunc(feqTransfer, func() {

		SendTransactionBackforth(nameFrom3, nameTo3, addrFrom3, addrTo3, seqFrom3, seqTo3)

		seqFrom3++
		seqTo3++

		log.Println("transfer cron2 running:")

	})
	//@every  second send one staking tx
	c.AddFunc(feqDelegate, func() {

		DelegateTransaction(nameDel, valFrom, valTo)

		log.Println("delegate cron running:")

	})

	c.AddFunc(feqGov, func() {

		VoteOnProposal(nameVoter, addrVoter, "Yes")

		log.Println("cron running:")

	})
	c.Start()

	select {} //阻塞主线程不退出

	fmt.Println("Terminating the application...")

}
