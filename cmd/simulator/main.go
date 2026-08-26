package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type client struct{ base, token string }

func (c client) call(method, path string, input, output any) (int, error) {
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return 0, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if output != nil {
		if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
			return resp.StatusCode, err
		}
	} else {
		_, _ = io.Copy(io.Discard, resp.Body)
	}
	return resp.StatusCode, nil
}
func require(name string, got, want int) {
	if got != want {
		fmt.Printf("%s failed: %d\n", name, got)
		os.Exit(1)
	}
	fmt.Printf("%s: %d\n", name, got)
}
func main() {
	base := flag.String("base", "http://127.0.0.1:28090", "base")
	flag.Parse()
	c := client{base: strings.TrimSuffix(*base, "/"), token: "dev-token"}
	var g map[string]any
	status, err := c.call("POST", "/v1/graphs", map[string]any{"name": "dependency"}, &g)
	if err != nil {
		panic(err)
	}
	require("graph", status, 201)
	id := g["id"].(string)
	status, err = c.call("POST", "/v1/graphs/"+id+"/publish", nil, &g)
	if err != nil {
		panic(err)
	}
	require("publish", status, 200)
	var schema map[string]any
	status, err = c.call("POST", "/v1/graphs/"+id+"/schema", map[string]any{"vertices": []any{map[string]any{"name": "service"}}}, &schema)
	if err != nil {
		panic(err)
	}
	require("schema", status, 200)
	for _, v := range []string{"a", "b", "c"} {
		status, err = c.call("POST", "/v1/vertices", map[string]any{"graph_id": id, "id": v, "type": "service", "properties": map[string]any{"name": v}}, &schema)
		if err != nil {
			panic(err)
		}
		require("vertex", status, 201)
	}
	for n, e := range []map[string]string{{"id": "ab", "from": "a", "to": "b"}, {"id": "bc", "from": "b", "to": "c"}} {
		status, err = c.call("POST", "/v1/edges", map[string]any{"graph_id": id, "id": e["id"], "type": "depends", "from_id": e["from"], "to_id": e["to"]}, &schema)
		if err != nil {
			panic(err)
		}
		require(fmt.Sprintf("edge-%d", n), status, 201)
	}
	status, err = c.call("POST", "/v1/queries", map[string]any{"graph_id": id, "start_vertex": "a", "depth": 3}, &schema)
	if err != nil {
		panic(err)
	}
	require("query", status, 200)
	status, err = c.call("POST", "/v1/algorithms", map[string]any{"graph_id": id, "name": "bfs", "parameters": map[string]any{"start": "a"}}, &schema)
	if err != nil {
		panic(err)
	}
	require("algorithm", status, 202)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("simulation completed")
}
