package weaviate_semdexer

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type Source struct {
	ID      xid.ID
	Kind    datagraph.Kind
	URL     url.URL
	Content string
}

// const fakeSDR = "https://sdr-dummy-domain.com/"

func mapObjectToSource(o WeaviateObject) (*Source, error) {
	id, err := xid.FromString(o.DatagraphID)
	if err != nil {
		return nil, err
	}

	kind, err := datagraph.NewKind(o.DatagraphType)
	if err != nil {
		return nil, err
	}

	fakeSDR, err := url.Parse(fmt.Sprintf("%s:%s/%s", datagraph.RefScheme, kind, id.String()))
	if err != nil {
		return nil, err
	}

	return &Source{
		ID:      id,
		Kind:    kind,
		URL:     *fakeSDR,
		Content: o.Content,
	}, nil
}

func mapObjectsToSources(objects []WeaviateObject) ([]*Source, error) {
	return dt.MapErr(objects, mapObjectToSource)
}

var AnswerPrompt = template.Must(template.New("").Parse(`
You are an expert assistant. Answer the user's question accurately and concisely using the provided sources. Cite the sources in a separate list at the end of your answer. 
Ensure that the source URLs (in "sdr" format) are kept exactly as they appear, without modification or breaking them across lines.

Sources:
{{- range .Context }}
- URL: {{ .URL.String }}
  Kind: {{ .Kind }}
  Content: {{ .Content }}
{{- end }}

Question: {{ .Question }}

Answer:
1. Provide your answer here in clear and concise paragraphs.
2. Use information from the sources above to support your answer, but do not include citations inline.
3. Include a "References" section with the source URLs listed, like this:

References:
- (the url to the source): (Short description of the source content)
`))

const maxContextForRAG = 10

func (s *weaviateSemdexer) Ask(ctx context.Context, q string) (chan string, chan error) {
	objects, err := s.SearchChunks(ctx, q, pagination.NewPageParams(1, 200), searcher.Options{})
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	if len(objects) > maxContextForRAG {
		objects = objects[:maxContextForRAG]
	}

	sources, err := mapObjectsToSources(objects)
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	t := strings.Builder{}
	err = AnswerPrompt.Execute(&t, map[string]any{
		"Context":  sources,
		"Question": q,
	})
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	chch, ech := s.ai.PromptStream(ctx, t.String())

	return chch, ech
}
