package generative

import (
	"context"
	"html/template"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

var TitlePrompt = template.Must(template.New("").Parse(`Generate 1 to 3 concise title suggestions for the following content. These titles will be used as page titles AND in URL slugs, so they must be URL-friendly and focus on the primary SUBJECT of the content, not descriptive topics or commentary.

Guidelines:
- If a URL is provided, extract the brand/product/subject name from the domain or path (e.g., "thedesignsystem.guide" → "Design System Guide")
- If an Original Title is provided, use it as a strong signal but simplify and clean it up (remove marketing language, "the ultimate", etc.)
- Focus on the main subject, entity, or thing being discussed (e.g., "React Hooks" not "Understanding React Hooks")
- Keep titles simple, direct, and factual
- Avoid creative or promotional language
- Limit to 60 characters maximum
- Prefer nouns and noun phrases over verb phrases
- Avoid question formats unless the content is explicitly a question
- Remove unnecessary articles (a, an, the) when they don't add clarity

<example>
URL: https://thedesignsystem.guide/
Original Title: The Design System Guide - Learn Design Systems

Content:
A comprehensive resource for learning design systems, covering everything from design tokens to component libraries...

Good titles:
- Design System Guide
- Design Systems Resource
- Design System Learning

Bad titles:
- The Ultimate Design System Guide
- Learn Everything About Design Systems
- Complete Design System Tutorial
</example>

<example>
URL: https://selfh.st/
Original Title: Selfh.st - Modern Self-Hosting Made Easy

Content:
Selfh.st is a modern self-hosting platform that makes it easy to deploy applications...

Good titles:
- Selfh.st
- Selfh.st Platform
- Self-Hosting Platform

Bad titles:
- Modern Self-Hosting Made Easy
- The Future of Self-Hosting
- Revolutionary Self-Hosting Platform
</example>

<example>
Content: "I've been working with Docker for a while now and wanted to share some best practices I've learned about container optimization and layer caching..."

Good titles:
- Docker Container Optimization
- Container Layer Caching
- Docker Best Practices

Bad titles:
- How I Learned to Optimize Docker Containers
- My Journey with Docker: Tips and Tricks
- Everything You Need to Know About Docker
</example>

<example>
Content: "Here's my implementation of a binary search tree in Go with some interesting performance characteristics..."
Good titles:
- Binary Search Tree in Go
- Go BST Implementation
- Binary Search Tree Performance

Bad titles:
- How to Build a Binary Search Tree
- My Take on Binary Search Trees
- The Ultimate Guide to BST in Go
</example>

<example>
Content: "What's the best way to handle authentication in a Next.js application? I'm trying to decide between NextAuth and custom JWT..."
Good titles:
- Next.js Authentication Options
- NextAuth vs Custom JWT
- Next.js Auth Implementation

Bad titles:
- Choosing the Right Authentication for Your Next.js App
- A Deep Dive into Next.js Authentication
- Which Auth Solution Should You Use?
</example>

Content:

{{ .Content }}
`))

type SuggestTitleResultSchema struct {
	Titles []string `json:"titles" jsonschema:"title=Titles,description=List of suggested titles,items=string"`
}

func (g *generator) SuggestTitle(ctx context.Context, content datagraph.Content) ([]string, error) {
	template := strings.Builder{}
	err := TitlePrompt.Execute(&template, map[string]any{
		"Content": content.Plaintext(),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := ai.PromptObject(ctx, g.prompter, "Suggest titles for content", template.String(), SuggestTitleResultSchema{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return result.Titles, nil
}
