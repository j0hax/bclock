// Utility functions to work with human-readable binary time representations
package bcd

import (
	"math"
	"time"
)

// Represents an integer [0,100) with two bit fields
type DoubleDigit struct {
	Tens, Ones [4]bool
}

// Represents individual digits of hours, minutes and seconds
type BCDTime struct {
	// Represents the hours as a binary value
	Hours DoubleDigit

	// Represents the minutes as a binary value
	Minutes DoubleDigit

	// Represents the seconds as a binary value
	Seconds DoubleDigit
}

// Sets the hour bits
func (ts *BCDTime) SetHour(h int) {
	ts.Hours.Tens = ToNibble(DigitAt(h, 1))
	ts.Hours.Ones = ToNibble(DigitAt(h, 0))
}

// Sets the minute bits
func (ts *BCDTime) SetMinute(m int) {
	ts.Minutes.Tens = ToNibble(DigitAt(m, 1))
	ts.Minutes.Ones = ToNibble(DigitAt(m, 0))
}

// Sets the seconds bits
func (ts *BCDTime) SetSeconds(s int) {
	ts.Seconds.Tens = ToNibble(DigitAt(s, 1))
	ts.Seconds.Ones = ToNibble(DigitAt(s, 0))
}

// Convert a given time to BCD
func ToBCD(t time.Time) *BCDTime {
	var out BCDTime

	out.SetHour(t.Hour())
	out.SetMinute(t.Minute())
	out.SetSeconds(t.Second())

	return &out
}

// Return the current time in BCD format
func Now() *BCDTime {
	return ToBCD(time.Now())
}

// Convert an integer [0,16) to an array of booleans.
// This function is limited to 4 bits as more are not needed to represent the digits [0,9]
func ToNibble(n int) [4]bool {
	var result [4]bool
	for i := 0; i < 4; i++ {
		result[i] = (n >> i & 1) != 0
	}
	return result
}

// Return the digit at a certain location [0, n)
func DigitAt(num int, n int) int {
	r := num / int(math.Pow10(n))
	digit := r % 10
	//fmt.Printf("%dth digit of %d is %d\n", n, num, digit)
	return digit
}
