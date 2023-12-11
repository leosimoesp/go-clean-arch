package sendtransfer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/lbsti/go-clean-arch/core/entity"
)

type TransationInputDTO struct {
	Date      time.Time `json:"date"`
	Sender    int64     `json:"sender"`
	Recipient int64     `json:"recipient"`
	Value     int64     `json:"value"` //value in cents
}

type TransationOutputDTO struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Status       string `json:"status"`
	ID           int64  `json:"id"`
}

type Authorizer interface {
	IsAuthorized(key string) (bool, error)
}

type Messager interface {
	Send(email, message string) error
}

type SendTransfer struct {
	transactionRepo entity.TransactionRepository
	authorizer      Authorizer
	messager        Messager
	mu              sync.Mutex
}

func NewSendTransfer(tRepo entity.TransactionRepository, auth Authorizer, msg Messager) *SendTransfer {
	return &SendTransfer{transactionRepo: tRepo, authorizer: auth, messager: msg}
}

func (s *SendTransfer) Execute(in TransationInputDTO) (TransationOutputDTO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sender, recipient, err := s.findSenderAndRecipient(in)

	if err != nil {
		return TransationOutputDTO{Status: "rejected"}, err
	}
	if e := sender.Debit(in.Value); e != nil {
		return TransationOutputDTO{Status: "rejected"}, e
	}
	if authorized, e := s.authorizer.IsAuthorized("5794d450-d2e2-4412-8131-73d0293ac1cc"); !authorized {
		return TransationOutputDTO{Status: "rejected"}, entity.AuthorizationDeniedErr
	} else if e != nil {
		return TransationOutputDTO{Status: "rejected"}, e
	}

	transaction := entity.NewTransaction()
	transaction.Sender = sender.ID
	transaction.Recipient = recipient.ID
	transaction.Value = in.Value
	transaction.Date = time.Now().UTC()

	if e := transaction.IsValid(); e != nil {
		return TransationOutputDTO{Status: "rejected"}, e
	}

	transactionID, err := s.transactionRepo.Save(*transaction)
	if err != nil {
		return TransationOutputDTO{Status: "rejected"}, err
	}

	if err := s.messager.Send(recipient.Email, "You have received a transfer"); err != nil {
		slog.Error("Error sending transaction message %+v\n", err)
	}
	return TransationOutputDTO{Status: "success", ID: transactionID}, nil
}

func (s *SendTransfer) findSenderAndRecipient(in TransationInputDTO) (*entity.User, *entity.User, error) {
	sender, err := s.transactionRepo.LoadUser(in.Sender)
	if err != nil {
		return nil, nil, err
	}
	recipient, err := s.transactionRepo.LoadUser(in.Recipient)
	if err != nil {
		return nil, nil, err
	}
	return sender, recipient, nil
}
