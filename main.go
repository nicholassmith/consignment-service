package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro"

	pb "github.com/nicholassmith/consignment-service/proto/consignment"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

//Repository is a faux data store
type Repository struct {
	consigments []*pb.Consignment
}

//Create adds a new assignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consigments, consignment)
	repo.consigments = updated
	return consignment, nil
}

//GetAll returns all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consigments
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(cts context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment

	return nil
}

func (s *service) GetConsignments(cts context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("consignment.service"),
	)

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
