package osint

import (
	"strings"
	"testing"
)

func TestExtractNorad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid format with parentheses",
			input:    "ISS (ZARYA) (25544)",
			expected: "ZARYA",
		},
		{
			name:     "Simple format",
			input:    "Satellite Name (12345)",
			expected: "12345",
		},
		{
			name:     "Multiple parentheses - takes first",
			input:    "Name (NORAD_ID) (extra)",
			expected: "NORAD_ID",
		},
		{
			name:     "No parentheses",
			input:    "Satellite Name",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only opening parenthesis",
			input:    "Name (NORAD",
			expected: "",
		},
		{
			name:     "Only closing parenthesis",
			input:    "Name NORAD)",
			expected: "",
		},
		{
			name:     "Reversed parentheses",
			input:    "Name )NORAD(",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNorad(tt.input)
			if result != tt.expected {
				t.Errorf("extractNorad(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenRowString(t *testing.T) {
	tests := []struct {
		name        string
		intro       string
		input       string
		checkPrefix bool
		checkSuffix bool
		checkLength bool
	}{
		{
			name:        "Normal case",
			intro:       "Name",
			input:       "Test",
			checkPrefix: true,
			checkSuffix: true,
			checkLength: true,
		},
		{
			name:        "Empty input",
			intro:       "Field",
			input:       "",
			checkPrefix: true,
			checkSuffix: true,
			checkLength: true,
		},
		{
			name:        "Short intro and input",
			intro:       "ID",
			input:       "123",
			checkPrefix: true,
			checkSuffix: true,
			checkLength: true,
		},
		{
			name:        "Long intro",
			intro:       "Very Long Field Name That Takes Up Space",
			input:       "Value",
			checkPrefix: true,
			checkSuffix: true,
			checkLength: false, // May exceed 67 chars if too long
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenRowString(tt.intro, tt.input)
			// Check that result has correct structure
			if tt.checkPrefix && !strings.HasPrefix(result, "║ ") {
				t.Errorf("GenRowString() should start with '║ '")
			}
			if tt.checkSuffix && !strings.HasSuffix(result, " ║") {
				t.Errorf("GenRowString() should end with ' ║'")
			}
			if !strings.Contains(result, tt.intro+": "+tt.input) {
				t.Errorf("GenRowString() should contain intro and input")
			}
			// Check total length only if not too long
			if tt.checkLength && len(result) != 67 {
				t.Errorf("GenRowString() length = %d, want 67", len(result))
			}
		})
	}
}

// Benchmark tests
func BenchmarkExtractNorad(b *testing.B) {
	testCases := []string{
		"ISS (ZARYA) (25544)",
		"Satellite Name (12345)",
		"Name (NORAD_ID)",
		"Invalid Format",
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			extractNorad(tc)
		}
	}
}

func BenchmarkGenRowString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRowString("Field Name", "Field Value")
	}
}

func TestBuildSatcatQuery(t *testing.T) {
	tests := []struct {
		name        string
		searchName  string
		country     string
		objectType  string
		launchYear  string
		page        int
		pageSize    int
		wantContain []string
	}{
		{
			name:        "Basic query with no filters",
			searchName:  "",
			country:     "",
			objectType:  "",
			launchYear:  "",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/class/satcat", "/orderby/SATNAME%20asc", "/limit/20,0", "/emptyresult/show"},
		},
		{
			name:        "Query with country filter",
			searchName:  "",
			country:     "US",
			objectType:  "",
			launchYear:  "",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/class/satcat", "/COUNTRY/US", "/orderby/SATNAME%20asc"},
		},
		{
			name:        "Query with object type filter",
			searchName:  "",
			country:     "",
			objectType:  "PAYLOAD",
			launchYear:  "",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/class/satcat", "/OBJECT_TYPE/PAYLOAD", "/orderby/SATNAME%20asc"},
		},
		{
			name:        "Query with launch year filter",
			searchName:  "",
			country:     "",
			objectType:  "",
			launchYear:  "2020",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/class/satcat", "/LAUNCH_YEAR/2020", "/orderby/SATNAME%20asc"},
		},
		{
			name:        "Query with name search (should fetch more)",
			searchName:  "ISS",
			country:     "",
			objectType:  "",
			launchYear:  "",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/class/satcat", "/limit/500"},
		},
		{
			name:        "Query with pagination",
			searchName:  "",
			country:     "",
			objectType:  "",
			launchYear:  "",
			page:        3,
			pageSize:    20,
			wantContain: []string{"/limit/20,40"},
		},
		{
			name:        "Query with multiple filters",
			searchName:  "",
			country:     "US",
			objectType:  "PAYLOAD",
			launchYear:  "2020",
			page:        1,
			pageSize:    20,
			wantContain: []string{"/COUNTRY/US", "/OBJECT_TYPE/PAYLOAD", "/LAUNCH_YEAR/2020"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSatcatQuery(tt.searchName, tt.country, tt.objectType, tt.launchYear, tt.page, tt.pageSize)
			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("buildSatcatQuery() = %q, should contain %q", result, want)
				}
			}
		})
	}
}

func TestFilterSatellitesByName(t *testing.T) {
	sats := []Satellite{
		{SATNAME: "ISS (ZARYA)", NORAD_CAT_ID: "25544"},
		{SATNAME: "STARLINK-1007", NORAD_CAT_ID: "44700"},
		{SATNAME: "NOAA 15", NORAD_CAT_ID: "25338"},
		{SATNAME: "STARLINK-1008", NORAD_CAT_ID: "44701"},
		{SATNAME: "Hubble Space Telescope", NORAD_CAT_ID: "20580"},
	}

	tests := []struct {
		name       string
		searchName string
		wantCount  int
		wantNames  []string
	}{
		{
			name:       "Empty search returns all",
			searchName: "",
			wantCount:  5,
			wantNames:  []string{"ISS (ZARYA)", "STARLINK-1007", "NOAA 15", "STARLINK-1008", "Hubble Space Telescope"},
		},
		{
			name:       "Search for STARLINK",
			searchName: "STARLINK",
			wantCount:  2,
			wantNames:  []string{"STARLINK-1007", "STARLINK-1008"},
		},
		{
			name:       "Search case insensitive",
			searchName: "iss",
			wantCount:  1,
			wantNames:  []string{"ISS (ZARYA)"},
		},
		{
			name:       "Search partial match",
			searchName: "100",
			wantCount:  2,
			wantNames:  []string{"STARLINK-1007", "STARLINK-1008"},
		},
		{
			name:       "Search with no matches",
			searchName: "XYZ",
			wantCount:  0,
			wantNames:  []string{},
		},
		{
			name:       "Search for Hubble",
			searchName: "Hubble",
			wantCount:  1,
			wantNames:  []string{"Hubble Space Telescope"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterSatellitesByName(sats, tt.searchName)
			if len(result) != tt.wantCount {
				t.Errorf("filterSatellitesByName() returned %d results, want %d", len(result), tt.wantCount)
			}
			for i, wantName := range tt.wantNames {
				if i < len(result) && result[i].SATNAME != wantName {
					t.Errorf("filterSatellitesByName() result[%d].SATNAME = %q, want %q", i, result[i].SATNAME, wantName)
				}
			}
		})
	}
}

func BenchmarkBuildSatcatQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildSatcatQuery("ISS", "US", "PAYLOAD", "2020", 1, 20)
	}
}

func BenchmarkFilterSatellitesByName(b *testing.B) {
	sats := []Satellite{
		{SATNAME: "ISS (ZARYA)", NORAD_CAT_ID: "25544"},
		{SATNAME: "STARLINK-1007", NORAD_CAT_ID: "44700"},
		{SATNAME: "NOAA 15", NORAD_CAT_ID: "25338"},
		{SATNAME: "STARLINK-1008", NORAD_CAT_ID: "44701"},
		{SATNAME: "Hubble Space Telescope", NORAD_CAT_ID: "20580"},
	}

	for i := 0; i < b.N; i++ {
		filterSatellitesByName(sats, "STARLINK")
	}
}

