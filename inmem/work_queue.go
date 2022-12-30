package inmem

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Lambels/go-csb"
)

// defaultBufSize represents the default buffer used as a queue to accumulate transactions.
//
// limited since most of the times the work queue will be used to process large requests
// for engage, and we dont want to spam engage.
const defaultBufSize int = 50

// WorkQueue represents an in memory implementation of a work queue.
//
// Note that this work queue implementation runs transactions synchronously using only
// one worker.
type WorkQueue struct {
	idCount int64

	done  chan struct{}
	queue chan *csb.Transaction

	// handler handels the message synchronously.
	handler func(*csb.Transaction) error

	statesMu sync.RWMutex
	states   map[int64]*state

	once sync.Once // used to close done only once.
}

// NewWorkQueue creates a new in memory work queue.
func NewWorkQueue(handler func(*csb.Transaction) error) *WorkQueue {
	w := &WorkQueue{
		done:    make(chan struct{}),
		queue:   make(chan *csb.Transaction, defaultBufSize),
		handler: handler,
		states:  make(map[int64]*state),
	}

	go w.listen()

	return w
}

// listen pulls transactions of the work queue and hands them in to the handler.
func (w *WorkQueue) listen() {
	for {
		select {
		case <-w.done:
			return
		case val := <-w.queue:
			select {
			case <-val.Ctx.Done():
				continue
			default:
			}

			state := w.states[val.Id]
			state.newSatus <- csb.Status{State: csb.Processing}

			status := csb.Status{State: csb.Done}
			if err := w.handler(val); err != nil {
				status.Error = err
			}

			state.newSatus <- status
		}
	}
}

// Publish pushes the transaction on the work queue, if the work queue is full it returns
// an error.
//
// If the work queue is closed, the call is no-op.
func (w *WorkQueue) Publish(transaction *csb.Transaction) error {
	w.statesMu.Lock()
	defer w.statesMu.Unlock()

	select {
	case <-w.done:
		return nil
	default:
	}

	w.idCount++
	transaction.Id = w.idCount

	s := &state{
		currStatus:    csb.Status{State: csb.Queued},
		subscriptions: make(map[int64]*Subscription),
	}
	w.states[transaction.Id] = s
	go s.bind(transaction) // bind the state to the transaction.

	select {
	case w.queue <- transaction:
		return nil
	default:
		return fmt.Errorf("publish: transaction queue is full")
	}
}

// Subscribe subscribes to the transaction with id = id.
//
// If the work queue is closed the call is no-op.
//
// If the transcation doesent exist ENOTFOUND is returned.
func (w *WorkQueue) Subscribe(ctx context.Context, id int64) (csb.Subscription, error) {
	w.statesMu.RLock()
	defer w.statesMu.RUnlock()

	select {
	case <-w.done:
		return nil, nil
	default:
	}

	state, ok := w.states[id]
	if !ok {
		return nil, csb.Errorf(csb.ENOTFOUND, "subscribe: no transaction was found with id: %v", id)
	}
	sub := &Subscription{
		state: state,
		c:     make(chan csb.Status, 1),
	}

	select {
	case state.newSub <- sub:
		return sub, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close closes the work queue, it is the callers responsability to cancel the context
// of the current transaction.
func (w *WorkQueue) Close() error {
	w.statesMu.Lock()
	defer w.statesMu.Unlock()

	w.once.Do(func() { close(w.done) })
	return nil
}

// Subscription represents
type Subscription struct {
	Id    int64
	state *state // parent state.
	c     chan csb.Status

	// used to synchronise calls from the client with calls from the parent state.
	closed atomic.Bool
}

// C returns a stream of status updates, the channel always has a status update
// when C is called indicating the current status of the message.
//
// When the channel is closed the previous status will indicate why.
func (s *Subscription) C() <-chan csb.Status {
	return s.c
}

// Close closes the subscription.
func (s *Subscription) Close() error {
	if s.closed.CompareAndSwap(false, true) {
		s.state.delSub <- s
		close(s.c)
	}

	return nil
}

// state ties a transaction to subscription, it constantly communicates with the transaction
// the current status and forwards the status to the subscriptions.
//
// It is responsible for handling new and perishing subscriptions as well as closing them
// when the transaction is complete or the context is cancelled.
type state struct {
	idCount int64

	// currStatus indicates the current status of the subscription.
	//
	// It is sent to new subscription by default.
	currStatus    csb.Status
	subscriptions map[int64]*Subscription

	// transaction is the transaction to which the state is binded to.
	transaction *csb.Transaction
	w           *WorkQueue // parent work queue.

	newSub   chan *Subscription
	delSub   chan *Subscription
	newSatus chan csb.Status // used to push new statuses.
}

// bind binds the state to the transaction and handels all status updates.
func (s *state) bind(transaction *csb.Transaction) {
	s.transaction = transaction
	defer func() {
		s.w.statesMu.Lock()
		defer s.w.statesMu.Unlock()

		s.closeSubscriptions()
		delete(s.w.states, s.transaction.Id)
	}()

	for {
		select {
		case <-transaction.Ctx.Done(): // transaction cancelled.
			// if we are currently in the processing state skip this case
			// since there will be a new status shortly detailing any error from the context
			// passed from the handler.
			if s.currStatus.State == csb.Processing {
				continue
			}

			s.currStatus = csb.Status{State: csb.Cancelled, Error: transaction.Ctx.Err()}
			s.broadcast()
			return
		case status := <-s.newSatus: // new status.
			s.currStatus = status
			s.broadcast()

			if status.State == csb.Done {
				return
			}
		case sub := <-s.newSub: // new subscriber.
			s.idCount++
			sub.Id = s.idCount
			s.subscriptions[sub.Id] = sub

			sub.c <- s.currStatus // c has buffer of 1, no point to handle blocking.
		case sub := <-s.delSub: // remove subscriber.
			delete(s.subscriptions, sub.Id)
		}
	}
}

// broadcast broadcasts the current status to all the subscribers.
func (s *state) broadcast() {
	for _, v := range s.subscriptions {
		select {
		case <-s.transaction.Ctx.Done():
			return
		case v.c <- s.currStatus:
		}
	}
}

// closeSubscriptions closes all subscribers, no-op if they are already closed.
func (s *state) closeSubscriptions() {
	for _, v := range s.subscriptions {
		if v.closed.CompareAndSwap(false, true) {
			close(v.c)
		}
	}
}
