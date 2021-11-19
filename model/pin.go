package model

import (
	"context"
	"log"

	"github.com/qiniu/qmgo"
)

type db struct {
	*qmgo.Database
	cli *qmgo.Client
}

var store *db

func OpenDb(addr string){
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: addr})
	if err != nil {
		log.Fatal(err)
	}
	d := client.Database("pin")
	store = &db{d, client}
}

func Add(){

}

func Get(){

}

func List(){

}

func Update(){

}

func Delete(){

}


func CloseDb() {
	if store != nil {
		err := store.cli.Close(context.Background())
		log.Println(err)
	}
}