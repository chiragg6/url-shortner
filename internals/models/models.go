package models

import (
	"databse/sql"
)

type ShortnerData struct {
	OriginalURL, ShortendURL string
	Clicks                   int
}

type ShortnerDataModel struct {
	DB *sql.DB
}

func (m *ShortnerDataModel) Latest() ([]*ShortnerData, error) {
	stmt := `SELECT original_url, shortened_url, clicks FROM urls`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	urls := []*ShortnerData{}
	for rows.Next() {
		url := &ShortnerData{}
		err = rows.Scan(&url.OriginalURL, &url.ShortendURL, &url.Clicks)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil

}
