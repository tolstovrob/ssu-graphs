#!/usr/bin/env python3
import json
import random
from datetime import datetime

def generate_large_graph():
    """Generate a graph with 10000 nodes and 20000 edges"""

    graph = {
        "nodes": {},
        "edges": {},
        "adjacencyMap": {},
        "options": {
            "isMulti": True,
            "IsDirected": True
        }
    }

    print("Generating 10000 nodes...")
    for i in range(1, 10001):
        graph["nodes"][str(i)] = {
            "key": i,
            "label": f"Node_{i}_Label"
        }
        graph["adjacencyMap"][str(i)] = []

    print("Generating 20000 edges...")
    edge_id = 1

    edge_id = create_connected_components(graph, edge_id)

    edge_id = create_multi_edges(graph, edge_id)

    edge_id = create_random_edges(graph, edge_id)

    rebuild_adjacency_map(graph)

    return graph

def create_connected_components(graph, start_edge_id):
    edge_id = start_edge_id
    component_sizes = [2000, 1500, 1000, 800, 500, 300, 200, 100, 50, 25, 10]

    node_offset = 1
    for size in component_sizes:
        if node_offset + size > 10000:
            break

        for i in range(size - 1):
            if edge_id > 20000:
                return edge_id

            src = node_offset + i
            dst = node_offset + i + 1

            graph["edges"][str(edge_id)] = {
                "key": edge_id,
                "source": src,
                "destination": dst,
                "weight": random.randint(1, 100),
                "label": f"edge_{edge_id}_type{random.choice(['A','B','C'])}"
            }
            edge_id += 1

        extra_edges = size // 4
        for _ in range(extra_edges):
            if edge_id > 20000:
                return edge_id

            src = random.randint(node_offset, node_offset + size - 1)
            dst = random.randint(node_offset, node_offset + size - 1)

            if src != dst:
                graph["edges"][str(edge_id)] = {
                    "key": edge_id,
                    "source": src,
                    "destination": dst,
                    "weight": random.randint(1, 50),
                    "label": f"edge_{edge_id}_internal"
                }
                edge_id += 1

        node_offset += size

    return edge_id

def create_multi_edges(graph, start_edge_id):
    edge_id = start_edge_id

    for _ in range(100):
        if edge_id > 20000:
            return edge_id

        src = random.randint(1, 500)
        dst = random.randint(1, 500)

        if src != dst:
            num_edges = random.randint(2, 5)
            for j in range(num_edges):
                if edge_id > 20000:
                    return edge_id

                graph["edges"][str(edge_id)] = {
                    "key": edge_id,
                    "source": src,
                    "destination": dst,
                    "weight": random.randint(1, 30),
                    "label": f"multi_edge_{src}_{dst}_{j+1}"
                }
                edge_id += 1

    return edge_id

def create_random_edges(graph, start_edge_id):
    """Fill remaining edges with random connections"""
    edge_id = start_edge_id

    while edge_id <= 20000:
        src = random.randint(1, 10000)
        dst = random.randint(1, 10000)

        if src != dst:
            graph["edges"][str(edge_id)] = {
                "key": edge_id,
                "source": src,
                "destination": dst,
                "weight": random.randint(1, 100),
                "label": f"random_edge_{edge_id}"
            }
            edge_id += 1

    return edge_id

def rebuild_adjacency_map(graph):
    """Rebuild adjacency map based on edges"""
    print("Rebuilding adjacency map...")

    for node_id in graph["nodes"]:
        graph["adjacencyMap"][node_id] = []

    for edge in graph["edges"].values():
        src = str(edge["source"])
        dst = str(edge["destination"])

        if dst not in graph["adjacencyMap"][src]:
            graph["adjacencyMap"][src].append(edge["destination"])

        # For undirected graphs, uncomment this:
        # if src not in graph["adjacencyMap"][dst]:
        #     graph["adjacencyMap"][dst].append(edge["source"])

def save_graph_to_file(graph, filename):
    """Save graph to JSON file"""
    print(f"Saving graph to {filename}...")

    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(graph, f, indent=2, ensure_ascii=False)

    print(f"Graph saved successfully!")
    print(f"Statistics:")
    print(f"  - Nodes: {len(graph['nodes'])}")
    print(f"  - Edges: {len(graph['edges'])}")
    print(f"  - Is Directed: {graph['options']['IsDirected']}")
    print(f"  - Is Multi: {graph['options']['isMulti']}")

def main():
    """Main function"""
    print("=== Large Graph Generator ===")
    print("Generating graph with 10000 nodes and 20000 edges...")
    print("This may take a few seconds...")

    start_time = datetime.now()

    graph = generate_large_graph()

    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"large_graph_{timestamp}.json"

    save_graph_to_file(graph, filename)

    end_time = datetime.now()
    duration = (end_time - start_time).total_seconds()

    print(f"\nGeneration completed in {duration:.2f} seconds")
    print(f"File: {filename}")

if __name__ == "__main__":
    main()
