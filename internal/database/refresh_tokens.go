package database

import "time"

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DBStructure.RefreshTokens Methods --

func (db *DB) SaveRefreshToken(userID int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour), // 1 hour
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserForRefreshToken(token string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return User{}, ErrorUserDoesNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, ErrorUserDoesNotExist
	}

	user, err := db.GetUser(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
