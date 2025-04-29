package reflectlib

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

// SampleRecord  is a sample struct used to test reflection.
type SampleRecord struct {
	BoolField   bool
	FloatField  float64
	IntField    int
	StringField string
}

func TestParseStructFields(t *testing.T) {
	// expected metadata for SampleRecord
	want := []FieldMetadata{
		{Kind: reflect.Bool, Index: 0, Name: "BoolField"},
		{Kind: reflect.Float64, Index: 1, Name: "FloatField"},
		{Kind: reflect.Int, Index: 2, Name: "IntField"},
		{Kind: reflect.String, Index: 3, Name: "StringField"},
	}

	tests := []struct {
		name    string
		input   any
		want    []FieldMetadata
		wantErr bool
	}{
		{
			name:    "value type",
			input:   SampleRecord{},
			want:    want,
			wantErr: false,
		},
		{
			name:    "pointer type",
			input:   &SampleRecord{},
			want:    want,
			wantErr: false,
		},
		{
			name:    "reflect.Type",
			input:   reflect.TypeOf(SampleRecord{}),
			want:    want,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   3.14,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStructFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStructFields(%T) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseStructFields(%T) mismatch (-want +got):\n%s", tt.input, diff)
			}
		})
	}
}
