package avatar_gen

import (
	"context"
	"fmt"
	"image/png"
	"math"
	"os"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/stretchr/testify/require"
)

func Test_service_Generate(t *testing.T) {
	r := require.New(t)
	s := service{}

	img, err := s.Generate(context.Background(), "southclaws")
	r.NoError(err)
	r.NotNil(img)

	file, err := os.Create("test/gradient.png")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	err = png.Encode(file, img)
	r.NoError(err)
}

func Test_hashfunction(t *testing.T) {
	input := []string{
		"southclaws",
		"JustMichael",
		"iAmir",
		"Hual",
		"J0sh_ES",
		"maddinat0r",
		"Dobby",
		"Y_Less",
		"Kyle_Smith",
		"Cheaterman",
		"aymel",
		"hiddos",
		"54m",
		"addetz",
		"blj",
		"bskyus",
		"catt",
		"discovery",
		"drstanton",
		"dvf",
		"flicknow",
		"gpte",
		"ivanbaldoino",
		"jasonchan",
		"jpren",
		"mediciners",
		"megumin",
		"miranda",
		"msrodrigo",
		"nexus",
		"philipbrown",
		"rathbone",
		"rui",
		"shad",
		"sho",
		"trippy",
		"yaf",
	}

	outputs := []uint16{}

	for _, v := range input {
		outputs = append(outputs, hashfunction(v))
		fmt.Println(v, hashfunction(v))
	}

	sum := dt.Reduce(outputs, func(c uint16, n uint16) uint16 {
		return c + n
	}, 0)
	avg := float64(sum) / float64(len(input))

	dev := dt.Reduce(outputs, func(c float64, n uint16) float64 {
		distance := float64(n) - avg
		sq := distance * distance
		return c + sq
	}, 0)

	stdev := math.Sqrt(dev / float64(len(input)))

	// we want the stdev of the hash function to be fairly high so we get
	// colours that are different enough for different usernames.
	fmt.Println(avg, stdev)
}
