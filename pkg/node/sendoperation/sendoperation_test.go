package sendoperation

import (
	"encoding/base64"
	"testing"

	"github.com/massalabs/station/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TransactionOpType = uint64(0)
	BuyRollOpType     = uint64(1)
	SellRollOpType    = uint64(2)
	ExecuteSCOpType   = uint64(3)
	CallSCOpType      = uint64(4)
)

func TestSerializeDeserializeCallSCMessage(t *testing.T) {
	assert := assert.New(t)

	testestcaseases := []struct {
		expiry     uint64
		fee        uint64
		address    string
		function   string
		parameters []byte
		maxGas     uint64
		coins      uint64
	}{
		{
			expiry:     uint64(123456),
			fee:        uint64(789),
			address:    "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5",
			function:   "exampleFunction",
			parameters: []byte("exampleParameters"),
			maxGas:     uint64(1000000),
			coins:      uint64(12345),
		},
	}

	for _, testcase := range testestcaseases {
		// Create a new CallSC operation
		operation, err := callsc.New(testcase.address, testcase.function, testcase.parameters,
			testcase.maxGas, testcase.coins)
		require.NoError(t, err, "Failed to create CallSC operation")

		// Serialize the operation
		msg := message(testcase.expiry, testcase.fee, operation)
		msgB64 := base64.StdEncoding.EncodeToString(msg)

		// Simulate decoding and deserialization
		decodedMsg, fee, expiry, err := DecodeMessage64(msgB64)
		require.NoError(t, err, "Error decoding message")

		// verify the fee and expiry
		assert.Equal(testcase.fee, fee, "fee mismatch")
		assert.Equal(testcase.expiry, expiry, "expiry mismatch")

		callSC, err := callsc.DecodeMessage(decodedMsg)
		require.NoError(t, err, "Error decoding CallSC")

		// Verify the fields
		assert.Equal(CallSCOpType, callSC.OperationType, "Operation type mismatch")
		assert.Equal(testcase.address, callSC.Address, "Address mismatch")
		assert.Equal(testcase.function, callSC.Function, "Function mismatch")
		assert.Equal(testcase.parameters, callSC.Parameters, "Parameters mismatch")
		assert.Equal(testcase.maxGas, callSC.MaxGas, "maxGas mismatch")
		assert.Equal(testcase.coins, callSC.Coins, "Coins mismatch")
	}
}

func TestSerializeDeserializeExecuteSCMessage(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		expiry    uint64
		fee       uint64
		maxGas    uint64
		maxCoins  uint64
		byteCode  []byte
		dataStore []byte
	}{
		{ // no dataStore
			expiry:   uint64(123456),
			fee:      uint64(789),
			maxGas:   uint64(100000),
			maxCoins: uint64(50000),
			byteCode: []byte("exampleByteCode"),
		},
		{ // with dataStore
			expiry:    uint64(123456),
			fee:       uint64(789),
			maxGas:    uint64(100000),
			maxCoins:  uint64(50000),
			byteCode:  []byte(" exampleByteCode "),
			dataStore: []byte(" exampleDataStore "),
		},
		{ // expiry, fees, maxGas and maxCoins are 0
			byteCode:  []byte("exampleByteCode"),
			dataStore: []byte("exampleDataStore"),
		},
	}

	for _, testCase := range testCases {
		// Create a new ExecuteSC operation
		operation := executesc.New(testCase.byteCode, testCase.maxGas, testCase.maxCoins, testCase.dataStore)

		// Serialize the operation
		msg := message(testCase.expiry, testCase.fee, operation)
		msgB64 := base64.StdEncoding.EncodeToString(msg)

		// Simulate decoding and deserialization
		decodedMsg, fee, expiry, err := DecodeMessage64(msgB64)
		require.NoError(t, err, "Error decoding message")

		// verify the fee and expiry
		assert.Equal(testCase.fee, fee, "fee mismatch")
		assert.Equal(testCase.expiry, expiry, "expiry mismatch")

		executeSC, err := executesc.DecodeMessage(decodedMsg)
		require.NoError(t, err, "Error decoding ExecuteSC")

		// Verify the fields
		assert.Equal(ExecuteSCOpType, executeSC.OperationType, "Operation type mismatch")
		assert.Equal(testCase.maxGas, executeSC.MaxGas, "MaxGas mismatch")
		assert.Equal(testCase.maxCoins, executeSC.MaxCoins, "MaxCoins mismatch")
		assert.Equal(testCase.byteCode, executeSC.ByteCode, "ByteCode mismatch")
		assert.Equal(testCase.dataStore, executeSC.DataStore, "DataStore mismatch")
	}
}

func TestSerializeDeserializeBuyRollsMessage(t *testing.T) {
	assert := assert.New(t)

	testcases := []struct {
		countRolls uint64
	}{
		{
			countRolls: 5,
		},
	}

	for _, testcase := range testcases {
		// Create a new BuyRolls operation
		operation := buyrolls.New(testcase.countRolls)

		// Simulate decoding and deserialization
		buyRolls, err := RollDecodeMessage(operation.Message())
		require.NoError(t, err, "Error decoding BuyRolls")

		// Verify the countRolls field
		assert.Equal(BuyRollOpType, buyRolls.OperationType, "Operation type mismatch")
		assert.Equal(testcase.countRolls, buyRolls.RollCount, "CountRolls mismatch")
	}
}

func TestSerializeDeserializeSellRollsMessage(t *testing.T) {
	assert := assert.New(t)

	testcases := []struct {
		countRolls uint64
	}{
		{
			countRolls: 10,
		},
	}

	for _, testcase := range testcases {
		// Create a new SellRolls operation
		operation := sellrolls.New(testcase.countRolls)

		// Simulate decoding and deserialization
		sellRolls, err := RollDecodeMessage(operation.Message())
		require.NoError(t, err, "Error decoding SellRolls")

		// Verify the countRolls field
		assert.Equal(SellRollOpType, sellRolls.OperationType, "Operation type mismatch")
		assert.Equal(testcase.countRolls, sellRolls.RollCount, "CountRolls mismatch")
	}
}

func TestSerializeDeserializeTransactionMessage(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		recipientAddress string
		amount           uint64
	}{
		{
			recipientAddress: "AU1MPDRXuR22mwYDFCeZUDgYjcTAF1co6xujx2X6ugoHeYeGY3B5",
			amount:           uint64(12345),
		},
	}

	for _, testCase := range testCases {
		// Create a new Transaction
		myTx, err := transaction.New(testCase.recipientAddress, testCase.amount)
		require.NoError(t, err, "Failed to create Transaction")

		decodedOpType, err := DecodeOperationType(myTx.Message())
		require.NoError(t, err, "Failed to retrieve operationID")

		// Simulate decoding and deserialization
		decodedTransaction, err := transaction.DecodeMessage(myTx.Message())
		require.NoError(t, err, "Error decoding message")

		// Verify the fields
		assert.Equal(TransactionOpType, decodedOpType, "Operation type mismatch")
		assert.Equal(testCase.recipientAddress, decodedTransaction.RecipientAddress, "RecipientAddress mismatch")
		assert.Equal(testCase.amount, decodedTransaction.Amount, "Amount mismatch")
	}
}

func TestDecodeOperationID(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		msg         []byte
		expectedID  uint64
		expectedErr bool
	}{
		{
			msg:         []byte{0x00}, // Encoded operation ID = 0
			expectedID:  0,
			expectedErr: false,
		},
		{
			msg:         []byte{0x01}, // Encoded operation ID = 1
			expectedID:  1,
			expectedErr: false,
		},
		{
			msg:         []byte{0x02}, // Encoded operation ID = 2
			expectedID:  2,
			expectedErr: false,
		},
	}

	for _, testCase := range testCases {
		decodedID, err := DecodeOperationType(testCase.msg)
		if testCase.expectedErr {
			require.Error(t, err, "Expected error")
		} else {
			require.NoError(t, err, "Unexpected error")
			assert.Equal(testCase.expectedID, decodedID, "Operation ID mismatch")
		}
	}
}
