package service

import (
    "context"
    pb "github.com/rustagram/template-service/genproto"
    "reflect"
    "testing"
)

func TestTaskService_Create(t *testing.T) {
    tests := []struct{
        name string
        input pb.Task
        want pb.Task
    } {
        {
            name: "testOne",
            input: pb.Task{
                Assignee: "aaa",
				Title: "bbb",
				Summary: "ccc",
				Deadline: "2021-01-01",
				Status: "active",
            },
            want: pb.Task{
                Assignee: "aaa",
				Title: "bbb",
				Summary: "ccc",
				Deadline: "2021-01-01",
				Status: "active",
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, funcTestTaskService_Create(t *testing.T) {
            got, err := client.Create(context.Background(), &tc.input)
            if err != nil {
                t.Error("failed to create task", err)
            }
            got.Id = ""
            if !reflect.DeepEqual(tc.want, *got) {
                t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
            }
        })
    }
}
