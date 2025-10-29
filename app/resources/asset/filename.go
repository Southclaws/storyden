package asset

import (
	"fmt"
	"path"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/rs/xid"
)

const AssetsSubdirectory = "assets"

var errInvalidFormat = fault.New("invalid format")

type Filename struct {
	id    opt.Optional[xid.ID]
	name  string
	hasID bool
}

func (f Filename) IsKnown() bool { return f.hasID }

func (a Filename) String() string {
	return a.name
}

func (f Filename) GetID() xid.ID {
	return f.id.Or(xid.New())
}

func NewFilename(name string) Filename {
	return NewExistingFilename(xid.New(), name)
}

func NewExistingFilename(id xid.ID, name string) Filename {
	return Filename{
		id:    opt.New(id),
		name:  mark.Slugify(formatFilename(id, name)),
		hasID: true,
	}
}

func NewFilepathFilename(name string) Filename {
	return Filename{name: mark.Slugify(name)}
}

func ParseAssetFilename(s string) (*Filename, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return nil, errInvalidFormat
	}

	id, err := xid.FromString(parts[0])
	if err != nil {
		return nil, fault.Wrap(err)
	}

	if !mark.IsSlug(parts[1]) {
		return nil, fault.Wrap(errInvalidFormat, fmsg.With("name is not a valid slug"))
	}

	return &Filename{
		id:   opt.New(id),
		name: fmt.Sprintf("%s-%s", parts[0], parts[1]),
	}, nil
}

func formatFilename(id xid.ID, name string) string {
	if name == "" {
		return id.String()
	}
	return fmt.Sprintf("%s-%s", id.String(), name)
}

func BuildAssetPath(name Filename) string {
	return path.Join(AssetsSubdirectory, name.String())
}
