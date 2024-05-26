import React, { useEffect, useRef, useState, useCallback } from 'react';
import * as d3 from 'd3';
import debounce from 'lodash.debounce';
import { useTheme } from 'next-themes';

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

interface GraphData {
    nodes: Node[];
    links: Link[];
}

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

const WhatIfGraph = ({ data }: { data: GraphData }) => {
    const svgRef = useRef<SVGSVGElement | null>(null);
    const { theme } = useTheme();
    const isDarkMode = theme === 'dark';
    const [hoveredNode, setHoveredNode] = useState<Node | null>(null);

    const renderGraph = useCallback((nodes: Node[], links: Link[]) => {
        const svg = d3.select(svgRef.current);
        svg.selectAll("*").remove(); // Clear previous elements

        const width = svg.node()?.clientWidth || 800;
        const height = svg.node()?.clientHeight || 600;

        const g = svg.append('g');

        const zoom = d3.zoom<SVGSVGElement, unknown>()
            .scaleExtent([0.1, 2])
            .on('zoom', (event) => {
                g.attr('transform', event.transform.toString());
            });

        svg.call(zoom as any);

        const simulation = d3.forceSimulation<Node>(nodes)
            .force('link', d3.forceLink<Node, Link>(links).id((d: any) => d.id).distance(100))
            .force('charge', d3.forceManyBody().strength(-100))
            .force('center', d3.forceCenter(width / 2, height / 2));

        const link = g.selectAll('.link')
            .data(links)
            .enter().append('line')
            .attr('class', 'link')
            .attr('stroke', '#666')
            .attr('stroke-width', 2);

        const node = g.selectAll('.node')
            .data(nodes)
            .enter().append('circle')
            .attr('class', 'node')
            .attr('r', 10)
            .attr('fill', d => d.kind === 'ClusterRoleBinding' ? 'orange' : d.kind === 'RoleBinding' ? 'green' : 'pink')
            .call(d3.drag<SVGCircleElement, Node>()
                .on('start', (event, d) => {
                    if (!event.active) simulation.alphaTarget(0.3).restart();
                    d.fx = d.x;
                    d.fy = d.y;
                })
                .on('drag', (event, d) => {
                    d.fx = event.x;
                    d.fy = event.y;
                })
                .on('end', (event, d) => {
                    if (!event.active) simulation.alphaTarget(0);
                    d.fx = null;
                    d.fy = null;
                }))
            .on('mouseover', debounce((_event, d) => setHoveredNode(d), 50))
            .on('mouseout', debounce(() => setHoveredNode(null), 50));

        const text = g.selectAll('.label')
            .data(nodes)
            .enter().append('text')
            .attr('class', 'label')
            .attr('x', d => d.x!)
            .attr('y', d => d.y!)
            .text(d => d.label)
            .style('font-size', '13px')
            .style('fill', isDarkMode ? '#FFF' : '#000')
            .style('text-anchor', 'middle')
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
        if (data && Array.isArray(data.nodes) && Array.isArray(data.links)) {
            const { nodes, links } = data;
            renderGraph(nodes, links);
        }
    }, [data, renderGraph]);

    return (
        <div style={{ position: 'relative', height: '100%' }}>
            <svg ref={svgRef} style={{ width: '100%', height: '100%' }}></svg>
            {hoveredNode && <Tooltip node={hoveredNode} isDarkMode={isDarkMode} />}
        </div>
    );
};

export default WhatIfGraph;
