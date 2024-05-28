"use client";
import Editor from '@monaco-editor/react';
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Badge, Button } from "@nextui-org/react";
import { useTheme } from 'next-themes';
import { useState } from 'react';
import axios from 'axios';
import DisjointGraph from '@/components/graph';
import { IoInformationCircle } from "react-icons/io5";
import {Tooltip} from "@nextui-org/tooltip";

export default function WhatIfPage() {
    const { theme } = useTheme();
    const isDarkMode = theme === 'dark';
    const editorTheme = isDarkMode ? 'vs-dark' : 'light';

    const [yamlContent, setYamlContent] = useState('');
    const [graphData, setGraphData] = useState<{ nodes: any[]; links: any[] } | null>(null);

    const handleEditorChange = (value: string | undefined) => {
        setYamlContent(value || '');
    };

    const handleGenerateClick = async () => {
        try {
            const response = await axios.post('/api/what-if', { yaml: yamlContent });
            setGraphData(response.data);
        } catch (error) {
            console.error('Error generating graph:', error);
        }
    };

    return (
        <section style={{ position: 'relative', width: '100%', height: '100vh', display: 'flex', flexDirection: 'column' }}>
            <div style={{ display: 'flex', flexDirection: 'row', flexGrow: 1 }}>
                <Card className="p-4" style={{ flex: 1, marginRight: '5px' }}>
                    <Badge content="Beta" color="warning" size="md" className="absolute">
                        <CardHeader className="p-2">
                                <h2>Editor</h2>
                                <Tooltip
                                    color="secondary"
                                    delay={1000}
                                    content={
                                        <div className="px-1 py-2">
                                            <div className="text-large font-bold">What is `What If?`</div>
                                            <div className="text-small">`What If?` helps you easily add one of your ClusterRoleBinding or RoleBinding manifests. When you click the `Generate` button, it visualizes your binding in a map format.</div>
                                            <br />
                                            <div className="text-tiny">⚠️Please note that this feature is still in beta. If you encounter any issues, please report them on our GitHub page.</div>
                                        </div>
                                    }
                                >
                                    <Button isIconOnly color="secondary" variant="shadow" size="sm" className="ml-2 rounded-b-small">
                                        <IoInformationCircle size={16} />
                                    </Button>
                                </Tooltip>
                        </CardHeader>
                    </Badge>
                    <CardBody style={{ height: '100%', padding: 10 }}>
                        <Editor
                            height="100%"
                            defaultLanguage="yaml"
                            defaultValue={``}
                            className="p-1 rounded-b-lg"
                            theme={editorTheme}
                            onChange={handleEditorChange}
                        />
                    </CardBody>
                </Card>
                <Card className="p-4" style={{ flex: 1, marginLeft: '5px' }}>
                    <CardHeader>
                        <h2>Graph</h2>
                    </CardHeader>
                    <CardBody style={{ height: '100%', padding: 10 }}>
                        {graphData && graphData.nodes && graphData.links && (
                            <DisjointGraph data={graphData} />
                        )}
                    </CardBody>
                </Card>
            </div>
            <Button onClick={handleGenerateClick} color="primary" variant="shadow" style={{ marginTop: '10px', marginBottom: '10px', width: '100%' }}>Generate</Button>
        </section>
    );
}
