import MainTable from "@/components/table";
import {Card, CardBody, CardHeader} from "@heroui/card";
import DisjointGraph from "@/components/graph";

export default function Home() {
	return (
		<section>
			<div>
				<MainTable/>
			</div>
			<br/>
			<div style={{width: '100%', height: '100vh'}}>
				<Card className="p-4" style={{height: '100%'}}>
					<CardHeader>
						<h2>RBAC Map</h2>
					</CardHeader>
					<CardBody style={{height: '100%', padding: 0}}>
						<DisjointGraph/>
					</CardBody>
				</Card>
			</div>
			<br/>
		</section>
	);
}
