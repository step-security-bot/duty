package duty

import (
	gocontext "context"

	"github.com/flanksource/duty/context"
	"github.com/flanksource/duty/models"
	"github.com/flanksource/duty/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

// deprecated use query.CheckSummaryByID
func CheckSummary(ctx DBContext, checkID string) (*models.CheckSummary, error) {
	return query.CheckSummaryByID(context.NewContext(ctx).WithDB(ctx.DB(), ctx.Pool()), checkID)
}

// deprecated use query.CheckSummary
func QueryCheckSummary(ctx gocontext.Context, dbpool *pgxpool.Pool, opts ...query.CheckSummaryOptions) ([]models.CheckSummary, error) {
	return query.CheckSummary(context.NewContext(ctx).WithDB(nil, dbpool), opts...)
}
