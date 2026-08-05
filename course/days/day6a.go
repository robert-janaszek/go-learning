package days

import (
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"

	_ "github.com/lib/pq"

	"github.com/robert-janaszek/go-learning/course/config"
	userPackage "github.com/robert-janaszek/go-learning/course/user"
)

func Day6a() {
	// ex 1
	cfg := config.AppConfig{
		Port: 1234,
	}

	fmt.Println(cfg)

	// ex 2
	config.Load()
	// config.parseEnv() -- name parseEnv not exported by package config

	// ex 3
	u := userPackage.User{
		// email: "test@test1", -- cannot refer to unexported field email in struct literal of type user.User
	}
	u.SetEmail("test@test.com")
	fmt.Println(u.Email())

	// ex 4
	crandText := crand.Text()
	mrandFloat := mrand.Float64()

	fmt.Printf("crand %s, mrand %f\n", crandText, mrandFloat)
}
