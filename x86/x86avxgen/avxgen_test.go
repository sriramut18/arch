// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"golang.org/x/arch/x86/xeddata"
)

func newTestContext(t testing.TB) *context {
	ctx := &context{xedPath: filepath.Join("testdata", "xedpath")}
	db, err := xeddata.NewDatabase(ctx.xedPath)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	ctx.db = db
	return ctx
}

func newStringSet(keys ...string) map[string]bool {
	set := make(map[string]bool)
	for _, k := range keys {
		set[k] = true
	}
	return set
}

func generateToString(t *testing.T) string {
	ctx := newTestContext(t)
	buildTables(ctx)
	var buf bytes.Buffer
	writeTables(&buf, ctx)
	return buf.String()
}

func TestOutput(t *testing.T) {
	// Ytab lists and optabs output checks.
	//
	// These tests are very fragile.
	// Slight changes can invalidate them.
	// It is better to keep testCases count at the minimum.

	type testCase struct {
		opcode     string
		ytabs      string
		optabLines string
	}
	var testCases []testCase
	{
		opcodeRE := regexp.MustCompile(`as: ([A-Z][A-Z0-9]*)`)
		data, err := ioutil.ReadFile(filepath.Join("testdata", "golden.txt"))
		if err != nil {
			t.Fatalf("read golden file: %v", err)
		}
		for _, entry := range bytes.Split(data, []byte("======")) {
			parts := bytes.Split(entry, []byte("----"))
			ytabs := parts[0]
			optabLines := parts[1]
			opcode := opcodeRE.FindSubmatch(optabLines)[1]
			testCases = append(testCases, testCase{
				ytabs:      strings.TrimSpace(string(ytabs)),
				optabLines: strings.TrimSpace(string(optabLines)),
				opcode:     string(opcode)[len("A"):],
			})
		}
	}

	output := generateToString(t)
	for _, tc := range testCases {
		if !strings.Contains(output, tc.ytabs) {
			t.Errorf("%s: ytabs not matched", tc.opcode)
		}
		if !strings.Contains(output, tc.optabLines) {
			t.Errorf("%s: optab lines not matched", tc.opcode)
		}
	}
}

func TestOutputStability(t *testing.T) {
	// Generate output count+1 times and check that every time
	// it is exactly the same string.
	//
	// The output should be deterministic to avoid unwanted diffs
	// between each code generation.
	const count = 8

	want := generateToString(t)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			if want != generateToString(t) {
				t.Errorf("output #%d mismatches", i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestOpcodeCoverage(t *testing.T) {
	// Check that generator produces all expected opcodes from testdata files.
	// All opcodes are in Go syntax.

	// VEX/EVEX opcodes collected from XED-based x86.csv.
	expectedOpcodes := newStringSet(
		"ANDNL",
		"ANDNQ",
		"BEXTRL",
		"BEXTRQ",
		"BLSIL",
		"BLSIQ",
		"BLSMSKL",
		"BLSMSKQ",
		"BLSRL",
		"BLSRQ",
		"BZHIL",
		"BZHIQ",
		"KADDB",
		"KADDD",
		"KADDQ",
		"KADDW",
		"KANDB",
		"KANDD",
		"KANDNB",
		"KANDND",
		"KANDNQ",
		"KANDNW",
		"KANDQ",
		"KANDW",
		"KMOVB",
		"KMOVD",
		"KMOVQ",
		"KMOVW",
		"KNOTB",
		"KNOTD",
		"KNOTQ",
		"KNOTW",
		"KORB",
		"KORD",
		"KORQ",
		"KORTESTB",
		"KORTESTD",
		"KORTESTQ",
		"KORTESTW",
		"KORW",
		"KSHIFTLB",
		"KSHIFTLD",
		"KSHIFTLQ",
		"KSHIFTLW",
		"KSHIFTRB",
		"KSHIFTRD",
		"KSHIFTRQ",
		"KSHIFTRW",
		"KTESTB",
		"KTESTD",
		"KTESTQ",
		"KTESTW",
		"KUNPCKBW",
		"KUNPCKDQ",
		"KUNPCKWD",
		"KXNORB",
		"KXNORD",
		"KXNORQ",
		"KXNORW",
		"KXORB",
		"KXORD",
		"KXORQ",
		"KXORW",
		"MULXL",
		"MULXQ",
		"PDEPL",
		"PDEPQ",
		"PEXTL",
		"PEXTQ",
		"RORXL",
		"RORXQ",
		"SARXL",
		"SARXQ",
		"SHLXL",
		"SHLXQ",
		"SHRXL",
		"SHRXQ",
		"V4FMADDPS",
		"V4FMADDSS",
		"V4FNMADDPS",
		"V4FNMADDSS",
		"VADDPD",
		"VADDPS",
		"VADDSD",
		"VADDSS",
		"VADDSUBPD",
		"VADDSUBPS",
		"VAESDEC",
		"VAESDECLAST",
		"VAESENC",
		"VAESENCLAST",
		"VAESIMC",
		"VAESKEYGENASSIST",
		"VALIGND",
		"VALIGNQ",
		"VANDNPD",
		"VANDNPS",
		"VANDPD",
		"VANDPS",
		"VBLENDMPD",
		"VBLENDMPS",
		"VBLENDPD",
		"VBLENDPS",
		"VBLENDVPD",
		"VBLENDVPS",
		"VBROADCASTF128",
		"VBROADCASTF32X2",
		"VBROADCASTF32X4",
		"VBROADCASTF32X8",
		"VBROADCASTF64X2",
		"VBROADCASTF64X4",
		"VBROADCASTI128",
		"VBROADCASTI32X2",
		"VBROADCASTI32X4",
		"VBROADCASTI32X8",
		"VBROADCASTI64X2",
		"VBROADCASTI64X4",
		"VBROADCASTSD",
		"VBROADCASTSS",
		"VCMPPD",
		"VCMPPS",
		"VCMPSD",
		"VCMPSS",
		"VCOMISD",
		"VCOMISS",
		"VCOMPRESSPD",
		"VCOMPRESSPS",
		"VCVTDQ2PD",
		"VCVTDQ2PS",
		"VCVTPD2DQ",
		"VCVTPD2DQX",
		"VCVTPD2DQY",
		"VCVTPD2PS",
		"VCVTPD2PSX",
		"VCVTPD2PSY",
		"VCVTPD2QQ",
		"VCVTPD2UDQ",
		"VCVTPD2UDQX",
		"VCVTPD2UDQY",
		"VCVTPD2UQQ",
		"VCVTPH2PS",
		"VCVTPS2DQ",
		"VCVTPS2PD",
		"VCVTPS2PH",
		"VCVTPS2QQ",
		"VCVTPS2UDQ",
		"VCVTPS2UQQ",
		"VCVTQQ2PD",
		"VCVTQQ2PS",
		"VCVTQQ2PSX",
		"VCVTQQ2PSY",
		"VCVTSD2SI",
		"VCVTSD2SIQ",
		"VCVTSD2SS",
		"VCVTSD2USIL",
		"VCVTSD2USIQ",
		"VCVTSI2SDL",
		"VCVTSI2SDQ",
		"VCVTSI2SSL",
		"VCVTSI2SSQ",
		"VCVTSS2SD",
		"VCVTSS2SI",
		"VCVTSS2SIQ",
		"VCVTSS2USIL",
		"VCVTSS2USIQ",
		"VCVTTPD2DQ",
		"VCVTTPD2DQX",
		"VCVTTPD2DQY",
		"VCVTTPD2QQ",
		"VCVTTPD2UDQ",
		"VCVTTPD2UDQX",
		"VCVTTPD2UDQY",
		"VCVTTPD2UQQ",
		"VCVTTPS2DQ",
		"VCVTTPS2QQ",
		"VCVTTPS2UDQ",
		"VCVTTPS2UQQ",
		"VCVTTSD2SI",
		"VCVTTSD2SIQ",
		"VCVTTSD2USIL",
		"VCVTTSD2USIQ",
		"VCVTTSS2SI",
		"VCVTTSS2SIQ",
		"VCVTTSS2USIL",
		"VCVTTSS2USIQ",
		"VCVTUDQ2PD",
		"VCVTUDQ2PS",
		"VCVTUQQ2PD",
		"VCVTUQQ2PS",
		"VCVTUQQ2PSX",
		"VCVTUQQ2PSY",
		"VCVTUSI2SDL",
		"VCVTUSI2SDQ",
		"VCVTUSI2SSL",
		"VCVTUSI2SSQ",
		"VDBPSADBW",
		"VDIVPD",
		"VDIVPS",
		"VDIVSD",
		"VDIVSS",
		"VDPPD",
		"VDPPS",
		"VEXP2PD",
		"VEXP2PS",
		"VEXPANDPD",
		"VEXPANDPS",
		"VEXTRACTF128",
		"VEXTRACTF32X4",
		"VEXTRACTF32X8",
		"VEXTRACTF64X2",
		"VEXTRACTF64X4",
		"VEXTRACTI128",
		"VEXTRACTI32X4",
		"VEXTRACTI32X8",
		"VEXTRACTI64X2",
		"VEXTRACTI64X4",
		"VEXTRACTPS",
		"VFIXUPIMMPD",
		"VFIXUPIMMPS",
		"VFIXUPIMMSD",
		"VFIXUPIMMSS",
		"VFMADD132PD",
		"VFMADD132PS",
		"VFMADD132SD",
		"VFMADD132SS",
		"VFMADD213PD",
		"VFMADD213PS",
		"VFMADD213SD",
		"VFMADD213SS",
		"VFMADD231PD",
		"VFMADD231PS",
		"VFMADD231SD",
		"VFMADD231SS",
		"VFMADDPD",
		"VFMADDPS",
		"VFMADDSD",
		"VFMADDSS",
		"VFMADDSUB132PD",
		"VFMADDSUB132PS",
		"VFMADDSUB213PD",
		"VFMADDSUB213PS",
		"VFMADDSUB231PD",
		"VFMADDSUB231PS",
		"VFMADDSUBPD",
		"VFMADDSUBPS",
		"VFMSUB132PD",
		"VFMSUB132PS",
		"VFMSUB132SD",
		"VFMSUB132SS",
		"VFMSUB213PD",
		"VFMSUB213PS",
		"VFMSUB213SD",
		"VFMSUB213SS",
		"VFMSUB231PD",
		"VFMSUB231PS",
		"VFMSUB231SD",
		"VFMSUB231SS",
		"VFMSUBADD132PD",
		"VFMSUBADD132PS",
		"VFMSUBADD213PD",
		"VFMSUBADD213PS",
		"VFMSUBADD231PD",
		"VFMSUBADD231PS",
		"VFMSUBADDPD",
		"VFMSUBADDPS",
		"VFMSUBPD",
		"VFMSUBPS",
		"VFMSUBSD",
		"VFMSUBSS",
		"VFNMADD132PD",
		"VFNMADD132PS",
		"VFNMADD132SD",
		"VFNMADD132SS",
		"VFNMADD213PD",
		"VFNMADD213PS",
		"VFNMADD213SD",
		"VFNMADD213SS",
		"VFNMADD231PD",
		"VFNMADD231PS",
		"VFNMADD231SD",
		"VFNMADD231SS",
		"VFNMADDPD",
		"VFNMADDPS",
		"VFNMADDSD",
		"VFNMADDSS",
		"VFNMSUB132PD",
		"VFNMSUB132PS",
		"VFNMSUB132SD",
		"VFNMSUB132SS",
		"VFNMSUB213PD",
		"VFNMSUB213PS",
		"VFNMSUB213SD",
		"VFNMSUB213SS",
		"VFNMSUB231PD",
		"VFNMSUB231PS",
		"VFNMSUB231SD",
		"VFNMSUB231SS",
		"VFNMSUBPD",
		"VFNMSUBPS",
		"VFNMSUBSD",
		"VFNMSUBSS",
		"VFPCLASSPDX",
		"VFPCLASSPDY",
		"VFPCLASSPDZ",
		"VFPCLASSPSX",
		"VFPCLASSPSY",
		"VFPCLASSPSZ",
		"VFPCLASSSD",
		"VFPCLASSSS",
		"VGATHERDPD",
		"VGATHERDPS",
		"VGATHERPF0DPD",
		"VGATHERPF0DPS",
		"VGATHERPF0QPD",
		"VGATHERPF0QPS",
		"VGATHERPF1DPD",
		"VGATHERPF1DPS",
		"VGATHERPF1QPD",
		"VGATHERPF1QPS",
		"VGATHERQPD",
		"VGATHERQPS",
		"VGETEXPPD",
		"VGETEXPPS",
		"VGETEXPSD",
		"VGETEXPSS",
		"VGETMANTPD",
		"VGETMANTPS",
		"VGETMANTSD",
		"VGETMANTSS",
		"VGF2P8AFFINEINVQB",
		"VGF2P8AFFINEQB",
		"VGF2P8MULB",
		"VHADDPD",
		"VHADDPS",
		"VHSUBPD",
		"VHSUBPS",
		"VINSERTF128",
		"VINSERTF32X4",
		"VINSERTF32X8",
		"VINSERTF64X2",
		"VINSERTF64X4",
		"VINSERTI128",
		"VINSERTI32X4",
		"VINSERTI32X8",
		"VINSERTI64X2",
		"VINSERTI64X4",
		"VINSERTPS",
		"VLDDQU",
		"VLDMXCSR",
		"VMASKMOVDQU",
		"VMASKMOVPD",
		"VMASKMOVPS",
		"VMAXPD",
		"VMAXPS",
		"VMAXSD",
		"VMAXSS",
		"VMINPD",
		"VMINPS",
		"VMINSD",
		"VMINSS",
		"VMOVAPD",
		"VMOVAPS",
		"VMOVD",
		"VMOVDDUP",
		"VMOVDQA",
		"VMOVDQA32",
		"VMOVDQA64",
		"VMOVDQU",
		"VMOVDQU16",
		"VMOVDQU32",
		"VMOVDQU64",
		"VMOVDQU8",
		"VMOVHLPS",
		"VMOVHPD",
		"VMOVHPS",
		"VMOVLHPS",
		"VMOVLPD",
		"VMOVLPS",
		"VMOVMSKPD",
		"VMOVMSKPS",
		"VMOVNTDQ",
		"VMOVNTDQA",
		"VMOVNTPD",
		"VMOVNTPS",
		"VMOVQ",
		"VMOVSD",
		"VMOVSHDUP",
		"VMOVSLDUP",
		"VMOVSS",
		"VMOVUPD",
		"VMOVUPS",
		"VMPSADBW",
		"VMULPD",
		"VMULPS",
		"VMULSD",
		"VMULSS",
		"VORPD",
		"VORPS",
		"VP4DPWSSD",
		"VP4DPWSSDS",
		"VPABSB",
		"VPABSD",
		"VPABSQ",
		"VPABSW",
		"VPACKSSDW",
		"VPACKSSWB",
		"VPACKUSDW",
		"VPACKUSWB",
		"VPADDB",
		"VPADDD",
		"VPADDQ",
		"VPADDSB",
		"VPADDSW",
		"VPADDUSB",
		"VPADDUSW",
		"VPADDW",
		"VPALIGNR",
		"VPAND",
		"VPANDD",
		"VPANDN",
		"VPANDND",
		"VPANDNQ",
		"VPANDQ",
		"VPAVGB",
		"VPAVGW",
		"VPBLENDD",
		"VPBLENDMB",
		"VPBLENDMD",
		"VPBLENDMQ",
		"VPBLENDMW",
		"VPBLENDVB",
		"VPBLENDW",
		"VPBROADCASTB",
		"VPBROADCASTD",
		"VPBROADCASTMB2Q",
		"VPBROADCASTMW2D",
		"VPBROADCASTQ",
		"VPBROADCASTW",
		"VPCLMULQDQ",
		"VPCMPB",
		"VPCMPD",
		"VPCMPEQB",
		"VPCMPEQD",
		"VPCMPEQQ",
		"VPCMPEQW",
		"VPCMPESTRI",
		"VPCMPESTRM",
		"VPCMPGTB",
		"VPCMPGTD",
		"VPCMPGTQ",
		"VPCMPGTW",
		"VPCMPISTRI",
		"VPCMPISTRM",
		"VPCMPQ",
		"VPCMPUB",
		"VPCMPUD",
		"VPCMPUQ",
		"VPCMPUW",
		"VPCMPW",
		"VPCOMPRESSB",
		"VPCOMPRESSD",
		"VPCOMPRESSQ",
		"VPCOMPRESSW",
		"VPCONFLICTD",
		"VPCONFLICTQ",
		"VPDPBUSD",
		"VPDPBUSDS",
		"VPDPWSSD",
		"VPDPWSSDS",
		"VPERM2F128",
		"VPERM2I128",
		"VPERMB",
		"VPERMD",
		"VPERMI2B",
		"VPERMI2D",
		"VPERMI2PD",
		"VPERMI2PS",
		"VPERMI2Q",
		"VPERMI2W",
		"VPERMIL2PD",
		"VPERMIL2PS",
		"VPERMILPD",
		"VPERMILPS",
		"VPERMPD",
		"VPERMPS",
		"VPERMQ",
		"VPERMT2B",
		"VPERMT2D",
		"VPERMT2PD",
		"VPERMT2PS",
		"VPERMT2Q",
		"VPERMT2W",
		"VPERMW",
		"VPEXPANDB",
		"VPEXPANDD",
		"VPEXPANDQ",
		"VPEXPANDW",
		"VPEXTRB",
		"VPEXTRD",
		"VPEXTRQ",
		"VPEXTRW",
		"VPGATHERDD",
		"VPGATHERDQ",
		"VPGATHERQD",
		"VPGATHERQQ",
		"VPHADDD",
		"VPHADDSW",
		"VPHADDW",
		"VPHMINPOSUW",
		"VPHSUBD",
		"VPHSUBSW",
		"VPHSUBW",
		"VPINSRB",
		"VPINSRD",
		"VPINSRQ",
		"VPINSRW",
		"VPLZCNTD",
		"VPLZCNTQ",
		"VPMADD52HUQ",
		"VPMADD52LUQ",
		"VPMADDUBSW",
		"VPMADDWD",
		"VPMASKMOVD",
		"VPMASKMOVQ",
		"VPMAXSB",
		"VPMAXSD",
		"VPMAXSQ",
		"VPMAXSW",
		"VPMAXUB",
		"VPMAXUD",
		"VPMAXUQ",
		"VPMAXUW",
		"VPMINSB",
		"VPMINSD",
		"VPMINSQ",
		"VPMINSW",
		"VPMINUB",
		"VPMINUD",
		"VPMINUQ",
		"VPMINUW",
		"VPMOVB2M",
		"VPMOVD2M",
		"VPMOVDB",
		"VPMOVDW",
		"VPMOVM2B",
		"VPMOVM2D",
		"VPMOVM2Q",
		"VPMOVM2W",
		"VPMOVMSKB",
		"VPMOVQ2M",
		"VPMOVQB",
		"VPMOVQD",
		"VPMOVQW",
		"VPMOVSDB",
		"VPMOVSDW",
		"VPMOVSQB",
		"VPMOVSQD",
		"VPMOVSQW",
		"VPMOVSWB",
		"VPMOVSXBD",
		"VPMOVSXBQ",
		"VPMOVSXBW",
		"VPMOVSXDQ",
		"VPMOVSXWD",
		"VPMOVSXWQ",
		"VPMOVUSDB",
		"VPMOVUSDW",
		"VPMOVUSQB",
		"VPMOVUSQD",
		"VPMOVUSQW",
		"VPMOVUSWB",
		"VPMOVW2M",
		"VPMOVWB",
		"VPMOVZXBD",
		"VPMOVZXBQ",
		"VPMOVZXBW",
		"VPMOVZXDQ",
		"VPMOVZXWD",
		"VPMOVZXWQ",
		"VPMULDQ",
		"VPMULHRSW",
		"VPMULHUW",
		"VPMULHW",
		"VPMULLD",
		"VPMULLQ",
		"VPMULLW",
		"VPMULTISHIFTQB",
		"VPMULUDQ",
		"VPOPCNTB",
		"VPOPCNTD",
		"VPOPCNTQ",
		"VPOPCNTW",
		"VPOR",
		"VPORD",
		"VPORQ",
		"VPROLD",
		"VPROLQ",
		"VPROLVD",
		"VPROLVQ",
		"VPRORD",
		"VPRORQ",
		"VPRORVD",
		"VPRORVQ",
		"VPSADBW",
		"VPSCATTERDD",
		"VPSCATTERDQ",
		"VPSCATTERQD",
		"VPSCATTERQQ",
		"VPSHLDD",
		"VPSHLDQ",
		"VPSHLDVD",
		"VPSHLDVQ",
		"VPSHLDVW",
		"VPSHLDW",
		"VPSHRDD",
		"VPSHRDQ",
		"VPSHRDVD",
		"VPSHRDVQ",
		"VPSHRDVW",
		"VPSHRDW",
		"VPSHUFB",
		"VPSHUFBITQMB",
		"VPSHUFD",
		"VPSHUFHW",
		"VPSHUFLW",
		"VPSIGNB",
		"VPSIGND",
		"VPSIGNW",
		"VPSLLD",
		"VPSLLDQ",
		"VPSLLQ",
		"VPSLLVD",
		"VPSLLVQ",
		"VPSLLVW",
		"VPSLLW",
		"VPSRAD",
		"VPSRAQ",
		"VPSRAVD",
		"VPSRAVQ",
		"VPSRAVW",
		"VPSRAW",
		"VPSRLD",
		"VPSRLDQ",
		"VPSRLQ",
		"VPSRLVD",
		"VPSRLVQ",
		"VPSRLVW",
		"VPSRLW",
		"VPSUBB",
		"VPSUBD",
		"VPSUBQ",
		"VPSUBSB",
		"VPSUBSW",
		"VPSUBUSB",
		"VPSUBUSW",
		"VPSUBW",
		"VPTERNLOGD",
		"VPTERNLOGQ",
		"VPTEST",
		"VPTESTMB",
		"VPTESTMD",
		"VPTESTMQ",
		"VPTESTMW",
		"VPTESTNMB",
		"VPTESTNMD",
		"VPTESTNMQ",
		"VPTESTNMW",
		"VPUNPCKHBW",
		"VPUNPCKHDQ",
		"VPUNPCKHQDQ",
		"VPUNPCKHWD",
		"VPUNPCKLBW",
		"VPUNPCKLDQ",
		"VPUNPCKLQDQ",
		"VPUNPCKLWD",
		"VPXOR",
		"VPXORD",
		"VPXORQ",
		"VRANGEPD",
		"VRANGEPS",
		"VRANGESD",
		"VRANGESS",
		"VRCP14PD",
		"VRCP14PS",
		"VRCP14SD",
		"VRCP14SS",
		"VRCP28PD",
		"VRCP28PS",
		"VRCP28SD",
		"VRCP28SS",
		"VRCPPS",
		"VRCPSS",
		"VREDUCEPD",
		"VREDUCEPS",
		"VREDUCESD",
		"VREDUCESS",
		"VRNDSCALEPD",
		"VRNDSCALEPS",
		"VRNDSCALESD",
		"VRNDSCALESS",
		"VROUNDPD",
		"VROUNDPS",
		"VROUNDSD",
		"VROUNDSS",
		"VRSQRT14PD",
		"VRSQRT14PS",
		"VRSQRT14SD",
		"VRSQRT14SS",
		"VRSQRT28PD",
		"VRSQRT28PS",
		"VRSQRT28SD",
		"VRSQRT28SS",
		"VRSQRTPS",
		"VRSQRTSS",
		"VSCALEFPD",
		"VSCALEFPS",
		"VSCALEFSD",
		"VSCALEFSS",
		"VSCATTERDPD",
		"VSCATTERDPS",
		"VSCATTERPF0DPD",
		"VSCATTERPF0DPS",
		"VSCATTERPF0QPD",
		"VSCATTERPF0QPS",
		"VSCATTERPF1DPD",
		"VSCATTERPF1DPS",
		"VSCATTERPF1QPD",
		"VSCATTERPF1QPS",
		"VSCATTERQPD",
		"VSCATTERQPS",
		"VSHUFF32X4",
		"VSHUFF64X2",
		"VSHUFI32X4",
		"VSHUFI64X2",
		"VSHUFPD",
		"VSHUFPS",
		"VSQRTPD",
		"VSQRTPS",
		"VSQRTSD",
		"VSQRTSS",
		"VSTMXCSR",
		"VSUBPD",
		"VSUBPS",
		"VSUBSD",
		"VSUBSS",
		"VTESTPD",
		"VTESTPS",
		"VUCOMISD",
		"VUCOMISS",
		"VUNPCKHPD",
		"VUNPCKHPS",
		"VUNPCKLPD",
		"VUNPCKLPS",
		"VXORPD",
		"VXORPS",
		"VZEROALL",
		"VZEROUPPER")

	// AMD-specific VEX opcodes.
	// Excluded from x86avxgen output for now.
	amdOpcodes := newStringSet(
		"VFMADDPD",
		"VFMADDPS",
		"VFMADDSD",
		"VFMADDSS",
		"VFMADDSUBPD",
		"VFMADDSUBPS",
		"VFMSUBADDPD",
		"VFMSUBADDPS",
		"VFMSUBPD",
		"VFMSUBPS",
		"VFMSUBSD",
		"VFMSUBSS",
		"VFNMADDPD",
		"VFNMADDPS",
		"VFNMADDSD",
		"VFNMADDSS",
		"VFNMSUBPD",
		"VFNMSUBPS",
		"VFNMSUBSD",
		"VFNMSUBSS",
		"VPERMIL2PD",
		"VPERMIL2PS")

	ctx := newTestContext(t)
	buildTables(ctx)

	for op := range amdOpcodes {
		delete(expectedOpcodes, op)
	}
	for op := range ctx.optabs {
		delete(expectedOpcodes, op)
	}

	for op := range expectedOpcodes {
		t.Errorf("missing opcode: %s", op)
	}
}
