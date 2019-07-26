package main

import (
	"context"
	"fmt"
	"math/rand"

	proto "github.com/mdevilliers/cache-service/proto/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func registerClientCommand(root *cobra.Command) {

	var address string

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client to exercise the service",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return errors.Wrap(err, "error connecting to endpoint")
			}

			// nolint: errcheck
			defer conn.Close()

			client := proto.NewCacheClient(conn)

			ctx := context.Background()

			fmt.Println("happy path")

			setResponse, err := client.Set(ctx, &proto.SetRequest{
				Key:      "foo",
				Contents: "hello",
				Ttl:      10000,
			})

			if err != nil {
				return errors.Wrap(err, "error caching item")
			}

			fmt.Println("setting value : response", setResponse)

			getResponse, err := client.GetByKey(ctx, &proto.GetByKeyRequest{
				Key: "foo",
			})

			if err != nil {
				return errors.Wrap(err, "error getting item")
			}

			fmt.Println("getting value : response", getResponse)

			purgeResponse, err := client.Purge(ctx, &proto.PurgeRequest{
				Key: "foo",
			})

			if err != nil {
				return errors.Wrap(err, "error getting item")
			}

			fmt.Println("purge value by key : response", purgeResponse)

			for i := 0; i < 100; i++ {

				// nolint: govet
				_, err := client.Set(ctx, &proto.SetRequest{
					Key:      fmt.Sprintf("foo:%d", rand.Intn(10000)),
					Contents: "hello",
				})

				if err != nil {
					return errors.Wrap(err, "error caching item")
				}

			}

			randomNResponse, err := client.GetRandomN(ctx, &proto.GetRandomNRequest{
				Count: 10,
			})

			if err != nil {
				return errors.Wrap(err, "error getting lastn items")
			}

			for n, i := range randomNResponse.GetKeys() {
				fmt.Println(fmt.Sprintf("%d : %s", n, i))
			}

			fmt.Println("sad path")

			setResponse, err = client.Set(ctx, &proto.SetRequest{
				Contents: "hello",
				Ttl:      10000,
			})

			if err != nil {
				return errors.Wrap(err, "error caching item")
			}

			fmt.Println("no key set : ", setResponse)

			setResponse, err = client.Set(ctx, &proto.SetRequest{
				Key: "hello",
			})

			if err != nil {
				return errors.Wrap(err, "error caching item")
			}

			fmt.Println("no content set : ", setResponse)

			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", address, "address of the service to call")

	root.AddCommand(cmd)
}
