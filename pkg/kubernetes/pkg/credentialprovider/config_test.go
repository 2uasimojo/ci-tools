package credentialprovider

import (
	"encoding/json"
	"testing"
)

func TestDockerConfigJSONUNmarshal(t *testing.T) {
	testCase := []struct {
		name                string
		in                  []byte
		expectedErrorString string
	}{
		{
			name: "Duplicate data inside the auth field, error",
			in: []byte(`{
  "auths": {
    "registry.ci.openshift.org": {
      "auth": "c2VydmljZWFjY291bnQ6ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrRndTekF0YjBaNGJXMUZURXRHTVMwMFVEa3djbEEwUTJWQlRUZERNMGRXUkZwdmJGOVllaTFEUW5NaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpoYkhaaGNtOHRkR1Z6ZENJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKa1pXWmhkV3gwTFhSdmEyVnVMV1EwT1d4aUlpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltUmxabUYxYkhRaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNTFhV1FpT2lJM05tVTRZMlpsTmkxbU1HWXhMVFF5WlRNdFlqUm1NQzFoTXpjM1pUbGhOemxrWWpRaUxDSnpkV0lpT2lKemVYTjBaVzA2YzJWeWRtbGpaV0ZqWTI5MWJuUTZZV3gyWVhKdkxYUmxjM1E2WkdWbVlYVnNkQ0o5LnMyajh6X2JfT3NMOHY5UGlLR1NUQmFuZDE0MHExMHc3VTlMdU9JWmZlUG1SeF9OMHdKRkZPcVN0MGNjdmtVaUVGV0x5QWNSU2k2cUt3T1FSVzE2MVUzSU52UEY4Q0pDZ2d2R3JHUnMzeHp6N3hjSmgzTWRpcXhzWGViTmNmQmlmWWxXUTU2U1RTZDlUeUh1RkN6c1poNXBlSHVzS3hOa2hJRTNyWHp5ZHNoMkhCaTZMYTlYZ1l4R1VjM0x3NWh4RnB5bXFyajFJNzExbWZLcUV2bUN0a0J4blJtMlhIZmFKalNVRkswWWdoY0lMbkhuWGhMOEx2MUl0bnU4SzlvWFRfWVZIQWY1R3hlaERjZ3FBMmw1NUZyYkJMTGVfNi1DV2V2N2RQZU5PbFlaWE5xbEtkUG5KbW9BREdsOEktTlhKN2x5ZXl2a2hfZ3JkanhXdVVqQ3lQUQ==c2VydmljZWFjY291bnQ6ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrRndTekF0YjBaNGJXMUZURXRHTVMwMFVEa3djbEEwUTJWQlRUZERNMGRXUkZwdmJGOVllaTFEUW5NaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpoYkhaaGNtOHRkR1Z6ZENJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKa1pXWmhkV3gwTFhSdmEyVnVMVFpzTW0xcUlpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltUmxabUYxYkhRaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNTFhV1FpT2lJM05tVTRZMlpsTmkxbU1HWXhMVFF5WlRNdFlqUm1NQzFoTXpjM1pUbGhOemxrWWpRaUxDSnpkV0lpT2lKemVYTjBaVzA2YzJWeWRtbGpaV0ZqWTI5MWJuUTZZV3gyWVhKdkxYUmxjM1E2WkdWbVlYVnNkQ0o5LnM5YS1Ucy1sbkZuSUYzSmNPSHM4TTJZS3VaX0dXcjZ1eTdaZTZVTF82YmZub1BDNTNwSHVBYzBDdTRDREVKQ3MxQm56S29BNlNrYW5qWFBzRHJnWlpYQndGMXJQOWduQ0p2d3N5REE0eTdBQmtyNW5pSEZfU2djSDNyQUxKNkhLYU1jTnhyQUd1d2hZOWFzREEzNlZOVmJGYVl6ejl5WW1GcHJFWnhiYlZKbE9CZVJZSDJ3eDA1WEVqQnhRWG9vS1JVTi02a3pXSnZ3WGlIdXhFd3ZKOFcwekc4c2ZmYXJFeW9oZG40ZVp2bkd1eW5XelFEcHNaUktsZGtmTmhIeVQxMm50cG5nS1g1TEhRd3J2bEF6S2k2OGlOZWc3SkhRY29XcmNXaS1uU1ljWDhpZnZaQlNHam8tazFDbXRYUHJqY08weWx1eVFjcS1XaDQzdGYycmZxUQ=="
    }
  }
}
`),
			expectedErrorString: "illegal base64 data at input byte 1236",
		},
		{
			name: "Parseable data, no error",
			in: []byte(`{
  "auths": {
    "registry.ci.openshift.org": {
      "auth": "c2VydmljZWFjY291bnQ6ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrRndTekF0YjBaNGJXMUZURXRHTVMwMFVEa3djbEEwUTJWQlRUZERNMGRXUkZwdmJGOVllaTFEUW5NaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpoYkhaaGNtOHRkR1Z6ZENJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKa1pXWmhkV3gwTFhSdmEyVnVMV1EwT1d4aUlpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltUmxabUYxYkhRaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNTFhV1FpT2lJM05tVTRZMlpsTmkxbU1HWXhMVFF5WlRNdFlqUm1NQzFoTXpjM1pUbGhOemxrWWpRaUxDSnpkV0lpT2lKemVYTjBaVzA2YzJWeWRtbGpaV0ZqWTI5MWJuUTZZV3gyWVhKdkxYUmxjM1E2WkdWbVlYVnNkQ0o5LnMyajh6X2JfT3NMOHY5UGlLR1NUQmFuZDE0MHExMHc3VTlMdU9JWmZlUG1SeF9OMHdKRkZPcVN0MGNjdmtVaUVGV0x5QWNSU2k2cUt3T1FSVzE2MVUzSU52UEY4Q0pDZ2d2R3JHUnMzeHp6N3hjSmgzTWRpcXhzWGViTmNmQmlmWWxXUTU2U1RTZDlUeUh1RkN6c1poNXBlSHVzS3hOa2hJRTNyWHp5ZHNoMkhCaTZMYTlYZ1l4R1VjM0x3NWh4RnB5bXFyajFJNzExbWZLcUV2bUN0a0J4blJtMlhIZmFKalNVRkswWWdoY0lMbkhuWGhMOEx2MUl0bnU4SzlvWFRfWVZIQWY1R3hlaERjZ3FBMmw1NUZyYkJMTGVfNi1DV2V2N2RQZU5PbFlaWE5xbEtkUG5KbW9BREdsOEktTlhKN2x5ZXl2a2hfZ3JkanhXdVVqQ3lQUQ=="
    }
  }
}
`),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			var errMsg string
			if err := json.Unmarshal(tc.in, &DockerConfigJSON{}); err != nil {
				errMsg = err.Error()
			}
			if tc.expectedErrorString != errMsg {
				t.Errorf("expected error %s, got error %s", tc.expectedErrorString, errMsg)
			}
		})
	}
}