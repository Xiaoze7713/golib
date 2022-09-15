package es

import (
	"github.com/olivere/elastic/v7"
)

type esLogic struct {
	Cli *elastic.Client
}

var Handler *esLogic

func (m *esLogic) Close() {
	m.Cli.Stop()
}

func Init(urls []string) (err error) {
	client, err := elastic.NewClient(elastic.SetURL(urls...), elastic.SetSniff(false))
	if err != nil {
		return err
	}
	Handler = &esLogic{Cli: client}
	return err
}
