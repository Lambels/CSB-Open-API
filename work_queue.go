package csb

import "context"

// Status represents the status of a transaction in the work queue.
type Status struct {
	// State of the transaction.
	//
	// Either: Queued, Processing, Done or Cancelled.
	State int `json:"state"`
	// Any error associated with the state. Should check for any error when the state is either
	// Done or Cancelled.
	Error error `json:"error"`
}

const (
	// Queued means the transaction is on a run queue.
	Queued = iota
	// Processing means the transaction is currently being worked on.
	Processing
	// Done means that the transaction was previously processing and now its done.
	//
	// This doesnt mean that the transaction was valid, just means that its done.
	Done
	// Cancelled means that the transaction was cancelled before it got a chance to run.
	Cancelled
)

// Transcation represents a transaction working through the work queue.
type Transaction struct {
	// Id of the transaction.
	Id int64 `json:"id"`
	// Data of the transaction.
	Data interface{} `json:"data"`
	// Ctx of the transaction, used to cancel the transaction.
	Ctx context.Context
}

// Subscription represents a closable one way flow of updates from the work queue service to the
// consumer.
type Subscription interface {
	// C returns the flow of status updates for the transaction you are subscribed to.
	C() <-chan Status

	// Close closes the stream of updates.
	Close() error
}

// WorkQueue represents a system where you publish your transaction to run in the future
// and have the option to subscribe to your transaction via the transaction id to get updates
// on the state of your transaction.
type WorkQueue interface {
	// Publish publishes a transaction to the work queue.
	Publish(transaction *Transaction) error

	// Subscribe returns a subscription to read updates on your transaction with id = id.
	//
	// You can have multiple subscribers on the same transaction.
	Subscribe(ctx context.Context, id int64) (Subscription, error)

	// Close closes the work queue.
	Close() error
}
