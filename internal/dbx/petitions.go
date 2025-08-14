package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const petitionsTable = "petitions"

type GeoPoint struct {
	Lat float64
	Lng float64
}

type Petition struct {
	ID          uuid.UUID `db:"id"`
	CityID      uuid.UUID `db:"city_id"`
	CreatorID   uuid.UUID `db:"creator_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	Signatures  int       `db:"signatures"`
	Goal        int       `db:"goal"`
	Reply       string    `db:"reply"`
	EndDate     time.Time `db:"end_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PetitionsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPetitionsQ(db *sql.DB) PetitionsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Явно выбираем колонки + вычисляем lat/lng из geometry
	selectCols := []string{
		"id",
		"city_id",
		"creator_id",
		"title",
		"description",
		"status",
		"signatures",
		"goal",
		"reply",
		"end_date",
		"created_at",
		"updated_at",
	}

	return PetitionsQ{
		db:       db,
		selector: builder.Select(selectCols...).From(petitionsTable),
		inserter: builder.Insert(petitionsTable),
		updater:  builder.Update(petitionsTable),
		deleter:  builder.Delete(petitionsTable),
		counter:  builder.Select("COUNT(*) AS count").From(petitionsTable),
	}
}

func (q PetitionsQ) New() PetitionsQ {
	return NewPetitionsQ(q.db)
}

func (q PetitionsQ) Insert(ctx context.Context, input Petition) error {
	values := map[string]interface{}{
		"id":          input.ID,
		"city_id":     input.CityID,
		"creator_id":  input.CreatorID,
		"title":       input.Title,
		"description": input.Description,
		"status":      input.Status,
		"signatures":  input.Signatures,
		"goal":        input.Goal,
		"reply":       input.Reply,
		"end_date":    input.EndDate,
		"created_at":  input.CreatedAt,
		"updated_at":  input.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building inserter query for table %s: %w", petitionsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PetitionsQ) Get(ctx context.Context) (Petition, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Petition{}, fmt.Errorf("building selector query for table %s: %w", petitionsTable, err)
	}

	var p Petition
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	err = row.Scan(
		&p.ID,
		&p.CityID,
		&p.CreatorID,
		&p.Title,
		&p.Description,
		&p.Status,
		&p.Signatures,
		&p.Goal,
		&p.Reply,
		&p.EndDate,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

func (q PetitionsQ) Select(ctx context.Context) ([]Petition, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table %s: %w", petitionsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Petition
	for rows.Next() {
		var p Petition
		if err := rows.Scan(
			&p.ID,
			&p.CityID,
			&p.CreatorID,
			&p.Title,
			&p.Description,
			&p.Status,
			&p.Signatures,
			&p.Goal,
			&p.Reply,
			&p.EndDate,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, nil
}

type UpdatePetitionInput struct {
	Status  *string
	Reply   *string
	EndDate *time.Time
}

func (q PetitionsQ) Update(ctx context.Context, in UpdatePetitionInput) error {
	updates := map[string]interface{}{}

	if in.Reply != nil {
		updates["reply"] = *in.Reply
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}
	if in.EndDate != nil {
		updates["end_date"] = *in.EndDate
	}

	if len(updates) == 0 {
		return nil
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table %s: %w", petitionsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PetitionsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table %s: %w", petitionsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PetitionsQ) FilterID(id uuid.UUID) PetitionsQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})

	return q
}

func (q PetitionsQ) FilterCityID(cityID uuid.UUID) PetitionsQ {
	q.selector = q.selector.Where(sq.Eq{"city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})

	return q
}

func (q PetitionsQ) FilterCreatorID(userID uuid.UUID) PetitionsQ {
	q.selector = q.selector.Where(sq.Eq{"creator_id": userID})
	q.counter = q.counter.Where(sq.Eq{"creator_id": userID})
	q.updater = q.updater.Where(sq.Eq{"creator_id": userID})
	q.deleter = q.deleter.Where(sq.Eq{"creator_id": userID})

	return q
}

func (q PetitionsQ) FilterStatus(status string) PetitionsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})

	return q
}

func (q PetitionsQ) FilterStatusIn(statuses ...string) PetitionsQ {
	if len(statuses) == 0 {
		return q
	}

	q.selector = q.selector.Where(sq.Eq{"status": statuses})
	q.counter = q.counter.Where(sq.Eq{"status": statuses})
	q.updater = q.updater.Where(sq.Eq{"status": statuses})
	q.deleter = q.deleter.Where(sq.Eq{"status": statuses})

	return q
}

func (q PetitionsQ) TitleLike(s string) PetitionsQ {
	p := fmt.Sprintf("%%%s%%", s)
	q.selector = q.selector.Where("title ILIKE ?", p)
	q.counter = q.counter.Where("title ILIKE ?", p)

	return q
}

func (q PetitionsQ) FilterCreatedAt(t time.Time, after bool) PetitionsQ {
	query := "created_at > ?"
	if !after {
		query = "created_at < ?"
	}

	q.selector = q.selector.Where(query, t)
	q.counter = q.counter.Where(query, t)
	q.updater = q.updater.Where(query, t)
	q.deleter = q.deleter.Where(query, t)

	return q
}

func (q PetitionsQ) FilterEndDate(t time.Time, after bool) PetitionsQ {
	query := "end_date > ?"
	if !after {
		query = "end_date < ?"
	}
	q.selector = q.selector.Where(query, t)
	q.counter = q.counter.Where(query, t)
	q.updater = q.updater.Where(query, t)
	q.deleter = q.deleter.Where(query, t)

	return q
}

func (q PetitionsQ) OrderByCreated(ascending bool) PetitionsQ {
	if ascending {
		q.selector = q.selector.OrderBy("created_at ASC")
	} else {
		q.selector = q.selector.OrderBy("created_at DESC")
	}

	return q
}

func (q PetitionsQ) OrderBySignatures(ascending bool) PetitionsQ {
	if ascending {
		q.selector = q.selector.OrderBy("signatures ASC")
	} else {
		q.selector = q.selector.OrderBy("signatures DESC")
	}

	return q
}

func (q PetitionsQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for table %s: %w", petitionsTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}

	return count, err
}

func (q PetitionsQ) Page(limit, offset uint64) PetitionsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)

	return q
}

func (q PetitionsQ) applyCondition(cond sq.Sqlizer) PetitionsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
