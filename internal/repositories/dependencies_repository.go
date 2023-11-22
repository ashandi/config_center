package repositories

import (
	"config_center/api/types"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DependenciesRepository struct {
	db *sqlx.DB
}

func NewDependenciesRepository(db *sqlx.DB) *DependenciesRepository {
	return &DependenciesRepository{
		db: db,
	}
}

func (r *DependenciesRepository) FindByMajor(table, platform string, major int) (types.Dependency, error) {
	query := fmt.Sprintf(`
			SELECT t1.version, t1.hash
			FROM %s t1
			JOIN (
			  SELECT major, minor, MAX(patch) as max_patch
			  FROM %s
			  WHERE platform = ? AND major = ?
			  GROUP BY major, minor
			  ORDER BY minor DESC
			  LIMIT 1
			) AS t2
			ON t1.major = t2.major AND t1.minor = t2.minor AND t1.patch = t2.max_patch
		`,
		table,
		table,
	)

	var dep types.Dependency
	err := r.db.QueryRowx(query, platform, major).StructScan(&dep)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Dependency{}, nil
		}
		return types.Dependency{}, err
	}

	urls, err := r.getUrls(table + "_urls")
	if err != nil {
		return types.Dependency{}, err
	}

	dep.Urls = urls

	return dep, nil
}

func (r *DependenciesRepository) FindByMajorMinor(table, platform string, major, minor int) (types.Dependency, error) {
	query := fmt.Sprintf(`
			SELECT t1.version, t1.hash
			FROM %s t1
			JOIN (
			  SELECT major, minor, MAX(patch) as max_patch
			  FROM %s
			  WHERE platform = ? AND major = ? AND minor = ?
			  GROUP BY major, minor
			  LIMIT 1
			) AS t2
			ON t1.major = t2.major AND t1.minor = t2.minor AND t1.patch = t2.max_patch
		`,
		table,
		table,
	)

	var dep types.Dependency
	err := r.db.QueryRowx(query, platform, major, minor).StructScan(&dep)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Dependency{}, nil
		}
		return types.Dependency{}, err
	}

	urls, err := r.getUrls(table + "_urls")
	if err != nil {
		return types.Dependency{}, err
	}

	dep.Urls = urls

	return dep, nil
}

func (r *DependenciesRepository) FindByMajorMinorPatch(table, platform string, major, minor, patch int) (types.Dependency, error) {
	query := fmt.Sprintf(`
			SELECT version, hash
			FROM %s
			WHERE platform = ? AND major = ? AND minor = ? AND patch = ?
		`,
		table,
	)

	var dep types.Dependency
	err := r.db.QueryRowx(query, platform, major, minor, patch).StructScan(&dep)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Dependency{}, nil
		}
		return types.Dependency{}, err
	}

	urls, err := r.getUrls(table + "_urls")
	if err != nil {
		return types.Dependency{}, err
	}

	dep.Urls = urls

	return dep, nil
}

func (r *DependenciesRepository) getUrls(table string) ([]string, error) {
	query := fmt.Sprintf(`
			SELECT url
			FROM %s
		`,
		table,
	)

	rows, err := r.db.Queryx(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, nil
		}
		return []string{}, err
	}
	defer func() { _ = rows.Close() }()

	var urls []string
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return []string{}, err
		}

		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return []string{}, err
	}

	return urls, err
}
