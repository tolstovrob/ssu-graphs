import json
import random


def generate_graph(graph_id, directed, weighted, has_loops=True, has_isolated=True):
    """Генерация графа с заданными параметрами"""

    num_nodes = random.randint(5, 15)
    nodes = {}
    edges = {}
    adjacency_map = {}

    for i in range(1, num_nodes + 1):
        nodes[str(i)] = {"key": i, "label": f"Node_{i}"}
        adjacency_map[str(i)] = []

    edge_key = 1

    connected_nodes = list(range(1, num_nodes + 1))
    if not directed:
        for i in range(1, num_nodes):
            if i + 1 <= num_nodes:
                weight = random.randint(1, 10) if weighted else 0
                edges[str(edge_key)] = {
                    "key": edge_key,
                    "source": i,
                    "destination": i + 1,
                    "weight": weight,
                    "label": f"Edge_{i}_{i + 1}",
                }
                adjacency_map[str(i)].append(i + 1)
                if not directed:
                    adjacency_map[str(i + 1)].append(i)
                edge_key += 1
    else:
        for i in range(1, num_nodes):
            if i + 1 <= num_nodes:
                weight = random.randint(1, 10) if weighted else 0
                edges[str(edge_key)] = {
                    "key": edge_key,
                    "source": i,
                    "destination": i + 1,
                    "weight": weight,
                    "label": f"Edge_{i}_{i + 1}",
                }
                adjacency_map[str(i)].append(i + 1)
                edge_key += 1

    num_additional_edges = random.randint(3, 6)
    for _ in range(num_additional_edges):
        source = random.randint(1, num_nodes)
        destination = random.randint(1, num_nodes)

        edge_exists = False
        for edge in edges.values():
            if edge["source"] == source and edge["destination"] == destination:
                edge_exists = True
                break
            if (
                not directed
                and edge["source"] == destination
                and edge["destination"] == source
            ):
                edge_exists = True
                break

        if not edge_exists:
            weight = random.randint(1, 10) if weighted else 0
            edges[str(edge_key)] = {
                "key": edge_key,
                "source": source,
                "destination": destination,
                "weight": weight,
                "label": f"Edge_{source}_{destination}",
            }
            adjacency_map[str(source)].append(destination)
            if not directed and source != destination:
                adjacency_map[str(destination)].append(source)
            edge_key += 1

    if has_loops:
        num_loops = random.randint(1, 2)
        for _ in range(num_loops):
            node = random.randint(1, num_nodes)
            weight = random.randint(1, 10) if weighted else 0
            edges[str(edge_key)] = {
                "key": edge_key,
                "source": node,
                "destination": node,
                "weight": weight,
                "label": f"Loop_{node}",
            }
            adjacency_map[str(node)].append(node)
            edge_key += 1

    if has_isolated:
        num_isolated = random.randint(1, 2)
        # Создаем новые изолированные вершины
        for i in range(num_isolated):
            node_id = num_nodes + i + 1
            nodes[str(node_id)] = {"key": node_id, "label": f"Isolated_Node_{node_id}"}
            adjacency_map[str(node_id)] = []

    options = {"isMulti": False, "IsDirected": directed}

    graph = {
        "nodes": nodes,
        "edges": edges,
        "adjacencyMap": adjacency_map,
        "options": options,
    }

    return graph


def generate_test_graphs():
    graphs = []

    print("Генерация неориентированного невзвешенного графа...")
    graph1 = generate_graph(1, directed=False, weighted=False)
    graphs.append(("undirected_unweighted.json", graph1))

    print("Генерация ориентированного невзвешенного графа...")
    graph2 = generate_graph(2, directed=True, weighted=False)
    graphs.append(("directed_unweighted.json", graph2))

    print("Генерация неориентированного взвешенного графа...")
    graph3 = generate_graph(3, directed=False, weighted=True)
    graphs.append(("undirected_weighted.json", graph3))

    print("Генерация ориентированного взвешенного графа...")
    graph4 = generate_graph(4, directed=True, weighted=True)
    graphs.append(("directed_weighted.json", graph4))

    print("Генерация ориентированного взвешенного графа с большим количеством ребер...")
    graph5 = generate_graph(5, directed=True, weighted=True)

    return graphs


def save_graphs(graphs):
    for filename, graph in graphs:
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(graph, f, indent=2, ensure_ascii=False)
        print(f"Сохранен: {filename}")

        nodes_count = len(graph["nodes"])
        edges_count = len(graph["edges"])
        directed = graph["options"]["IsDirected"]
        weighted = any(edge["weight"] > 0 for edge in graph["edges"].values())

        print(f"  Узлы: {nodes_count}, Ребра: {edges_count}")
        print(f"  Ориентированный: {directed}, Взвешенный: {weighted}")

        loops = sum(
            1
            for edge in graph["edges"].values()
            if edge["source"] == edge["destination"]
        )
        isolated = sum(
            1
            for node_id, neighbors in graph["adjacencyMap"].items()
            if len(neighbors) == 0
        )

        print(f"  Петли: {loops}, Изолированные вершины: {isolated}")
        print()


if __name__ == "__main__":
    print("Генерация тестовых графов...")
    print("=" * 50)

    graphs = generate_test_graphs()
    save_graphs(graphs)

    print("Successfully generated test graphs")
