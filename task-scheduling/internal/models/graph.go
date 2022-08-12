package models

import (
	"errors"
)

/*
	通过依赖关系构成图并进行解析二维数组
*/
//key  为Topic:taskID
type nodeset map[string]struct{}

type depmap map[string]nodeset

type Graph struct {

	//节点集合
	nodes nodeset

	//节点边使用 p -> c
	dependencyP2C depmap

	//c -> p
	depenDencyC2P depmap
}

//初始化
func GraphNew() *Graph {
	return &Graph{
		nodes:         make(nodeset),
		dependencyP2C: make(depmap),
		depenDencyC2P: make(depmap),
	}
}

//添加节点
func (g *Graph) DependOn(child, parent string) error {
	if parent == child {
		return errors.New("self-dependence not allowed")
	}

	//判断是否成环
	if g.DependsOn(parent, child) {
		return errors.New("circular dependencies not allowed")
	}

	g.nodes[parent] = struct{}{}
	g.nodes[child] = struct{}{}

	addNodeToNodeset(g.dependencyP2C, parent, child)
	addNodeToNodeset(g.depenDencyC2P, child, parent)
	return nil
}

//没有需要依赖的 C2P中没有的元素
func (g *Graph) Leaves() []string {
	leaves := []string{}

	for nodes := range g.nodes {
		if _, ok := g.depenDencyC2P[nodes]; !ok {
			leaves = append(leaves, nodes)
		}
	}
	return leaves
}

func (g *Graph) TopoSortedLayers() [][]string {
	out := [][]string{}

	copyGraph := g.Clone()

	for len(copyGraph.nodes) > 0 {
		leaves := copyGraph.Leaves()

		out = append(out, leaves)

		for _, leafNode := range leaves {
			copyGraph.remove(leafNode)
		}
	}
	return out

}

func (g *Graph) remove(node string) {
	for p2cNode := range g.dependencyP2C[node] {
		removeFromDepmap(g.depenDencyC2P, p2cNode, node)
	}
	delete(g.dependencyP2C, node)
	delete(g.nodes, node)
}

func removeFromDepmap(dm depmap, key, node string) {
	nodes := dm[key]
	if len(nodes) == 1 {
		//删除c2p中没有需要依赖的元素
		delete(dm, key)
	} else {
		delete(nodes, node)
	}
}

//复制一个新的图
func (g *Graph) Clone() *Graph {
	return &Graph{
		nodes:         copyNodeset(g.nodes),
		dependencyP2C: copyDepmap(g.dependencyP2C),
		depenDencyC2P: copyDepmap(g.depenDencyC2P),
	}
}

//判断依赖child任务中有没有parents
func (g *Graph) DependsOn(child, parent string) bool {
	out := g.dependC2P(child)
	_, ok := out[parent]
	return ok
}

func (g *Graph) dependC2P(child string) nodeset {
	return g.buildTransitive(child, g.dependC2PImp)
}

func (g *Graph) dependC2PImp(node string) nodeset {
	return g.depenDencyC2P[node]
}

func addNodeToNodeset(dep depmap, key, node string) {
	nodes, ok := dep[key]
	if !ok {
		nodes = make(nodeset)
		dep[key] = nodes
	}
	nodes[node] = struct{}{}
}

func (g *Graph) buildTransitive(root string, nextFn func(string) nodeset) nodeset {
	if _, ok := g.nodes[root]; !ok {
		return nil
	}

	out := make(nodeset)

	searchNodes := []string{root}

	for len(searchNodes) > 0 {
		discovery := []string{}
		for _, node := range searchNodes {
			for nextNode := range nextFn(node) {
				if _, ok := out[nextNode]; !ok {
					out[nextNode] = struct{}{}
					discovery = append(discovery, nextNode)
				}
			}
		}
		searchNodes = discovery

	}

	return out
}

func copyNodeset(s nodeset) nodeset {
	newNodes := make(nodeset, len(s))

	for key, value := range s {
		newNodes[key] = value
	}
	return newNodes
}

func copyDepmap(m depmap) depmap {
	newMap := make(depmap, len(m))

	for key, value := range m {
		newMap[key] = copyNodeset(value)
	}

	return newMap
}
