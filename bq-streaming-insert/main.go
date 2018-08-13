package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
)

type T struct {
	Name string
}

func main() {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "vg-adingo-fluct-dlv")
	if err != nil {
		log.Fatal(err)
	}

	jstr := "{\"Name\": \"aaa\"}"
	var t interface{}

	if err := json.Unmarshal([]byte(jstr), &t); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", t)

	schema, err := bigquery.InferSchema(t)
	if err != nil {
		log.Fatal(err)
	}
	table := client.Dataset("sandbox_saxsir").Table("moge")
	if _, err := table.Metadata(ctx); err != nil {
		log.Println("create table")

		if err := table.Create(ctx, schema); err != nil {
			log.Fatal(err)
		}
	}

	// items := []T{
	// 	T{Name: "aaa"},
	// 	T{Name: "bbb"},
	// 	T{Name: "ccc"},
	// }
	// items := []*bigquery.StructSaver{
	// 	{Struct: T{Name: "aaa"}, Schema: schema, InsertID: "id1"},
	// 	{Struct: T{Name: "bbb"}, Schema: schema, InsertID: "id2"},
	// 	{Struct: T{Name: "ccc"}, Schema: schema, InsertID: "id3"},
	// }

	// u := client.Dataset("sandbox_saxsir").Table("hoge").Uploader()
	// if err := u.Put(ctx, items); err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("success")
}
