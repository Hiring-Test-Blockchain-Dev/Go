# VeriFi AI

Wallet-native finance with optional on-chain loans. This repo is a **10–15 minute Go hiring screen**.

## Your task

The loan registry can create, approve, and disburse a loan, but it cannot record **partial repayments**. Finish that in `loan/registry.go`.

Do not edit the tests.

```bash
go test ./...
```

You should see failing tests under `RecordPayment` / `RemainingOf`. Make them pass.

### Implement

1. **`RecordPayment(caller string, loanID uint64, amountCents uint64) error`**
   - Platform **or** the borrower may call it
   - Loan must be `StatusDisbursed`
   - `amountCents` must be `> 0` and `<=` remaining
   - Decrease remaining
   - If remaining hits 0, set status to `StatusRepaid`
   - Append a `Payment` to `Registry.Payments` (`LoanID`, `Payer`, `AmountCents`, `RemainingCents`)
   - On failure, return one of the sentinel errors in `loan/errors.go` and do not append a payment

2. **`RemainingOf(loanID uint64) (uint64, error)`**
   - Return `ErrLoanNotFound` if the loan does not exist

3. **Fix `MarkDefaulted`**
   - It currently defaults an `Approved` loan. Only `Disbursed` is allowed.
   - Anything else (including `Repaid`) must return `ErrInvalidStatusTransition`

Use the sentinel errors already in the package. Do not change function signatures, the `Loan` struct, the `Payment` struct, or those errors.

## Submit

Work on your own branch and commit locally. Do not push.

```bash
git checkout -b <your-name>/verifi-loan
git add loan/registry.go
git commit -m "Record partial loan repayments"
```

## Optional: run the app

```bash
cp .env.example .env
npm install
npm run db:push
npm run db:seed
npm run dev
```

On Windows PowerShell, create the env file with `Copy-Item .env.example .env`.

App: http://localhost:5173 · API: http://localhost:3001
