import {HIT} from "../App";

interface ResultTableProps {
    hits: HIT[]
}

export default function ResultTable(props: ResultTableProps) {
    return (
        <table>
            <thead>
            <tr>
                <th>URL</th>
                <th>TFIDF</th>
            </tr>
            </thead>
            <tbody>
            {props.hits.map(((hit, index) => (
                <tr key={index}>
                    <td><a className="text-blue" href={hit.URL}>{hit.URL}</a></td>
                    <td>{hit.TFIDF}</td>
                </tr>
            )))}
            </tbody>
        </table>
    );
}