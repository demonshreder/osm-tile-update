package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jinzhu/gorm"
)

var workDir, _ = os.Getwd()

var avail bool

var links []string

type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

type Tile struct {
	Z    int
	X    int
	Y    int
	Lat  float64
	Long float64
}

type Url struct {
	gorm.Model
	oUrl string
	x    int
	y    int
	z    int
}

func Deg2num(Lat, Long float64, Z int) (x int, y int) {
	x = int(math.Floor((Long + 180.0) / 360.0 * (math.Exp2(float64(Z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(Lat*math.Pi/180.0)+1.0/math.Cos(Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(Z)))))
	return x, y
}

func (*Tile) Num2deg(t *Tile) (lat float64, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, long
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.Create(filepath.Join(workDir, "links"))
	check(err)
	defer f.Close()

	site := "a"
	curBB := BoundingBox{MinLat: 11.8680, MinLon: 79.7374, MaxLat: 11.9755, MaxLon: 79.8493}
	// fmt.Println(a)
	count := 0

	curLat := curBB.MinLat

	for curLon := curBB.MinLon; curLon <= curBB.MaxLon; curLon = curLon + 0.0001 {
		for Z := 0; Z < 15; Z++ {
			avail = true
			x, y := Deg2num(curLat, curLon, Z)

			fullUrl := fmt.Sprint(".tile.openstreetmap.org/" + strconv.Itoa(Z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png")
			for link := range links {
				if links[link] == fullUrl {
					avail = false
				}
			}
			if avail {
				count++
				links = append(links, fullUrl)
				fmt.Println(count, fullUrl)
			}

		}
	}

	for ; curLat <= curBB.MaxLat; curLat = curLat + 0.0001 {
		for curLon := curBB.MinLon; curLon <= curBB.MaxLon; curLon = curLon + 0.0001 {
			for Z := 15; Z < 20; Z++ {
				avail = true
				x, y := Deg2num(curLat, curLon, Z)

				fullUrl := fmt.Sprint(".tile.openstreetmap.org/" + strconv.Itoa(Z) + "/" + strconv.Itoa(x) + "/" + strconv.Itoa(y) + ".png")
				for link := range links {
					if links[link] == fullUrl {
						avail = false
					}
				}
				if avail {
					// links[count] = fullUrl
					count++
					links = append(links, fullUrl)
					fmt.Println(count, fullUrl)
				}
			}
		}
	}

	fmt.Println("Writing the file links")
	for i := range links {
		links[i] = "http://" + site + links[i] + "\n"
		_, err := f.WriteString(links[i])
		check(err)

		switch site {
		case "a":
			site = "b"
		case "b":
			site = "c"
		case "c":
			site = "a"
		}
	}
	f.Sync()
}
