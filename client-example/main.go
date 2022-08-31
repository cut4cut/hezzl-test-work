package main

import (
	"context"
	"log"
	"time"

	pb "github.com/cut4cut/hezzl-test-work/internal/controller/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr  string   = "localhost:50051"
	names []string = []string{"Bob", "Anna", "John", "Ivan", "Lee", "Zora"}
)

func main() {

	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect gRPC, error: %v", err)
	}
	defer conn.Close()
	c := pb.NewServiceUserClient(conn)

	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Check RPC method create and redis cashe
	log.Print("1. Check RPC method Create() and redis cashe")
	for idx, name := range names {
		// create method
		user, err := c.Create(ctx, &pb.UserName{Name: name})
		if err != nil {
			log.Fatalf("could not create username=%s, error : %v", name, err)
		}
		log.Printf("  %v. successfully created user with ID: %d", idx, user.GetId())

		// getList method and redis cache
		userList, err := c.GetList(ctx, &pb.Pagination{Page: 0, DescName: false, DescCreated: false})
		if err != nil {
			log.Fatalf("could not got users list, error: %v", err)
		}

		if len(userList.Users) == 1 {
			log.Printf("  %v. successfully cashed user list with correct len: %d", idx, len(userList.Users))
		} else {
			log.Fatalf("incorrect len of cashed user list")
		}

		if idx == len(names)-1 {
			log.Printf("Cashed user list: %v", userList.Users)
		}
	}

	// Check RPC method getList
	log.Print("2. Check RPC method GetList()")
	userList, err := c.GetList(ctx, &pb.Pagination{Page: 0, DescName: true, DescCreated: false})
	if err != nil {
		log.Fatalf("could not get users list, error: %v", err)
	}

	if userList.Users[0].GetName() == "Zora" && len(userList.Users) == 5 {
		log.Printf("  correct behavior: first Zora name")
	} else {
		log.Fatalf("incorrect behavior")
	}

	userList, err = c.GetList(ctx, &pb.Pagination{Page: 1, DescName: true, DescCreated: false})
	if err != nil {
		log.Fatalf("could not get users list, error: %v", err)
	}

	if userList.Users[0].GetName() == "Anna" && len(userList.Users) == 1 {
		log.Printf("  correct behavior: first Anna name")
	} else {
		log.Fatalf("incorrect behavior")
	}

	// Check RPC method Delete
	log.Print("3. Check RPC method Delete()")
	for idx, name := range names {
		deletedUser, err := c.Delete(ctx, &pb.UserId{Id: int64(idx + 1)})
		if err != nil {
			log.Fatalf("could not delete username=%s, error : %v", name, err)
		}
		log.Printf("  %v. successfully deleted user with ID: %d", idx, deletedUser.GetId())
	}

	userList, err = c.GetList(ctx, &pb.Pagination{Page: 0, DescName: true, DescCreated: true})
	if err != nil {
		log.Fatalf("could not get users list, error: %v", err)
	}

	if len(userList.Users) == 0 {
		log.Printf("  correct behavior: empty user list")
	} else {
		log.Fatalf("incorrect behavior")
	}
}
