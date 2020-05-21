# wooCommerce-connector
wooCommerce go client for shoptype

![](https://github.com/shoptype/wooCommerce-connector/workflows/Test/badge.svg)

### Run tests
```bash
go test ./...
```

### Example Usage
```go
import "github.com/shoptype/wooCommerce-connector"

func main() {
    client := NewClient("https://example.com/", "ck_xxxxxx", "cs_xxxxxx")
    products, err := client.QueryProducts(
                    url.Values{
                        "page":     []string{"2"},
                        "per_page": []string{"20"}}
                    )
    //...
```