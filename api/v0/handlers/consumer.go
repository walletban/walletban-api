package handlers

import (
	json2 "encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
	"walletban-api/api/v0/presenter"
	"walletban-api/internal/entities"
	"walletban-api/internal/services"
	"walletban-api/internal/utils"
)

// TODO: CreateWalletFunction
type ConsumserRequest struct {
	Token    string
	Password string
}

type UsernameRequest struct {
	Username string `json:"username"`
}

type ContractInvoke struct {
	ContractID        string        `json:"contract_id"`
	ContractFunction  string        `json:"contract_function"`
	SecretKey         string        `json:"secret_key"`
	UserId            string        `json:"user_id"`
	ContractArguments []interface{} `json:"contract_arguments"`
}

type UserkeyResult struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		GenerateAccountResponse struct {
			Username string `json:"username"`
			Res      struct {
				Links struct {
					Account struct {
						Href string `json:"href"`
					} `json:"account"`
					Effects struct {
						Href      string `json:"href"`
						Templated bool   `json:"templated"`
					} `json:"effects"`
					Ledger struct {
						Href string `json:"href"`
					} `json:"ledger"`
					Operations struct {
						Href      string `json:"href"`
						Templated bool   `json:"templated"`
					} `json:"operations"`
					Precedes struct {
						Href string `json:"href"`
					} `json:"precedes"`
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Succeeds struct {
						Href string `json:"href"`
					} `json:"succeeds"`
					Transaction struct {
						Href string `json:"href"`
					} `json:"transaction"`
				} `json:"_links"`
				CreatedAt      time.Time `json:"created_at"`
				EnvelopeXdr    string    `json:"envelope_xdr"`
				FeeAccount     string    `json:"fee_account"`
				FeeCharged     string    `json:"fee_charged"`
				FeeMetaXdr     string    `json:"fee_meta_xdr"`
				Hash           string    `json:"hash"`
				ID             string    `json:"id"`
				Ledger         int       `json:"ledger"`
				MaxFee         string    `json:"max_fee"`
				MemoType       string    `json:"memo_type"`
				OperationCount int       `json:"operation_count"`
				PagingToken    string    `json:"paging_token"`
				Preconditions  struct {
					Timebounds struct {
						MinTime string `json:"min_time"`
					} `json:"timebounds"`
				} `json:"preconditions"`
				ResultMetaXdr         string    `json:"result_meta_xdr"`
				ResultXdr             string    `json:"result_xdr"`
				Signatures            []string  `json:"signatures"`
				SourceAccount         string    `json:"source_account"`
				SourceAccountSequence string    `json:"source_account_sequence"`
				Successful            bool      `json:"successful"`
				ValidAfter            time.Time `json:"valid_after"`
			} `json:"res"`
			PrivateKey string `json:"private_key"`
			PublicKey  string `json:"public_key"`
		} `json:"GenerateAccountResponse"`
	} `json:"data"`
}

type InvokeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		InvokeContractResponse struct {
			Result string `json:"result"`
		} `json:"InvokeContractResponse"`
	} `json:"data"`
}

func RegisterConsumer(service services.ApplicationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data ConsumserRequest
		err := c.BodyParser(&data)
		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		consumer := entities.Consumer{}
		intVal, err := strconv.Atoi(data.Token)
		uintVal := uint(intVal)
		consumer.ID = uintVal
		userName := UsernameRequest{Username: "abc"}
		json, err := json2.Marshal(userName)
		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		resp, err := utils.MakePOSTRequest("https://backend.sorobix.xyz/api/account", json)
		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		var userkey UserkeyResult
		parsedData, err := json2.Marshal(resp)
		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		err = json2.Unmarshal(parsedData, &userkey)
		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		consumer.WalletGKey = userkey.Data.GenerateAccountResponse.PublicKey
		consumer.WalletEncryptedSKey = userkey.Data.GenerateAccountResponse.PrivateKey
		consumer.Token = data.Password
		respoli, err := service.ConsumerRepository.Update(c.Context(), consumer)

		if err != nil {
			return handleError(c, err, "consumer cannot be registered")
		}
		return c.JSON(presenter.Success(respoli, "user updated!"))
	}
}

func InvokeContract(service services.ApplicationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data ContractInvoke
		err := c.BodyParser(&data)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		consumer := entities.Consumer{}
		intVal, err := strconv.Atoi(data.UserId)
		uintVal := uint(intVal)
		consumer.ID = uintVal
		consumerUpdate, err := service.ConsumerRepository.FindOne(c.Context(), consumer)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		if consumerUpdate.Token != data.SecretKey {
			return handleError(c, errors.New("invalid authentication"), "consumer cannot invoke contract")
		}
		data.SecretKey = consumerUpdate.WalletEncryptedSKey

		json, err := json2.Marshal(data)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		resp, err := utils.MakePOSTRequest("https://backend.sorobix.xyz/api/invoke", json)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		var userkey InvokeResponse
		parsedData, err := json2.Marshal(resp)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		err = json2.Unmarshal(parsedData, &userkey)
		if err != nil {
			return handleError(c, err, "consumer cannot invoke contract")
		}
		//Make Post request and return the result
		return c.JSON(presenter.Success(userkey, "user updated!"))
	}
}
