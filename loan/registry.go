package loan

import (
	"errors"
	"time"
)

// Status is the lifecycle of a loan. Do not renumber existing values.
type Status int

const (
	StatusPending Status = iota
	StatusApproved
	StatusDisbursed
	StatusRepaid
	StatusDefaulted
	StatusRejected
)

// Loan is one registered loan. Do not add, remove, or reorder fields.
type Loan struct {
	OffChainID          string
	Borrower            string
	AmountCents         uint64
	InterestRateBps     uint64
	TermMonths          uint64
	MonthlyPaymentCents uint64
	Status              Status
	CreatedAt           time.Time
	ApprovedAt          time.Time
	DisbursedAt         time.Time
}

// Payment is one successful repayment. RecordPayment appends these to Registry.Payments.
type Payment struct {
	LoanID         uint64
	Payer          string
	AmountCents    uint64
	RemainingCents uint64
}

// Registry is an in-memory loan book.
// Hiring task (10–15 min): implement RecordPayment and RemainingOf, and fix MarkDefaulted.
// Instructions are in the repo README under "Your task". Do not edit the tests.
type Registry struct {
	owner      string
	platform   string
	nextLoanID uint64
	loans      map[uint64]*Loan
	remaining  map[uint64]uint64
	byOffChain map[string]uint64
	byBorrower map[string][]uint64

	// Payments is the repayment log, in call order.
	Payments []Payment
}

// New returns a registry owned by owner. platform may create and transition loans.
func New(owner, platform string) (*Registry, error) {
	if owner == "" || platform == "" {
		return nil, ErrInvalidAddress
	}
	return &Registry{
		owner:      owner,
		platform:   platform,
		nextLoanID: 1,
		loans:      make(map[uint64]*Loan),
		remaining:  make(map[uint64]uint64),
		byOffChain: make(map[string]uint64),
		byBorrower: make(map[string][]uint64),
	}, nil
}

// CreateLoan registers a pending loan. caller must be the owner or the platform.
func (r *Registry) CreateLoan(caller, offChainID, borrower string, amountCents, interestRateBps, termMonths, monthlyPaymentCents uint64) (uint64, error) {
	if !r.isPlatform(caller) {
		return 0, ErrUnauthorized
	}
	if borrower == "" {
		return 0, ErrInvalidAddress
	}
	if amountCents == 0 || termMonths == 0 {
		return 0, ErrInvalidAmount
	}
	if _, exists := r.byOffChain[offChainID]; exists {
		return 0, ErrLoanExists
	}

	id := r.nextLoanID
	r.nextLoanID++
	r.loans[id] = &Loan{
		OffChainID:          offChainID,
		Borrower:            borrower,
		AmountCents:         amountCents,
		InterestRateBps:     interestRateBps,
		TermMonths:          termMonths,
		MonthlyPaymentCents: monthlyPaymentCents,
		Status:              StatusPending,
		CreatedAt:           time.Now().UTC(),
	}
	r.byOffChain[offChainID] = id
	r.byBorrower[borrower] = append(r.byBorrower[borrower], id)
	r.remaining[id] = amountCents
	return id, nil
}

// ApproveLoan moves Pending -> Approved.
func (r *Registry) ApproveLoan(caller string, loanID uint64) error {
	if !r.isPlatform(caller) {
		return ErrUnauthorized
	}
	if err := r.transition(loanID, StatusPending, StatusApproved); err != nil {
		return err
	}
	r.loans[loanID].ApprovedAt = time.Now().UTC()
	return nil
}

// DisburseLoan moves Approved -> Disbursed.
func (r *Registry) DisburseLoan(caller string, loanID uint64) error {
	if !r.isPlatform(caller) {
		return ErrUnauthorized
	}
	if err := r.transition(loanID, StatusApproved, StatusDisbursed); err != nil {
		return err
	}
	r.loans[loanID].DisbursedAt = time.Now().UTC()
	return nil
}

// RejectLoan moves Pending -> Rejected.
func (r *Registry) RejectLoan(caller string, loanID uint64) error {
	if !r.isPlatform(caller) {
		return ErrUnauthorized
	}
	return r.transition(loanID, StatusPending, StatusRejected)
}

// MarkRepaid moves Disbursed -> Repaid and clears the outstanding balance.
func (r *Registry) MarkRepaid(caller string, loanID uint64) error {
	if !r.isPlatform(caller) {
		return ErrUnauthorized
	}
	if err := r.transition(loanID, StatusDisbursed, StatusRepaid); err != nil {
		return err
	}
	r.remaining[loanID] = 0
	return nil
}

// MarkDefaulted moves a disbursed loan to Defaulted.
// BUG: also allows Defaulted from Approved. Only Disbursed should work.
func (r *Registry) MarkDefaulted(caller string, loanID uint64) error {
	if !r.isPlatform(caller) {
		return ErrUnauthorized
	}
	ln, ok := r.loans[loanID]
	if !ok || ln.Borrower == "" {
		return ErrLoanNotFound
	}
	if ln.Status != StatusDisbursed && ln.Status != StatusApproved {
		return ErrInvalidStatusTransition
	}
	ln.Status = StatusDefaulted
	return nil
}

var errTODO = errors.New("TODO")

// RecordPayment records a repayment against a disbursed loan. See README "Your task".
func (r *Registry) RecordPayment(caller string, loanID uint64, amountCents uint64) error {
	// TODO: candidate implements
	return errTODO
}

// RemainingOf returns the outstanding balance in cents.
func (r *Registry) RemainingOf(loanID uint64) (uint64, error) {
	// TODO: candidate implements
	return 0, errTODO
}

// GetLoan returns a copy of the loan.
func (r *Registry) GetLoan(loanID uint64) (Loan, error) {
	ln, ok := r.loans[loanID]
	if !ok || ln.Borrower == "" {
		return Loan{}, ErrLoanNotFound
	}
	return *ln, nil
}

func (r *Registry) isPlatform(caller string) bool {
	return caller == r.platform || caller == r.owner
}

func (r *Registry) transition(loanID uint64, expected, next Status) error {
	ln, ok := r.loans[loanID]
	if !ok || ln.Borrower == "" {
		return ErrLoanNotFound
	}
	if ln.Status != expected {
		return ErrInvalidStatusTransition
	}
	ln.Status = next
	return nil
}
