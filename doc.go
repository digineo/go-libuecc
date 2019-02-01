// Package libuecc is a very small generic-purpose Elliptic Curve
// Cryptography library compatible with Ed25519.
//
// There are two sets of functions converting between libuecc's internal point
// representation and coordinates or compressed representation. The functions
// ending with Ed25519 use the same representation as original Ed25519
// implementation and should be used by new software. The functions with the
// suffix Legacy are provided for compatibility with libuecc version before
// v6.
//
// Ed25519 and the legacy representation are isomorphic, they use a Twisted
// Edwards Curve
//
//     ax^2 + y^2 = 1 + dx^2y^2
//
// over the prime field for p = 2^255 - 19.
//
// Ed25519 uses the parameters
//
//     a = -1 and
//     d = -(121665/121666),
//
// while the legacy curve has
//
//     a = 486664
//     d = 486660.
package libuecc
