package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/3n0ugh/grpc-test-sample/pb"
	"golang.org/x/net/context"
)

func runGetContact(client pb.TelephoneClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.GetContact(ctx, &pb.GetContactRequest{
		Number: "11111111111",
	})
	if err != nil {
		log.Printf("Error getting contact: %v", err)
		return
	}

	fmt.Printf("Got contact: %+v\n", res)
}

func runListContacts(client pb.TelephoneClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.ListContacts(ctx, &pb.ListContactsRequest{})
	if err != nil {
		log.Printf("Error getting contacts: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("All contacts delivered")
				return
			}
			log.Printf("Error getting contacts: %v", err)
			return
		}
		fmt.Printf("Got contact: %+v\n", res)
		time.Sleep(time.Second)
	}
}

func runRecordCallHistory(client pb.TelephoneClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	var wait = make(chan struct{})

	stream, err := client.RecordCallHistory(ctx)
	if err != nil {
		log.Printf("Error involing RecordCallHistory function: %v", err)
		return
	}

	requests := []*pb.RecordCallHistoryRequest{
		{Number: "11111111111"},
		{Number: "1111111111"},
		{Number: "77777777777"},
		{Number: "22222222222"},
		{Number: "33333333333"},
	}

	for _, r := range requests {
		time.Sleep(time.Second)
		err := stream.Send(r)
		if err != nil {
			log.Printf("Error sending call: %v", err)
		}
	}

	go func() {
		for {
			in, err := stream.CloseAndRecv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					close(wait)
					return
				}
				log.Printf("Error getting call history: %v", err)
				return
			}
			fmt.Printf("Got call history: %+v\n", in)
		}
	}()

	<-wait
	stream.CloseSend()
}

func runSendMessage(client pb.TelephoneClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 3030*time.Second)
	defer cancel()

	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Error making call: %v", err)
		return
	}

	var wait = make(chan struct{})

	go func() {
		for {
			r, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					close(wait)
					return
				}
				log.Printf("Error at get call recieved: %v", err)
				return
			}

			fmt.Printf("Got call response: %+v\n", r)
		}
	}()

	request1 := []*pb.SendMessageRequest{
		{
			Msg: []byte("Hi!"),
		},
		{
			Msg: []byte("How are you?"),
		},
		{
			Msg: []byte("."),
		},
	}

	request2 := []*pb.SendMessageRequest{
		{
			Msg: []byte("Thanks"),
		},
		{
			Msg: []byte("See you later"),
		},
		{
			Msg: []byte("."),
		},
	}

	for _, r := range request1 {
		time.Sleep(time.Second)
		err := stream.Send(r)
		if err != nil {
			log.Printf("Error sending request: %v", err)
		}
	}

	time.Sleep(3 * time.Second)

	for _, r := range request2 {
		time.Sleep(time.Second)
		err := stream.Send(r)
		if err != nil {
			log.Printf("Error sending request: %v", err)
		}
	}
	<-wait
	stream.CloseSend()
}
