package spam_checker

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetector_Detect(t *testing.T) {
	threshold := 0.0003
	d := repeatedContentDetector{
		threshold: threshold,
	}

	check := func(f string, below float64) {
		t.Run(f, func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			reader, err := os.Open(f)
			r.NoError(err)

			got, err := d.getRatio(reader)
			r.NoError(err)

			a.Less(got, below)
		})
	}

	check("./post01_spam.txt" /**/, threshold)
	check("./post02.txt" /*     */, 0.03)
	check("./post03.txt" /*     */, 0.0008)
	check("./post04.txt" /*     */, 0.0014)
	check("./post05.txt" /*     */, 0.0014)
	check("./post06.txt" /*     */, 0.009)
	check("./post07.txt" /*     */, 0.005)
	check("./post08.txt" /*     */, 0.003)
	check("./post09_spam.txt" /**/, 0.003) // false negative
	check("./post10_spam.txt" /**/, 0.003) // false negative
}
