package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// PopulateDB inserts some test users into the database.
func (db *Database) PopulateUsers() error {
	ctx := context.Background()

	users := []struct {
		Username   string
		Email      string
		FirstName  string
		LastName   string
		Gender     string
		Age        string
		Avatar     string
		Identifier string
		Password   string
	}{
		{
			Username:   "alice",
			Email:      "alice@example.com",
			FirstName:  "Alice",
			LastName:   "Smith",
			Gender:     "female",
			Age:        "25",
			Avatar:     "A",
			Identifier: "alice@example.com",
			Password:   "p",
		},
		{
			Username:   "bob",
			Email:      "bob@example.com",
			FirstName:  "Bob",
			LastName:   "Jones",
			Gender:     "male",
			Age:        "30",
			Avatar:     "B",
			Identifier: "bob@example.com",
			Password:   "ps",
		},
		{
			Username:   "charlie",
			Email:      "charlie@example.com",
			FirstName:  "Charlie",
			LastName:   "Brown",
			Gender:     "male",
			Age:        "28",
			Avatar:     "C",
			Identifier: "charlie@example.com",
			Password:   "p",
		},
		{
			Username:   "diana",
			Email:      "diana@example.com",
			FirstName:  "Diana",
			LastName:   "Prince",
			Gender:     "female",
			Age:        "27",
			Avatar:     "D",
			Identifier: "diana@example.com",
			Password:   "p",
		},
		{
			Username:   "edward",
			Email:      "edward@example.com",
			FirstName:  "Edward",
			LastName:   "Norton",
			Gender:     "male",
			Age:        "35",
			Avatar:     "E",
			Identifier: "edward@example.com",
			Password:   "p",
		},
		{
			Username:   "fiona",
			Email:      "fiona@example.com",
			FirstName:  "Fiona",
			LastName:   "Gallagher",
			Gender:     "female",
			Age:        "22",
			Avatar:     "F",
			Identifier: "fiona@example.com",
			Password:   "shameless",
		},
		{
			Username:   "george",
			Email:      "george@example.com",
			FirstName:  "George",
			LastName:   "Miller",
			Gender:     "male",
			Age:        "40",
			Avatar:     "G",
			Identifier: "george@example.com",
			Password:   "madmax",
		},
		{
			Username:   "hannah",
			Email:      "hannah@example.com",
			FirstName:  "Hannah",
			LastName:   "Baker",
			Gender:     "female",
			Age:        "19",
			Avatar:     "H",
			Identifier: "hannah@example.com",
			Password:   "thirteen",
		},
		{
			Username:   "ian",
			Email:      "ian@example.com",
			FirstName:  "Ian",
			LastName:   "Curtis",
			Gender:     "male",
			Age:        "29",
			Avatar:     "I",
			Identifier: "ian@example.com",
			Password:   "joydivision",
		},
		{
			Username:   "julia",
			Email:      "julia@example.com",
			FirstName:  "Julia",
			LastName:   "Roberts",
			Gender:     "female",
			Age:        "33",
			Avatar:     "J",
			Identifier: "julia@example.com",
			Password:   "prettywoman",
		},
		{
			Username:   "kevin",
			Email:      "kevin@example.com",
			FirstName:  "Kevin",
			LastName:   "Hart",
			Gender:     "male",
			Age:        "38",
			Avatar:     "K",
			Identifier: "kevin@example.com",
			Password:   "funnyguy",
		},
		{
			Username:   "lisa",
			Email:      "lisa@example.com",
			FirstName:  "Lisa",
			LastName:   "Simpson",
			Gender:     "female",
			Age:        "8",
			Avatar:     "L",
			Identifier: "lisa@example.com",
			Password:   "saxophone",
		},
		{
			Username:   "michael",
			Email:      "michael@example.com",
			FirstName:  "Michael",
			LastName:   "Jordan",
			Gender:     "male",
			Age:        "57",
			Avatar:     "M",
			Identifier: "michael@example.com",
			Password:   "air23",
		},
		{
			Username:   "nina",
			Email:      "nina@example.com",
			FirstName:  "Nina",
			LastName:   "Dobrev",
			Gender:     "female",
			Age:        "35",
			Avatar:     "N",
			Identifier: "nina@example.com",
			Password:   "vampire",
		},
		{
			Username:   "oliver",
			Email:      "oliver@example.com",
			FirstName:  "Oliver",
			LastName:   "Queen",
			Gender:     "male",
			Age:        "32",
			Avatar:     "O",
			Identifier: "oliver@example.com",
			Password:   "arrow",
		},
		{
			Username:   "paula",
			Email:      "paula@example.com",
			FirstName:  "Paula",
			LastName:   "Patton",
			Gender:     "female",
			Age:        "41",
			Avatar:     "P",
			Identifier: "paula@example.com",
			Password:   "mission",
		},
		{
			Username:   "quentin",
			Email:      "quentin@example.com",
			FirstName:  "Quentin",
			LastName:   "Tarantino",
			Gender:     "male",
			Age:        "61",
			Avatar:     "Q",
			Identifier: "quentin@example.com",
			Password:   "pulpfiction",
		},
		{
			Username:   "rachel",
			Email:      "rachel@example.com",
			FirstName:  "Rachel",
			LastName:   "Green",
			Gender:     "female",
			Age:        "34",
			Avatar:     "R",
			Identifier: "rachel@example.com",
			Password:   "fashion",
		},
		{
			Username:   "steven",
			Email:      "steven@example.com",
			FirstName:  "Steven",
			LastName:   "Spielberg",
			Gender:     "male",
			Age:        "77",
			Avatar:     "S",
			Identifier: "steven@example.com",
			Password:   "jaws",
		},
		{
			Username:   "tina",
			Email:      "tina@example.com",
			FirstName:  "Tina",
			LastName:   "Turner",
			Gender:     "female",
			Age:        "83",
			Avatar:     "T",
			Identifier: "tina@example.com",
			Password:   "simplythebest",
		},
	}

	for _, u := range users {
		hashedPassword, err := security.HashPassword(u.Password)
		if err != nil {
			log.Printf("error hashing password for %s: %v", u.Username, err)
			continue
		}

		req := models.RegisterDbRequest{
			Ctx:          ctx,
			Username:     &u.Username,
			Email:        &u.Email,
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			Gender:       &u.Gender,
			Age:          &u.Age,
			Avatar:       u.Avatar,
			Identifier:   u.Identifier,
			PasswordHash: hashedPassword,
		}

		res, err := db.AddAuthUser(req)
		if err != nil {
			log.Printf("error adding user %s: %v", u.Username, err)
			continue
		}

		fmt.Printf("Inserted user %s with ID %d\n", res.RegisterResponseUser.UserName, res.Id)
	}

	return nil
}

// PopulatePosts inserts 20 demo posts into the database.
func (db *Database) PopulatePosts() error {

	ctx := context.Background()

	// Example post titles and bodies (feel free to expand or randomize)
	postTitles := []string{
		"Go Concurrency Explained",
		"Understanding SQL Transactions",
		"Building a REST API in Go",
		"Why Context is Important in Go",
		"How to Write Unit Tests",
		"PostgreSQL Tips and Tricks",
		"Working with JSON in Go",
		"Building a Simple CLI Tool",
		"Exploring Interfaces in Go",
		"Using Channels for Communication",
		"Optimizing Database Queries",
		"Error Handling Best Practices",
		"Deploying Go Apps to Production",
		"Explaining Goroutines",
		"Go vs Rust: A Quick Comparison",
		"How to Hash Passwords Securely",
		"Creating Middleware in Go",
		"Understanding defer, panic, and recover",
		"Building a WebSocket Server",
		"Working with Files in Go",
	}

	postBodies := []string{
		"This post explains the concept of concurrency in Go using goroutines and channels.",
		"A deep dive into database transactions, isolation levels, and best practices.",
		"Step-by-step guide on how to create a REST API using the net/http package.",
		"Context helps control cancellations, deadlines, and passing values through requests.",
		"Learn how to write clean and maintainable unit tests for your Go applications.",
		"Improve your database performance with these PostgreSQL optimization tips.",
		"An example of working with JSON serialization and deserialization in Go.",
		"A guide to building a command-line tool using flags and Cobra.",
		"Interfaces are a core feature of Go. This post explores their use cases.",
		"A practical example of using channels to synchronize goroutines.",
		"Tips on writing faster queries and reducing N+1 problems.",
		"Learn to handle errors in Go idiomatically with wrapped errors.",
		"Deploying Go applications on Docker and Kubernetes made easy.",
		"Goroutines are lightweight threads. This post explores their internals.",
		"A high-level comparison of Go and Rust for backend development.",
		"Learn how to hash passwords securely using bcrypt.",
		"Build your own HTTP middleware for logging and authentication.",
		"How to safely use defer, panic, and recover for error handling.",
		"A hands-on guide to building a simple WebSocket chat server.",
		"Reading and writing files safely and efficiently in Go.",
	}

	// User IDs (from your SELECT query)
	userIDs := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < len(postTitles); i++ {
		postCreator := userIDs[rng.Intn(len(userIDs))]
		title := postTitles[i]
		body := postBodies[i]

		// Randomly assign 0–1 categories (IDs 1, 2, 3 for demo)
		var categories []int
		if rng.Intn(2) == 1 { // 50% chance
			categories = []int{rng.Intn(3) + 1}
		}

		id, createdAt, err := db.CreatePost(ctx, title, body, postCreator, categories)
		if err != nil {
			log.Printf("error creating post %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Created post %d by user %d at %s\n", id, postCreator, createdAt.Format(time.RFC3339))
	}

	return nil
}

// PopulateConversations creates several direct message conversations with Alice.
func (db *Database) PopulateConversations() error {
	ctx := context.Background()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Alice's user ID
	const aliceID int64 = 1

	// Other users that will start conversations with Alice
	senders := []int64{2, 3, 5, 6, 7, 8, 9, 10}

	// Example messages
	openers := []string{
		"Hey Alice, how are you?",
		"Hi Alice! Long time no see.",
		"Hello Alice, I wanted to ask you something.",
		"Good morning Alice!",
		"Hey! Are you free this weekend?",
	}
	replies := []string{
		"I'm good, thanks!",
		"Sure, what’s up?",
		"That sounds interesting.",
		"Yes, I have some time to chat.",
		"Let’s plan something soon!",
	}

	for _, sender := range senders {
		// First message starts the conversation (conversationID = 0)
		convID, msgID, createdAt, err := db.CreateMessage(ctx, 0, sender, aliceID, openers[rng.Intn(len(openers))])
		if err != nil {
			log.Printf("error creating initial message from user %d: %v", sender, err)
			continue
		}
		fmt.Printf("Started conversation %d (msg %d) between user %d -> Alice at %s\n",
			convID, msgID, sender, createdAt.Format(time.RFC3339))

		// Add 3-4 follow-up messages in the same conversation
		numFollowUps := rng.Intn(2) + 3 // 3 or 4
		currentSender := aliceID

		for i := 0; i < numFollowUps; i++ {
			// Alternate sender/receiver
			receiver := sender
			if currentSender == sender {
				receiver = aliceID
			}

			text := replies[rng.Intn(len(replies))]

			_, msgID, createdAt, err := db.CreateMessage(ctx, convID, currentSender, receiver, text)
			if err != nil {
				log.Printf("error creating follow-up message in conversation %d: %v", convID, err)
				break
			}
			fmt.Printf("Added message %d to conversation %d (%d -> %d) at %s\n",
				msgID, convID, currentSender, receiver, createdAt.Format(time.RFC3339))

			// Alternate sender for next message
			if currentSender == aliceID {
				currentSender = sender
			} else {
				currentSender = aliceID
			}
		}
	}

	return nil
}
