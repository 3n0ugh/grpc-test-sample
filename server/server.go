package main

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/3n0ugh/grpc-test-sample/pb"
)

type telephoneServer struct {
	pb.UnimplementedTelephoneServer
	contacts []*pb.ListContactsReply
	calls    []string
}

func NewServer() pb.TelephoneServer {
	t := &telephoneServer{}
	t.loadContact()
	return t
}

func (t *telephoneServer) GetContact(ctx context.Context, telNum *pb.GetContactRequest) (*pb.GetContactReply, error) {
	_, err := strconv.ParseUint(telNum.Number, 10, 64)

	if len(telNum.Number) != 11 || telNum.Number == "" || err != nil {
		return &pb.GetContactReply{}, errors.New("invalid number")
	}

	for _, c := range t.contacts {
		if telNum.Number == c.Number {
			return &pb.GetContactReply{
				Name:     c.Name,
				Lastname: c.Lastname,
				Number:   c.Number,
			}, nil
		}
	}
	return &pb.GetContactReply{}, errors.New("no contact found")
}

func (t *telephoneServer) ListContacts(_ *pb.ListContactsRequest, stream pb.Telephone_ListContactsServer) error {
	for _, c := range t.contacts {
		if err := stream.Send(&pb.ListContactsReply{
			Name:     c.Name,
			Lastname: c.Lastname,
			Number:   c.Number,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (t *telephoneServer) RecordCallHistory(stream pb.Telephone_RecordCallHistoryServer) error {
	var callCount int32

	start := time.Now()

	for {
		_, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&pb.RecordCallHistoryReply{
					CallCount:   callCount,
					ElapsedTime: int32(time.Since(start)),
				})
			}
			return err
		}
		callCount++
	}
}

func (t *telephoneServer) SendMessage(stream pb.Telephone_SendMessageServer) error {
	for {
		rec, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		//log.Printf("Received: %+v", rec)

		t.calls = append(t.calls, string(rec.Msg))

		if t.calls[len(t.calls)-1] == "." {
			for _, m := range t.calls {
				//time.Sleep(time.Second)
				switch m {
				case ".":
				case "Hi!":
					stream.Send(&pb.SendMessageReply{
						Msg: []byte("Hello!"),
					})
				case "How are you?":
					stream.Send(&pb.SendMessageReply{
						Msg: []byte("Fine, you?"),
					})
				case "See you later":
					stream.Send(&pb.SendMessageReply{
						Msg: []byte("See you!"),
					})
				default:
					stream.Send(&pb.SendMessageReply{
						Msg: []byte("Sorry, I don't understand :/"),
					})
				}
			}
			t.calls = t.calls[:0]
		}

	}
}

func (t *telephoneServer) loadContact() {
	t.contacts = []*pb.ListContactsReply{
		{
			Name:     "Nukhet",
			Lastname: "Duru",
			Number:   "11111111111",
		},
		{
			Name:     "Zeki",
			Lastname: "Muren",
			Number:   "22222222222",
		},
		{
			Name:     "Sebnem",
			Lastname: "Ferah",
			Number:   "33333333333",
		},
	}
}
