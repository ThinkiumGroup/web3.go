package dto

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThinkiumGroup/go-common"
	"github.com/ThinkiumGroup/go-common/trie"
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
	"github.com/ThinkiumGroup/web3.go/web3/complex/types"
	"github.com/ThinkiumGroup/web3.go/web3/constants"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type RequestResult struct {
	// ID      int         `json:"id"`
	// Version string      `json:"jsonrpc"`
	Result interface{} `json:"result"`
	Error  *Error      `json:"error,omitempty"`
	Data   string      `json:"data,omitempty"`
}

type SendTxResult struct {
	TXhash string `json:"TXhash,omitempty"`
	ErrMsg string `json:"ErrMsg,omitempty"`
}
type RpcMakeVccProofJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

type MakeCCCExistenceProofJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

//GetCCCRelativeTx
type GetCCCRelativeTxJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}
type CompileContractJson struct {
	Test   map[string]interface{} `json:"test,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

type TransactionResult struct {
	ChainId   int      `json:"chainid"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     int      `json:"nonce"`
	Value     *big.Int `json:"value"`
	Input     string   `json:"input"`
	Hash      string   `json:"hash"`
	UseLocal  bool     `json:"uselocal"`
	Extra     string   `json:"extra"` // It is currently used to save transaction types. If it does not exist, it is a normal transaction. Otherwise, it will correspond to special operations
	Timestamp uint64   `json:"timestamp"`
}

type TxResult struct {
	Transaction     TransactionResult `json:"tx"`
	Root            string            `json:"root"`
	Status          int               `json:"status"`
	Logs            interface{}       `json:"logs"`
	TransactionHash string            `json:"transactionHash"`
	ContractAddress string            `json:"contractAddress"`
	Out             string            `json:"out"`
	GasFee          string            `json:"gasFee"`
	GasUsed         int               `json:"gasUsed"`
	BlockHeight     int               `json:"blockHeight"`
	Error           string            `json:"errorMsg"`
	ErrMsg          string            `json:"ErrMsg,omitempty"`
}

type GetBlockResult struct {
	Hash          string `json:"hash"`          // Hash of this block
	Previoushash  string `json:"previoushash"`  // Hash of parent block
	ChainId       int    `json:"chainid"`       //
	Height        int    `json:"height"`        // The block height of the query block
	Empty         bool   `json:"empty"`         // Is it an empty block
	RewardAddress string `json:"rewardaddress"` // Receiving address
	Mergeroot     string `json:"mergeroot"`     // Merge other chain block data hash
	Deltaroot     string `json:"deltaroot"`     // Cross chain transfer data hash
	Stateroot     string `json:"stateroot"`     // State hash
	RREra         int    `json:"rrera"`
	RRCurrent     string `json:"rrcurrent"`
	RRNext        string `json:"rrnext"`
	Txcount       int    `json:"txcount"`
	Timestamp     int64  `json:"timestamp"`
	ErrMsg        string `json:"ErrMsg,Omitempty"`
}

type NodeInfo struct {
	NodeId        string      `json:"nodeId"`        // Node ID
	Version       string      `json:"version"`       // edition
	IsDataNode    bool        `json:"isDataNode"`    // Is it a data node
	DataNodeOf    int         `json:"dataNodeOf"`    // Data node
	LastMsgTime   int64       `json:"lastMsgTime"`   // Last message time
	LastEventTime int64       `json:"lastEventTime"` // Last event time
	LastBlockTime int64       `json:"lastBlockTime"` // Last block time
	Overflow      bool        `json:"overflow"`      // overflow
	LastBlocks    interface{} `json:"lastBlocks"`    // The last block
	OpTypes       interface{} `json:"opTypes"`       // type
	ErrMsg        string      `json:"ErrMsg,Omitempty"`
}

type DataNode struct {
	DataNodeId   string `json:"dataNodeId"`
	DataNodeIp   string `json:"dataNodeIp"`
	DataNodePort int    `json:"dataNodePort"`
}
type GetChainInfo struct {
	ChainId   int        `json:"chainId"`
	DataNodes []DataNode `json:"datanodes"`
	Mode      int        `json:"mode"`
	Parent    int        `json:"parent"`
}

type BlockTxs struct {
	Elections      interface{}         `json:"elections"`
	AccountChanges []TransactionResult `json:"accountchanges"`
	ErrMsg         string              `json:"errMsg,omitempty"`
}

type GetTransactions struct {
	ChainId   int    `json:"chainId"`
	From      string `json:"from"`
	To        string `json:"to"`
	Nonce     int    `json:"nonce"`
	Value     int    `json:"value"`
	Input     string `json:"input"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
}

type GetChainStats struct {
	ChainId           int      `json:"chainId"`
	CurrentHeight     int      `json:"currentheight"`
	EpochDuration     int      `json:"epochduration"`
	EpochLength       int      `json:"epochlength"`
	GasLimit          int      `json:"gaslimit"`
	GasPrice          string   `json:"gasprice"`
	LastEpochDuration int      `json:"lastepochduration"`
	Lives             int      `json:"lives"`
	Tps               int      `json:"tps"`
	TpsLastEpoch      int      `json:"tpsLastEpoch"`
	N                 int      `json:"n"`
	TpsLastN          int      `json:"tpsLastN"`
	LastNduration     int      `json:"lastNduration"`
	TxCount           int      `json:"txcount"`
	AccountCount      int      `json:"accountcount"`
	CurrentComm       []string `json:"currentcomm"`
}

type GetMultiStatsResult struct {
	ErrMsg string `json:"ErrMsg,Omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TxProof struct {
	TxReceipt
	Proof  MerkleItems `json:"proof"`
	Error  string      `json:"errorMsg"`
	ErrMsg string      `json:"ErrMsg,omitempty"`
}
type MerkleItem struct {
	HashVal   hexutil.Bytes `json:"hash"`
	Direction uint8         `json:"direction"`
}

func (m MerkleItem) Proof(toBeProof []byte) ([]byte, error) {
	order := true
	if m.Direction != 0 {
		order = false
	}
	return common.HashPairOrder(order, m.HashVal, toBeProof)
}

func (ms MerkleItems) Proof(toBeProof []byte) ([]byte, error) {
	r := toBeProof
	var err error
	for _, item := range ms {
		r, err = item.Proof(r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

type MerkleItems []MerkleItem
type TxReceipt struct {
	Transaction     *Transaction   `json:"tx"`                                  // Transaction data object
	Sig             *PubAndSig     `json:"signature"`                           // transaction signature
	PostState       []byte         `json:"root"`                                // It is used to record the information of transaction execution in JSON format, such as gas, cost "gas", and world state "root" after execution.
	Status          uint64         `json:"status"`                              // Transaction execution status, 0: failed, 1: successful. (refers to whether the execution is abnormal)
	Logs            []Log          `json:"logs" gencodec:"required"`            // The log written by the contract during execution
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"` // Transaction Hash
	ContractAddress common.Address `json:"contractAddress"`                     // If you are creating a contract, save the address of the created contract here
	Out             hexutil.Bytes  `json:"out"`                                 // Return value of contract execution
	Height          common.Height  `json:"blockHeight"`                         // The block where the transaction is packaged is high and will not be returned when calling
	GasUsed         uint64         `json:"gasUsed"`                             // The gas value consumed by transaction execution is not returned in call
	GasFee          string         `json:"gasFee"`                              // The gas cost of transaction execution is not returned in call
	PostRoot        []byte         `json:"postroot"`                            // World state root after transaction execution (never return, always empty)
	Error           string         `json:"errorMsg"`                            // Error message in case of transaction execution failure
}

type Transaction struct {
	ChainId   uint32   `json:"chainid"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     uint64   `json:"nonce"`
	Value     *big.Int `json:"value"`
	Input     string   `json:"input"`
	Hash      string   `json:"hash"`
	UseLocal  bool     `json:"uselocal"`
	Extra     string   `json:"extra"` // 目前用来存交易类型，不存在时为普通交易，否则会对应特殊操作
	Timestamp uint64   `json:"timestamp"`
	Version   uint16   `json:"version"` // Version number used to distinguish different execution methods when the transaction execution is incompatible due to upgrade
}
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber" gencodec:"required"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// index of the log in the receipt
	Index uint `json:"logIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash *common.Hash `json:"blockHash"`
}

type BlockDetail struct {
	BlockHeader *BlockHeader
	BlockBody   *BlockBody
	BlockPass   PubAndSigs
	ErrMsg      string `json:"ErrMsg,Omitempty"`
}

type (
	Public    [65]byte
	Signature [65]byte
)

func (p Public) MarshalText() ([]byte, error) {
	return hexutil.Bytes(p[:]).MarshalText()
}

var (
	PublicT    = reflect.TypeOf(Public{})
	SignatureT = reflect.TypeOf(Signature{})
)

func (p *Public) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(PublicT, input, p[:])
}

func (p *Public) Bytes() []byte {
	return p[:]

}

func (s *Signature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(SignatureT, input, s[:])
}

func (s *Signature) Bytes() []byte {
	return s[:]
}

func (s Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(s[:]).MarshalText()
}

func (b *BlockDetail) Hash() common.Hash {
	if b == nil || b.BlockHeader == nil {
		return common.Hash{}
	}
	return b.BlockHeader.Hash()
}

type BlockHeader struct {
	PreviousHash   common.Hash    `json:"previoushash"` // the hash of the previous block header on current chain
	HashHistory    common.Hash    `json:"history"`      // hash of the history tree of hash for each block recorded in height order
	ChainID        common.ChainID `json:"chainid"`      // current chain id
	Height         common.Height  `json:"height"`       // height of current block
	Empty          bool           `json:"empty"`        // empty block
	ParentHeight   common.Height  // height of parent height, is 0 if current is main chain
	ParentHash     *common.Hash   // block hash of main chain block at ParentHeight, nil if current is main chain
	RewardAddress  common.Address // reward to
	AttendanceHash *common.Hash   // The current epoch attendance record hash
	RewardedCursor *common.Height // If the current chain is the reward chain, record start height of main chain when next reward issues

	CommitteeHash    *common.Hash   // current epoch Committee member trie root hash
	ElectedNextRoot  *common.Hash   // root hash of the election result of next epoch committee members
	NewCommitteeSeed *common.Seed   `json:"seed"` // Current election seeds, only in the main chain
	RREra            *common.EraNum // the era corresponding to the root of the current Required Reserve tree. When this value is inconsistent with the height of main chain, it indicates that a new RR tree needs to be calculated
	RRRoot           *common.Hash   // root hash of the Required Reserve tree in current era. Only in the reward chain and the main chain
	RRNextRoot       *common.Hash   // root hash of the Required Reserve tree in next era. Only in the reward chain and the main chain
	RRChangingRoot   *common.Hash   // changes waiting to be processed in current era

	MergedDeltaRoot  *common.Hash `json:"mergeroot"` // Root hash of the merged delta sent from other shards
	BalanceDeltaRoot *common.Hash `json:"deltaroot"` // Root hash of the generated deltas by this block which needs to be sent to the other shards
	StateRoot        common.Hash  `json:"stateroot"` // account on current chain state trie root hash
	ChainInfoRoot    *common.Hash // for main chain only: all chain info trie root hash
	WaterlinesRoot   *common.Hash // since v2.3.0, the waterlines of other shards to current chain after the execution of this block. nil represent all zeros. Because the value of the previous block needs to be inherited when the block is empty, values after block execution recorded.
	VCCRoot          *common.Hash // Root hash of transfer out check tree in business chain
	CashedRoot       *common.Hash // Root hash of transfer in check tree in business chain
	TransactionRoot  *common.Hash // transactions in current block trie root hash
	ReceiptRoot      *common.Hash // receipts for transactions in current block trie root hash
	HdsRoot          *common.Hash // if there's any child chain of current chain, this is the Merkle trie root hash generated by the reported block header information of the child chain in order

	TimeStamp uint64 `json:"timestamp"`

	ElectResultRoot *common.Hash // Since v1.5.0, Election result hash root (including pre election and ordinary election, ordinary one has not been provided yet)
	PreElectRoot    *common.Hash // Since v1.5.0, the root hash of current preelecting list sorted by (Expire, ChainID), only in the main chain
	FactorRoot      *common.Hash // since v2.0.0, seed random factor hash
	RRReceiptRoot   *common.Hash // since v2.11.0, receipts of RRActs applied in current block
	Version         uint16       // since v2.11.0
}

func (h *BlockHeader) Hash() common.Hash {
	hashOfHeader, err := h.HashValue()
	if err != nil {
		panic(fmt.Sprintf("BlockHeader %s merkle tree hash failed: %v", h, err))
	}
	return common.BytesToHash(hashOfHeader)
}

func (h *BlockHeader) HashValue() ([]byte, error) {
	hashList, err := h.hashList()
	if err != nil {
		return nil, fmt.Errorf("BlockHeader %s hash failed: %v", h, err)
	}
	ret, err := common.MerkleHashComplete(hashList, 0, nil)
	return ret, err
}

// Hash value and its corresponding position are generated together to generate hash, which can
// prove that this value is the value in this position
func hashIndexProperty(posBuffer [13]byte, index byte, h []byte) []byte {
	indexHash := common.HeaderIndexHash(posBuffer, index)
	return common.HashPair(indexHash, h)
}

func hashPointerHash(h *common.Hash) []byte {
	if h == nil {
		return common.NilHashSlice
	} else {
		return h[:]
	}
}
func (h *BlockHeader) hashList() ([][]byte, error) {
	if h == nil {
		return nil, common.ErrNil
	}
	posBuffer := common.ToHeaderPosHashBuffer(h.ChainID, h.Height)
	hashlist := make([][]byte, 0, 30)
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 0, h.PreviousHash[:]))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 1, h.HashHistory[:]))
	hh, err := h.ChainID.HashValue()
	if err != nil {
		return nil, err
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 2, hh))
	hh, err = h.Height.HashValue()
	if err != nil {
		return nil, err
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 3, hh))
	var b byte = 0
	if h.Empty {
		b = 1
	}
	hh, err = common.Hash256s([]byte{b})
	if err != nil {
		return nil, err
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 4, hh))
	hh, err = h.ParentHeight.HashValue()
	if err != nil {
		return nil, err
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 5, hh))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 6, hashPointerHash(h.ParentHash)))
	hh, err = common.Hash256s(h.RewardAddress[:])
	if err != nil {
		return nil, err
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 7, hh))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 8, hashPointerHash(h.CommitteeHash)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 9, hashPointerHash(h.ElectedNextRoot)))
	if h.Version == 0 {
		if h.NewCommitteeSeed == nil {
			hh = common.NilHashSlice
		} else {
			hh = h.NewCommitteeSeed[:]
		}
		hh, err = common.Hash256s(hh)
		if err != nil {
			return nil, err
		}
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 10, hh))
	} else {
		if h.NewCommitteeSeed == nil {
			hh = common.NilHashSlice
		} else {
			hh = h.NewCommitteeSeed[:]
			hh, err = common.Hash256s(hh)
			if err != nil {
				return nil, err
			}
		}
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 10, hh))
	}
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 11, hashPointerHash(h.MergedDeltaRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 12, hashPointerHash(h.BalanceDeltaRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 13, h.StateRoot[:]))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 14, hashPointerHash(h.ChainInfoRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 15, hashPointerHash(h.WaterlinesRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 16, hashPointerHash(h.VCCRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 17, hashPointerHash(h.CashedRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 18, hashPointerHash(h.TransactionRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 19, hashPointerHash(h.ReceiptRoot)))
	hashlist = append(hashlist, hashIndexProperty(posBuffer, 20, hashPointerHash(h.HdsRoot)))
	{
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, h.TimeStamp)
		hh, err = common.Hash256s(bs)
		if err != nil {
			return nil, err
		}
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 21, hh))
	}
	if h.Version == 0 {
		// // TODO: should remove conditions when restart the chain with new version
		// // v1.5.0: Because each leaf of merkle tree is not the field value of the block header, nil data is not NilHash
		if h.AttendanceHash != nil || h.RewardedCursor != nil ||
			h.RREra != nil || h.RRRoot != nil || h.RRNextRoot != nil || h.RRChangingRoot != nil {
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 22, hashPointerHash(h.AttendanceHash)))
			if h.RewardedCursor == nil {
				hh = common.NilHashSlice
			} else {
				hh, err = common.HashObject(h.RewardedCursor)
				if err != nil {
					return nil, err
				}
			}
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 23, hh))
			if h.RREra == nil {
				hh = common.NilHashSlice
			} else {
				hh, err = common.HashObject(h.RREra)
				if err != nil {
					return nil, err
				}
			}
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 24, hh))
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 25, hashPointerHash(h.RRRoot)))
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 26, hashPointerHash(h.RRNextRoot)))
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 27, hashPointerHash(h.RRChangingRoot)))
			// add by v1.5.0
			if h.ElectResultRoot != nil || h.PreElectRoot != nil {
				hashlist = append(hashlist, hashIndexProperty(posBuffer, 28, hashPointerHash(h.ElectResultRoot)))
				hashlist = append(hashlist, hashIndexProperty(posBuffer, 29, hashPointerHash(h.PreElectRoot)))
			}
			// add by v2.0.0 newSeed
			if h.FactorRoot != nil {
				hashlist = append(hashlist, hashIndexProperty(posBuffer, 30, hashPointerHash(h.FactorRoot)))
			}
		}
	} else {
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 22, hashPointerHash(h.AttendanceHash)))
		if h.RewardedCursor == nil {
			hh = common.NilHashSlice
		} else {
			hh, err = common.HashObject(h.RewardedCursor)
			if err != nil {
				return nil, err
			}
		}
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 23, hh))
		if h.RREra == nil {
			hh = common.NilHashSlice
		} else {
			hh, err = common.HashObject(h.RREra)
			if err != nil {
				return nil, err
			}
		}
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 24, hh))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 25, hashPointerHash(h.RRRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 26, hashPointerHash(h.RRNextRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 27, hashPointerHash(h.RRChangingRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 28, hashPointerHash(h.ElectResultRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 29, hashPointerHash(h.PreElectRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 30, hashPointerHash(h.FactorRoot)))
		hashlist = append(hashlist, hashIndexProperty(posBuffer, 31, hashPointerHash(h.RRReceiptRoot)))
		{
			bs := make([]byte, 2)
			binary.BigEndian.PutUint16(bs, h.Version)
			hh, err = common.Hash256s(bs)
			if err != nil {
				return nil, err
			}
			hashlist = append(hashlist, hashIndexProperty(posBuffer, 32, hh))
		}
	}
	return hashlist, err
}

type BlockBody struct {
	NextCommittee     *Committee        // election results of the next committee
	NCMsg             []*ElectMessage   // election requests for chains (in main chain)
	DeltaFroms        DeltaFroms        // deltas merged to current shard
	Txs               []*BlockTx        // transactions
	TxsPas            []*PubAndSig      // signatures corresponding to packaged transactions
	Deltas            []*AccountDelta   // the delta generated by packaged transactions on current shard needs to be sent to other shards
	Hds               []*BlockSummary   // block summary reported by children chains
	Attendance        *AttendanceRecord // attendance table of the current epoch
	RewardReqs        RewardRequests    // self-proving reward request of each chain received on the main chain
	ElectingResults   ChainElectResults // Since v1.5.0, a list of election results, it's a preelection when Epoch.IsNil()==true, others are local election
	PreElectings      PreElectings      // Since v1.5.0, the list of preselections in progress, sorted by (expire, chainid)
	NextRealCommittee *Committee        // Since v1.5.0, when election finished, the result will be put into NextCommittee. If the election is failed, the current committee will continue to be used in the next epoch. At this time, the current committee needs to be written into this field, which can be brought with it when reporting.
	SeedFactor        SeedFactor        // Since v2.0.0, random factor of seed
}

type BlockTx struct {
	ChainID   common.ChainID  `json:"chainID"`   // The chain ID that needs to process this transaction
	From      *common.Address `json:"from"`      // Address of transaction transmitter
	To        *common.Address `json:"to"`        // Address of transaction receiver
	Nonce     uint64          `json:"nonce"`     // Nonce of sender account
	UseLocal  bool            `json:"uselocal"`  // true: local currency，false: basic currency; default false
	Val       *big.Int        `json:"value"`     // Amount of the transaction
	Input     hexutil.Bytes   `json:"input"`     // Contract code/initial parameters when creating a contract, or input parameters when calling a contract
	Extra     hexutil.Bytes   `json:"extra"`     // Store transaction additional information
	Version   uint16          `json:"version"`   // Version number used to distinguish different execution methods when the transaction execution is incompatible due to upgrade
	MultiSigs PubAndSigs      `json:"multiSigs"` // The signatures used to sign this transaction will only be used when there are multiple signatures. The signature of the transaction sender is not here. Not included in Hash
}

type (
	PubAndSig struct {
		PublicKey Public
		Signature Signature
	}

	Committee struct {
		Members   []common.NodeID
		indexMap  map[common.NodeID]common.CommID
		indexLock sync.Mutex
	}

	ElectMessage struct {
		// EpochNum is the current epoch number
		// I.e., the elected committee is for epoch EpochNum+1
		EpochNum     common.EpochNum `json:"epoch"` // the epoch when election starts
		ElectChainID common.ChainID  `json:"chainid"`
	}
	DeltaFromKey struct {
		ShardID common.ChainID
		Height  common.Height
	}

	AccountDelta struct {
		Addr          common.Address
		Delta         *big.Int // Balance modification
		CurrencyDelta *big.Int // LocalCurrency modification (if has)
	}

	DeltaFrom struct {
		Key    DeltaFromKey
		Deltas []*AccountDelta
	}

	DeltaFroms []DeltaFrom

	PubAndSigs []*PubAndSig

	BlockSummary struct {
		ChainId   common.ChainID
		Height    common.Height
		BlockHash common.Hash
		// since v1.5.0, the election result of the next committee whill be packaged together.
		// Because only the data and comm node will receive the report and record the next committee
		// of the sub chain. Since the new elected node has already been synchronizing the main chain,
		// it will not synchronize the data again, then it will not be able to synchronize all the sub
		// chain committee information, resulting in the nodes missing the corresponding information
		// when the new epoch begins.
		NextComm *EpochCommittee
		// V0's BlockSummary.Hash is same with blockhash, which can't reflect the location information
		// of the block, and can't complete the proof of cross chain. V1 adds chainid and height to hash
		Version uint16
	}

	EpochCommittee struct {
		Result *Committee // actual election results
		Real   *Committee // the final result, if Result.IsAvailable()==false, then Real is the actual Committee. Otherwise, it is nil
	}

	AttendanceRecord struct {
		Epoch      common.EpochNum // current epoch
		Attendance *big.Int        // Indicates by bit whether the corresponding data block is empty, Attendance.Bit(BlockNum)==1 is normal block and ==0 is empty block
		DataNodes  common.NodeIDs  // List of datanode nodeid in ascending order
		Stats      []int           // Stats of alive data nodes

		nodeIdxs map[common.NodeID]int // cache data node id -> index of Stats
	}

	RewardRequest struct {
		ChainId      common.ChainID
		CommitteePks [][]byte          // The public key list of the members of the current committee in the order of proposing
		Epoch        common.EpochNum   // Epoch where the reward is declared
		LastHeader   *BlockHeader      // The block header of the last block of the epoch declared
		Attendance   *AttendanceRecord // The attendance table of the last block, which contains the attendance records of the entire epoch
		PASs         PubAndSigs        // Signature list for the last block
	}

	RewardRequests []*RewardRequest

	// expireEra >= (Withdrawing.Demand + WithdrawDelayEras)
	// Withdrawing.Demand >= (DepositIndex.Era + MinDepositEras)
	Withdrawing struct {
		// since v2.11.0, change to the era of withdraw request execution, will cause the
		// generated withdraws to be delayed by one more WithdrawDelayEras.
		Demand common.EraNum `json:"demand"`
		// Withdraw amount, if nil, it means withdrawing all
		Amount *big.Int `json:"amount,omitempty"`
		// since v2.11.0, mining pool sub-account address
		PoolAddr *common.Address `json:"addr,omitempty"`
	}

	Withdrawings []*Withdrawing

	RRInfo struct {
		// The hash value of the NodeID of the node is used to store information in a more
		// private way. It can also reduce storage capacity
		NodeIDHash common.Hash
		// The main chain block height at the time of the last deposit
		Height common.Height
		// Which type of node, supports common.Consensus/common.Data
		Type common.NodeType
		// If it is not nil, it means that this deposit has been applied for withdrawing and
		// will no longer participate in the calculation. When the value >= current era, execute
		// the withdrawing. Redemption will be executed at the end of the era.
		WithdrawDemand *common.EraNum
		// Record the number of penalties, initially 0, +1 after each Penalty execution
		PenalizedTimes int
		// Depositing: sum of all the deposits of the node
		Amount *big.Int
		// The percentage of the effective pledge amount of the current node in the total
		// effective pledge. If it is nil, it indicates that the current pledge does not
		// account for the proportion. It may be waiting for withdrawing at this time.
		Ratio *big.Rat
		// Reward binding address
		RewardAddr common.Address
		// Since v1.3.4. When WithdrawDemand!=nil, record all pending withdrawing records. If it
		// exists, the withdrawing due in the list will be executed every era.
		Withdrawings Withdrawings
		// since v1.5.0. Version number, used for compatible
		Version uint16
		// since v1.5.0。Used to record a total of valid pledged consensus nodes, only valid
		// when Type==common.Consensus, others are 0
		NodeCount uint32
		// since v2.9.17, node status
		Status uint16
		// since v2.11.0, available amount of the node, use for election and settle
		Avail *big.Int
		// since v2.11.0, voted data node id hash
		Voted *common.Hash
		// since v2.11.0, voted amount of current data node
		VotedAmount *big.Int
		// since v2.11.0, if not nil means it's a pool node (only Type==common.Consensus supports pool mode)
		Settles *SettleInfo
	}

	SettleInfo struct {
		ChargeRatio *big.Rat     // pool owner charge ratio
		Root        *common.Hash // root hash of the trie build by all deposit from each account to this node (Address->SettleValue)
	}

	RRProofs struct {
		Info  *RRInfo
		Proof trie.ProofChain
	}

	NodeResult struct {
		NodeID     common.NodeID // The ID of the node participating in the election. For ManagedComm, only this field is needed, and the other fields are empty
		Sorthash   *common.Hash  // The result of the VRF algorithm
		Proof      []byte        // Proof of VRF algorithm results
		RRProof    *RRProofs     // The proof of the deposit of the nodes participating in the election
		FactorHash *common.Hash  // since2.0.0 The node declares the hash of the random factor participating in the seed calculation
		RandNum    uint32        // since version 100
	}

	NodeResults []*NodeResult

	// The compound data structure packed in the block, the memory and the form of the data set in the block
	ChainElectResult struct {
		ChainID common.ChainID  // Election chain
		Epoch   common.EpochNum // The Epoch where the election took place, the value of the pre-election is NilEpoch
		Results NodeResults
	}

	ChainElectResults []*ChainElectResult

	PreElectPhase byte

	PreElecting struct {
		// Chain of pre-election
		ChainID common.ChainID
		// Current execution stage
		Phase PreElectPhase
		// Seed of main chain when pre-electing
		Seed *common.Seed
		// Count the number of election retrys, because the election may not be successful, and the
		// election can be automatically started again (3 times in total)
		Count int
		// The height of the main chain when the pre-election starts. Because the Hash value of the
		// current block is required when creating PreElecting, it cannot be stored in the object and
		// needs to be obtained from the data node when synchronizing data
		Start common.Height
		// The Hash of the main chain height block at startup has a value in the cache and is nil in
		// the BlockBody
		CachedHash *common.Hash
		// When the new chain is a ManagedComm chain, NidHashes saves the hash values of all authorized
		// node IDs, which are the basis for the pre-election. The election type can also be judged
		// based on whether this field is empty
		NidHashes []common.Hash
		// Electing phase: the height of the main chain at which the pre-election ends;
		// Starting phase: the height of the main chain at which consensus is initiated
		Expire common.Height
	}

	SeedFactor []byte

	PreElectings []*PreElecting
)

func (pointer *RequestResult) ToStringArray() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.([]interface{})
	new := make([]string, len(result))
	for i, v := range result {
		new[i] = v.(string)
	}
	return new, nil
}

func (pointer *RequestResult) ToComplexString() (types.ComplexString, error) {
	if err := pointer.checkResponse(); err != nil {
		return "", err
	}
	result := (pointer).Result.(interface{})
	return types.ComplexString(result.(string)), nil
}

func (pointer *RequestResult) ToString() (string, error) {
	if err := pointer.checkResponse(); err != nil {
		return "", err
	}
	result := (pointer).Result.(interface{})
	return result.(string), nil
}

func (pointer *RequestResult) ToInt() (int64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}
	result := (pointer).Result.(interface{})
	hex := result.(string)
	numericResult, err := strconv.ParseInt(hex, 16, 64)
	return numericResult, err
}

func (pointer *RequestResult) ToBigInt() (*big.Int, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	res := (pointer).Result.(interface{})
	ret, success := big.NewInt(0).SetString(res.(string)[2:], 16)
	if !success {
		return nil, errors.New(fmt.Sprintf("Failed to convert %s to BigInt", res.(string)))
	}
	return ret, nil
}

func (pointer *RequestResult) ToComplexIntResponse() (types.ComplexIntResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return types.ComplexIntResponse(0), err
	}
	result := (pointer).Result.(interface{})
	var hex string
	switch v := result.(type) {
	// Testrpc returns a float64
	case float64:
		hex = strconv.FormatFloat(v, 'E', 16, 64)
		break
	default:
		hex = result.(string)
	}
	cleaned := strings.TrimPrefix(hex, "0x")
	return types.ComplexIntResponse(cleaned), nil
}

func (pointer *RequestResult) ToBoolean() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}
	result := (pointer).Result.(interface{})
	return result.(bool), nil
}

func (pointer *RequestResult) ToSignTransactionResponse() (*SignTransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	signTransactionResponse := &SignTransactionResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, signTransactionResponse)
	return signTransactionResponse, err
}

func (pointer *RequestResult) ToTransactionResponse() (*TransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	transactionResponse := &TransactionResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, transactionResponse)
	return transactionResponse, err
}

func (pointer *RequestResult) ToTransactionReceipt() (*TransactionReceipt, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	transactionReceipt := &TransactionReceipt{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, transactionReceipt)
	return transactionReceipt, err
}

func (pointer *RequestResult) ToBlock() (*Block, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	block := &Block{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, block)
	return block, err
}

func (pointer *RequestResult) ToSyncingResponse() (*SyncingResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	switch (pointer).Result.(type) {
	case bool:
		return &SyncingResponse{}, nil
	case map[string]interface{}:
		result = (pointer).Result.(map[string]interface{})
	default:
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	syncingResponse := &SyncingResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	json.Unmarshal(marshal, syncingResponse)
	return syncingResponse, nil
}

// To avoid a conversion of a nil interface
func (pointer *RequestResult) checkResponse() error {
	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}
	if pointer.Result == nil {
		return customerror.EMPTYRESPONSE
	}
	return nil
}
