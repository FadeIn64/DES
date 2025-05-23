package repositories

import (
	"DAS/internal/repositories/db"
	"DAS/models"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db      *pgxpool.Pool
	manager trm.Manager
}

func NewTeamRepository(db *pgxpool.Pool, manager trm.Manager) *TeamRepository {
	return &TeamRepository{db: db, manager: manager}
}

// GetByID возвращает команду по ID
func (r *TeamRepository) GetByID(ctx context.Context, teamKey int) (*models.Team, error) {
	q := db.New(r.db)
	dbTeam, err := q.GetTeamByID(ctx, int64(teamKey))
	if err != nil {
		return nil, err
	}

	return &models.Team{
		TeamKey:     int(dbTeam.TeamKey),
		Name:        dbTeam.Name,
		Description: dbTeam.Description,
		Color:       dbTeam.Color.String,
		Country:     dbTeam.Country.String,
	}, nil
}

// GetAll возвращает все команды, отсортированные по названию
func (r *TeamRepository) GetAll(ctx context.Context) ([]*models.Team, error) {
	q := db.New(r.db)
	dbTeams, err := q.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}

	teams := make([]*models.Team, 0, len(dbTeams))
	for _, dbTeam := range dbTeams {
		teams = append(teams, &models.Team{
			TeamKey:     int(dbTeam.TeamKey),
			Name:        dbTeam.Name,
			Description: dbTeam.Description,
			Color:       dbTeam.Color.String,
			Country:     dbTeam.Country.String,
		})
	}

	return teams, nil
}
