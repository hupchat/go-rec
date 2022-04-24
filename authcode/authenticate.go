package main

import (
	"context"
	"fmt"
	spotify "go-rec"
	spotifyauth "go-rec/auth"
	"go-rec/helpers"
	"go-rec/lib/cache"
	"log"
	"net/http"
	"os"
	"sync"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI()), spotifyauth.WithScopes(spotifyauth.ScopeUserFollowRead, spotifyauth.ScopeUserLibraryModify))
	state = "go-rec-state"
)

func redirectURI() string {
	if helpers.IsDev() {
		return "http://localhost:" + os.Getenv("PORT") + "/callback"
	}
	return "https://go-rec-app.herokuapp.com/callback"
}

func main() {
	env := helpers.VerifyEnv([]string{
		"PORT",
	})
	if env != nil {
		log.Fatal(env)
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	createClientCache(r.Context(), client, user)
}

func createClientCache(ctx context.Context, client *spotify.Client, user *spotify.PrivateUser) {
	clientToken, _ := client.Token()
	var wg sync.WaitGroup
	wg.Add(2)
	go cache.SetHash(ctx, &wg, user.ID, "token", clientToken.AccessToken)
	go cache.SetHash(ctx, &wg, user.ID, "refresh_token", clientToken.RefreshToken)
	wg.Wait()
}
