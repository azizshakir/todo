package service

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/azizshakir/todo/genproto"
)

func TestTaskService_Create(t *testing.T) {
	tests := []struct {
		name  string
		input pb.Task
		want  pb.Task
	}{
		{
			name: "testOne",
			input: pb.Task{
				Assignee: "aaa",
				Title:    "bbb",
				Summary:  "ccc",
				Deadline: "2021-01-01",
				Status:   "active",
			},
			want: pb.Task{
				Assignee: "aaa",
				Title:    "bbb",
				Summary:  "ccc",
				Deadline: "2021-01-01T00:00:00Z",
				Status:   "active",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := client.Create(context.Background(), &tc.input)
			if err != nil {
				t.Error("failed to create task", err)
			}
			got.Id = "0"
			if !reflect.DeepEqual(tc.want, *got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}
