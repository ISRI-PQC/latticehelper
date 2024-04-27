package pqdevkit

/*
	Returns x mod q, but centered around 0

Args:

	x (int): number to be modded
	q (int): modulus

Returns:

	int: x mod q, centered around 0
*/
func centeredModulo(x, q int64) int64 {
	ret := x % q
	if ret > (q >> 1) {
		ret -= q
	}
	return ret
}