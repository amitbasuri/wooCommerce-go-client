# wooCommerce-connector

![](https://github.com/shoptype/wooCommerce-connector/workflows/Test/badge.svg)

wooCommerce go client

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
}
```

### Run tests
```bash
go test ./...
```