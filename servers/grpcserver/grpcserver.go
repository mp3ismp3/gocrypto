package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/mp3ismp3/gocrypto/controllers"
	pb "github.com/mp3ismp3/gocrypto/proto"
	"google.golang.org/grpc"
)

type server struct {
	*pb.UnimplementedExchangeServer
}

func (s *server) OpenMatching(ctx context.Context, req *pb.OpenMatchingRequest) (*pb.Response, error) {
	openPrice := req.GetOpenPrice()
	symbol := req.GetSymbol()
	status, message := controllers.AddEngine(symbol, openPrice)
	response := &pb.Response{Code: status, Message: message}
	return response, nil

}

func (s *server) CloseMatching(ctx context.Context, req *pb.CloseMatchingRequest) (*pb.Response, error) {
	symbol := req.GetSymbol()
	status, message := controllers.DeleteEngine(symbol)
	response := &pb.Response{Code: status, Message: message}
	return response, nil

}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderLogResponse, error) {
	side := req.GetOrderSide()
	ordertype := req.GetOrderType()
	qty := req.GetQty()
	symbol := req.GetSymbol()
	price := req.GetPrice()
	orderNode := controllers.AddOrder(symbol, int32(*side.Enum()), price, qty, int32(*ordertype.Enum()))
	fmt.Println("新訂單:", orderNode)
	// // orderList, ok := EngineList[orderNode.Symbol]
	// // if !ok {
	// // 	log.Println("沒有這個幣種")
	// // }
	// trade := orderList.BuyOrSell(*orderNode)
	// fmt.Println("交易結果:", trade)
	response := &pb.OrderLogResponse{OrderId: string(strconv.Itoa(int(orderNode.ID)))}
	return response, nil
}

func (s *server) CancelOrder(ctx context.Context, request *pb.CancelOrderRequest) (*pb.OrderLogResponse, error) {
	return nil, nil
}

func Init() {
	fmt.Println("starting gRPC server...")
	s := grpc.NewServer() // grpc.NewServer(ServerOption) 創建一個未啟動服務和尚未接受請求的 gRPC server
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v \n", err)
	}

	pb.RegisterExchangeServer(s, &server{}) // 再調用服務前，呼叫 RegisterService 註冊服務使其實現到gRPC服務器
	if err := s.Serve(lis); err != nil {    //Serve為每一個連接創建 new ServerTransport 和 service goroutine
		log.Fatalf("faild to serve: %v \n", err)
	}
}
