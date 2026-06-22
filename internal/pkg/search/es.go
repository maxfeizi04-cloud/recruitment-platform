package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// JobDoc ES 中的职位文档
type JobDoc struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
	City         string   `json:"city"`
	Province     string   `json:"province"`
	SalaryMin    int      `json:"salary_min"`
	SalaryMax    int      `json:"salary_max"`
	CompanyName  string   `json:"company_name"`
	Status       string   `json:"status"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
}

// Client ES 搜索客户端
type Client struct {
	es  *elasticsearch.Client
	idx string
}

// NewClient 创建 ES 客户端
func NewClient(addrs []string, indexName string) (*Client, error) {
	cfg := elasticsearch.Config{Addresses: addrs}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("es client: %w", err)
	}
	// Ping
	_, err = es.Ping()
	if err != nil {
		return nil, fmt.Errorf("es ping: %w", err)
	}
	// Create index if not exists
	if err := ensureIndex(es, indexName); err != nil {
		return nil, err
	}
	return &Client{es: es, idx: indexName}, nil
}

func ensureIndex(es *elasticsearch.Client, name string) error {
	res, err := es.Indices.Exists([]string{name})
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		mapping := `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 0,
				"analysis": {
					"analyzer": {
						"zh_analyzer": {
							"type": "standard"
						}
					}
				}
			},
			"mappings": {
				"properties": {
					"title":       {"type": "text", "analyzer": "zh_analyzer", "boost": 3.0},
					"description": {"type": "text", "analyzer": "zh_analyzer", "boost": 1.0},
					"skills":      {"type": "keyword"},
					"city":        {"type": "keyword"},
					"province":    {"type": "keyword"},
					"salary_min":  {"type": "integer"},
					"salary_max":  {"type": "integer"},
					"company_name":{"type": "text", "boost": 2.0},
					"status":      {"type": "keyword"}
				}
			}
		}`
		req := esapi.IndicesCreateRequest{Index: name, Body: strings.NewReader(mapping)}
		res, err := req.Do(context.Background(), es)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.IsError() {
			return fmt.Errorf("create index: %s", res.String())
		}
	}
	return nil
}

// IndexJob 索引或更新一个职位
func (c *Client) IndexJob(ctx context.Context, doc JobDoc) error {
	body, _ := json.Marshal(doc)
	req := esapi.IndexRequest{
		Index:      c.idx,
		DocumentID: doc.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("index job: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("index job: %s", res.String())
	}
	return nil
}

// DeleteJob 从索引中删除职位
func (c *Client) DeleteJob(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{Index: c.idx, DocumentID: id}
	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("delete job: %w", err)
	}
	defer res.Body.Close()
	return nil
}

// SearchJobs 搜索职位，返回 ID 列表（按相关性排序）
func (c *Client) SearchJobs(ctx context.Context, query, city string, limit, offset int) ([]SearchResult, int, error) {
	var must []map[string]interface{}

	if query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^3", "description^1", "company_name^2", "skills"},
			},
		})
	}

	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{"status": "active"},
	})

	if city != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{"city": city},
		})
	}

	body := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"must": must},
		},
		"from": offset,
		"size": limit,
	}

	buf, _ := json.Marshal(body)
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.idx),
		c.es.Search.WithBody(bytes.NewReader(buf)),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Total struct{ Value int } `json:"total"`
			Hits  []struct {
				ID     string  `json:"_id"`
				Score  float64 `json:"_score"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("decode: %w", err)
	}

	var ids []SearchResult
	for _, h := range result.Hits.Hits {
		ids = append(ids, SearchResult{ID: h.ID, Score: h.Score})
	}
	return ids, result.Hits.Total.Value, nil
}
