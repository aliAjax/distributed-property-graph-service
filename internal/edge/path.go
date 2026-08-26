package edge

import (
	"context"
	"sort"
)

func SimplePaths(ctx context.Context, repo Repository, g, start, target string, maxDepth int) [][]string {
	if maxDepth < 1 {
		maxDepth = 1
	}
	paths := [][]string{}
	var walk func([]string)
	walk = func(path []string) {
		select {
		case <-ctx.Done():
			return
		default:
		}
		current := path[len(path)-1]
		if current == target {
			paths = append(paths, append([]string(nil), path...))
			return
		}
		if len(path)-1 >= maxDepth {
			return
		}
		for _, e := range repo.ListFrom(ctx, g, current) {
			seen := false
			for _, id := range path {
				if id == e.ToID {
					seen = true
				}
			}
			if !seen {
				walk(append(append([]string(nil), path...), e.ToID))
			}
		}
	}
	walk([]string{start})
	sort.Slice(paths, func(i, j int) bool { return len(paths[i]) < len(paths[j]) })
	return paths
}
