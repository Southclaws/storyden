package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/local"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/glebarez/go-sqlite"
)

// Sample content for realistic distribution
var shortReplies = []string{
	"thanks!", "lol", "nice", "cool", "awesome", "great post", "+1", "agreed",
	"thanks for sharing", "interesting", "good point", "makes sense", "I see",
	"definitely", "exactly", "yep", "sure", "ok", "got it", "understood",
}

var mediumSentences = []string{
	"This is a really interesting perspective on the topic.",
	"I completely agree with your point about this.",
	"Thanks for taking the time to write this up.",
	"This reminds me of something similar I experienced.",
	"I think there's definitely merit to this approach.",
	"This is exactly what I was looking for, thanks!",
	"Great explanation, this really helped me understand.",
	"I had never thought about it this way before.",
}

var longParagraphs = []string{
	`I think this is a really important topic that deserves more discussion. 
    In my experience, I've found that taking the time to really understand the 
    underlying principles makes a huge difference in implementation success. 
    Too often, people rush into solutions without fully grasping the problem space.`,

	`This approach has worked well for me in several projects. The key is to 
    start small and iterate based on feedback. I've seen teams try to implement 
    everything at once and end up with a mess that's hard to maintain. Better 
    to get something working first, then enhance it gradually.`,

	`I disagree with this assessment. While it's true that simplicity is 
    important, sometimes complex problems require sophisticated solutions. 
    The trick is finding the right balance between simplicity and functionality. 
    I've seen oversimplified solutions that create more problems than they solve.`,
}

var threadTopics = []string{
	"Best practices for", "How to handle", "Thoughts on", "Experience with",
	"Question about", "Discussion: ", "Advice needed for", "Tips for",
	"Strategies for", "Approaches to", "Methods for", "Techniques for",
}

var threadSubjects = []string{
	"database optimization", "API design", "user authentication", "error handling",
	"performance tuning", "code organization", "testing strategies", "deployment",
	"monitoring", "logging", "caching", "security", "scalability", "architecture",
	"frontend frameworks", "backend services", "microservices", "DevOps",
	"team collaboration", "project management", "documentation", "code reviews",
}

func generateContent(lengthType string) string {
	switch lengthType {
	case "short":
		return shortReplies[rand.Intn(len(shortReplies))]
	case "medium":
		return mediumSentences[rand.Intn(len(mediumSentences))]
	case "long":
		if rand.Float32() < 0.7 {
			return longParagraphs[rand.Intn(len(longParagraphs))]
		}
		// Multiple paragraphs
		numParagraphs := rand.Intn(2) + 2 // 2-3 paragraphs
		paragraphs := make([]string, numParagraphs)
		for i := 0; i < numParagraphs; i++ {
			paragraphs[i] = longParagraphs[rand.Intn(len(longParagraphs))]
		}
		result := ""
		for i, p := range paragraphs {
			if i > 0 {
				result += "\n\n"
			}
			result += p
		}
		return result
	default: // "thread" - opening posts are longer
		numParagraphs := rand.Intn(4) + 2 // 2-5 paragraphs
		paragraphs := make([]string, numParagraphs)
		for i := 0; i < numParagraphs; i++ {
			paragraphs[i] = longParagraphs[rand.Intn(len(longParagraphs))]
		}
		result := ""
		for i, p := range paragraphs {
			if i > 0 {
				result += "\n\n"
			}
			result += p
		}
		return result
	}
}

func generateTitle() string {
	topic := threadTopics[rand.Intn(len(threadTopics))]
	subject := threadSubjects[rand.Intn(len(threadSubjects))]
	return fmt.Sprintf("%s %s", topic, subject)
}

func generateHTMLContent(textContent string) string {
	paragraphs := splitParagraphs(textContent)
	htmlParagraphs := make([]string, 0, len(paragraphs))
	for _, p := range paragraphs {
		if p = cleanString(p); p != "" {
			htmlParagraphs = append(htmlParagraphs, fmt.Sprintf("<p>%s</p>", p))
		}
	}
	return fmt.Sprintf("<body>%s</body>", joinStrings(htmlParagraphs))
}

func splitParagraphs(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '\n' && s[i+1] == '\n' {
			result = append(result, s[start:i])
			start = i + 2
		}
	}
	result = append(result, s[start:])
	return result
}

func cleanString(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if s != "" {
			if i > 0 {
				result += "\n\n"
			}
			result += s
		}
	}
	return result
}

type seeder struct {
	db            *ent.Client
	accountWriter *account_writer.Writer
}

func newSeeder(db *ent.Client) *seeder {
	store, err := local.New()
	if err != nil {
		panic(err)
	}

	roleHydrator := role_repo.New(db, store)
	accountQuerier := account_querier.New(db, roleHydrator)
	return &seeder{
		db:            db,
		accountWriter: account_writer.New(db, accountQuerier, roleHydrator),
	}
}

func (s *seeder) createRandomAccount(ctx context.Context) (*ent.Account, error) {
	// Generate random handle using petname
	handle := petname.Generate(2, "-") + "-" + strconv.Itoa(rand.Intn(1000))

	bioContent, err := datagraph.NewRichText("<body><p>Random user generated for testing</p></body>")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountWriter.Create(ctx, handle,
		account_writer.WithName(strings.Title(strings.ReplaceAll(handle, "-", " "))),
		account_writer.WithBio(bioContent))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.db.Account.Get(ctx, xid.ID(acc.ID))
}

func (s *seeder) getOrCreateCategories(ctx context.Context) ([]*ent.Category, error) {
	categories, err := s.db.Category.Query().Limit(20).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	if len(categories) == 0 {
		// Create test categories if none exist
		for i := 0; i < 10; i++ {
			cat, err := s.db.Category.Create().
				SetName(fmt.Sprintf("Category %d", i+1)).
				SetSlug(fmt.Sprintf("category-%d", i+1)).
				SetDescription(fmt.Sprintf("Test category %d", i+1)).
				SetColour("rgba(59, 130, 246, 1)").
				SetSort(i).
				SetAdmin(false).
				SetMetadata(map[string]any{}).
				Save(ctx)
			if err != nil {
				return nil, fault.Wrap(err)
			}
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (s *seeder) seedData(ctx context.Context, numThreads int, numAccounts int) error {
	fmt.Printf("Seeding %d threads with replies using %d accounts...\n", numThreads, numAccounts)

	// Create multiple test accounts
	accounts := make([]*ent.Account, 0, numAccounts)
	for i := 0; i < numAccounts; i++ {
		acc, err := s.createRandomAccount(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		accounts = append(accounts, acc)
		if i%10 == 0 {
			fmt.Printf("  Created %d accounts\r", i+1)
		}
	}
	fmt.Printf("  Created %d accounts\n", numAccounts)

	// Get or create categories
	categories, err := s.getOrCreateCategories(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	fmt.Printf("Using %d accounts and %d categories\n", len(accounts), len(categories))

	// Generate threads in batches for better performance
	batchSize := 1000
	totalReplies := 0

	startTime := time.Now()

	for batchStart := 0; batchStart < numThreads; batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > numThreads {
			batchEnd = numThreads
		}

		fmt.Printf("Processing batch %d/%d\n", batchStart/batchSize+1, (numThreads-1)/batchSize+1)

		// Generate threads for this batch
		for i := batchStart; i < batchEnd; i++ {
			title := generateTitle()
			bodyText := generateContent("thread")
			bodyHTML := generateHTMLContent(bodyText)
			shortText := bodyText
			if len(bodyText) > 200 {
				shortText = bodyText[:200] + "..."
			}

			category := categories[rand.Intn(len(categories))]
			createdAt := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)
			slug := fmt.Sprintf("%s-%s", title, xid.New().String()[:8])

			// Create thread
			account := accounts[rand.Intn(len(accounts))]
			thread, err := s.db.Post.Create().
				SetTitle(title).
				SetSlug(slug).
				SetBody(bodyHTML).
				SetShort(shortText).
				SetLastReplyAt(time.Now()).
				SetVisibility("published").
				SetCreatedAt(createdAt).
				SetUpdatedAt(createdAt).
				SetAccountPosts(account.ID).
				SetCategoryID(category.ID).
				SetMetadata(map[string]any{}).
				Save(ctx)
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

			if i%100 == 0 {
				fmt.Printf("  Thread %d/%d\r", i+1, numThreads)
			}

			// Generate replies for this thread
			numReplies := rand.Intn(191) + 10 // 10-200 replies
			totalReplies += numReplies

			for j := 0; j < numReplies; j++ {
				// Content distribution: 40% short, 40% medium, 20% long
				contentRand := rand.Float32()
				var bodyText string
				if contentRand < 0.4 {
					bodyText = generateContent("short")
				} else if contentRand < 0.8 {
					bodyText = generateContent("medium")
				} else {
					bodyText = generateContent("long")
				}

				bodyHTML := generateHTMLContent(bodyText)
				shortText := bodyText
				if len(bodyText) > 200 {
					shortText = bodyText[:200] + "..."
				}

				// Replies are created after the thread, but never in the future
				minutesAfterThread := rand.Intn(10080) // 1 minute to 1 week later
				replyCreated := createdAt.Add(time.Duration(minutesAfterThread) * time.Minute)

				// Ensure reply is not in the future
				if replyCreated.After(time.Now()) {
					replyCreated = time.Now().Add(-time.Duration(rand.Intn(60)) * time.Minute) // Recent reply within last hour
				}

				// Create reply
				account := accounts[rand.Intn(len(accounts))]
				_, err := s.db.Post.Create().
					SetBody(bodyHTML).
					SetShort(shortText).
					SetVisibility("published").
					SetCreatedAt(replyCreated).
					SetUpdatedAt(replyCreated).
					SetAccountPosts(account.ID).
					SetRootPostID(thread.ID).
					SetMetadata(map[string]any{}).
					Save(ctx)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				if j%50 == 0 {
					fmt.Printf("    Thread %d: reply %d/%d\r", i+1, j+1, numReplies)
				}
			}
		}

		// Progress update
		elapsed := time.Since(startTime)
		threadsPerSec := float64(batchEnd) / elapsed.Seconds()
		eta := time.Duration(float64(numThreads-batchEnd)/threadsPerSec) * time.Second

		fmt.Printf("  Batch completed. Threads: %d/%d (%.1f/sec, ETA: %.1fmin)\n",
			batchEnd, numThreads, threadsPerSec, eta.Minutes())
	}

	// Final statistics
	totalTime := time.Since(startTime)
	fmt.Printf("\nSeeding completed!\n")
	fmt.Printf("Total threads: %d\n", numThreads)
	fmt.Printf("Total replies: %d\n", totalReplies)
	fmt.Printf("Total posts: %d\n", numThreads+totalReplies)
	fmt.Printf("Time taken: %.1f minutes\n", totalTime.Minutes())
	fmt.Printf("Average rate: %.1f posts/second\n", float64(numThreads+totalReplies)/totalTime.Seconds())

	// Verify data
	threadCount, err := s.db.Post.Query().Where(post.RootPostIDIsNil()).Count(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	replyCount, err := s.db.Post.Query().Where(post.RootPostIDNotNil()).Count(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	fmt.Printf("\nVerification:\n")
	fmt.Printf("Threads in DB: %d\n", threadCount)
	fmt.Printf("Replies in DB: %d\n", replyCount)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <database_path>")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("Database %s does not exist.\n", dbPath)
		os.Exit(1)
	}

	// Create database connection manually with proper foreign keys
	sqlDB, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqlDB.Close()

	// Create ent client
	opts := []ent.Option{}
	opts = append(opts, ent.Driver(entsql.OpenDB(dialect.SQLite, sqlDB)))

	client := ent.NewClient(opts...)

	// Run schema migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	fmt.Printf("Connected to database: %s\n", dbPath)

	// Create seeder
	seeder := newSeeder(client)

	// Seed data
	if err := seeder.seedData(context.Background(), 10000, 50); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	fmt.Println("\nDatabase seeding complete!")
}
