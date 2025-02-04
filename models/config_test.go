package models

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func ptr[T any](v T) *T {
	return &v
}

func TestConfig_AsMap(t *testing.T) {
	id := uuid.New()
	fmt.Println(id.String())
	tests := []struct {
		name         string
		canary       ConfigItem
		removeFields []string
		want         map[string]any
	}{
		{
			name: "remove single field",
			canary: ConfigItem{
				ID:        id,
				Namespace: ptr("canary"),
				Name:      ptr("dummy-canary"),
			},
			removeFields: []string{"updated_at", "created_at", "config_class"},
			want: map[string]any{
				"name":      "dummy-canary",
				"namespace": "canary",
				"agent_id":  "00000000-0000-0000-0000-000000000000",
				"id":        id.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.canary.AsMap(tt.removeFields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Canary.AsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
