# go-libuecc

[![CircleCI](https://circleci.com/gh/digineo/go-libuecc/tree/master.svg?style=shield)](https://circleci.com/gh/digineo/go-libuecc/tree/master)
[![Codecov](http://codecov.io/github/digineo/go-libuecc/coverage.svg?branch=master)](http://codecov.io/github/digineo/go-libuecc?branch=master)


This is a port of the [libuecc v7](https://git.universe-factory.net/libuecc) C library to Go. It is used to
reduce the C call overhead in our [fastd](https://github.com/digineo/fastd) implementation.

**WARNING:** While extra care was taken while porting the code, this
was not crucifyingly reviewed by the original author, nor by other
security experts. Expect some nasty bugs!

## Notes

- Where possible, the ported code adapts an idiomatic Go style:

  ```c
  int ecc_25519_load_xy_ed25519(ecc_25519_work_t *out, const ecc_int256_t *x, const ecc_int256_t *y) { /* ... */ }
  ```

  vs.

  ```go
  func loadXYEd25519(x, y *int256) (out *point) { /* ... */ }
  ```

- Data structures are handled immutable:

  ```c
  /* squeeze modifies a */
  static void squeeze(uint32_t a[32]) { /* ... */ }
  ```

  vs.

  ```go
  // squeeze does not modify a, and returns new unpacked object
  func (a unpacked) squeeze() unpacked { /* ... */ }
  ```

  This will result in a slightly higher resource consumption and
  (presumably) in slower code execution.

- When feasable, tests rely on generated data. We'll use `testdata/gen.c`
  to generate *expected behaviour* and run the Go tests against these
  precomputed results.

  To achieve this, the upstream libuecc is bundled as submodule in
  `testdata/libuecc`. Just run `make` to run the tests.

  Note: This also means that you'll need GCC to run the tests on your
  machine (`go test` will fail otherwise, because it can't find
  `testdata/cases/*`).
