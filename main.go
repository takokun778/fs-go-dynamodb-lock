package main

import (
	"context"
	"lock/db"
	_ "lock/ddl"
	"lock/gateway"
	"lock/model"
	"log"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	db, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	g := gateway.NewGateway(db)

	src := model.Model{
		ID:   uuid.NewString(),
		Name: "model",
	}

	if err := g.Save(ctx, src); err != nil {
		log.Fatal(err.Error())
	}

	ctx = gateway.SetRIDCtx(ctx)

	m, err := g.Find(ctx, src.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("src: %+v\n", src)

	m.Name = "updated"

	if err := g.Update(ctx, m); err != nil {
		log.Fatal(err.Error())
	}

	dst, err := g.Find(ctx, m.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("dst: %+v\n", dst)
}
