package login

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"
	tokenHandler "github.com/notamedia-org/backend/internal/api/jwt"
	"github.com/notamedia-org/backend/internal/api/user"
	"github.com/notamedia-org/backend/internal/config"
)

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.User true "User"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Router /api/v1/user/register [post]
func Auth(db *pg.DB, env *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		tokenString := req.Header.Get("Authorization")
		if status, err := tokenHandler.VerifyToken(env, tokenString); status != true || err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		uToken := user.UserToken{Token: tokenString}
		user, err := uToken.FindUserByToken(db)
		if err != nil {
			http.Error(w, "Couldnt find user", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Ошибка преобразования в JSON", http.StatusInternalServerError)
			return
		}

		w.Write(userJson)
	}
}
