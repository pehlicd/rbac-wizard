"use client";
import { useEffect, useRef, useState } from 'react';
import * as d3 from 'd3';

interface Node extends d3.SimulationNodeDatum {
    id: string;
    kind?: string;
    label: string;
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

const DisjointGraph = () => {
    const svgRef = useRef<SVGSVGElement | null>(null);
    const [hoveredNode, setHoveredNode] = useState<Node | null>(null);
    const [isDarkMode, setIsDarkMode] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/data');
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: BindingData[] = await response.json();
                if (!data) {
                    throw new Error('Data is null or undefined');
                }
                console.log('Fetched data:', data);
                renderGraph(data);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        const renderGraph = (data: BindingData[]) => {
            if (!data || !Array.isArray(data)) {
                console.error('Invalid data format:', data);
                return;
            }

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

            const svg = d3.select(svgRef.current);
            svg.selectAll("*").remove(); // Clear previous elements

            const width = svg.node()?.clientWidth || 800;
            const height = svg.node()?.clientHeight || 600;

            const g = svg.append('g');

            const zoom = d3.zoom<SVGSVGElement, unknown>().on('zoom', (event) => {
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
                .attr('fill', d => d.kind === 'ClusterRoleBinding' ? 'orange' : 'green')
                .call(drag(simulation) as any)
                .on('mouseover', (_event, d) => setHoveredNode(d))
                .on('mouseout', () => setHoveredNode(null));

            const text = g.selectAll('.label')
                .data(nodes)
                .enter().append('text')
                .attr('class', 'label')
                .attr('x', d => d.x!)
                .attr('y', d => d.y!)
                .text(d => d.label)
                .style('font-size', '12px')
                .style('fill', '#4E46C1')
                .style('text-anchor', '-moz-initial');

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

            return () => {
                simulation.stop();
            };
        };

        fetchData();
    }, [isDarkMode]);

    useEffect(() => {
        const handleThemeChange = (e: MediaQueryListEvent) => {
            setIsDarkMode(e.matches);
        };

        const darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        darkModeMediaQuery.addEventListener('change', handleThemeChange);

        // Set initial dark mode state
        setIsDarkMode(darkModeMediaQuery.matches);

        return () => {
            darkModeMediaQuery.removeEventListener('change', handleThemeChange);
        };
    }, []);

    // Additional effect to update text color on theme change
    useEffect(() => {
        const text = d3.selectAll('.label');
        text.style('fill', isDarkMode ? '#FFF' : '#000');
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

    return (
        <div style={{ width: '100%', height: '100%' }}>
            <svg ref={svgRef} style={{ width: '100%', height: '100%' }}></svg>
            {hoveredNode && (
                <div style={{
                    position: 'absolute',
                    left: hoveredNode.x,
                    top: hoveredNode.y,
                    backgroundColor: isDarkMode ? '#333' : 'white',
                    color: isDarkMode ? 'white' : 'black',
                    padding: '2px 5px',
                    border: `1px solid ${isDarkMode ? '#555' : '#ccc'}`,
                    borderRadius: '3px'
                }}>
                    {hoveredNode.label}
                </div>
            )}
        </div>
    );
};

export default DisjointGraph;
