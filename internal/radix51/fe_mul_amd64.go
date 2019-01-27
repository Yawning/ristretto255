// Copyright (c) 2017 George Tankersley. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!noasm

package radix51

// Mul sets out = a * b.
func (v *FieldElement) Mul(x, y *FieldElement) {
	feMul(v, x, y)
}

// go:noescape
func feMul(out, a, b *FieldElement)