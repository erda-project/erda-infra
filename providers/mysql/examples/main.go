// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/mysql"
	"github.com/jinzhu/gorm"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"mysql"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	DB    *gorm.DB        // autowired
	MySQL mysql.Interface // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	fmt.Println(p.DB)
	fmt.Println(p.MySQL)
	// do something
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
