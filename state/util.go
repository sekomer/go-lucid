package state

import "github.com/holiman/uint256"

func NBitsToTarget(nBits uint32) *uint256.Int {
	/*
		In order to reduce block size, we store nBits(32 bits) instead of target(256 bits) in the block header.
			nbits:
			  [1 byte exponent][3 bytes coefficient]

			e.g. 0x1e000fff
				exponent: 0x1e
				coefficient: 0x000fff

			target:
				coefficient * 256**(exponent-3)
	*/

	exponent := nBits >> 24
	coefficient := nBits & 0xffffff

	target := uint256.NewInt(uint64(coefficient))

	// Optimize the calculation by shifting left or right
	if exponent > 3 {
		target.Lsh(target, uint(8*(exponent-3)))
	} else if exponent < 3 {
		target.Rsh(target, uint(8*(3-exponent)))
	}

	return target
}

// TODO(optional): Implement this
// func TargetToNBits(target *uint256.Int) uint32 {
// 	return 0
// }
