package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/samualhalder/go-social/internal/store"
)

var usernames = []string{
	"SwiftTiger821",
	"SilentFalcon504",
	"CleverWizard192",
	"BraveSamurai377",
	"WittyPanther946",
	"WildNinja623",
	"LuckyKnight230",
	"GentlePhoenix784",
	"JollyDragon105",
	"FierceRobot699",
	"RapidCobra482",
	"BrightWolf316",
	"ChillRaven258",
	"MightyGhost730",
	"SharpBear194",
	"CoolMonkey565",
	"EpicHawk811",
	"SlyShark089",
	"NimbleFox327",
	"ZanyKoala471",
}
var titles = []string{
	"Scaling Microservices with Go and Kubernetes",
	"Understanding Prisma Migrations in Production",
	"Debugging EPERM Errors on Windows with Go",
	"Building REST APIs with Gin vs Fiber",
	"TypeScript vs JavaScript: When to Choose What",
	"Optimizing MongoDB Queries for Real-Time Apps",
	"Setting Up Prettier for Seamless VS Code Workflow",
	"Mastering Go Concurrency for Web Backends",
	"Handling Dirty Database States Gracefully",
	"The Hidden Costs of Inefficient API Design",
	"Unit Testing Node.js Apps the Right Way",
	"What Every Full-Stack Dev Should Know About CORS",
	"Implementing Secure Authentication with JWT",
	"Using Makefiles to Automate Go Dev Tasks",
	"Monitoring Go Services with Prometheus",
	"Exploring Go’s Standard Library: Gems and Gotchas",
	"Creating Reusable Components in React with TypeScript",
	"Migrating Legacy APIs to Express.js",
	"Understanding Market Demand for Go Frameworks",
	"Secrets to Clean Code in Full-Stack Projects",
}
var tags = []string{
	"golang",
	"nextjs",
	"prisma",
	"typescript",
	"mongodb",
	"expressjs",
	"jwt",
	"vscode",
	"makefile",
	"gin",
	"fiber",
	"nodejs",
	"reactjs",
	"api",
	"rest",
	"fullstack",
	"backend",
	"frontend",
	"debugging",
	"performance",
}
var contents = []string{
	"Learn how to scale Go-based microservices using Kubernetes orchestration.",
	"Dive into Prisma's migration system and handle schema changes safely.",
	"Troubleshoot EPERM permission errors on Windows when using Go apps.",
	"Compare Gin and Fiber frameworks for building RESTful APIs efficiently.",
	"Understand use cases for TypeScript vs JavaScript in modern dev stacks.",
	"Improve MongoDB performance by applying advanced query optimization.",
	"Configure Prettier in VS Code for a streamlined coding experience.",
	"Use goroutines and channels to implement high-concurrency Go backends.",
	"Detect and fix dirty database states during automated Go migrations.",
	"Identify bottlenecks in API design and how to resolve them cleanly.",
	"Write effective unit tests in Node.js using best practices and tools.",
	"Handle cross-origin requests correctly to avoid CORS-related issues.",
	"Secure web apps with JWT authentication and token verification flows.",
	"Use Makefiles to automate builds, tests, and routine Go development.",
	"Set up Prometheus to monitor performance metrics of Go applications.",
	"Explore powerful, lesser-known packages in Go’s standard library.",
	"Build modular React components using TypeScript for safety and reuse.",
	"Refactor legacy APIs and migrate to modern Express.js architecture.",
	"Analyze trends in Go frameworks and choose based on market relevance.",
	"Discover habits and patterns that lead to clean full-stack codebases.",
}
var commentsTmp = []string{
	"This post nailed the core challenges with Prisma migrations—great insights!",
	"Gin vs Fiber comparison helped clarify performance tradeoffs, thanks!",
	"Loved the breakdown on Go concurrency—goroutines explained beautifully.",
	"Using Makefiles for Go tasks was a game changer in my workflow.",
	"The JWT section made authentication finally click for me!",
	"Appreciate the real-world EPERM fix—saved me hours of debugging.",
	"Great walkthrough on setting up Prettier—code feels cleaner already.",
	"MongoDB optimization tips were exactly what my app needed.",
	"Sharp take on TypeScript vs JavaScript—concise and practical.",
	"React component patterns made my UI dev so much smoother.",
	"Prometheus setup steps were clear—monitoring feels less intimidating now.",
	"Your REST API advice should be mandatory reading for junior devs!",
	"Finally understood CORS after reading your explanation. 🙌",
	"Express.js migration guide is pure gold for legacy projects.",
	"Loved the market analysis on Go frameworks—very thoughtful!",
	"Secure auth practices with JWT? Absolutely essential nowadays.",
	"Your clean code habits are inspiring—makes me rethink my structure.",
	"Go standard library gems post had me bookmarking every other line!",
	"Node.js testing techniques are spot-on—applied them immediately.",
	"Handling dirty databases felt chaotic before—this made it manageable.",
}

func Seed(store store.Store) {
	ctx := context.Background()
	users := generateUsers(50)
	for _, user := range users {
		if err := store.User.Create(ctx, user); err != nil {
			log.Fatal("Error while inserting user")
		}
	}
	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := store.Post.Create(ctx, post); err != nil {
			log.Fatal("Error while inserting post")
		}
	}
	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comment.Create(ctx, comment); err != nil {
			log.Fatal("Error while inserting comments", err)
		}
	}

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@gmail.com",
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Password: "samual1234",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags:    []string{tags[rand.Intn(len(tags))], tags[rand.Intn(len(tags))]},
			UserId:  user.Id,
		}
	}
	return posts
}
func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			PostId:  post.Id,
			Content: commentsTmp[rand.Intn(len(commentsTmp))],
			UserId:  user.Id,
		}
	}
	return comments
}
