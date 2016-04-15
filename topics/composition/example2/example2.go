// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program demonstrating decoupling with interfaces.
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// EOD represents the end of the data stream.
var EOD = errors.New("EOD")

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

// =============================================================================

// Puller declares behavior for pulling data.
type Puller interface {
	Pull(d *Data) error
}

// Storer declares behavior for storing data.
type Storer interface {
	Store(d Data)
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct{}

// Pull knows how to pull data out of Xenia.
func (Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return EOD
	case 5:
		return errors.New("Error reading data from Xenia")
	}

	d.Line = "Data"
	fmt.Println("In:", d.Line)

	return nil
}

// Pillar is a system we need to store data into.
type Pillar struct{}

// Store knows how to store data into Pillar.
func (Pillar) Store(d Data) {
	fmt.Println("Out:", d.Line)
}

// =============================================================================

// System wraps Xenia and Pillar together into a single system.
type System struct {
	Xenia
	Pillar
}

// =============================================================================

// IO provides support to copy bulk data.
type IO struct{}

// pull knows how to pull bulks of data from any Puller.
func (IO) pull(p Puller, data []Data) error {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return err
		}
	}

	return nil
}

// store knows how to store bulks of data from any Storer.
func (IO) store(s Storer, data []Data) {
	for _, d := range data {
		s.Store(d)
	}
}

// Copy knows how to pull and store data from the System.
func (io IO) Copy(sys *System, batch int) error {
	for {
		data := make([]Data, batch)
		if err := io.pull(&sys.Xenia, data); err != nil {
			return err
		}

		io.store(&sys.Pillar, data)
	}
}

// =============================================================================

func main() {

	// Initialize the system for use.
	sys := System{
		Xenia:  Xenia{},
		Pillar: Pillar{},
	}

	var io IO
	if err := io.Copy(&sys, 3); err != nil && err != EOD {
		fmt.Println(err)
	}
}
