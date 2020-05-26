package main

/*****************************************************************************
 *  EMSMiner mines ectocopial Mandelbrot seeds used to create Anthropobrots. *
 *  Copyright © 2020 Daïm Aggott-Hönsch                                      *
 *                                                                           *
 *  This program is free software: you can redistribute it and/or modify     *
 *  it under the terms of the GNU General Public License as published by     *
 *  the Free Software Foundation, either version 3 of the License, or        *
 *  (at your option) any later version.                                      *
 *                                                                           *
 *  This program is distributed in the hope that it will be useful,          *
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of           *
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the            *
 *  GNU General Public License for more details.                             *
 *                                                                           *
 *  You should have received a copy of the GNU General Public License        *
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.   *
 *****************************************************************************/

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func main() {

	min := flag.Int("min", 100, "minimum depth of seeds to mine")
	max := flag.Int("max", 1000, "maximum depth of seeds to mine")
	howmany := flag.Int("howmany", 1000000, "number of seeds to mine")
	flag.Parse()

	fmt.Println("\nEMSMiner v0.1 Copyright (C) 2020 Daïm Aggott-Hönsch. This program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("This is free software, and you are welcome to redistribute it under the conditions specified by")
	fmt.Println("the GNU General Public License 3 (https://www.gnu.org/licenses/gpl-3.0).")

	fmt.Println("\nUsage: " + filepath.Base(os.Args[0]) + " -min [minimum_depth] -max [maximum_depth] -howmany [number_of_seeds_wanted]")

	fmt.Println("")
	rand.Seed(time.Now().UTC().UnixNano())
	seeds, realmin, realmax := Mine(*howmany, *min, *max)
	SaveEMSFile(seeds, realmin, realmax)
}

// .EMS file handling

func SaveEMSFile(seeds seedpack, min, max int) {
	buf := new(bytes.Buffer)
	buf.Reset()

	seeds = seeds.Sort()

	for _, c := range seeds {
		binary.Write(buf, binary.LittleEndian, c)
	}

	md5 := md5.Sum(buf.Bytes())

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	outfilename := filepath.Join(dir, strconv.Itoa(min)+"-"+strconv.Itoa(max)+"_"+fmt.Sprintf("%x", string(md5[:]))+".ems")
	outfile, err := os.OpenFile(outfilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	outfile = outfile

	buf.Reset()
	binary.Write(buf, binary.LittleEndian, []byte("@DM.EMS{codex.apeirography.art} "))
	for _, c := range seeds {
		binary.Write(buf, binary.LittleEndian, c)
	}
	outfile.Write(buf.Bytes())

	return
}

// Optimized Mining Function

func Mine(howmany, min, max int) (seedpack, int, int) {

	/**** Initialization ****/

	if howmany < 1 {
		panic("Number of seeds sought is less than one.")
	}

	if max < min {
		panic("Maximum seed depth is less than minimum seed depth.")
	}

	if min < 2 {
		panic("Minimum seed depth is less than 2.")
	}

	seeds := NewSeedpack(howmany)
	sidx := 0
	guidemap := GenerateGuidemap(51)
	found := 0
	relfound := 0

	realmin, realmax := max, min

	startTime := time.Now()
	relstartTime := time.Now()
	updateInterval := 1
	fmt.Println("Commencing mining of "+strconv.Itoa(howmany)+" seeds with depths between "+strconv.Itoa(min)+" - "+strconv.Itoa(max)+":")

	b := 2.00 * 2.00

	var z, c, oldz complex128
	var l, i, j int
	var repcheck, repcheckstart int

	/**** Outer Loop Begins ****/
	j = 0
CheckNewC:

	z = complex(0, 0)
	c = complex(rand.Float64()*4-2, rand.Float64()*2)
	l = max + 2
	i = 0
	repcheckstart = 2
	repcheck = repcheckstart
	oldz = z

	/**** Inner Loop Begins ****/
	i = 0
IterateZ:
	z = z*z + c
	if repcheck == 0 {
		if oldz == z {
			i = -1
			goto IterateZDone
		}
		oldz = z
		if i%8 == 0 {
			repcheckstart = repcheckstart + 2
			if !guidemap.Check(c) && i%64 != 0 {
				i = -1
				goto IterateZDone
			}
		} else {
			repcheckstart = repcheckstart + 1
		}
		repcheck = repcheckstart
	}
	repcheck--

	i++
	if i < l && (real(z)*real(z))+(imag(z)*imag(z)) <= b {
		goto IterateZ
	}
	/**** Inner Loop Ceases ****/

IterateZDone:
	if i >= min && i <= max {
		if i < realmin {
			realmin = i
		}
		if i > realmax {
			realmax = i
		}
		found++
		relfound++
		seeds[sidx] = c
		sidx++
		guidemap.Mark(c)
		if relfound % updateInterval == 0 {
			if time.Since(relstartTime).Seconds() < 45 {
				if updateInterval > 5 && time.Since(relstartTime).Seconds() > 0 {
					updateInterval = updateInterval * int(float64(90)/float64(time.Since(relstartTime).Seconds()))
				}
				updateInterval++
				relfound = 0
				relstartTime = time.Now()
			} else {
				totalseconds := int(math.Floor(time.Since(startTime).Seconds()))
				sps := float64(found)/float64(totalseconds)
				totalseconds = int((float64(howmany) - float64(found)) / float64(sps))
				hours := totalseconds / 3600
				minutes := (totalseconds - (hours * 3600)) / 60
				seconds := totalseconds - (hours * 3600) - (minutes * 60)

				fmt.Println(strconv.Itoa(found) + " seeds with depths between "+strconv.Itoa(min) + " - " + strconv.Itoa(max)+" have been found so far. "+ strconv.Itoa(hours) +"h "+strconv.Itoa(minutes)+"m " +strconv.Itoa(seconds) +"s"+" left at current speed of "+strconv.Itoa(int(sps*60*60))+" sph.")
				relfound = 0
				relstartTime = time.Now()
			}
		}
	}

	j++
	if found < howmany {
		goto CheckNewC
	}
	/**** Outer Loop Ceases ****/

	totalseconds := int(math.Floor(time.Since(startTime).Seconds()))
	hours := totalseconds / 3600
	minutes := (totalseconds - (hours * 3600)) / 60
	seconds := totalseconds - (hours * 3600) - (minutes * 60)
	sps := int(math.Round(float64(found)/float64(totalseconds)))

	fmt.Println(strconv.Itoa(found) + " seeds with depths between "+strconv.Itoa(min) + " - " + strconv.Itoa(max)+" have been found after "+ strconv.Itoa(hours) +"h "+strconv.Itoa(minutes)+"m " +strconv.Itoa(seconds) +"s"+" with an overall speed of "+strconv.Itoa(int(sps*60*60))+" sph.")

	return seeds, realmin, realmax
}

// Seedpack

type seedpack []complex128

func NewSeedpack(howmany int) seedpack {
	return seedpack(make([]complex128, howmany))
}

func (this seedpack) Sort() seedpack {
	sort.SliceStable(this, func(i, j int) bool {
		if real(this[i]) != real(this[j]) {
			return real(this[i]) < real(this[j])
		} else {
			return imag(this[i]) < imag(this[j])
		}
	})
	return this
}

// Guidemap

type Guidemap struct {
	itsWidth, itsHeight int
	itsMinR, itsMaxR float64
	itsMinI, itsMaxI float64
	itsDelR, itsDelI float64
	itsData []bool
}

func GenerateGuidemap(size int) *Guidemap {

	fmt.Print("Generating guidemap... ")

	this := new(Guidemap)

	this.itsWidth = size
	this.itsHeight = size

	this.itsMinR, this.itsMaxR = -2.00, 2.00
	this.itsMinI, this.itsMaxI = -2.00, 2.00

	this.itsDelR = (this.itsMaxR - this.itsMinR) / float64(this.itsWidth)
	this.itsDelI = (this.itsMaxI - this.itsMinI) / float64(this.itsHeight)

	this.itsData = make([]bool, this.itsWidth * this.itsHeight)

	for idx := 0; idx < len(this.itsData); idx++ {
		this.itsData[idx] = false
	}

	startTime := time.Now()
	found := 0
	limmin := 32
	limmax := limmin * 2
	for time.Since(startTime).Seconds() < 60 {

		z := complex(0.00, 0.00)
		c := complex(rand.Float64()*4-2, rand.Float64()*2)

		for idx := 0; idx < limmax+2; idx++ {
			z = z*z + c
			if cmplx.Abs(z) > 2 {
				if idx >= limmin {
					found++
					if found % (1000) == 0 {
						limmax *= 2
						limmin *= 2
					}
					this.Mark(c)
				}
				break
			}
		}
	}

	fmt.Println("done.")

	/*
	for idx := 0; idx < len(this.itsData); idx++ {
		if idx % this.itsWidth == 0 {
			fmt.Print("\n")
		}
		if this.itsData[idx] {
			fmt.Print("O")
		} else {
			fmt.Print("-")
		}
	}
	fmt.Print("\n")
	*/

	return this
}

func (this *Guidemap) Mark(c complex128) {
	x := int(math.Round((real(c) - this.itsMinR) / this.itsDelR))
	y := int(math.Round((imag(c) - this.itsMinI) / this.itsDelI))
	if x < 0 {
		x = 0
	}
	if x > this.itsWidth-1 {
		x = this.itsWidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y > this.itsHeight-1 {
		y = this.itsHeight - 1
	}
	this.itsData[y*this.itsWidth+x] = true
}

func (this *Guidemap) Check(c complex128) bool {
	x := int(math.Round((real(c) - this.itsMinR) / this.itsDelR))
	y := int(math.Round((imag(c) - this.itsMinI) / this.itsDelI))
	if x < 0 {
		x = 0
	}
	if x > this.itsWidth-1 {
		x = this.itsWidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y > this.itsHeight-1 {
		y = this.itsHeight - 1
	}
	return this.itsData[y*this.itsWidth+x]
}
