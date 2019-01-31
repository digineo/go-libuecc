/* gen.c generates a bunch of test cases for the Go port. */

#include <errno.h>
#include <error.h>
#include <inttypes.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "../../libuecc/src/ec25519.c"

void saveQuad(FILE *f, uint32_t word) {
	for (int pos = 0; pos < 4; ++pos) {
		int c = (unsigned char)(word >> (pos*8));
		if (fputc(c, f) != c) {
			error(1, 0, "failed to write byte 0x%02x", c);
		}
	}
}

void saveUnpacked(const char *filename, const uint32_t a[32]) {
	FILE *f;
	if ((f = fopen(filename, "w")) == NULL) {
		error(1, errno, "cannot open %s", filename);
	}

	for (int i = 0; i < 32; ++i) {
		saveQuad(f, a[i]);
	}
	fclose(f);
}

void saveInt256(const char *filename, const ecc_int256_t *in) {
	FILE *f;
	if ((f = fopen(filename, "w")) == NULL) {
		error(1, errno, "cannot open %s", filename);
	}

	for (int i = 0; i < 32; ++i) {
		int c = (unsigned char)(in->p[i]);
		if (fputc(c, f) != c) {
			error(1, 0, "failed to write byte 0x%02x", c);
		}
	}
	fclose(f);
}

void saveWork(const char *filename, const ecc_25519_work_t *p) {
	FILE *f;
	if ((f = fopen(filename, "w")) == NULL) {
		error(1, errno, "cannot open %s", filename);
	}

	#define saveElem(arr) do { \
		for (int i = 0; i < 32; ++i) { \
			saveQuad(f, arr[i]); \
		} \
	} while(0);

	saveElem(p->X);
	saveElem(p->Y);
	saveElem(p->Z);
	saveElem(p->T);
	fclose(f);

	#undef saveElem
}

void copy(uint32_t dst[32], const uint32_t src[32]) {
	for (int i = 0; i < 32; ++i) {
		dst[i] = src[i];
	}
}

int main(int argc, char const *argv[]) {
	#define save(name, code) do { \
		if (access("cases/"name, F_OK) == -1) { \
			uint32_t t[32] = {0}; \
			code ; \
			fprintf(stderr, " %s", name); \
			saveUnpacked("cases/"name, t); \
		} \
	} while(0);

	puts("generating... ");

	save("one",    copy(t, one));
	save("zero",   copy(t, zero));
	save("minusp", copy(t, minusp));

	save("add_0_0",      add(t, zero,   zero));
	save("add_0_1",      add(t, zero,    one));
	save("add_1_0",      add(t,  one,   zero));
	save("add_1_1",      add(t,  one,    one));
	save("add_0_minusp", add(t, zero, minusp));
	save("add_1_minusp", add(t,  one, minusp));

	save("sub_0_0",      sub(t, zero,   zero));
	save("sub_0_1",      sub(t, zero,    one));
	save("sub_1_0",      sub(t,  one,   zero));
	save("sub_1_1",      sub(t,  one,    one));
	save("sub_0_minusp", sub(t, zero, minusp));
	save("sub_1_minusp", sub(t,  one, minusp));

	save("sub_add_0_0_0", add(t, zero, zero); sub(t, t, zero));
	save("sub_add_0_0_1", add(t, zero, zero); sub(t, t,  one));
	save("sub_add_0_1_0", add(t, zero,  one); sub(t, t, zero));
	save("sub_add_0_1_1", add(t, zero,  one); sub(t, t,  one));
	save("sub_add_1_0_0", add(t,  one, zero); sub(t, t, zero));
	save("sub_add_1_0_1", add(t,  one, zero); sub(t, t,  one));
	save("sub_add_1_1_0", add(t,  one,  one); sub(t, t, zero));
	save("sub_add_1_1_1", add(t,  one,  one); sub(t, t,  one));

	save("add_sub_0_0_0", sub(t, zero, zero); add(t, t, zero));
	save("add_sub_0_0_1", sub(t, zero, zero); add(t, t,  one));
	save("add_sub_0_1_0", sub(t, zero,  one); add(t, t, zero));
	save("add_sub_0_1_1", sub(t, zero,  one); add(t, t,  one));
	save("add_sub_1_0_0", sub(t,  one, zero); add(t, t, zero));
	save("add_sub_1_0_1", sub(t,  one, zero); add(t, t,  one));
	save("add_sub_1_1_0", sub(t,  one,  one); add(t, t, zero));
	save("add_sub_1_1_1", sub(t,  one,  one); add(t, t,  one));

	{
		uint32_t z[32] = {0};
		save("squeeze_zero", squeeze(z); *t = *z);
	}{
		uint32_t z[32] = {1};
		save("squeeze_one", squeeze(z); *t = *z);
	}
	save("squeeze_sub_0_1", sub(t, zero, one); squeeze(t))

	{
		uint32_t z[32] = {0};
		save("freeze_zero", freeze(z); *t = *z);
	}{
		uint32_t z[32] = {1};
		save("freeze_one", freeze(z); *t = *z);
	}
	save("freeze_sub_0_1", sub(t, zero, one); freeze(t))

	{
		uint32_t z[32] = {0};
		squeeze(z);
		printf(" parity_zero=%"PRIu32, parity(z));
	}{
		uint32_t z[32] = {1};
		squeeze(z);
		printf(" parity_one=%"PRIu32, parity(z));
	}{
		uint32_t z[32];
		copy(z, minusp);
		squeeze(z);
		printf(" parity_minusp=%"PRIu32, parity(z));
	}{
		uint32_t z[32] = {0};
		sub(z, zero, one);
		squeeze(z);
		printf(" parity_sub_0_1=%"PRIu32, parity(z));
	}

	save("mult_0_0", mult(t, zero, zero));
	save("mult_0_1", mult(t, zero,  one));
	save("mult_1_0", mult(t,  one, zero));
	save("mult_1_1", mult(t,  one,  one));
	{
		uint32_t mp[32];
		copy(mp, minusp);
		save("mult_minusp_minusp", mult(t, mp, mp));
	}

	save("mult_int_0_0",   mult_int(t, UINT32_C(0), zero));
	save("mult_int_1_0",   mult_int(t, UINT32_C(0),  one));
	save("mult_int_0_1",   mult_int(t, UINT32_C(1), zero));
	save("mult_int_1_1",   mult_int(t, UINT32_C(1),  one));
	save("mult_int_0_max", mult_int(t, UINT32_MAX,  zero));
	save("mult_int_1_max", mult_int(t, UINT32_MAX,   one));

	save("square_0",      square(t,   zero));
	save("square_1",      square(t,    one));
	save("square_minusp", square(t, minusp));

	printf(" equal_0_0=%d",      check_equal(zero,   zero));
	printf(" equal_0_1=%d",      check_equal(zero,    one));
	printf(" equal_1_0=%d",      check_equal( one,   zero));
	printf(" equal_1_1=%d",      check_equal( one,    one));
	printf(" equal_0_minusp=%d", check_equal(zero, minusp));

	save("select_0_1_0",      select(t, zero,    one, 0));
	save("select_0_1_1",      select(t, zero,    one, 1));
	save("select_0_minusp_0", select(t, zero, minusp, 0));
	save("select_0_minusp_1", select(t, zero, minusp, 1));
	save("select_1_minusp_0", select(t,  one, minusp, 0));
	save("select_1_minusp_1", select(t,  one, minusp, 1));

	const ecc_int256_t testKeys[4] = {
		// 83369beddca777585167520fb54a7fb059102bf4e0a46dd5fb1c633d83db77a2
		{{
			0x83, 0x36, 0x9b, 0xed, 0xdc, 0xa7, 0x77, 0x58,
			0x51, 0x67, 0x52, 0x0f, 0xb5, 0x4a, 0x7f, 0xb0,
			0x59, 0x10, 0x2b, 0xf4, 0xe0, 0xa4, 0x6d, 0xd5,
			0xfb, 0x1c, 0x63, 0x3d, 0x83, 0xdb, 0x77, 0xa2
		}},
		// b4dbdb0c05dd28204534fa27c5afca4dcda5397d833e3064f7a7281b249dc7c7
		{{
			0xb4, 0xdb, 0xdb, 0x0c, 0x05, 0xdd, 0x28, 0x20,
			0x45, 0x34, 0xfa, 0x27, 0xc5, 0xaf, 0xca, 0x4d,
			0xcd, 0xa5, 0x39, 0x7d, 0x83, 0x3e, 0x30, 0x64,
			0xf7, 0xa7, 0x28, 0x1b, 0x24, 0x9d, 0xc7, 0xc7
		}},
		// 346a11a8bd8fcedfcde2e19c996b6e4497d0dafc3f5af7096c915bd0f9fe4fe9
		{{
			0x34, 0x6a, 0x11, 0xa8, 0xbd, 0x8f, 0xce, 0xdf,
			0xcd, 0xe2, 0xe1, 0x9c, 0x99, 0x6b, 0x6e, 0x44,
			0x97, 0xd0, 0xda, 0xfc, 0x3f, 0x5a, 0xf7, 0x09,
			0x6c, 0x91, 0x5b, 0xd0, 0xf9, 0xfe, 0x4f, 0xe9
		}},
		// 3bac2ada2fbfa1ea75b2cb214490d5d718f1bbe5b226184488c07cf1a551e8d9
		{{
			0x3b, 0xac, 0x2a, 0xda, 0x2f, 0xbf, 0xa1, 0xea,
			0x75, 0xb2, 0xcb, 0x21, 0x44, 0x90, 0xd5, 0xd7,
			0x18, 0xf1, 0xbb, 0xe5, 0xb2, 0x26, 0x18, 0x44,
			0x88, 0xc0, 0x7c, 0xf1, 0xa5, 0x51, 0xe8, 0xd9
		}},
	};

	{
		ecc_25519_work_t w = ecc_25519_work_identity;

		ecc_25519_double(&w, &ecc_25519_work_base_legacy);
		saveWork("cases/ecc_point_double", &w);
	}

	{
		ecc_25519_work_t w = {0};

		ecc_25519_add(&w, &ecc_25519_work_identity, &ecc_25519_work_base_legacy);
		saveWork("cases/ecc_point_add", &w);
	}

	const size_t flen = 32*sizeof(char); // [sic]
	char *filename = malloc(flen);
	for (int i = 0; i < 4; ++i) {
		ecc_int256_t key = testKeys[i];

		snprintf(filename, flen, "cases/ecc_key_%d", i);
		if (access(filename, F_OK) == -1) {
			fprintf(stderr, " ecc_key_%d", i);
			saveInt256(filename, &key);
		}

		ecc_25519_work_t p;
		if (ecc_25519_load_packed_legacy(&p, &key) == 0) {
			puts("");
			error(1, 0, "failed to unpack key");
		}

		snprintf(filename, flen, "cases/ecc_key_unpacked_%d", i);
		if (1||access(filename, F_OK) == -1) {
			fprintf(stderr, " ecc_key_unpacked_%d", i);
			saveWork(filename, &p);
		}

		snprintf(filename, flen, "cases/ecc_key_derived_public_%d", i);
		{
			ecc_25519_work_t work;
			ecc_int256_t pub;
			ecc_25519_scalarmult_bits(&work, &key, &ecc_25519_work_base_legacy, 256);
			ecc_25519_store_packed_legacy(&pub, &work);
			saveInt256(filename, &pub);
		}
	}
	free(filename);

	puts("\ndone.");
	return 0;
}

