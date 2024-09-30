import React from 'react';
import './App.css';
import Search from "./components/search/Search"
import ResultTable from "./components/ResultTable";

export interface HIT {
    TFIDF: number
    URL: string
}

function App() {
    const [searchTerm, setSearchTerm] = React.useState('');
    const [results, setResults] = React.useState<HIT[]>([]);


    return (
        <html lang="en">
        <body className="flex flex-col items-center justify-center bg-gray-100 pt-32">
            <div className="text-center p-8 bg-white shadow-md rounded-md max-w-md w-full">
                <h1 className="text-2xl font-bold mb-4 text-gray-800">Welcome to my search engine!</h1>
                <h2 className="text-lg mb-6 text-gray-600">Search for a term</h2>
                <Search
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                    setResults={setResults}
                />
            </div>
            {results !== undefined && results.length > 0 && (
                <ResultTable hits={results} />
            )}
            </body>
        </html>
    );
}

export default App;
