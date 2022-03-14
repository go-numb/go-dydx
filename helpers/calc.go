package helpers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/umbracle/ethgo/jsonrpc"

	"github.com/go-numb/go-dydx/types"
)

var (
	Eip712OnboardingActionStruct = []map[string]string{
		{"type": "string", "name": "action"},
		{"type": "string", "name": "onlySignOn"},
	}
	Eip712OnboardingActionStructString = "dYdX(string action,string onlySignOn)"

	Eip712OnboardingActionStructTestnet = []map[string]string{
		{"type": "string", "name": "action"},
	}
	Eip712OnboardingActionStructStringTestnet = "dYdX(string action)"
	Eip712StructName                          = "dYdX"
	OnlySignOnDomainMainnet                   = "https://trade.dydx.exchange"
)

type SignOnboardingAction struct {
	Signer    EthSigner
	NetworkId int
}

func NewSigner(signer EthSigner, networkId int) *SignOnboardingAction {
	return &SignOnboardingAction{signer, networkId}
}

func (a *SignOnboardingAction) Sign(signerAddress string, message map[string]interface{}) string {
	eip712Message := a.GetEIP712Message(message)
	action := message["action"].(string)
	messageHash := a.GetHash(action)
	typedSignature := a.Signer.sign(eip712Message, messageHash, signerAddress)
	return typedSignature
}

func (a *SignOnboardingAction) GetEIP712Message(message map[string]interface{}) map[string]interface{} {
	structName := a.GetEIP712StructName()
	eip712Message := map[string]interface{}{
		"types": map[string]interface{}{
			"EIP712Domain": []map[string]string{
				{
					"name": "name",
					"type": "string",
				},
				{
					"name": "version",
					"type": "string",
				},
				{
					"name": "chainId",
					"type": "uint256",
				},
			},
			structName: a.GetEIP712Struct(),
		},
		"domain": map[string]interface{}{
			"name":    types.Domain,
			"version": types.Version,
			"chainId": a.NetworkId,
		},
		"primaryType": structName,
		"message":     message,
	}
	if a.NetworkId == types.NetworkIdMainnet {
		msg := eip712Message["message"].(map[string]interface{})
		msg["onlySignOn"] = OnlySignOnDomainMainnet
	}

	return eip712Message
}

func (a *SignOnboardingAction) GetEip712Hash(structHash string) string {
	fact := solsha3.SoliditySHA3(
		[]string{"bytes2", "bytes32", "bytes32"},
		[]interface{}{"0x1901", a.GetDomainHash(), structHash},
	)
	return fmt.Sprintf("0x%x", fact)
}

func (a *SignOnboardingAction) GetDomainHash() string {
	fact := solsha3.SoliditySHA3(
		[]string{"bytes32", "bytes32", "bytes32", "uint256"},
		[]interface{}{HashString(types.Eip712DomainStringNoContract), HashString(types.Domain), HashString(types.Version), a.NetworkId},
	)
	return fmt.Sprintf("0x%x", fact)
}

func (a *SignOnboardingAction) GetEIP712Struct() []map[string]string {
	if a.NetworkId == types.NetworkIdMainnet {
		return Eip712OnboardingActionStruct
	} else {
		return Eip712OnboardingActionStructTestnet
	}
}

func (a *SignOnboardingAction) GetEIP712StructName() string {
	return Eip712StructName
}

func (a *SignOnboardingAction) GetHash(action string) string {
	var eip712StructStr string
	if a.NetworkId == types.NetworkIdMainnet {
		eip712StructStr = Eip712OnboardingActionStructString
	} else {
		eip712StructStr = Eip712OnboardingActionStructStringTestnet
	}
	data := [][]string{
		{"bytes32", "bytes32"},
		{HashString(eip712StructStr), HashString(action)},
	}
	if a.NetworkId == types.NetworkIdMainnet {
		data[0] = append(data[0], "bytes32")
		data[1] = append(data[1], HashString(OnlySignOnDomainMainnet))
	}
	structHash := solsha3.SoliditySHA3(data[0], data[1])
	return a.GetEip712Hash(hexutil.Encode(structHash))
}

type EthSigner interface {
	sign(eip712Message map[string]interface{}, messageHash, optSingerAddress string) string
}

type EthWeb3Signer struct {
	Web3 *jsonrpc.Client
}

func (web3Singer *EthWeb3Signer) sign(eip712Message map[string]interface{}, messageHash, address string) string {
	rawSignature := signTypedData(eip712Message, web3Singer, address)
	return CreateTypedSignature(rawSignature, types.SignatureTypeNoPrepend)
}

func signTypedData(eip712Message map[string]interface{}, web3Singer *EthWeb3Signer, address string) string {
	var out string
	if err := web3Singer.Web3.Call("eth_signTypedData", &out, address, eip712Message); err != nil {
		panic(err)
	}
	return out
}

type EthKeySinger struct {
	Address    string
	PrivateKey string
}

func (keySinger EthKeySinger) sign(eip712Message map[string]interface{}, messageHash, optSingerAddress string) string {
	panic("implement me")
}
