// Package est implements a custom file format that supports:
//   - uint64, which is encoded as little-endian 8 bytes
//   - strings and bytes, which is encoded as 8 length in bytes, followed by the data
//   - messages, which is encoded as concatenated field values
//   - sub-messages, which is encoded into a bytes and then written as a byte slice
//
// For example:
//   - a message containing field `uint64(42)`
//     is encoded as `[42 00 00 00 00 00 00 00]`
//   - a message containing field `"a"`
//     is encoded as `[01 00 00 00 00 00 00 00 97]`
//   - a message containing a sub-message containing fields `uint64(42)` and `"a"`
//     is encoded as `[09 00 00 00 00 00 00 00 42 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 97]`
package est
