import React, { useEffect, useState } from 'react';
import './App.css';

import DataViz, {
    VizType,
    BarChartDatum,
} from 'react-fast-charts';


function App() {

    const [dataInFormat, setDataInFormat] = useState<BarChartDatum[][]>([])

    useEffect(() => {

        fetch('http://localhost:8080/ws')
        const ws = new WebSocket("ws://localhost:8081/ws-consume");

        ws.addEventListener("message", event => {
            const newData = JSON.parse(event.data) as number[];
            const tempData: BarChartDatum[] = newData.map((val, idx) => ({
                x: idx.toString(),
                y: val,
                fill: "lightblue"
            }));
            setDataInFormat([tempData]);
        })


        return () => {
            ws.close(); // Cleanup on unmount
        };
    }, [])

    return (
        <div className="App">
            <header className="App-header">
                <DataViz
                    id={'example-bar-chart'}
                    vizType={VizType.BarChart}
                    data={dataInFormat}
                    hideAxis={{ left: true, bottom: true }}
                />
            </header>
        </div>
    );
}

export default App;
