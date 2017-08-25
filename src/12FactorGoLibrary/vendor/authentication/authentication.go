package authentication
//
//Author: C Neuhardt
//Last Updated: 8/17/2017
//Last Updated By: Ricky Clevinger
import (
	"net/http"
	"time"
	"fmt"
	"context"
	"jwt-go"
	"helper"
	"member"
)

type Key int

const MyKey Key = 0

// JWT schema of the data it will store
type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// create a JWT and put in the clients cookie
func SetToken(res http.ResponseWriter, req *http.Request) {
	
	mail := helper.HTMLClean(req.FormValue("email"))
	pass := helper.HTMLClean(req.FormValue("pass"))

	var members []member.Member

	members = member.MemberExist(mail,pass)
	role := members[0].Role

	if (len(members) > 0){
			
		expireToken := time.Now().Add(time.Hour * 1).Unix()
		expireCookie := time.Now().Add(time.Hour * 1)

		claims := Claims{
			role,
			jwt.StandardClaims{
				ExpiresAt: expireToken,
				Issuer:    "localhost:8080",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, _ := token.SignedString([]byte("secret"))

		cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
		http.SetCookie(res, &cookie)

		http.Redirect(res, req, "/index.html", 307)
	}else{
		http.Redirect(res, req, "/login.html", 307)
	}
	
}

// deletes the cookie
func Logout(res http.ResponseWriter, req *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(res, &deleteCookie)
	fmt.Fprintf(res, "Logged Out")
}

// only viewable if the client has a valid token
func ProtectedProfile(res http.ResponseWriter, req *http.Request) {
	claims, ok := req.Context().Value(MyKey).(Claims)
	if !ok {
		http.NotFound(res, req)
		return
	}

	fmt.Fprintf(res, "Hello %s", claims.Role)
}


// middleware to protect private pages
func Validate(page http.HandlerFunc, role string) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.NotFound(res, req)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(res, req)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid && claims.Role == role{
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			http.NotFound(res, req)
			return
		}
	})
}

// middleware to protect private pages
func Validate2(page http.HandlerFunc, role1, role2 string) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.NotFound(res, req)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(res, req)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid && claims.Role == role1 || claims.Role == role2{
			ctx := context.WithValue(req.Context(), MyKey, *claims)
			page(res, req.WithContext(ctx))
		} else {
			http.NotFound(res, req)
			return
		}
	})
}


