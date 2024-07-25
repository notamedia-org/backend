package register

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-pg/pg"
	tokenHandler "github.com/notamedia-org/backend/internal/api/jwt"
	user "github.com/notamedia-org/backend/internal/api/user"
	"github.com/notamedia-org/backend/internal/config"
)

// Register godoc
// @Summary Create user
// @Description Create a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body user.User true "User"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Router /api/v1/user/register [post]
func Register(db *pg.DB, env *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Unable to read request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		var u user.User
		if err = json.Unmarshal(body, &u); err != nil {
			http.Error(w, "Unable to parse JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		token, err := tokenHandler.CreateToken(env)
		if err != nil {
			http.Error(w, "Unable to create jwt token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		u.GetUUID()
		if err := u.CreateUser(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userToken := &user.UserToken{Token: token}
		userToken.GetUUID()
		if err := userToken.CreateUserToken(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
	}
}
