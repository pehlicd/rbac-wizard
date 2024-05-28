"use client";
import Editor from '@monaco-editor/react';
import { Card, CardBody, CardHeader } from "@nextui-org/card";
import { Badge, Button } from "@nextui-org/react";
import { useTheme } from 'next-themes';
import { useState } from 'react';
import axios from 'axios';
import DisjointGraph from '@/components/graph';

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
                        <CardHeader>
                            <h2>Editor</h2>
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
