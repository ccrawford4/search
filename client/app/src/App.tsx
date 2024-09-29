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
      <html lang="english">
      <body>
          <h1>Welcome to my search engine!</h1>
          <h2>Search for a term</h2>
          <Search
              searchTerm={searchTerm}
              setSearchTerm={setSearchTerm}
              setResults={setResults}
          />
          {results.length > 0 && (
              <ResultTable hits={results}/>
          )}
      </body>
      </html>
    );
}

export default App;
