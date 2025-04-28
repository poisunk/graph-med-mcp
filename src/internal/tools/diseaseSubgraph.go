package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"graph-med-mcp/src/internal/config"
	"graph-med-mcp/src/internal/types"
	"log"
)

type DiseaseSubgraphTool struct {
	c config.Config
}

func NewDiseaseSubgraphTool(c config.Config) *DiseaseSubgraphTool {
	return &DiseaseSubgraphTool{
		c: c,
	}
}

func (t *DiseaseSubgraphTool) Name() string {
	return "get_disease_subgraph"
}

func (t *DiseaseSubgraphTool) Definition() mcp.Tool {
	return mcp.NewTool(t.Name(),
		mcp.WithDescription("当你想查询指定疾病的相关知识图谱时，输入疾病名称，返回该疾病的知识图谱。"),
		mcp.WithString(
			"disease_name",
			mcp.Required(),
			mcp.Description("疾病名称"),
		),
	)
}

func (t *DiseaseSubgraphTool) ToolHandlerFunc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["disease_name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	// 链接neo4j
	var neo4jDB neo4j.DriverWithContext

	neo4jDB, err := neo4j.NewDriverWithContext(t.c.Neo4j.Addr, neo4j.BasicAuth(t.c.Neo4j.Username, t.c.Neo4j.Password, ""))
	if err != nil {
		return nil, err
	}

	err = neo4jDB.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	// 查询知识图谱
	subgraph, err := t.getSubgraph(neo4jDB, "疾病", name, 100)
	if err != nil {
		return nil, err
	}

	if subgraph == nil {
		return mcp.NewToolResultText(fmt.Sprintf("没有找到疾病 %s 的相关知识图谱", name)), nil
	}

	result, err := json.Marshal(subgraph)
	if err != nil {
		return nil, err
	}

	//result := ""
	//idToName := make(map[string]string)
	//
	//for _, node := range subgraph.Nodes {
	//	idToName[node.ID] = node.Name
	//}
	//
	//for _, edge := range subgraph.Edges {
	//	result += fmt.Sprintf("%s ==%s==> %s\n\n", idToName[edge.Target], edge.Type, idToName[edge.Source])
	//}

	return mcp.NewToolResultText(string(result)), nil
}

func (t *DiseaseSubgraphTool) getSubgraph(neo4jDB neo4j.DriverWithContext, label string, name string, limit int) (*types.KGGraph, error) {
	nodeList, relaList, err := t.getNodeSubGraph(neo4jDB, label, name, limit)
	if err != nil {
		return nil, err
	}

	subgraph := &types.KGGraph{
		Nodes: make([]types.KGNode, len(nodeList)),
		Edges: make([]types.KGEdge, len(relaList)),
	}

	for i, node := range nodeList {
		subgraph.Nodes[i] = types.KGNode{
			ID:       node.ElementId,
			Label:    node.Labels[0],
			Name:     node.Props["name"].(string),
			NodeType: node.Labels[0],
		}
	}

	for i, rela := range relaList {
		subgraph.Edges[i] = types.KGEdge{
			ID:     rela.ElementId,
			Target: rela.StartElementId,
			Source: rela.EndElementId,
			Type:   rela.Type,
			Label:  rela.Type,
		}
	}

	return subgraph, nil
}

func (t *DiseaseSubgraphTool) getNodeSubGraph(neo4jDB neo4j.DriverWithContext, label, name string, limit int) ([]neo4j.Node, []neo4j.Relationship, error) {
	cypher := "MATCH path=(s:`%s` {name: '%s'})-[*..1]-(connected) " +
		"RETURN nodes(path) AS nodes, relationships(path) AS relationships LIMIT %d"

	cypher = fmt.Sprintf(cypher, label, name, limit)

	nodeIdentitySet := make(map[string]struct{})
	var nodeList []neo4j.Node
	var relaList []neo4j.Relationship
	ctx := context.Background()
	session := neo4jDB.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			if value, ok := record.Get("nodes"); ok {
				list := value.([]interface{})
				for _, ele := range list {
					node := ele.(neo4j.Node)
					if _, ok := nodeIdentitySet[node.ElementId]; !ok {
						nodeIdentitySet[node.ElementId] = struct{}{}
						nodeList = append(nodeList, node)
					}
				}
			}

			if value, ok := record.Get("relationships"); ok {
				list := value.([]interface{})
				for _, ele := range list {
					relationship := ele.(neo4j.Relationship)
					if _, ok := nodeIdentitySet[relationship.ElementId]; !ok {
						nodeIdentitySet[relationship.ElementId] = struct{}{}
						relaList = append(relaList, relationship)
					}
				}
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return nodeList, relaList, err
}
