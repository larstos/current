package current

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestBatchDo(t *testing.T) {
	param := []int32{1, 2, 3, 4, 5, 6, 7}
	dunc := func(ctx context.Context, input int32) (int32, error) {
		if input < 5 {
			time.Sleep(time.Millisecond * time.Duration(input))
			return input + 1, nil
		}
		if input < 7 {
			if input == 6 {
				time.Sleep(time.Millisecond * time.Duration(input))
			}
			return 0, fmt.Errorf("failed")
		} else {
			panic("value equal 7")
		}
	}
	res := BatchDo(context.Background(), param, dunc, WithSyncLimit(2))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	res = BatchDo(context.Background(), param, dunc, WithSyncLimit(2), WithTimeout(3*time.Millisecond))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	res = BatchDo(context.Background(), param, dunc, WithSyncLimit(2), WithTimeout(3*time.Millisecond), WithFastFails(true))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
}
