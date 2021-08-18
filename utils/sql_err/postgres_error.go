package sql_err

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/muchlist/sagasql/utils/rest_err"
	"strings"
)

func ParseError(err error) rest_err.APIError {
	if err == pgx.ErrNoRows {
		return rest_err.NewBadRequestError(fmt.Sprintf("tidak ada data yang sesuai dengan id yang diberikan"))
	}

	if strings.Contains(err.Error(), "23505") {
		return rest_err.NewBadRequestError("input yang diberikan mengalami konflik dengan data existing")
	}

	if strings.Contains(err.Error(), "42703") {
		return rest_err.NewInternalServerError("galat pada query database", err)
	}

	if strings.Contains(err.Error(), "must equal") {
		return rest_err.NewInternalServerError("galat pada query database", err)
	}

	return rest_err.NewInternalServerError("galat", err)
}
