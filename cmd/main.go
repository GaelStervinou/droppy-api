package main

import (
	"go-api/cmd/fixtures"
	"go-api/internal/repositories"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	repo := repositories.Setup(&wg)
	wg.Add(1)
	defer repo.Disconnect()

	fixtures.PopulateUsers(repo.UserRepository)
}
