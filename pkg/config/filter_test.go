package config_test

import (
	"testing"

	"github.com/chemapolo/aws-nuke/pkg/config"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalFilter(t *testing.T) {

	cases := []struct {
		yaml            string
		match, mismatch []string
	}{
		{
			yaml:     `foo`,
			match:    []string{"foo"},
			mismatch: []string{"fo", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"exact","value":"foo"}`,
			match:    []string{"foo"},
			mismatch: []string{"fo", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"glob","value":"b*sh"}`,
			match:    []string{"bish", "bash", "bosh", "bush", "boooooosh", "bsh"},
			mismatch: []string{"woooosh", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"glob","value":"b?sh"}`,
			match:    []string{"bish", "bash", "bosh", "bush"},
			mismatch: []string{"woooosh", "fooo", "o", "fo", "boooooosh", "bsh"},
		},
		{
			yaml:     `{"type":"regex","value":"b[iao]sh"}`,
			match:    []string{"bish", "bash", "bosh"},
			mismatch: []string{"woooosh", "fooo", "o", "fo", "boooooosh", "bsh", "bush"},
		},
		{
			yaml:     `{"type":"contains","value":"mba"}`,
			match:    []string{"bimbaz", "mba", "bi mba z"},
			mismatch: []string{"bim-baz"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.yaml, func(t *testing.T) {
			var filter config.Filter

			err := yaml.Unmarshal([]byte(tc.yaml), &filter)
			if err != nil {
				t.Fatal(err)
			}

			for _, o := range tc.match {
				match, err := filter.Match(o)
				if err != nil {
					t.Fatal(err)
				}

				if !match {
					t.Fatalf("'%v' should match", o)
				}
			}

			for _, o := range tc.mismatch {
				match, err := filter.Match(o)
				if err != nil {
					t.Fatal(err)
				}

				if match {
					t.Fatalf("'%v' should not match", o)
				}
			}
		})
	}

}
