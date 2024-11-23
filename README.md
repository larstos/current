# current

example:
````
    param := []int32{1, 2, 3, 4, 5, 6, 7}
	handler := func(ctx context.Context, input int32) (int32, error) {
			time.Sleep(time.Millisecond * time.Duration(input))
			return input + 1, nil
	}
	res := BatchDo(context.Background(), param, handler, WithSyncLimit(2))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
````