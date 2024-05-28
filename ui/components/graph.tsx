"use client";
import React, { useEffect, useRef, useState, useCallback } from 'react';
import * as d3 from 'd3';
import { useTheme } from 'next-themes';
import debounce from 'lodash.debounce';
import axios from "axios";
import { Select, SelectItem, Button } from '@nextui-org/react';

interface Node extends d3.SimulationNodeDatum {
    id: string;
    kind?: string;
    label: string;
    x?: number;
    y?: number;
}

interface Link extends d3.SimulationLinkDatum<Node> {
    source: string | Node;
    target: string | Node;
}

type Subject = {
    kind: string;
    apiGroup: string;
    name: string;
};

type RoleRef = {
    kind: string;
    apiGroup: string;
    name: string;
};

type BindingData = {
    id: number;
    name: string;
    kind: string;
    subjects: Subject[];
    roleRef: RoleRef;
    details?: string;
};

const Tooltip = ({ node, isDarkMode }: { node: Node; isDarkMode: boolean }) => (
    <div style={{
        position: 'absolute',
        left: node.x,
        top: node.y,
        backgroundColor: isDarkMode ? '#333' : 'white',
        color: isDarkMode ? 'white' : 'black',
        padding: '2px 5px',
        border: `1px solid ${isDarkMode ? '#555' : '#ccc'}`,
        borderRadius: '3px',
        pointerEvents: 'none',
        transform: 'translate(-50%, -100%)',
    }}>
        {node.label}
    </div>
);

const DisjointGraph = ({ data }: { data?: { nodes: Node[], links: Link[] } }) => {
    const svgRef = useRef<SVGSVGElement | null>(null);
    const [hoveredNode, setHoveredNode] = useState<Node | null>(null);
    const [selectedNodes, setSelectedNodes] = useState<Set<string>>(new Set());
    const [allNodes, setAllNodes] = useState<Node[]>([]);
    const [allLinks, setAllLinks] = useState<Link[]>([]);
    const [bindingData, setBindingData] = useState<BindingData[]>([]);
    const { theme } = useTheme();
    const isDarkMode = theme === 'dark';

    const processGraphData = (data: BindingData[]) => {
        const nodes: Node[] = [];
        const links: Link[] = [];

        data.forEach(binding => {
            if (!binding.name || !binding.kind || !binding.subjects || !binding.roleRef) {
                console.error('Invalid binding data:', binding);
                return;
            }

            nodes.push({ id: binding.name, kind: binding.kind, label: binding.name });

            binding.subjects.forEach(subject => {
                if (!subject.kind || !subject.apiGroup || !subject.name) {
                    console.error('Invalid subject data:', subject);
                    return;
                }
                const subjectId = `${subject.kind}-${subject.name}`;
                if (!nodes.find(n => n.id === subjectId)) {
                    nodes.push({ id: subjectId, label: `${subject.kind} - ${subject.name}` });
                }
                links.push({ source: binding.name, target: subjectId });
            });

            const roleRefId = `${binding.roleRef.kind}-${binding.roleRef.name}`;
            if (!nodes.find(n => n.id === roleRefId)) {
                nodes.push({ id: roleRefId, label: `${binding.roleRef.kind} - ${binding.roleRef.name}` });
            }
            links.push({ source: binding.name, target: roleRefId });
        });

        return { nodes, links };
    };

    const renderGraph = useCallback((nodes: Node[], links: Link[], selectedIds: Set<string>) => {
        const svg = d3.select(svgRef.current);
        svg.selectAll("*").remove(); // Clear previous elements

        const width = svg.node()?.clientWidth || 800;
        const height = svg.node()?.clientHeight || 600;

        const filteredNodes = new Set<Node>();
        const filteredLinks = new Set<Link>();

        if (selectedIds.size > 0) {
            selectedIds.forEach(id => {
                const mainNode = nodes.find(node => node.id === id);
                if (mainNode) {
                    filteredNodes.add(mainNode);

                    links.forEach(link => {
                        if (link.source === id || link.target === id) {
                            filteredLinks.add(link);
                            filteredNodes.add(nodes.find(node => node.id === link.source) as Node);
                            filteredNodes.add(nodes.find(node => node.id === link.target) as Node);
                        }
                    });
                }
            });
        } else {
            nodes.forEach(node => filteredNodes.add(node));
            links.forEach(link => filteredLinks.add(link));
        }

        const nodesToRender = Array.from(filteredNodes);
        const linksToRender = Array.from(filteredLinks);

        const g = svg.append('g');

        const zoom = d3.zoom<SVGSVGElement, unknown>()
            .scaleExtent([0.1, 2])
            .on('zoom', (event) => {
                g.attr('transform', event.transform.toString());
            });

        svg.call(zoom as any);

        const simulation = d3.forceSimulation<Node>(nodesToRender)
            .force('link', d3.forceLink<Node, Link>(linksToRender).id((d: any) => d.id).distance(100))
            .force('charge', d3.forceManyBody().strength(-100))
            .force('center', d3.forceCenter(width / 2, height / 2));

        const link = g.selectAll('.link')
            .data(linksToRender)
            .enter().append('line')
            .attr('class', 'link')
            .attr('stroke', '#666')
            .attr('stroke-width', 2);

        const node = g.selectAll('.node')
            .data(nodesToRender)
            .enter().append('circle')
            .attr('class', 'node')
            .attr('r', 10)
            .attr('fill', d => d.kind === 'ClusterRoleBinding' ? 'orange' : d.kind === 'RoleBinding' ? 'green' : 'pink')
            .call(drag(simulation) as any)
            .on('mouseover', debounce((_event, d) => setHoveredNode(d), 50))
            .on('mouseout', debounce(() => setHoveredNode(null), 50));

        const text = g.selectAll('.label')
            .data(nodesToRender)
            .enter().append('text')
            .attr('class', 'label')
            .attr('x', d => d.x!)
            .attr('y', d => d.y!)
            .text(d => d.label)
            .style('font-size', '13px')
            .style('fill', isDarkMode ? '#FFF' : '#000')
            .style('text-anchor', '-moz-initial')
            .style('font-weight', 'bold');

        simulation.on('tick', () => {
            link
                .attr('x1', d => (d.source as Node).x!)
                .attr('y1', d => (d.source as Node).y!)
                .attr('x2', d => (d.target as Node).x!)
                .attr('y2', d => (d.target as Node).y!);

            node
                .attr('cx', d => d.x!)
                .attr('cy', d => d.y!);

            text
                .attr('x', d => d.x!)
                .attr('y', d => d.y!);
        });

        // Legend
        const legend = svg.append('g')
            .attr('class', 'legend')
            .attr('transform', 'translate(20,20)');

        const legendData = [
            { label: 'ClusterRoleBinding', color: 'orange' },
            { label: 'RoleBinding', color: 'green' },
            { label: 'Other', color: 'pink' }
        ];

        const legendItem = legend.selectAll('.legend-item')
            .data(legendData)
            .enter().append('g')
            .attr('class', 'legend-item')
            .attr('transform', (_d, i) => `translate(0, ${i * 20})`);

        legendItem.append('rect')
            .attr('width', 18)
            .attr('height', 18)
            .attr('fill', d => d.color);

        legendItem.append('text')
            .attr('x', 24)
            .attr('y', 9)
            .attr('dy', '0.35em')
            .text(d => d.label)
            .style('font-size', '12px')
            .style('fill', isDarkMode ? '#FFF' : '#333');

        return () => {
            simulation.stop();
        };
    }, [isDarkMode]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get('/api/data');
                const data: BindingData[] | null = response.data;
                if (!data) {
                    console.error(new Error('Data is null or undefined'));
                    return;
                }
                const { nodes, links } = processGraphData(data);
                setAllNodes(nodes);
                setAllLinks(links);
                setBindingData(data);
                renderGraph(nodes, links, new Set());
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        if (!data) {
            fetchData().finally(() => { console.log('Data fetching completed'); });
        } else {
            renderGraph(data.nodes, data.links, new Set());
        }
    }, [isDarkMode, renderGraph, data]);

    useEffect(() => {
        const text = d3.selectAll('.label');
        text.style('fill', isDarkMode ? '#FFF' : '#000');

        const legendText = d3.selectAll('.legend-item text');
        legendText.style('fill', isDarkMode ? '#FFF' : '#333');
    }, [isDarkMode]);

    const drag = (simulation: d3.Simulation<Node, undefined>) => {
        const dragStarted = (event: d3.D3DragEvent<SVGCircleElement, Node, Node>, d: Node) => {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
        };

        const dragged = (event: d3.D3DragEvent<SVGCircleElement, Node, Node>, d: Node) => {
            d.fx = event.x;
            d.fy = event.y;
        };

        const dragEnded = (event: d3.D3DragEvent<SVGCircleElement, Node, Node>, d: Node) => {
            if (!event.active) simulation.alphaTarget(0);
            d.fx = null;
            d.fy = null;
        };

        return d3.drag<SVGCircleElement, Node>()
            .on('start', dragStarted)
            .on('drag', dragged)
            .on('end', dragEnded);
    };

    const handleSelectionChange = (keys: Set<React.Key> | any) => {
        const newSelectedNodes = new Set(Array.from(keys) as string[]);
        setSelectedNodes(newSelectedNodes);

        const selectedData = bindingData.filter(binding => newSelectedNodes.has(binding.name));
        const { nodes, links } = processGraphData(selectedData);
        renderGraph(nodes, links, newSelectedNodes);
    };

    const handleResetSelection = () => {
        setSelectedNodes(new Set());
        if (selectedNodes.size === 0) return; // Nothing to reset
        renderGraph(allNodes, allLinks, new Set());
    };

    return (
        <div style={{ width: '100%', height: '100%', position: 'relative' }}>
            <div style={{ display: 'flex', gap: '10px', alignItems: 'center', padding: '10px' }}>
                <Select
                    multiple
                    aria-label="Select bindings"
                    placeholder="Select bindings"
                    selectedKeys={Array.from(selectedNodes)}
                    onSelectionChange={handleSelectionChange}
                    selectionMode="multiple"
                    className="max-w-xs"
                >
                    {allNodes.filter(node => node.kind === 'ClusterRoleBinding' || node.kind === 'RoleBinding').map(node => (
                        <SelectItem key={node.id} value={node.id}>{node.label}</SelectItem>
                    ))}
                </Select>
                <Button color="secondary" onClick={handleResetSelection}>Reset</Button>
            </div>
            <svg ref={svgRef} style={{ width: '100%', height: '100%' }}></svg>
            {hoveredNode && (
                <Tooltip node={hoveredNode} isDarkMode={isDarkMode} />
            )}
        </div>
    );
};

export default DisjointGraph;