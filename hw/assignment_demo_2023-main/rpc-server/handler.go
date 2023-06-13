package main

import (
	"context"
	"math/rand"
	"fmt"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) { //sending of messages
	if err := validateSendRequest(req); err != nil {
		return nil, fmt.Errorf("invalid send request: %w", err)
	}
	timestamp := time.Now().Unix()
	message := &Message{
		Message:   req.Message.GetText(),
		Sender:    req.Message.GetSender(),
		Timestamp: timestamp,
	}
	roomID, err := getRoomID(req.Message.GetChat())
	if err != nil{
		return nil, fmt.Errorf("invalid room ID: %w", err)
	}
	err = rdb.SaveMessage(ctx, roomID, message)
	if err != nil{
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success"
	return resp, nil
func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}
	limit := int64(req.GetLimit())
	if limit == 0 {
		limit = 100 // default limit 100
	}
	start := req.GetCursor()
	end := start + limit
	
	messages, err := rdb.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	response := make([]*rpc.Message, 0)
	var counter int64 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages{
		if counter+1>limit{
			hasMore = true
			nextCursor = end
			break
		}
		temp := &rpc.Message{
			Chat: req.GetChat(),
			Text: msg.Message
			Sender: msg.Sender
			SendTime: msg.Timestamp
		}
		response = append(response, temp)
		counter +=1
	}

	resp := rpc.NewPullResponse()
	resp.Messages = response
	resp.Code = 0
	resp.Msg = "success"
	resp.HasMore = &hasMore
	resp.NextCuror = &nextCursor
	
	return resp, nil
}

func validateSendRequest(req *rpc.SendRequest) error{
	senders := string.Split(req.Message.Chat, ":")
	if len(senders) != 2{
		err := fmt.Errorf("wrong number of senders/should be in the format of sender1:sender2") //error interface
		return err
	}

	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("wrong sender")
		return err
	}

	return nil
}

func getRoomID(chat string)(string, error){
	var roomID string
	//get the ID, if cannot get then error
	senders := strings.Split(strings.ToLower(chat), ":")
	if len(senders)!=2{
		err := fmt.Errorf("invalid senders")
		return "", err //error type, hence Errorf
	}
	sender1, sender2 := senders[0], senders[1]
	if comp := strings.Compare(sender1, sender2);comp == 1{
		roomID := fmt.Sprintf("%s:%s", sender2, sender1)
	} else{
		roomID := fmt.Sprintf("%s:%s", sender2, sender1)
	}
	return roomID, nil //nil means no error
}
// func areYouLucky() (int32, string) {
// 	if rand.Int31n(2) == 1 {
// 		return 0, "success"
// 	} else {
// 		return 500, "oops"
// 	}
// }  replace with code for writing and reading messages from server
