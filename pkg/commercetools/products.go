package commercetools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"

	"github.com/shurcooL/graphql"
	"github.com/shurcooL/graphql/ident"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2/clientcredentials"
)

type SingleProductResponse struct {
	Data SingleProduct
}

type QueryRequest struct {
	Products struct {
		Results []ProductBasic
	}
}

type ProductBasic struct {
	Id         string
	MasterData struct {
		Current struct {
			Name          string `graphql:"name(locale: $lang)"`
			MasterVariant VariantBasic
		}
	}
}

type ProductFull struct {
	Id         string
	MasterData struct {
		Current struct {
			Name          string `graphql:"name(locale: $lang)"`
			MasterVariant VariantAttributes
		}
	}
}

type SingleProduct struct {
	Product ProductFull `graphql:"product(id: $id)"`
}

type VariantBasic struct {
	Images []struct {
		Url string
	}
	Price struct {
		Value struct {
			CentAmount int64
		}
	} `graphql:"price(currency: $curr)"`
}

type VariantAttributes struct {
	AttributesText []struct {
		Name  string
		Value interface{}
	} `graphql:"attributesText: attributesRaw (includeNames: $attributesText)"`
	AttributesEnum []struct {
		Name  string
		Value struct {
			Key   string
			Label string
		} `skipinner:"value"`
	} `graphql:"attributesEnum: attributesRaw (includeNames: $attributesEnum)"`
	Images []struct {
		Url string
	}
	Price struct {
		Value struct {
			CentAmount int64
		}
	} `graphql:"price(currency: $curr)"`
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

func GetProduct(id string, textAttributes []string, enumAttributes []string) SingleProduct {

	query := SingleProduct{}
	ctx := context.Background()
	conf := &clientcredentials.Config{
		ClientID:     "5CBbLCVsJQCTAuTQlsvpRYD-",
		ClientSecret: "fGHePZjoSTVWTsUILmymmgRd-YS21m8T",
		Scopes:       []string{"manage_project:flexy-commerce"},
		TokenURL:     "https://auth.sphere.io/oauth/token",
	}

	httpClient := conf.Client(ctx)

	//client := graphql.NewClient("https://api.sphere.io/flexy-commerce/graphql", httpClient)

	var attributesText []graphql.String

	for _, s := range textAttributes {
		attributesText = append(attributesText, graphql.String(s))
	}

	var attributesEnum []graphql.String

	for _, s := range enumAttributes {
		attributesEnum = append(attributesEnum, graphql.String(s))
	}

	variables := map[string]interface{}{
		"id":             graphql.String(id),
		"lang":           Locale("EN"),
		"curr":           Currency("USD"),
		"attributesText": attributesText,
		"attributesEnum": attributesEnum,
	}

	out := queryGraphql(context.Background(), &query, variables, httpClient, "https://api.sphere.io/flexy-commerce/graphql")
	b, _ := json.Marshal(out)

	fmt.Println(string(b))

	/*
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			fmt.Println("the error is: ", err)
		}
		fmt.Println("the query is: ", query)
	*/
	var singleProduct SingleProduct
	err := json.Unmarshal(b, &singleProduct)

	if err != nil {
		fmt.Println("the error is: ", err)
	}

	return singleProduct
}

// do executes a single GraphQL operation.
func queryGraphql(ctx context.Context, v interface{}, variables map[string]interface{}, client *http.Client, url string) interface{} {
	var query string

	query = constructQuery(v, variables)

	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}

	resp, err := ctxhttp.Post(ctx, client, url, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	var out struct {
		Data   *json.RawMessage
		Errors *json.RawMessage
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		// TODO: Consider including response body in returned error, if deemed helpful.
		return err
	}

	return out.Data
}

func constructQuery(v interface{}, variables map[string]interface{}) string {
	query := query(v)
	fmt.Println(query)
	if len(variables) > 0 {
		return "query(" + queryArguments(variables) + ")" + query
	}
	return query
}

func query(v interface{}) string {
	var buf bytes.Buffer
	writeQuery(&buf, reflect.TypeOf(v), false)
	return buf.String()
}

func writeQuery(w io.Writer, t reflect.Type, inline bool) {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice:
		writeQuery(w, t.Elem(), false)
	case reflect.Struct:
		// If the type implements json.Unmarshaler, it's a scalar. Don't expand it.
		if reflect.PtrTo(t).Implements(jsonUnmarshaler) {
			return
		}
		if !inline {
			io.WriteString(w, "{")
		}
		for i := 0; i < t.NumField(); i++ {
			skipinner := false
			if i != 0 {
				io.WriteString(w, ",")
			}
			f := t.Field(i)
			value, ok := f.Tag.Lookup("graphql")
			_, skipinner = f.Tag.Lookup("skipinner")
			inlineField := f.Anonymous && !ok
			if !inlineField {
				if ok {
					io.WriteString(w, value)
				} else {
					io.WriteString(w, ident.ParseMixedCaps(f.Name).ToLowerCamelCase())
				}
			}
			if !skipinner {
				writeQuery(w, f.Type, inlineField)
			}
		}
		if !inline {
			io.WriteString(w, "}")
		}
	}
}

func queryArguments(variables map[string]interface{}) string {
	// Sort keys in order to produce deterministic output for testing purposes.
	// TODO: If tests can be made to work with non-deterministic output, then no need to sort.
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		io.WriteString(&buf, "$")
		io.WriteString(&buf, k)
		io.WriteString(&buf, ":")
		writeArgumentType(&buf, reflect.TypeOf(variables[k]), true)
		// Don't insert a comma here.
		// Commas in GraphQL are insignificant, and we want minified output.
		// See https://facebook.github.io/graphql/October2016/#sec-Insignificant-Commas.
	}
	return buf.String()
}

func writeArgumentType(w io.Writer, t reflect.Type, value bool) {
	if t.Kind() == reflect.Ptr {
		// Pointer is an optional type, so no "!" at the end of the pointer's underlying type.
		writeArgumentType(w, t.Elem(), false)
		return
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		// List. E.g., "[Int]".
		io.WriteString(w, "[")
		writeArgumentType(w, t.Elem(), true)
		io.WriteString(w, "]")
	default:
		// Named type. E.g., "Int".
		name := t.Name()
		if name == "string" { // HACK: Workaround for https://github.com/shurcooL/githubv4/issues/12.
			name = "ID"
		}
		io.WriteString(w, name)
	}

	if value {
		// Value is a required type, so add "!" to the end.
		io.WriteString(w, "!")
	}
}

var jsonUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
