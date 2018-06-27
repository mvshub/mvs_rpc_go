package mvs_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RPCClient struct {
	sync.RWMutex
	Url         string
	sick        bool
	sickRate    int
	successRate int
	client      *http.Client
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}

func NewRPCClient(url, timeout string) *RPCClient {
	rpcClient := &RPCClient{Url: url}
	timeoutIntv := MustParseDuration(timeout)
	rpcClient.client = &http.Client{
		Timeout: timeoutIntv,
	}
	return rpcClient
}

type JSONRpcResp struct {
	Id     *json.RawMessage       `json:"id"`
	Result *json.RawMessage       `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func (r *RPCClient) doPost(url string, method string, params interface{}) (*JSONRpcResp, error) {
	jsonReq := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 0}
	data, _ := json.Marshal(jsonReq)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Length", (string)(len(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.markSick()
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp *JSONRpcResp
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		r.markSick()
		return nil, err
	}
	if rpcResp.Error != nil {
		r.markSick()
		return nil, errors.New(rpcResp.Error["message"].(string))
	}
	return rpcResp, err
}

func (r *RPCClient) markSick() {
	r.Lock()
	r.sickRate++
	r.successRate = 0
	if r.sickRate >= 5 {
		r.sick = true
	}
	r.Unlock()
}

// auto-generate code begin

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TOADDRESS(std::string): "Target address"
   :param: DIDSYMBOL(std::string): "Did symbol"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Didchangeaddress(ACCOUNTNAME string, ACCOUNTAUTH string, TOADDRESS string, DIDSYMBOL string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didchangeaddress"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TOADDRESS, DIDSYMBOL}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TRANSACTION(string of hexcode): "The input Base16 transaction to sign."
   :param: selfpublickey(std::string): "The private key of this public key will be used to sign."
   :param: broadcast(bool): "Broadcast the tx if it is fullly signed, disabled by default."
*/
func (r *RPCClient) Signmultisigtx(ACCOUNTNAME string, ACCOUNTAUTH string, TRANSACTION string, selfpublickey string, broadcast bool) (*JSONRpcResp, error) {
	cmd := "signmultisigtx"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TRANSACTION}

	optional := map[string]interface{}{}
	if broadcast == true {
		positional = append(positional, "--broadcast")
	}
	if selfpublickey != "" {
		optional["selfpublickey"] = selfpublickey
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: ADDRESS(std::string): "The address will be bound to, can change to other addresses later."
   :param: SYMBOL(std::string): "The symbol of global unique MVS Digital Identity Destination/Index, supports alphabets/numbers/(“@”, “.”, “_”, “-“), case-sensitive, maximum length is 64."
   :param: fee(uint64_t): "The fee of tx. defaults to 1 etp."
*/
func (r *RPCClient) Registerdid(ACCOUNTNAME string, ACCOUNTAUTH string, ADDRESS string, SYMBOL string, fee uint64) (*JSONRpcResp, error) {
	cmd := "registerdid"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, ADDRESS, SYMBOL}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: SYMBOL(std::string): "The asset symbol, global uniqueness, only supports UPPER-CASE alphabet and dot(.)"
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "The fee of tx. minimum is 10 etp."
*/
func (r *RPCClient) Issue(ACCOUNTNAME string, ACCOUNTAUTH string, SYMBOL string, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "issue"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, SYMBOL}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: WORD(list of string): "The set of words that that make up the mnemonic. If not specified the words are read from STDIN."
   :param: language(explorer::config::language): "The language identifier of the dictionary of the mnemonic. Options are 'en', 'es', 'ja', 'zh_Hans', 'zh_Hant' and 'any', defaults to 'any'."
   :param: accountname(std::string): Account name required.
   :param: password(std::string): Account password(authorization) required.
   :param: hd_index(std::uint32_t): "The HD index for the account."
*/
func (r *RPCClient) Importaccount(WORD []string, language string, accountname string, password string, hd_index uint32) (*JSONRpcResp, error) {
	cmd := "importaccount"
	positional := []interface{}{strings.Join(WORD, " ")}

	optional := map[string]interface{}{
		"accountname": accountname,
		"password":    password,
	}

	if language != "" {
		optional["language"] = language
	}
	if hd_index != 0 {
		optional["hd_index"] = hd_index
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Stopmining(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "stopmining"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FROMADDRESS(std::string): "Send from this address, must be a multi-signature script address."
   :param: TOADDRESS(std::string): "Send to this address"
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: symbol(std::string): "asset name, not specify this option for etp tx"
   :param: type(uint16_t): "Transaction type, defaults to 0. 0 -- transfer etp, 3 -- transfer asset"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Createmultisigtx(ACCOUNTNAME string, ACCOUNTAUTH string, FROMADDRESS string, TOADDRESS string, AMOUNT uint64, symbol string, type_ uint16, fee uint64) (*JSONRpcResp, error) {
	cmd := "createmultisigtx"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FROMADDRESS, TOADDRESS, AMOUNT}

	optional := map[string]interface{}{}

	if symbol != "" {
		optional["symbol"] = symbol
	}
	if type_ != 0 {
		optional["type"] = type_
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: ADDRESS(std::string): "Address."
*/
func (r *RPCClient) Getpublickey(ACCOUNTNAME string, ACCOUNTAUTH string, ADDRESS string) (*JSONRpcResp, error) {
	cmd := "getpublickey"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, ADDRESS}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: address(std::string): "The deposit target address."
   :param: deposit(uint16_t): "Deposits support [7, 30, 90, 182, 365] days. defaluts to 7 days"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Deposit(ACCOUNTNAME string, ACCOUNTAUTH string, AMOUNT uint64, address string, deposit uint16, fee uint64) (*JSONRpcResp, error) {
	cmd := "deposit"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, AMOUNT}

	optional := map[string]interface{}{}

	if address != "" {
		optional["address"] = address
	}
	if deposit != 0 {
		optional["deposit"] = deposit
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: SYMBOL(std::string): "Asset symbol."
   :param: cert(bool): "If specified, then only get related asset cert. Default is not specified."
*/
func (r *RPCClient) Getaccountasset(ACCOUNTNAME string, ACCOUNTAUTH string, SYMBOL string, cert bool) (*JSONRpcResp, error) {
	cmd := "getaccountasset"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, SYMBOL}

	optional := map[string]interface{}{}
	if cert == true {
		positional = append(positional, "--cert")
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TO_(std::string): "Asset receiver did/address."
   :param: ASSET(std::string): "Asset MST symbol."
   :param: AMOUNT(uint64_t): "Asset integer bits. see asset <decimal_number>."
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Didsendasset(ACCOUNTNAME string, ACCOUNTAUTH string, TO_ string, ASSET string, AMOUNT uint64, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didsendasset"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TO_, ASSET, AMOUNT}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: SYMBOL(std::string): "The asset will be burned."
   :param: AMOUNT(uint64_t): "Asset integer bits. see asset <decimal_number>."
*/
func (r *RPCClient) Burn(ACCOUNTNAME string, ACCOUNTAUTH string, SYMBOL string, AMOUNT uint64) (*JSONRpcResp, error) {
	cmd := "burn"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, SYMBOL, AMOUNT}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: height(uint32_t): "specify the starting point to pop out blocks. eg, if specified 1000000, then all blocks with height greater than or equal to 1000000 will be poped out."
*/
func (r *RPCClient) Popblock(height uint32) (*JSONRpcResp, error) {
	cmd := "popblock"
	positional := []interface{}{height}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: nozero(bool): "Defaults to false."
   :param: greater_equal(uint64_t): "Greater than ETP bits."
   :param: lesser_equal(uint64_t): "Lesser than ETP bits."
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Listbalances(nozero bool, greater_equal uint64, lesser_equal uint64, ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "listbalances"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}
	if nozero == true {
		positional = append(positional, "--nozero")
	}
	if greater_equal != 0 {
		optional["greater_equal"] = greater_equal
	}
	if lesser_equal != 0 {
		optional["lesser_equal"] = lesser_equal
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: rate(int32_t): "The percent threshold value when you secondary issue.              0,  not allowed to secondary issue;              -1,  the asset can be secondary issue freely;             [1, 100], the asset can be secondary issue when own percentage greater than or equal to this value.             Defaults to 0."
   :param: symbol(std::string): "The asset symbol, global uniqueness, only supports UPPER-CASE alphabet and dot(.), eg: CHENHAO.LAPTOP, dot separates prefix 'CHENHAO', It's impossible to create any asset named with 'CHENHAO' prefix, but this issuer."
   :param: issuer(std::string): "Issue must be specified as a DID symbol."
   :param: volume(non_negative_uint64): "The asset maximum supply volume, with unit of integer bits."
   :param: decimalnumber(uint32_t): "The asset amount decimal number, defaults to 0."
   :param: description(std::string): "The asset data chuck, defaults to empty string."
*/
func (r *RPCClient) Createasset(ACCOUNTNAME string, ACCOUNTAUTH string, rate int32, symbol string, issuer string, volume uint64, decimalnumber uint32, description string) (*JSONRpcResp, error) {
	cmd := "createasset"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"symbol": symbol,
		"issuer": issuer,
		"volume": volume,
	}

	if rate != 0 {
		optional["rate"] = rate
	}
	if decimalnumber != 0 {
		optional["decimalnumber"] = decimalnumber
	}
	if description != "" {
		optional["description"] = description
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TOADDRESS(std::string): "Send to this address"
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: memo(std::string): "Attached memo for this transaction."
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 etp bits"
*/
func (r *RPCClient) Send(ACCOUNTNAME string, ACCOUNTAUTH string, TOADDRESS string, AMOUNT uint64, memo string, fee uint64) (*JSONRpcResp, error) {
	cmd := "send"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TOADDRESS, AMOUNT}

	optional := map[string]interface{}{}

	if memo != "" {
		optional["memo"] = memo
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: password(std::string): "The new password."
*/
func (r *RPCClient) Changepasswd(ACCOUNTNAME string, ACCOUNTAUTH string, password string) (*JSONRpcResp, error) {
	cmd := "changepasswd"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"password": password,
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: type(uint16_t): "Transaction type. 0 -- transfer etp, 1 -- deposit etp, 3 -- transfer asset"
   :param: senders(list of string): "Send from addresses"
   :param: receivers(list of string): "Send to [address:amount]. amount is asset number if sybol option specified"
   :param: symbol(std::string): "asset name, not specify this option for etp tx"
   :param: deposit(uint16_t): "Deposits support [7, 30, 90, 182, 365] days. defaluts to 7 days"
   :param: mychange(std::string): "Mychange to this address, includes etp and asset change"
   :param: message(std::string): "Message/Information attached to this transaction"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Createrawtx(type_ uint16, senders []string, receivers []string, symbol string, deposit uint16, mychange string, message string, fee uint64) (*JSONRpcResp, error) {
	cmd := "createrawtx"
	positional := []interface{}{}

	optional := map[string]interface{}{
		"type":      type_,
		"senders":   senders,
		"receivers": receivers,
	}

	if symbol != "" {
		optional["symbol"] = symbol
	}
	if deposit != 0 {
		optional["deposit"] = deposit
	}
	if mychange != "" {
		optional["mychange"] = mychange
	}
	if message != "" {
		optional["message"] = message
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: PAYMENT_ADDRESS(std::string): "Valid payment address. If not specified the address is read from STDIN."
*/
func (r *RPCClient) Validateaddress(PAYMENT_ADDRESS string) (*JSONRpcResp, error) {
	cmd := "validateaddress"
	positional := []interface{}{PAYMENT_ADDRESS}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FROMADDRESS(std::string): "Send from this address"
   :param: TOADDRESS(std::string): "Send to this address"
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: memo(std::string): "The memo to descript transaction"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Sendfrom(ACCOUNTNAME string, ACCOUNTAUTH string, FROMADDRESS string, TOADDRESS string, AMOUNT uint64, memo string, fee uint64) (*JSONRpcResp, error) {
	cmd := "sendfrom"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FROMADDRESS, TOADDRESS, AMOUNT}

	optional := map[string]interface{}{}

	if memo != "" {
		optional["memo"] = memo
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: ADDRESS(std::string): "The multisig script corresponding address."
*/
func (r *RPCClient) Deletemultisig(ACCOUNTNAME string, ACCOUNTAUTH string, ADDRESS string) (*JSONRpcResp, error) {
	cmd := "deletemultisig"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, ADDRESS}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Listdids(ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "listdids"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getheight(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getheight"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TO_(std::string): "Send to this did/address"
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: memo(std::string): "Attached memo for this transaction."
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 etp bits"
*/
func (r *RPCClient) Didsend(ACCOUNTNAME string, ACCOUNTAUTH string, TO_ string, AMOUNT uint64, memo string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didsend"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TO_, AMOUNT}

	optional := map[string]interface{}{}

	if memo != "" {
		optional["memo"] = memo
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TODID(std::string): "Target did"
   :param: SYMBOL(std::string): "Asset cert symbol"
   :param: CERT(std::string): "Asset cert type name. eg. ISSUE, DOMAIN or NAMING"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Transfercert(ACCOUNTNAME string, ACCOUNTAUTH string, TODID string, SYMBOL string, CERT string, fee uint64) (*JSONRpcResp, error) {
	cmd := "transfercert"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TODID, SYMBOL, CERT}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: TRANSACTION(string of hexcode): "The input Base16 transaction to broadcast."
   :param: fee(uint64_t): "The max tx fee. default_value 10 etp"
*/
func (r *RPCClient) Sendrawtx(TRANSACTION string, fee uint64) (*JSONRpcResp, error) {
	cmd := "sendrawtx"
	positional := []interface{}{TRANSACTION}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TODID(std::string): "The DID will own this cert."
   :param: SYMBOL(std::string): "Asset Cert Symbol/Name."
   :param: CERT(std::string): "Asset cert type name can be: ISSUE: cert of issuing asset, generated by issuing asset and used in secondaryissue asset.  DOMAIN: cert of domain, generated by issuing asset, the symbol is same as asset symbol(if it does not contain dot) or the prefix part(that before the first dot) of asset symbol. NAMING: cert of naming right of domain. The owner of domain cert can issue this type of cert by issuecert with symbol like “domain.XYZ”(domain is the symbol of domain cert)."
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Issuecert(ACCOUNTNAME string, ACCOUNTAUTH string, TODID string, SYMBOL string, CERT string, fee uint64) (*JSONRpcResp, error) {
	cmd := "issuecert"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TODID, SYMBOL, CERT}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: NUMBER(std::string): "Block number, or earliest, latest or pending"
*/
func (r *RPCClient) Fetchheaderext(ACCOUNTNAME string, ACCOUNTAUTH string, NUMBER string) (*JSONRpcResp, error) {
	cmd := "fetchheaderext"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, NUMBER}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FROM_(std::string): "From did/address"
   :param: TO_(std::string): "Target did/address"
   :param: SYMBOL(std::string): "Asset symbol"
   :param: AMOUNT(uint64_t): "Asset integer bits. see asset <decimal_number>."
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Didsendassetfrom(ACCOUNTNAME string, ACCOUNTAUTH string, FROM_ string, TO_ string, SYMBOL string, AMOUNT uint64, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didsendassetfrom"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FROM_, TO_, SYMBOL, AMOUNT}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: receivers(list of string): "Send to [did/address:etp_bits]."
   :param: mychange(std::string): "Mychange to this did/address"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Didsendmore(ACCOUNTNAME string, ACCOUNTAUTH string, receivers []string, mychange string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didsendmore"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"receivers": receivers,
	}

	if mychange != "" {
		optional["mychange"] = mychange
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: receivers(list of string): "Send to [address:etp_bits]."
   :param: mychange(std::string): "Mychange to this address"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Sendmore(ACCOUNTNAME string, ACCOUNTAUTH string, receivers []string, mychange string, fee uint64) (*JSONRpcResp, error) {
	cmd := "sendmore"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"receivers": receivers,
	}

	if mychange != "" {
		optional["mychange"] = mychange
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: symbol(std::string): "The asset symbol/name. Global unique."
*/
func (r *RPCClient) Deletelocalasset(ACCOUNTNAME string, ACCOUNTAUTH string, symbol string) (*JSONRpcResp, error) {
	cmd := "deletelocalasset"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"symbol": symbol,
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: address(std::string): "Address."
   :param: height(a range expressed by 2 integers): "Get tx according height eg: -e start-height:end-height will return tx between [start-height, end-height)"
   :param: symbol(std::string): "Asset symbol."
   :param: limit(uint64_t): "Transaction count per page."
   :param: index(uint64_t): "Page index."
*/
func (r *RPCClient) Listtxs(ACCOUNTNAME string, ACCOUNTAUTH string, address string, height [2]uint64, symbol string, limit uint64, index uint64) (*JSONRpcResp, error) {
	cmd := "listtxs"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	if address != "" {
		optional["address"] = address
	}
	if height != [2]uint64{0, 0} {
		optional["height"] = strings.Join([]string{strconv.FormatUint(uint64(height[0]), 10), strconv.FormatUint(uint64(height[1]), 10)}, ":")
	}
	if symbol != "" {
		optional["symbol"] = symbol
	}
	if limit != 0 {
		optional["limit"] = limit
	}
	if index != 0 {
		optional["index"] = index
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: SYMBOL(std::string): "Asset symbol. If not specified then show whole network MIT symbols."
   :param: trace(bool): "If specified then trace the history. Default is not specified."
   :param: limit(uint32_t): "MIT count per page."
   :param: index(uint32_t): "Page index."
   :param: current(bool): "If specified then show the lastest information of specified MIT. Default is not specified."
*/
func (r *RPCClient) Getmit(SYMBOL string, trace bool, limit uint32, index uint32, current bool) (*JSONRpcResp, error) {
	cmd := "getmit"
	positional := []interface{}{SYMBOL}

	optional := map[string]interface{}{}
	if trace == true {
		positional = append(positional, "--trace")
	}
	if current == true {
		positional = append(positional, "--current")
	}
	if limit != 0 {
		optional["limit"] = limit
	}
	if index != 0 {
		optional["index"] = index
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: language(std::string): "Options are 'en', 'es', 'ja', 'zh_Hans', 'zh_Hant' and 'any', defaults to 'en'."
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Getnewaccount(language string, ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "getnewaccount"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	if language != "" {
		optional["language"] = language
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Listmits(ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "listmits"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): "admin name."
   :param: ADMINAUTH(std::string): "admin password/authorization."
*/
func (r *RPCClient) Shutdown(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "shutdown"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TRANSACTION(string of hexcode): "The input Base16 transaction to sign."
*/
func (r *RPCClient) Signrawtx(ACCOUNTNAME string, ACCOUNTAUTH string, TRANSACTION string) (*JSONRpcResp, error) {
	cmd := "signrawtx"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TRANSACTION}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: json(bool): "Json format or Raw format, default is Json(true)."
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getmemorypool(json bool, ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getmemorypool"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	if json != false {
		optional["json"] = json
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: hash(string of hash256): "The Base16 block hash."
   :param: height(uint32_t): "The block height."
*/
func (r *RPCClient) Getblockheader(hash string, height uint32) (*JSONRpcResp, error) {
	cmd := "getblockheader"
	positional := []interface{}{}

	optional := map[string]interface{}{}

	if hash != "" {
		optional["hash"] = hash
	}
	if height != 0 {
		optional["height"] = height
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: cert(bool): "If specified, then only get related asset cert. Default is not specified."
*/
func (r *RPCClient) Listassets(ACCOUNTNAME string, ACCOUNTAUTH string, cert bool) (*JSONRpcResp, error) {
	cmd := "listassets"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}
	if cert == true {
		positional = append(positional, "--cert")
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FROMADDRESS(std::string): "From address"
   :param: TOADDRESS(std::string): "Target address"
   :param: SYMBOL(std::string): "Asset symbol"
   :param: AMOUNT(uint64_t): "Asset integer bits. see asset <decimal_number>."
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Sendassetfrom(ACCOUNTNAME string, ACCOUNTAUTH string, FROMADDRESS string, TOADDRESS string, SYMBOL string, AMOUNT uint64, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "sendassetfrom"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FROMADDRESS, TOADDRESS, SYMBOL, AMOUNT}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: SYMBOL(std::string): "Asset symbol. If not specified, will show whole network asset symbols."
   :param: cert(bool): "If specified, then only get related asset cert. Default is not specified."
*/
func (r *RPCClient) Getasset(SYMBOL string, cert bool) (*JSONRpcResp, error) {
	cmd := "getasset"
	positional := []interface{}{SYMBOL}

	optional := map[string]interface{}{}
	if cert == true {
		positional = append(positional, "--cert")
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getinfo(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getinfo"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TODID(std::string): "target did to check and issue asset, fee from and mychange to the address of this did too."
   :param: SYMBOL(std::string): "issued asset symbol"
   :param: VOLUME(uint64_t): "The volume of asset, with unit of integer bits."
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "The fee of tx. default_value 10000 ETP bits"
*/
func (r *RPCClient) Secondaryissue(ACCOUNTNAME string, ACCOUNTAUTH string, TODID string, SYMBOL string, VOLUME uint64, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "secondaryissue"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TODID, SYMBOL, VOLUME}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADDRESS(std::string): "address"
   :param: cert(bool): "If specified, then only get related asset cert. Default is not specified."
*/
func (r *RPCClient) Getaddressasset(ADDRESS string, cert bool) (*JSONRpcResp, error) {
	cmd := "getaddressasset"
	positional := []interface{}{ADDRESS}

	optional := map[string]interface{}{}
	if cert == true {
		positional = append(positional, "--cert")
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: number(std::uint32_t): "The number of addresses to be generated, defaults to 1."
*/
func (r *RPCClient) Getnewaddress(ACCOUNTNAME string, ACCOUNTAUTH string, number uint32) (*JSONRpcResp, error) {
	cmd := "getnewaddress"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	if number != 0 {
		optional["number"] = number
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Getbalance(ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "getbalance"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: signaturenum(uint16_t): "Account multisig signature number."
   :param: publickeynum(uint16_t): "Account multisig public key number."
   :param: selfpublickey(std::string): "the public key belongs to this account."
   :param: publickey(list of string): "cosigner public key used for multisig"
   :param: description(std::string): "multisig record description."
*/
func (r *RPCClient) Getnewmultisig(ACCOUNTNAME string, ACCOUNTAUTH string, signaturenum uint16, publickeynum uint16, selfpublickey string, publickey []string, description string) (*JSONRpcResp, error) {
	cmd := "getnewmultisig"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{
		"signaturenum":  signaturenum,
		"publickeynum":  publickeynum,
		"selfpublickey": selfpublickey,
	}

	if publickey != nil {
		optional["publickey"] = publickey
	}
	if description != "" {
		optional["description"] = description
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TODID(std::string): "Target did"
   :param: SYMBOL(std::string): "Asset MIT symbol"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Transfermit(ACCOUNTNAME string, ACCOUNTAUTH string, TODID string, SYMBOL string, fee uint64) (*JSONRpcResp, error) {
	cmd := "transfermit"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TODID, SYMBOL}

	optional := map[string]interface{}{}

	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: LASTWORD(std::string): "The last word of your private-key phrase."
*/
func (r *RPCClient) Deleteaccount(ACCOUNTNAME string, ACCOUNTAUTH string, LASTWORD string) (*JSONRpcResp, error) {
	cmd := "deleteaccount"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, LASTWORD}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Listmultisig(ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "listmultisig"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: DidOrAddress(std::string): "Did symbol or standard address; If no input parameters, then display whole network DIDs."
*/
func (r *RPCClient) Getdid(DidOrAddress string) (*JSONRpcResp, error) {
	cmd := "getdid"
	positional := []interface{}{DidOrAddress}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: address(std::string): "The mining target address. Defaults to empty, means a new address will be generated."
   :param: number(uint16_t): "The number of mining blocks, useful for testing. Defaults to 0, means no limit."
*/
func (r *RPCClient) Startmining(ACCOUNTNAME string, ACCOUNTAUTH string, address string, number uint16) (*JSONRpcResp, error) {
	cmd := "startmining"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	if address != "" {
		optional["address"] = address
	}
	if number != 0 {
		optional["number"] = number
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getwork(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getwork"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FILE(string of file path): "key file path."
   :param: FILECONTENT(std::string): "key file content. this will omit the FILE argument if specified."
*/
func (r *RPCClient) Importkeyfile(ACCOUNTNAME string, ACCOUNTAUTH string, FILE string, FILECONTENT string) (*JSONRpcResp, error) {
	cmd := "importkeyfile"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FILE, FILECONTENT}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: TRANSACTION(string of hexcode): "The input Base16 transaction to sign."
*/
func (r *RPCClient) Decoderawtx(TRANSACTION string) (*JSONRpcResp, error) {
	cmd := "decoderawtx"
	positional := []interface{}{TRANSACTION}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: ADDRESS(std::string): "Asset receiver."
   :param: SYMBOL(std::string): "Asset symbol/name."
   :param: AMOUNT(uint64_t): "Asset integer bits. see asset <decimal_number>."
   :param: model(std::string): The token offering model by block height.
   TYPE=1 - fixed quantity model; TYPE=2 - specify parameters;
   LQ - Locked Quantity each period;
   LP - Locked Period, numeber of how many blocks;
   UN - Unlock Number, number of how many LPs;
   eg:
       TYPE=1;LQ=9000;LP=60000;UN=3
       TYPE=2;LQ=9000;LP=60000;UN=3;UC=20000,20000,20000;UQ=3000,3000,3000
   defaults to disable.
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Sendasset(ACCOUNTNAME string, ACCOUNTAUTH string, ADDRESS string, SYMBOL string, AMOUNT uint64, model string, fee uint64) (*JSONRpcResp, error) {
	cmd := "sendasset"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, ADDRESS, SYMBOL, AMOUNT}

	optional := map[string]interface{}{}

	if model != "" {
		optional["model"] = model
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: NONCE(std::string): "nonce. without leading 0x"
   :param: HEADERHASH(std::string): "header hash. with leading 0x"
   :param: MIXHASH(std::string): "mix hash. with leading 0x"
*/
func (r *RPCClient) Submitwork(NONCE string, HEADERHASH string, MIXHASH string) (*JSONRpcResp, error) {
	cmd := "submitwork"
	positional := []interface{}{NONCE, HEADERHASH, MIXHASH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: PAYMENT_ADDRESS(string of Base58-encoded public key address): "The payment address. If not specified the address is read from STDIN."
*/
func (r *RPCClient) Getaddressetp(PAYMENT_ADDRESS string) (*JSONRpcResp, error) {
	cmd := "getaddressetp"
	positional := []interface{}{PAYMENT_ADDRESS}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: json(bool): "Json/Raw format, default is '--json=true'."
   :param: HASH(string of hash256): "The Base16 transaction hash of the transaction to get. If not specified the transaction hash is read from STDIN."
*/
func (r *RPCClient) Gettx(json bool, HASH string) (*JSONRpcResp, error) {
	cmd := "gettx"
	positional := []interface{}{json, HASH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getmininginfo(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getmininginfo"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: TODID(std::string): "Target did"
   :param: SYMBOL(std::string): "MIT symbol"
   :param: content(std::string): "Content of MIT"
   :param: mits(list of string): "List of symbol and content pair. Symbol and content are separated by a ':'"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Registermit(ACCOUNTNAME string, ACCOUNTAUTH string, TODID string, SYMBOL string, content string, mits []string, fee uint64) (*JSONRpcResp, error) {
	cmd := "registermit"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, TODID}
	if SYMBOL != "" {
		positional = append(positional, SYMBOL)
	}
	optional := map[string]interface{}{}

	if content != "" {
		optional["content"] = content
	}
	if mits != nil {
		optional["mits"] = mits
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: PAYMENT_ADDRESS(string of Base58-encoded public key address): "the payment address of this account."
*/
func (r *RPCClient) Setminingaccount(ACCOUNTNAME string, ACCOUNTAUTH string, PAYMENT_ADDRESS string) (*JSONRpcResp, error) {
	cmd := "setminingaccount"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, PAYMENT_ADDRESS}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
*/
func (r *RPCClient) Listaddresses(ACCOUNTNAME string, ACCOUNTAUTH string) (*JSONRpcResp, error) {
	cmd := "listaddresses"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: LASTWORD(std::string): "The last word of your master private-key phrase."
   :param: DESTINATION(string of file path): "The keyfile storage path to."
   :param: data(bool): "If specified, the keyfile content will be append to the report, rather than to local file specified by DESTINATION."
*/
func (r *RPCClient) Dumpkeyfile(ACCOUNTNAME string, ACCOUNTAUTH string, LASTWORD string, DESTINATION string, data bool) (*JSONRpcResp, error) {
	cmd := "dumpkeyfile"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, LASTWORD}
	if DESTINATION != "" {
		positional = append(positional, DESTINATION)
	}
	optional := map[string]interface{}{}
	if data == true {
		positional = append(positional, "--data")
	}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ADMINNAME(std::string): Administrator required.(when administrator_required in mvs.conf is set true)
   :param: ADMINAUTH(std::string): Administrator password required.
*/
func (r *RPCClient) Getpeerinfo(ADMINNAME string, ADMINAUTH string) (*JSONRpcResp, error) {
	cmd := "getpeerinfo"
	positional := []interface{}{ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: FROM_(std::string): "Send from this did/address"
   :param: TO_(std::string): "Send to this did/address"
   :param: AMOUNT(uint64_t): "ETP integer bits."
   :param: memo(std::string): "The memo to descript transaction"
   :param: fee(uint64_t): "Transaction fee. defaults to 10000 ETP bits"
*/
func (r *RPCClient) Didsendfrom(ACCOUNTNAME string, ACCOUNTAUTH string, FROM_ string, TO_ string, AMOUNT uint64, memo string, fee uint64) (*JSONRpcResp, error) {
	cmd := "didsendfrom"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, FROM_, TO_, AMOUNT}

	optional := map[string]interface{}{}

	if memo != "" {
		optional["memo"] = memo
	}
	if fee != 0 {
		optional["fee"] = fee
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: ACCOUNTNAME(std::string): Account name required.
   :param: ACCOUNTAUTH(std::string): Account password(authorization) required.
   :param: LASTWORD(std::string): "The last word of your backup words."
*/
func (r *RPCClient) Getaccount(ACCOUNTNAME string, ACCOUNTAUTH string, LASTWORD string) (*JSONRpcResp, error) {
	cmd := "getaccount"
	positional := []interface{}{ACCOUNTNAME, ACCOUNTAUTH, LASTWORD}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: NODEADDRESS(std::string): "The target node address[x.x.x.x:port]."
   :param: ADMINNAME(std::string): "admin name."
   :param: ADMINAUTH(std::string): "admin password/authorization."
   :param: operation(std::string): "The operation[ add|ban ] to the target node address. default: add."
*/
func (r *RPCClient) Addnode(NODEADDRESS string, ADMINNAME string, ADMINAUTH string, operation string) (*JSONRpcResp, error) {
	cmd := "addnode"
	positional := []interface{}{NODEADDRESS, ADMINNAME, ADMINAUTH}

	optional := map[string]interface{}{}

	if operation != "" {
		optional["operation"] = operation
	}
	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}

/*
   :param: HASH_OR_HEIGH(std::string): "block hash or block height"
   :param: json(bool): "Json/Raw format, default is '--json=true'."
   :param: tx_json(bool): "Json/Raw format for txs, default is '--tx_json=true'."
*/
func (r *RPCClient) Getblock(HASH_OR_HEIGH string, json bool, tx_json bool) (*JSONRpcResp, error) {
	cmd := "getblock"
	positional := []interface{}{HASH_OR_HEIGH, json, tx_json}

	optional := map[string]interface{}{}

	args := append(positional, optional)
	return r.doPost(r.Url, cmd, args)
}
