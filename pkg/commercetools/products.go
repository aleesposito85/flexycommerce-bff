package commercetools

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2/clientcredentials"
)

type QueryRequest struct {
	Products struct {
		Results []struct {
			Id         string
			MasterData struct {
				Current struct {
					Name          string `graphql:"name(locale: $lang)"`
					MasterVariant struct {
						Images []struct {
							Url string
						}
						Price struct {
							Value struct {
								CentAmount int64
							}
						} `graphql:"price(currency: $curr)"`
					}
				}
			}
		}
	}
}

type Locale string
type Currency string

func GetProducts() QueryRequest {

	query := QueryRequest{}
	ctx := context.Background()
	conf := &clientcredentials.Config{
		ClientID:     "5CBbLCVsJQCTAuTQlsvpRYD-",
		ClientSecret: "fGHePZjoSTVWTsUILmymmgRd-YS21m8T",
		Scopes:       []string{"manage_project:flexy-commerce"},
		TokenURL:     "https://auth.sphere.io/oauth/token",
	}

	httpClient := conf.Client(ctx)

	client := graphql.NewClient("https://api.sphere.io/flexy-commerce/graphql", httpClient)

	variables := map[string]interface{}{
		"lang": Locale("EN"),
		"curr": Currency("USD"),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println("the error is: ", err)
	}

	return query
}
