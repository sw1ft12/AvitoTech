package repository

import (
    "avitoTech/internal/models"
    "context"
    "database/sql"
    "errors"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
    "strings"
)

type TenderRepo struct {
    conn *pgxpool.Pool
}

func NewTenderRepo(conn *pgxpool.Pool) *TenderRepo {
    return &TenderRepo{conn: conn}
}

func ParseError(err error) error {
    var pgErr *pgconn.PgError

    if errors.As(err, &pgErr) {
        if strings.Contains(pgErr.Message, "NO_USER_FOUND") {
            return ErrUserNotFound
        } else if strings.Contains(pgErr.Message, "NO_TENDER_FOUND") {
            return ErrTenderNotFound
        } else if strings.Contains(pgErr.Message, "ACCESS_DENIED") {
            return ErrNoAccess
        }
    }
    return err
}

func (r *TenderRepo) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
    query := `SELECT id, name, description, type, created_by, version, created_at, status 
				FROM create_tender(@created_by, @name, @description, @type, @organization_id)`
    row, _ := r.conn.Query(ctx, query, pgx.NamedArgs{
        "name":            tender.Name,
        "description":     tender.Description,
        "type":            tender.ServiceType,
        "organization_id": tender.OrganizationId,
        "created_by":      tender.CreatorUsername,
    })

    tender, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[models.Tender])
    if err != nil {
        return models.Tender{}, ParseError(err)
    }

    return tender, nil
}

func (r *TenderRepo) GetTenders(ctx context.Context, limit, offset int, serviceType []string) ([]models.Tender, error) {
    query := `SELECT * FROM tender WHERE status='Published' ORDER BY name LIMIT @limit OFFSET @offset`
    if len(serviceType) != 0 {
        query = `SELECT * FROM tender WHERE status='Published' ORDER BY name AND type IN @arr LIMIT @limit OFFSET @offset`
    }
    rows, err := r.conn.Query(ctx, query, pgx.NamedArgs{
        "arr":    serviceType,
        "limit":  limit,
        "offset": offset,
    })

    if err != nil {
        return []models.Tender{}, ParseError(err)
    }

    tenders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tender])
    if err != nil {
        return []models.Tender{}, err
    }

    return tenders, nil
}

func (r *TenderRepo) GetTendersByUsername(ctx context.Context, limit, offset int, username string) ([]models.Tender, error) {
    query2 := `SELECT t1.* FROM tender t1 JOIN organization_responsible t2 ON t1.organization_id=t2.organization_id
    	WHERE user_id=get_user(@username) LIMIT @limit OFFSET @offset`
    rows, err := r.conn.Query(ctx, query2, pgx.NamedArgs{
        "username": username,
        "limit":    limit,
        "offset":   offset,
    })

    if err != nil {
        return []models.Tender{}, ParseError(err)
    }

    tenders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Tender])
    if err != nil {
        return []models.Tender{}, ParseError(err)
    }

    return tenders, nil
}

func (r *TenderRepo) GetTenderStatus(ctx context.Context, id, username string) (string, error) {
    query := `SELECT status FROM tender WHERE id=check_tender(@username, @id)`
    var s string
    err := r.conn.QueryRow(ctx, query, pgx.NamedArgs{
        "username": username,
        "id":       id,
    }).Scan(&s)
    if err != nil {
        return "", ParseError(err)
    }
    return s, nil
}

func NewNullString(s string) sql.NullString {
    if len(s) == 0 {
        return sql.NullString{}
    }
    return sql.NullString{
        String: s,
        Valid:  true,
    }
}

func (r *TenderRepo) UpdateTender(ctx context.Context, username string, tender models.Tender) (models.Tender, error) {
    query := `SELECT id, name, description, status, type, version, created_at FROM update_tender(@username, @id, @name, @description, @type, @status)`
    err := r.conn.QueryRow(ctx, query, pgx.NamedArgs{
        "username":    NewNullString(username),
        "id":          NewNullString(tender.Id),
        "name":        NewNullString(tender.Name),
        "description": NewNullString(tender.Description),
        "type":        NewNullString(tender.ServiceType),
        "status":      NewNullString(tender.Status),
    }).Scan(&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt)
    if err != nil {
        return models.Tender{}, ParseError(err)
    }
    return tender, nil
}
