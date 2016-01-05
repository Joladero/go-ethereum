package filters

import (
	"testing"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/event"
)

func TestCallbacks(t *testing.T) {
	var mux event.TypeMux

	fs := NewFilterSystem(&mux)

	blockDone, txDone, logDone, removedLogDone := false, false, false, false
	blockFilter := &Filter{
		BlockCallback: func(*types.Block, vm.Logs) {
			blockDone = true
		},
	}
	txFilter := &Filter{
		TransactionCallback: func(*types.Transaction) {
			txDone = true
		},
	}
	logFilter := &Filter{
		LogCallback: func(*vm.Log, bool) {
			logDone = true
		},
	}

	removedLogFilter := &Filter{
		LogCallback: func(*vm.Log, bool) {
			removedLogDone = true
		},
	}

	fs.Add(blockFilter)
	fs.Add(txFilter)
	fs.Add(logFilter)
	fs.Add(removedLogFilter)

	mux.Post(core.ChainEvent{})
	mux.Post(core.TxPreEvent{})
	mux.Post(core.RemovedLogEvent{vm.Logs{&vm.Log{}}})
	mux.Post(vm.Logs{&vm.Log{}})

	if !blockDone {
		t.Error("block filter failed to trigger")
	}

	if !txDone {
		t.Error("transaction filter failed to trigger")
	}

	if !logDone {
		t.Error("log filter failed to trigger")
	}

	if !removedLogDone {
		t.Error("removed log filter failed to trigger")
	}
}
