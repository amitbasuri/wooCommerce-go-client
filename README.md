# wooCommerce-connector
wooCommerce go client for shoptype

### Run tests
```bash
go test ./...
```

### Example Usage
```go
client := NewClient("https://example.com/", "ck_xxxxxx", "cs_xxxxxx")
products, err := client.QueryProducts(
                    url.Values{
                        "page":     []string{"2"},
                        "per_page": []string{"20"}}
                    )
```