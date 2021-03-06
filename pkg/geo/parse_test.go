// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package geo

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/geo/geopb"
	"github.com/stretchr/testify/require"
)

func TestParseWKB(t *testing.T) {
	testCases := []struct {
		desc          string
		b             []byte
		defaultSRID   geopb.SRID
		expected      geopb.EWKB
		expectedError string
	}{
		{
			"EWKB should make this error",
			[]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
			4326,
			[]byte(""),
			"wkb: unknown type: 536870913",
		},
		{
			"Normal WKB should take the SRID",
			[]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
			4326,
			[]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret, err := ParseWKB(tc.b, tc.defaultSRID)
			if tc.expectedError != "" {
				require.Error(t, err)
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, ret)
			}
		})
	}
}

func TestParseEWKB(t *testing.T) {
	testCases := []struct {
		desc        string
		b           []byte
		defaultSRID geopb.SRID
		overwrite   defaultSRIDOverwriteSetting
		expected    geopb.EWKB
	}{
		{
			"SRID 4326 is hint; EWKB has 3857",
			[]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
			4326,
			DefaultSRIDIsHint,
			[]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
		},
		{
			"Overwrite SRID 3857 with 4326",
			[]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
			4326,
			DefaultSRIDShouldOverwrite,
			[]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret, err := ParseEWKB(tc.b, tc.defaultSRID, tc.overwrite)
			require.NoError(t, err)
			require.Equal(t, tc.expected, ret)
		})
	}
}

func TestParseEWKT(t *testing.T) {
	testCases := []struct {
		desc        string
		wkt         geopb.EWKT
		defaultSRID geopb.SRID
		overwrite   defaultSRIDOverwriteSetting
		expected    geopb.EWKB
	}{
		{
			"SRID 4326 is hint; EWKT has 3857",
			"SRID=3857;POINT(1.0 1.0)",
			4326,
			DefaultSRIDIsHint,
			[]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
		},
		{
			"Overwrite SRID 3857 with 4326",
			"SRID=3857;POINT(1.0 1.0)",
			4326,
			DefaultSRIDShouldOverwrite,
			[]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ret, err := ParseEWKT(tc.wkt, tc.defaultSRID, tc.overwrite)
			require.NoError(t, err)
			require.Equal(t, tc.expected, ret)
		})
	}
}

func TestParseGeometry(t *testing.T) {
	testCases := []struct {
		str          string
		expected     *Geometry
		expectedErr  string
		expectedSRID geopb.SRID
	}{
		{
			"0101000000000000000000F03F000000000000F03F",
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			0,
		},
		{
			"0101000020E6100000000000000000F03F000000000000F03F",
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"POINT(1.0 1.0)",
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			0,
		},
		{
			"SRID=3857;POINT(1.0 1.0)",
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			3857,
		},
		{
			"\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f",
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			0,
		},
		{
			`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [1.0, 1.0] }, "properties": { "name": "┳━┳ ヽ(ಠل͜ಠ)ﾉ" } }`,
			NewGeometry(geopb.EWKB([]byte("\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			0,
		},
		{
			"invalid",
			nil,
			"geos error: ParseException: Unknown type: 'INVALID'",
			0,
		},
		{
			"",
			nil,
			"geo: parsing empty string to geo type",
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			g, err := ParseGeometry(tc.str)
			if len(tc.expectedErr) > 0 {
				require.Equal(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, g)
				require.Equal(t, tc.expectedSRID, g.SRID())
			}
		})
	}
}

func TestParseGeography(t *testing.T) {
	testCases := []struct {
		str          string
		expected     *Geography
		expectedErr  string
		expectedSRID geopb.SRID
	}{
		{
			// Even forcing an SRID to 0 using EWKB will make it 4326.
			"0101000000000000000000F03F000000000000F03F",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"0101000020E6100000000000000000F03F000000000000F03F",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"0101000020110F0000000000000000F03F000000000000F03F",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			3857,
		},
		{
			"POINT(1.0 1.0)",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			// Even forcing an SRID to 0 using WKT will make it 4326.
			"SRID=0;POINT(1.0 1.0)",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"SRID=3857;POINT(1.0 1.0)",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			3857,
		},
		{
			"\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f",
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\x11\x0F\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			3857,
		},
		{
			`{ "type": "Feature", "geometry": { "type": "Point", "coordinates": [1.0, 1.0] }, "properties": { "name": "┳━┳ ヽ(ಠل͜ಠ)ﾉ" } }`,
			NewGeography(geopb.EWKB([]byte("\x01\x01\x00\x00\x20\xe6\x10\x00\x00\x00\x00\x00\x00\x00\x00\xf0\x3f\x00\x00\x00\x00\x00\x00\xf0\x3f"))),
			"",
			4326,
		},
		{
			"invalid",
			nil,
			"geos error: ParseException: Unknown type: 'INVALID'",
			0,
		},
		{
			"",
			nil,
			"geo: parsing empty string to geo type",
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			g, err := ParseGeography(tc.str)
			if len(tc.expectedErr) > 0 {
				require.Equal(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, g)
				require.Equal(t, tc.expectedSRID, g.SRID())
			}
		})
	}
}
