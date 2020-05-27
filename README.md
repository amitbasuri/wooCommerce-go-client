# wooCommerce-go-client

![](https://github.com/shoptype/wooCommerce-connector/workflows/Test/badge.svg)

A Golang wrapper for the WooCommerce REST API

### Example Usage
```go
import "github.com/shoptype/wooCommerce-go-client"

func main() {
    client := NewClient("https://example.com/", "ck_xxxxxx", "cs_xxxxxx")
    products, err := client.QueryProducts(
                    url.Values{
                        "page":     []string{"2"},
                        "per_page": []string{"20"}}
                    )
    //...
}
```

### Run tests
```bash
go test ./...
```